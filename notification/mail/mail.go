package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"os"
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
	body := m.BodyToBytes()
	//attachments := m.Attachments

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to
	headers["CC"] = cc
	headers["BCC"] = bcc
	headers["Subject"] = subj

	// Setup sms
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	//sms += "\r\n"

	contentData := bytes.Join([][]byte{
		[]byte(message),
		body,
	}, []byte(""))
	//contentData := body

	// Connect to the SMTP Server
	servername := s.SMTPConfig.ServerName

	host, _, _ := net.SplitHostPort(servername)

	auth := s.auth

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsConfig)
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

	_, err = w.Write(contentData)
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
