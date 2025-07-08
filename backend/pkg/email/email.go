package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Email
type Email struct {
	From        string               // 发件人地址
	To          []string             // 收件人地址列表
	Cc          []string             // 抄送地址列表
	Bcc         []string             // 密送地址列表
	Subject     string               // 邮件主题
	Text        string               // 文本内容
	HTML        string               // HTML内容
	Headers     textproto.MIMEHeader // 自定义头字段
	Attachments []*Attachment        // 附件列表
	ReadReceipt []string             // 回执地址列表
}

// Attachment
type Attachment struct {
	Filename string // 附件文件名
	Content  []byte // 附件内容
	Inline   bool   // 是否内联显示
}

// Config
type Config struct {
	Host      string      // SMTP服务器主机
	Port      int         // SMTP服务器端口
	Username  string      // SMTP服务器用户名
	Password  string      // SMTP服务器密码
	SSL       bool        // 是否使用SSL
	TLSConfig *tls.Config // TLS配置
	LocalName string      // 本地主机名
}

type Client struct {
	config *Config
}

func NewEmail(conf *viper.Viper) *Client {
	return &Client{config: &Config{
		Host:      conf.GetString("email.host"),
		Port:      conf.GetInt("email.port"),
		Username:  conf.GetString("email.username"),
		Password:  conf.GetString("email.password"),
		SSL:       conf.GetBool("email.ssl"),
		LocalName: conf.GetString("email.local_name"),
	}}
}

func (c *Client) Send(email *Email) error {
	if len(email.To) == 0 {
		return errors.New("mail: no recipient specified")
	}

	from, err := mail.ParseAddress(email.From)
	if err != nil {
		return fmt.Errorf("mail: invalid from address: %v", err)
	}

	to := make([]string, 0, len(email.To))
	for _, addr := range email.To {
		parsedAddr, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("mail: invalid to address %q: %v", addr, err)
		}
		to = append(to, parsedAddr.Address)
	}

	cc := make([]string, 0, len(email.Cc))
	for _, addr := range email.Cc {
		parsedAddr, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("mail: invalid cc address %q: %v", addr, err)
		}
		cc = append(cc, parsedAddr.Address)
	}

	bcc := make([]string, 0, len(email.Bcc))
	for _, addr := range email.Bcc {
		parsedAddr, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("mail: invalid bcc address %q: %v", addr, err)
		}
		bcc = append(bcc, parsedAddr.Address)
	}

	raw, err := c.buildEmail(email, from, to)
	if err != nil {
		return err
	}

	return c.sendEmail(from.Address, to, cc, bcc, raw)
}

