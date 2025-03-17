package email

import (
	"fmt"
	"log"
	"testing"
)

func TestFetchEventsFromEmail(t *testing.T) {
	// 邮箱服务器信息
	server := ImapServer{
		Username: "382295014@qq.com",
		Password: "grpmcasnwgszbihe",
		Server:   "imap.qq.com",
		Port:     993,
	}

	// 关键字
	keyword := "FFlow"

	// 获取事件列表
	events, err := FetchEventsFromEmail(server, keyword)
	if err != nil {
		log.Fatal(err)
	}

	// 打印事件列表
	for _, event := range events {
		fmt.Printf("Subject: %s\nFrom: %s\nDate: %s\nContent: %s\n", event.Subject, event.From, event.Date, event.Content)
	}
}
