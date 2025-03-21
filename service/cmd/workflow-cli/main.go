package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/fflow-tech/fflow/service/cmd/workflow-cli/factory"
	"github.com/fflow-tech/fflow/service/cmd/workflow-cli/service"
	"github.com/fflow-tech/fflow/service/cmd/workflow-cli/service/event"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/dto"
	"github.com/fflow-tech/fflow/service/internal/workflow-app/engine/domain/entity"
	"github.com/fflow-tech/fflow/service/pkg/config"
	"github.com/fflow-tech/fflow/service/pkg/k8s"
	"github.com/fflow-tech/fflow/service/pkg/log"
	"github.com/fflow-tech/fflow/service/pkg/registry"
	"github.com/fflow-tech/fflow/service/pkg/utils"
)

var (
	globalConfigName = flag.String("config-name", "app", "The global config name")
	globalConfigType = flag.String("config-type", "yaml", "The global config type")
	globalConfigPath = flag.String("config-dir", ".fflow/", "The global config path")
	definitionPath   = flag.String("def-dir", ".fflow/definitions", "Workflow definition history directory")
	instancePath     = flag.String("inst-dir", ".fflow/instances", "Workflow instance history directory")
	workflowFile     = flag.String("f", "", "Workflow definition file path, e.g. examples/example-http.json")
	inputFile        = flag.String("i", "", "Workflow input file path, e.g. examples/example-http-input.json")
	showHelp         = flag.Bool("h", false, "Show help information")
)

func main() {
	flag.Parse()

	// 显示帮助信息
	if *showHelp {
		printHelp()
		return
	}

	// 初始化环境
	if err := initializeEnvironment(); err != nil {
		log.Fatalf("Environment initialization failed: %v", err)
		panic(err)
	}

	// 初始化工作流服务
	workflowService, err := initializeWorkflowService()
	if err != nil {
		log.Fatalf("Initialize workflow service failed: %v", err)
		panic(err)
	}

	// 执行工作流
	instId, err := executeWorkflow(workflowService)
	if err != nil {
		log.Fatalf("Failed to execute workflow: %v", err)
		panic(err)
	}

	// 监控工作流执行
	monitorWorkflow(workflowService, instId)
}

// 初始化环境：工厂、数据库、事件服务器和目录
func initializeEnvironment() error {
	// 命令行模式下先创建一个空的 app.yaml 文件，如果存在这个文件则忽略
	configFilePath := fmt.Sprintf("%s/%s.%s", *globalConfigPath, *globalConfigName, *globalConfigType)

	// 检查文件是否已存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// 确保目录存在
		if err := ensureDir(*globalConfigPath); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// 创建空文件
		if err := os.WriteFile(configFilePath, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to create app.yaml file: %w", err)
		}
		log.Infof("Created empty config file: %s", configFilePath)
	}

	// 初始化本地工厂
	if err := factory.New(factory.WithRegistryClientType(registry.Kubernetes),
		factory.WithConfigClientType(config.Kubernetes),
		factory.WithK8sConfig(k8s.Config{
			GlobalConfigName: *globalConfigName,
			GlobalConfigType: *globalConfigType,
			GlobalConfigPath: *globalConfigPath,
		}),
	); err != nil {
		return fmt.Errorf("factory init failed: %w", err)
	}

	// 创建数据库表
	if err := factory.CreateTables(); err != nil {
		return fmt.Errorf("create tables failed: %w", err)
	}

	// 启动事件服务器
	eventServer := event.NewServer()
	if err := eventServer.Serve(); err != nil {
		return fmt.Errorf("event server not serve: %w", err)
	}

	// 创建工作流目录
	if err := ensureDir(*definitionPath); err != nil {
		return fmt.Errorf("failed to ensure definition directory: %w", err)
	}
	if err := ensureDir(*instancePath); err != nil {
		return fmt.Errorf("failed to ensure instance directory: %w", err)
	}

	return nil
}

// 初始化工作流服务
func initializeWorkflowService() (*service.WorkflowService, error) {
	return service.NewWorkflowService(*definitionPath, *instancePath)
}

// 执行工作流
func executeWorkflow(workflowService *service.WorkflowService) (string, error) {
	// 检查工作流文件
	if *workflowFile == "" {
		return "", fmt.Errorf("workflow file is required")
	}

	// 复制工作流定义文件到定义目录
	defJson, err := copyWorkflowFile(*workflowFile, *definitionPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy workflow definition file: %w", err)
	}

	// 读取输入文件
	inputMap, err := readInputFile(*inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to process input file: %w", err)
	}

	// 执行工作流
	instId, err := workflowService.ExecuteWorkflow(defJson, inputMap)
	if err != nil {
		return "", fmt.Errorf("failed to execute workflow: %w", err)
	}

	return instId, nil
}

