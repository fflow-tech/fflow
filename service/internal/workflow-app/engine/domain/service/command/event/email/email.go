package email

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"log"
	"mime/quotedprintable"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

const maxReadEmailNum = 10

// ImapServer 存储邮箱服务器信息
type ImapServer struct {
	Server   string
	Port     int
	Username string
	Password string
}

// Event 存储邮件事件信息
type Event struct {
	Subject string
	From    string
	Date    time.Time
	Content string
}

// FetchEventsFromEmail 根据指定的邮箱和关键字，收取邮件并生成事件列表
func FetchEventsFromEmail(server ImapServer, keyword string) ([]Event, error) {
	var events []Event

	// 连接到服务器
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", server.Server, server.Port), nil)
	if err != nil {
		return events, err
	}
	defer c.Logout()

	// 登录邮箱
	if err := c.Login(server.Username, server.Password); err != nil {
		return events, err
	}

	// 选择邮箱文件夹
	if _, err := c.Select("INBOX", false); err != nil {
		return events, err
	}

	// 设置搜索条件
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag} // 只搜索未读邮件
	criteria.Header.Set("Subject", keyword)         // 根据关键字搜索邮件的主题

	// 执行搜索
	uids, err := c.Search(criteria)
	if err != nil {
		return events, err
	}

	if len(uids) == 0 {
		return events, err
	}

	// 获取邮件内容
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)
	// 最多读取 10 封相关的邮件
	messages := make(chan *imap.Message, maxReadEmailNum)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, imap.FetchRFC822Text}, messages)
	}()

	// 处理收到的邮件
	for msg := range messages {
		event, err := parseEventFromMessage(msg)
		if err != nil {
			log.Println("Failed to parse message:", err)
			continue
		}
		events = append(events, event)

		// 标记邮件为已读
		if err := markMessageAsSeen(c, seqSet); err != nil {
			log.Println("Failed to mark message as seen:", err)
		}
	}

	if err := <-done; err != nil {
		return events, err
	}

	return events, nil
}

func markMessageAsSeen(c *client.Client, seqSet *imap.SeqSet) error {
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err := c.Store(seqSet, item, flags, nil)
	if err != nil {
		return fmt.Errorf("failed to mark message as seen: %w", err)
	}
	return nil
}

func parseEventFromMessage(msg *imap.Message) (Event, error) {
	event := Event{
		Subject: msg.Envelope.Subject,
		From:    msg.Envelope.From[0].Address(),
		Date:    msg.Envelope.Date,
	}

	content, err := decodeEmailContent(msg)
	if err != nil {
		log.Println("Failed to decode email content:", err)
	}
	event.Content = content

	return event, nil
}

func decodeEmailContent(msg *imap.Message) (string, error) {
	r := msg.GetBody(&imap.BodySectionName{BodyPartName: imap.BodyPartName{Specifier: imap.TextSpecifier}})
	if r == nil {
		return "", nil
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return "", fmt.Errorf("failed to read message body: %w", err)
	}

	rawEmail := buf.String()
	return decodeEmailContentFromString(rawEmail)
}

func decodeEmailContentFromString(rawEmail string) (string, error) {
	r := strings.NewReader(rawEmail)
	mr, err := mail.CreateReader(r)
	if err != nil {
		lines := strings.SplitN(rawEmail, "\n", 2)
		if len(lines) != 2 {
			return "", fmt.Errorf("invalid email content")
		}
		rawEmail = lines[1]
		return decodeEmailContentFromString(rawEmail)
	}

	var decodedContent strings.Builder
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("failed to read mail part: %w", err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			b, err := io.ReadAll(p.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read mail part body: %w", err)
			}

			encoding := h.Get("Content-Transfer-Encoding")
			switch encoding {
			case "quoted-printable":
				decoded, err := ioutil.ReadAll(quotedprintable.NewReader(strings.NewReader(string(b))))
				if err != nil {
					return "", fmt.Errorf("failed to decode quoted-printable: %w", err)
				}
				decodedContent.Write(decoded)
			default:
				decodedContent.Write(b)
			}
		case *mail.AttachmentHeader:
			// Ignore attachments
		default:
			// Ignore unknown headers
		}
	}

	return decodedContent.String(), nil
}
