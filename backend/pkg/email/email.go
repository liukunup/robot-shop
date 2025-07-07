package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Email 简化版邮件结构体
type Email struct {
	From    string   // 发件人地址
	To      []string // 收件人地址列表
	Subject string   // 邮件主题
	HTML    string   // HTML内容
}

// Config SMTP服务器配置
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	UseSSL   bool // 是否使用SSL
}

// Client 邮件客户端
type Client struct {
	config *Config
}

// NewClient 创建新的邮件客户端
func NewClient(conf *viper.Viper) *Client {
	return &Client{
		config: &Config{
			Host:     conf.GetString("email.host"),
			Port:     conf.GetInt("email.port"),
			Username: conf.GetString("email.username"),
			Password: conf.GetString("email.password"),
			UseSSL:   conf.GetBool("email.use_ssl"),
		},
	}
}

// Send 发送邮件
func (c *Client) Send(email *Email) error {
	if len(email.To) == 0 {
		return fmt.Errorf("no recipient specified")
	}

	// 构建邮件内容
	raw, err := c.buildEmail(email)
	if err != nil {
		return err
	}

	// 发送邮件
	return c.send(email.From, email.To, raw)
}

// buildEmail 构建邮件内容
func (c *Client) buildEmail(email *Email) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// 设置头部
	headers := make(textproto.MIMEHeader)
	headers.Set("From", email.From)
	headers.Set("To", strings.Join(email.To, ", "))
	headers.Set("Subject", email.Subject)
	headers.Set("Date", time.Now().Format(time.RFC1123Z))
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", "text/html; charset=UTF-8")
	headers.Set("Content-Transfer-Encoding", "quoted-printable")

	// 写入头部
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ", ")))
	}
	buf.WriteString("\r\n")

	// 写入HTML内容
	buf.WriteString(email.HTML)
	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

// send 实际发送邮件
func (c *Client) send(from string, to []string, raw []byte) error {
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// 创建认证
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	// 如果是SSL直接发送
	if c.config.UseSSL {
		tlsconfig := &tls.Config{
			ServerName: c.config.Host,
		}
		return smtp.SendMail(addr, auth, from, to, raw)
	}

	// 非SSL连接
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial smtp server: %w", err)
	}
	defer client.Close()

	// 尝试STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsconfig := &tls.Config{
			ServerName: c.config.Host,
		}
		if err = client.StartTLS(tlsconfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	}

	// 认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("mail from failed: %w", err)
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("rcpt to failed: %w", err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data failed: %w", err)
	}

	if _, err = w.Write(raw); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("close failed: %w", err)
	}

	return client.Quit()
}