// 读取输入文件并转换为map
func readInputFile(filePath string) (map[string]interface{}, error) {
	input, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	var inputMap map[string]interface{}
	if err := json.Unmarshal(input, &inputMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input file: %w", err)
	}

	return inputMap, nil
}

// 监控工作流执行
func monitorWorkflow(workflowService *service.WorkflowService, instId string) {
	// 等待中断信号来优雅地关闭服务
	quit := make(chan os.Signal, 1)

	go monitorWorkflowStatus(workflowService, instId, quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	// 打印提示信息
	fmt.Println("\nWorkflow executor has started, processing workflow tasks...")
	fmt.Println("Press Ctrl+C to exit fflow-cli")

	<-quit
}

// 监控工作流状态
func monitorWorkflowStatus(workflowService *service.WorkflowService, instId string, quit chan os.Signal) {
	for {
		time.Sleep(3 * time.Second)
		inst, err := workflowService.GetWorkflowStatus(instId)
		if err != nil {
			log.Errorf("Failed to get workflow status: %v", err)
			return
		}

		// 打印工作流实例状态
		fmt.Printf("Workflow instance execute path: %v\n", inst.ExecutePath)
		fmt.Printf("Workflow instance status: %v\n", inst.Status)

		// 如果工作流实例状态为完成，则退出
		if inst.Status == entity.InstSucceed || inst.Status == entity.InstFailed {
			saveWorkflowInstance(inst)
			close(quit)
		}
	}
}

// 保存工作流实例
func saveWorkflowInstance(inst *dto.WorkflowInstDTO) {
	workflowFileName := strings.TrimSuffix(filepath.Base(*workflowFile), filepath.Ext(filepath.Base(*workflowFile)))
	instanceFileName := fmt.Sprintf("%s_%s.json", workflowFileName, inst.InstID)
	instanceFilePath := filepath.Join(*instancePath, instanceFileName)

	instanceData, err := json.MarshalIndent(inst, "", "  ")
	if err != nil {
		log.Errorf("Failed to marshal instance data: %v", err)
		return
	}

	if err := os.WriteFile(instanceFilePath, instanceData, 0644); err != nil {
		log.Errorf("Failed to write instance file: %v", err)
		return
	}

	log.Infof("Instance file saved successfully: %s", instanceFilePath)
}

func copyWorkflowFile(srcPath, destDir string) (string, error) {
	// 读取源文件
	data, err := readDefinitionFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("Failed to read workflow definition file: %w", err)
	}

	destPath := filepath.Join(destDir, getDstFileName(srcPath))
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return "", fmt.Errorf("Failed to write workflow definition file: %w", err)
	}

	log.Infof("Finished copying workflow definition file: %s\n", destPath)

	return utils.BytesToJsonStr(data), nil
}

func readDefinitionFile(srcPath string) ([]byte, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read source file: %w", err)
	}

	// 检查文件扩展名，如果是yaml或yml，转换为json
	ext := strings.ToLower(filepath.Ext(srcPath))
	if ext == ".yaml" || ext == ".yml" {
		var yamlObj interface{}
		if err := yaml.Unmarshal(data, &yamlObj); err != nil {
			return nil, fmt.Errorf("Failed to parse YAML file: %w", err)
		}

		jsonData, err := json.Marshal(yamlObj)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert YAML to JSON: %w", err)
		}

		// 更新数据为JSON格式
		return jsonData, nil
	}

	return data, nil
}

func getDstFileName(srcPath string) string {
	fileName := filepath.Base(srcPath)
	fileName = fmt.Sprintf("%s_%s%s", strings.TrimSuffix(fileName, filepath.Ext(fileName)), utils.GetCurrentTimestamp(), filepath.Ext(fileName))
	return fileName
}

// 确保目录存在，如果不存在则创建
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// 优雅关闭服务
func shutdownGraceful(fs ...func(chan struct{}) error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f func(chan struct{}) error) {
			defer wg.Done()

			c := make(chan struct{}, 1)
			go f(c)
			select {
			case <-c:
			case <-ctx.Done():
			}
		}(f)
	}

	wg.Wait()
	fmt.Println("Service has been closed")
}

// 打印帮助信息
func printHelp() {
	fmt.Println("FFlow Workflow CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  fflow-cli [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  fflow-cli -f examples/example-http.json -i examples/example-http-input.json")
	fmt.Println("  fflow-cli -f examples/example-http.yaml -i examples/example-http-input.json")
}
