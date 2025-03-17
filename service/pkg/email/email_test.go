package email

import "testing"

func TestEmail_Send(t *testing.T) {
	type fields struct {
		To      []string
		Cc      []string
		Subject string
		Body    string
		Attach  string
	}
	type args struct {
		smtpServer SMTPServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Case1", fields{
			To:      []string{"fflow-tech@gmail.com"},
			Cc:      []string{"zhanghuan2@linklogis.com"},
			Subject: "测试主题",
			Body:    "测试内容",
		}, args{smtpServer: SMTPServer{
			Host:     "smtp.qq.com",
			Port:     465,
			From:     "382295014@qq.com",
			Password: "pwvkexqarhpvbhja",
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Email{
				To:      tt.fields.To,
				Cc:      tt.fields.Cc,
				Subject: tt.fields.Subject,
				Body:    tt.fields.Body,
				Attach:  tt.fields.Attach,
			}
			if err := e.Send(tt.args.smtpServer); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
