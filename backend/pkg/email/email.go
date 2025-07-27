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

type Message struct {
	From         string               // 发件人地址
	To           []string             // 收件人地址列表
	Cc           []string             // 抄送地址列表
	Bcc          []string             // 密送地址列表
	Subject      string               // 邮件主题
	Text         string               // 文本内容
	HTML         string               // HTML内容
	Headers      textproto.MIMEHeader // 自定义头字段
	Attachments  []*Attachment        // 附件列表
	ReadReceipts []string             // 回执地址列表
}

type Attachment struct {
	Filename string // 附件文件名
	Content  []byte // 附件内容
	Inline   bool   // 是否内联显示
}

// Mail Submission Agent
type MSA struct {
	Host      string      // 主机
	Port      int         // 端口
	Username  string      // 用户名
	Password  string      // 密码
	UseSSL    bool        // 是否使用 SSL
	UseTLS    bool        // 是否使用 TLS (SSL/TLS 二选一)
	TLSConfig *tls.Config // TLS 配置
	LocalName string      // 本地主机名
}

type Email struct {
	msa *MSA
}

func NewEmail(conf *viper.Viper) *Email {
	return &Email{msa: &MSA{
		Host:      conf.GetString("email.host"),
		Port:      conf.GetInt("email.port"),
		Username:  conf.GetString("email.username"),
		Password:  conf.GetString("email.password"),
		UseSSL:    conf.GetBool("email.use_ssl"),
		UseTLS:    conf.GetBool("email.use_tls"),
		LocalName: conf.GetString("email.local_name"),
	}}
}

func (e *Email) Send(msg *Message) error {

	// 必须指定至少一个收件人
	if len(msg.To) == 0 {
		return errors.New("mail: no recipient specified")
	}

	// 默认发件人与用户名相同
	if msg.From == "" {
		msg.From = e.msa.Username
	}

	from, err := mail.ParseAddress(msg.From)
	if err != nil {
		return fmt.Errorf("mail: invalid from address: %v", err)
	}

	to := make([]string, 0, len(msg.To))
	for _, addr := range msg.To {
		parsedAddr, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("mail: invalid to address %q: %v", addr, err)
		}
		to = append(to, parsedAddr.Address)
	}

	cc := make([]string, 0, len(msg.Cc))
	for _, addr := range msg.Cc {
		parsedAddr, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("mail: invalid cc address %q: %v", addr, err)
		}
		cc = append(cc, parsedAddr.Address)
	}

	bcc := make([]string, 0, len(msg.Bcc))
	for _, addr := range msg.Bcc {
		parsedAddr, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("mail: invalid bcc address %q: %v", addr, err)
		}
		bcc = append(bcc, parsedAddr.Address)
	}

	raw, err := e.buildEmail(msg, from, to)
	if err != nil {
		return err
	}

	return e.sendEmail(from.Address, to, cc, bcc, raw)
}

func (e *Email) buildEmail(msg *Message, from *mail.Address, to []string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// 设置头部
	headers := make(textproto.MIMEHeader)
	headers.Set("From", from.String())
	headers.Set("To", strings.Join(to, ", "))
	headers.Set("Subject", mime.QEncoding.Encode("utf-8", msg.Subject))
	headers.Set("Date", time.Now().Format(time.RFC1123Z))
	headers.Set("MIME-Version", "1.0")

	if len(msg.Cc) > 0 {
		headers.Set("Cc", strings.Join(msg.Cc, ", "))
	}

	if len(msg.ReadReceipts) > 0 {
		headers.Set("Disposition-Notification-To", strings.Join(msg.ReadReceipts, ", "))
	}

	// 添加自定义头部
	for k, v := range msg.Headers {
		headers[k] = v
	}

	// 根据是否有附件决定内容类型
	mixed := len(msg.Attachments) > 0
	var (
		related     bool
		alternative bool
	)

	if mixed {
		headers.Set("Content-Type", "multipart/mixed; boundary=MIXED_BOUNDARY")
		buf.WriteString("--MIXED_BOUNDARY\r\n")
	}

	alternative = msg.Text != "" && msg.HTML != ""
	if alternative {
		headers.Set("Content-Type", "multipart/alternative; boundary=ALTERNATIVE_BOUNDARY")
		buf.WriteString("--ALTERNATIVE_BOUNDARY\r\n")
	}

	related = len(msg.getInlineAttachments()) > 0
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
	if msg.Text != "" || msg.HTML == "" {
		textHeader := make(textproto.MIMEHeader)
		textHeader.Set("Content-Type", "text/plain; charset=utf-8")
		textHeader.Set("Content-Transfer-Encoding", "quoted-printable")

		writeHeaders(buf, textHeader)
		buf.WriteString(encodeText(msg.Text))
		buf.WriteString("\r\n")

		if alternative {
			buf.WriteString("--ALTERNATIVE_BOUNDARY\r\n")
		}
	}

	// 写入HTML内容
	if msg.HTML != "" {
		htmlHeader := make(textproto.MIMEHeader)
		htmlHeader.Set("Content-Type", "text/html; charset=utf-8")
		htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")

		writeHeaders(buf, htmlHeader)
		buf.WriteString(encodeText(msg.HTML))
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
		for _, attachment := range msg.Attachments {
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
func (e *Email) sendEmail(from string, to, cc, bcc []string, raw []byte) error {

	addr := net.JoinHostPort(e.msa.Host, strconv.Itoa(e.msa.Port))

	var conn net.Conn
	var err error

	if e.msa.UseSSL {
		tlsconfig := e.msa.TLSConfig
		if tlsconfig == nil {
			tlsconfig = &tls.Config{ServerName: e.msa.Host}
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

	client, err := smtp.NewClient(conn, e.msa.Host)
	if err != nil {
		return fmt.Errorf("mail: smtp new client error: %v", err)
	}
	defer client.Close()

	if e.msa.LocalName != "" {
		if err = client.Hello(e.msa.LocalName); err != nil {
			return fmt.Errorf("mail: helo error: %v", err)
		}
	}

	if e.msa.UseTLS && !e.msa.UseSSL {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsconfig := e.msa.TLSConfig
			if tlsconfig == nil {
				tlsconfig = &tls.Config{ServerName: e.msa.Host}
			}
			if err = client.StartTLS(tlsconfig); err != nil {
				return fmt.Errorf("mail: starttls error: %v", err)
			}
		}
	}

	if e.msa.Username != "" && e.msa.Password != "" {
		auth := smtp.PlainAuth("", e.msa.Username, e.msa.Password, e.msa.Host)
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
func (msg *Message) getInlineAttachments() []*Attachment {
	var inlines []*Attachment
	for _, a := range msg.Attachments {
		if a.Inline {
			inlines = append(inlines, a)
		}
	}
	return inlines
}