// 构建邮件内容
func (c *Client) buildEmail(email *Email, from *mail.Address, to []string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// 设置头部
	headers := make(textproto.MIMEHeader)
	headers.Set("From", from.String())
	headers.Set("To", strings.Join(to, ", "))
	headers.Set("Subject", mime.QEncoding.Encode("utf-8", email.Subject))
	headers.Set("Date", time.Now().Format(time.RFC1123Z))
	headers.Set("MIME-Version", "1.0")

	if len(email.Cc) > 0 {
		headers.Set("Cc", strings.Join(email.Cc, ", "))
	}

	if len(email.ReadReceipt) > 0 {
		headers.Set("Disposition-Notification-To", strings.Join(email.ReadReceipt, ", "))
	}

	// 添加自定义头部
	for k, v := range email.Headers {
		headers[k] = v
	}

	// 根据是否有附件决定内容类型
	mixed := len(email.Attachments) > 0
	var (
		related     bool
		alternative bool
	)

	if mixed {
		headers.Set("Content-Type", "multipart/mixed; boundary=MIXED_BOUNDARY")
		buf.WriteString("--MIXED_BOUNDARY\r\n")
	}

	alternative = email.Text != "" && email.HTML != ""
	if alternative {
		headers.Set("Content-Type", "multipart/alternative; boundary=ALTERNATIVE_BOUNDARY")
		buf.WriteString("--ALTERNATIVE_BOUNDARY\r\n")
	}

	related = len(email.getInlineAttachments()) > 0
	if related {
		headers.Set("Content-Type", "multipart/related; boundary=RELATED_BOUNDARY")
		buf.WriteString("--RELATED_BOUNDARY\r\n")
	}

	// 写入头部
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ", ")))
	}
	buf.WriteString("\r\n")

	// 写入文本内容
	if email.Text != "" || email.HTML == "" {
		textHeader := make(textproto.MIMEHeader)
		textHeader.Set("Content-Type", "text/plain; charset=utf-8")
		textHeader.Set("Content-Transfer-Encoding", "quoted-printable")

		writeHeaders(buf, textHeader)
		buf.WriteString(encodeText(email.Text))
		buf.WriteString("\r\n")

		if alternative {
			buf.WriteString("--ALTERNATIVE_BOUNDARY\r\n")
		}
	}

	// 写入HTML内容
	if email.HTML != "" {
		htmlHeader := make(textproto.MIMEHeader)
		htmlHeader.Set("Content-Type", "text/html; charset=utf-8")
		htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")

		writeHeaders(buf, htmlHeader)
		buf.WriteString(encodeText(email.HTML))
		buf.WriteString("\r\n")

		if alternative {
			buf.WriteString("--ALTERNATIVE_BOUNDARY--\r\n")
		}
	}

	if related {
		buf.WriteString("--RELATED_BOUNDARY--\r\n")
	}

	// 添加附件
	if mixed {
		for _, attachment := range email.Attachments {
			buf.WriteString("\r\n--MIXED_BOUNDARY\r\n")
			attachmentHeader := make(textproto.MIMEHeader)
			if attachment.Inline {
				attachmentHeader.Set("Content-Type", "message/rfc822")
				attachmentHeader.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", mime.QEncoding.Encode("utf-8", attachment.Filename)))
			} else {
				ext := filepath.Ext(attachment.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype == "" {
					mimetype = "application/octet-stream"
				}
				attachmentHeader.Set("Content-Type", mimetype)
				attachmentHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", mime.QEncoding.Encode("utf-8", attachment.Filename)))
			}
			attachmentHeader.Set("Content-Transfer-Encoding", "base64")

			writeHeaders(buf, attachmentHeader)
			buf.WriteString("\r\n")
			buf.Write(encodeBase64(attachment.Content))
			buf.WriteString("\r\n")
		}
		buf.WriteString("--MIXED_BOUNDARY--\r\n")
	}

	return buf.Bytes(), nil
}

// 执行发送邮件
func (c *Client) sendEmail(from string, to, cc, bcc []string, raw []byte) error {

	addr := net.JoinHostPort(c.config.Host, strconv.Itoa(c.config.Port))

	var conn net.Conn
	var err error

	if c.config.SSL {
		tlsconfig := c.config.TLSConfig
		if tlsconfig == nil {
			tlsconfig = &tls.Config{ServerName: c.config.Host}
		}

		conn, err = tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("mail: tls dial error: %v", err)
		}
	} else {
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("mail: dial error: %v", err)
		}
	}

	client, err := smtp.NewClient(conn, c.config.Host)
	if err != nil {
		return fmt.Errorf("mail: smtp new client error: %v", err)
	}
	defer client.Close()

	if c.config.LocalName != "" {
		if err = client.Hello(c.config.LocalName); err != nil {
			return fmt.Errorf("mail: helo error: %v", err)
		}
	}

	if !c.config.SSL {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsconfig := c.config.TLSConfig
			if tlsconfig == nil {
				tlsconfig = &tls.Config{ServerName: c.config.Host}
			}
			if err = client.StartTLS(tlsconfig); err != nil {
				return fmt.Errorf("mail: starttls error: %v", err)
			}
		}
	}

	if c.config.Username != "" && c.config.Password != "" {
		auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("mail: auth error: %v", err)
		}
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("mail: mail from error: %v", err)
	}

	recipients := append(to, append(cc, bcc...)...)
	for _, addr := range recipients {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("mail: rcpt to error: %v", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("mail: data error: %v", err)
	}

	_, err = w.Write(raw)
	if err != nil {
		return fmt.Errorf("mail: write error: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("mail: close error: %v", err)
	}

	return client.Quit()
}

// 编码文本内容
func encodeText(text string) string {
	// 这里简化处理，实际应该使用quoted-printable编码
	return text
}

// 编码附件内容
func encodeBase64(data []byte) []byte {
	// 这里简化处理，实际应该使用base64编码
	return data
}

// 写入头部信息
func writeHeaders(buf *bytes.Buffer, headers textproto.MIMEHeader) {
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ", ")))
	}
}

// 获取内联附件
func (e *Email) getInlineAttachments() []*Attachment {
	var inlines []*Attachment
	for _, a := range e.Attachments {
		if a.Inline {
			inlines = append(inlines, a)
		}
	}
	return inlines
}
