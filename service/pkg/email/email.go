package email

import (
	"github.com/fflow-tech/fflow/service/pkg/log"
	"gopkg.in/gomail.v2"
	"regexp"
)

// Email 邮件
type Email struct {
	To      []string // 收件人邮箱地址列表
	Cc      []string // 抄送人邮箱地址列表
	Subject string   // 邮件主题
	Body    string   // 邮件正文
	Attach  string   // 邮件附件
}

// SMTPServer 邮箱服务器
type SMTPServer struct {
	Host     string // SMTP服务器地址
	Port     int    // SMTP服务器端口
	From     string // 发件人邮箱地址
	Password string // 发件人邮箱密码
}

// Send 发送接口
func (e *Email) Send(smtpServer SMTPServer) error {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpServer.From)
	m.SetHeader("To", e.To...)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/html", e.Body)

	for _, c := range e.Cc {
		m.SetAddressHeader("Cc", c, "")
	}

	if len(e.Attach) > 0 {
		m.Attach(e.Attach)
	}

	// 拿到 token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer(smtpServer.Host, smtpServer.Port, smtpServer.From, smtpServer.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Errorf("DialAndSend err %v:", err)
		return err
	}

	return nil
}

// IsEmailValid 校验邮件地址
func IsEmailValid(email string) bool {
	if len(email) <= 0 {
		return false
	}

	// 邮件地址的正则表达式
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, _ := regexp.MatchString(regex, email)
	return match
}
