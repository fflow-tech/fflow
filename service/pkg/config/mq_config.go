package config

// TDMQConfig tdmq 的配置
type TDMQConfig struct {
	URL                 string `json:"url"`                                  // URL 接入点地址
	Authentication      string `yaml:"authentication" json:"authentication"` // Authentication 密钥
	NackRedeliveryDelay int    `json:"nackRedeliveryDelay"`                  // 重投时间间隔
	MaxDeliveries       uint32 `json:"maxDeliveries"`                        // 最大重试次数
	RetryEnable         bool   `json:"retryEnable"`                          // 是否开启重试
	RetryInitialDelay   uint64 `json:"retryInitialDelay"`                    // 重试的基础延迟时间, 以秒为单位
	RetryMaxDelay       uint64 `json:"retryMaxDelay"`                        // 重试的最大延迟时间, 以秒为单位
}

// KafkaConfig kafka 的配置
type KafkaConfig struct {
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	Network              string `json:"network"`
	WriteTimeout         int    `json:"writeTimeout"`
	ReadTimeout          int    `json:"readTimeout"`
	NumPartitions        int    `json:"numPartitions"`        // 消息分区的数量
	ReplicationFactor    int    `json:"replicationFactor"`    // 消息副本的数量
	ConsumerMaxWaitTime  int    `json:"consumerMaxWaitTime"`  // 消费者消费时的最长等待时间, 以 ms 为单位
	ProducerBatchTimeout int    `json:"producerBatchTimeout"` // 生产者发送消息时的最少等待时间，以 ms 为单位
}
