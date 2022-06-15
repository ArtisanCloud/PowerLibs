package mail

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
)

type SMTPConfig struct {
	Host       string
	UserName   mail.Address
	Password   string
	PortNumber string
	ServerName string
	TLSConfig  *tls.Config
}

var (
	smtpConfig = &SMTPConfig{
		Host:       os.Getenv("email_host"),
		UserName:   mail.Address{"", os.Getenv("email_username")},
		Password:   os.Getenv("email_password"),
		PortNumber: os.Getenv("email_port"),
	}
)

type Sender struct {
	auth       smtp.Auth
	SMTPConfig *SMTPConfig
}

type Message struct {
	To          []mail.Address
	CC          []mail.Address
	BCC         []mail.Address
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func NewSender(config *SMTPConfig) *Sender {
	if config == nil {
		config = smtpConfig
	}
	auth := smtp.PlainAuth("", config.UserName.Address, config.Password, config.Host)

	config.ServerName = fmt.Sprintf("%s:%s", config.Host, config.PortNumber)
	// TLS config
	config.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.Host,
	}

	return &Sender{
		auth,
		config,
	}
}

func (s *Sender) Send(m *Message) error {

	from := s.SMTPConfig.UserName
	to := ""
	for _, addr := range m.To {
		to += addr.String() + ";"
	}
	cc := ""
	for _, addr := range m.CC {
		cc += addr.String() + ";"
	}
	bcc := ""
	for _, addr := range m.CC {
		bcc += addr.String() + ";"
	}
	subj := m.Subject
	body := m.Body
	//attachments := m.Attachments

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to
	headers["CC"] = cc
	headers["BCC"] = bcc
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := s.SMTPConfig.ServerName

	host, _, _ := net.SplitHostPort(servername)

	auth := s.auth

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	for _, addr := range m.To {
		if err = c.Rcpt(addr.Address); err != nil {
			log.Panic(err)
		}
	}

	for _, addr := range m.CC {
		if err = c.Rcpt(addr.Address); err != nil {
			log.Panic(err)
		}
	}

	for _, addr := range m.BCC {
		if err = c.Rcpt(addr.Address); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	err = c.Quit()
	return err
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachData(fileName string, data []byte) error {
	m.Attachments[fileName] = data
	return nil
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

//func (m *Message) ToBytes() []byte {
//buf := bytes.NewBuffer(nil)
//withAttachments := len(m.Attachments) > 0
//buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
//for _, addr := range m.To {
//	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(addr.Address, ",")))
//}
//if len(m.CC) > 0 {
//	for _, addr := range m.To {
//
//		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(addr.Address, ",")))
//	}
//}
//
//if len(m.BCC) > 0 {
//	buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
//}
//
//buf.WriteString("MIME-Version: 1.0\n")
//writer := multipart.NewWriter(buf)
//boundary := writer.Boundary()
//if withAttachments {
//	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
//	buf.WriteString(fmt.Sprintf("--%s\n", boundary))
//} else {
//	buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
//}
//
//buf.WriteString(m.Body)
//if withAttachments {
//	for k, v := range m.Attachments {
//		buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
//		buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
//		buf.WriteString("Content-Transfer-Encoding: base64\n")
//		buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))
//
//		b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
//		base64.StdEncoding.Encode(b, v)
//		buf.Write(b)
//		buf.WriteString(fmt.Sprintf("\n--%s", boundary))
//	}
//
//	buf.WriteString("--")
//}
//
//return buf.Bytes()
//}
