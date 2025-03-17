package p

import (
	"encoding/json"

	"github.com/fflow-tech/fflow-sdk-go/faas"
	"github.com/go-resty/resty/v2"
)

type RobotMessage struct {
	Msgtype  string               `json:"msgtype"`
	Text     *TextMessageText     `json:"text"`
	Markdown *MarkdownMessageText `json:"markdown"`
}

type TextMessageText struct {
	Content string `json:"content"`
}

type MarkdownMessageText struct {
	Content string `json:"content"`
}

type SendMsgReq struct {
	APIKey  string `json:"apiKey"`
	Content string `json:"content"`
}

// fucntion name(handler) DO NOT MODIFY
// @ctx[context.RuntimeContext]: Some inner functions like GetLogger is on it
// @input[map[string]interface{}]: The input of the function
// @return [interface{}, error]: The return value should be allowed to convert to json
func handler(ctx faas.Context, input map[string]interface{}) (interface{}, error) {
	ctx.Logger().Infof("input: %v", input)
	bytes, err := json.Marshal(input)
	if err != nil {
		return err.Error(), nil
	}
	var req SendMsgReq
	if err := json.Unmarshal(bytes, &req); err != nil {
		return err.Error(), nil
	}

	ctx.Logger().Infof("API Key: %v", req.APIKey)
	ctx.Logger().Infof("Content: %v", req.Content)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&RobotMessage{
			Msgtype: "markdown",
			Markdown: &MarkdownMessageText{
				Content: req.Content,
			},
		}).
		Post("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + req.APIKey)
	if err != nil {
		return err.Error(), nil
	}

	return resp.String(), nil
}
