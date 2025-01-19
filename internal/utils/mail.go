package utils

import (
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/wt3022/github-release-notifier/internal/env"
)

type EmailRequest struct {
	To      string
	Subject string
	Body    string
}

func SendEmail(req EmailRequest, env env.Config) error {
	fmt.Println("run SendEmail")
	headers := make(map[string]string)
	headers["From"] = env.EmailConfig.FromEmail
	headers["To"] = req.To
	headers["Subject"] = req.Subject

	// メッセージ構築
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + req.Body

	if env.ProductionMode {
		fmt.Println("ProductionMode")
		// メール送信
		addr := fmt.Sprintf("%s:%s", env.EmailConfig.SMTPHost, env.EmailConfig.SMTPPort)
		d := net.Dialer{
			Timeout: 10 * time.Second,
		}

		conn, err := d.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("connection error: %v", err)
		}
		defer conn.Close()

		c, err := smtp.NewClient(conn, env.EmailConfig.SMTPHost)
		if err != nil {
			return fmt.Errorf("SMTP client error: %v", err)
		}
		defer c.Close()

		if err := c.Mail(env.EmailConfig.FromEmail); err != nil {
			return fmt.Errorf("MAIL FROM error: %v", err)
		}

		if err := c.Rcpt(req.To); err != nil {
			return fmt.Errorf("RCPT TO error: %v", err)
		}

		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("DATA error: %v", err)
		}
		defer w.Close()

		if _, err := w.Write([]byte(message)); err != nil {
			return fmt.Errorf("write error: %v", err)
		}

		fmt.Println("Email sent successfully")
		return nil
	}

	// テストモードの場合はメールを出力する
	fmt.Println("================== 送信メール ==================")
	fmt.Printf("To: %s\n", req.To)
	fmt.Printf("Subject: %s\n", req.Subject)
	fmt.Printf("Body: %s\n", req.Body)
	fmt.Println("===============================================")

	return nil
}
