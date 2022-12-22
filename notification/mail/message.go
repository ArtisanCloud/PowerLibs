package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/notification/contract"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/mail"
	"path/filepath"
)

type Message struct {
	contract.MessageInterface

	To          []mail.Address
	CC          []mail.Address
	BCC         []mail.Address
	Subject     string
	Body        string
	Attachments map[string][]byte
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

func (m *Message) BodyToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	//buf.WriteString(fmt.Sprintf("From: %s\n", "matrix-x@artisan-cloud.com"))
	//buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	//
	//if len(m.To) > 0 {
	//	addresses := ""
	//	for _, addr := range m.To {
	//		addresses += addr.Address + ","
	//	}
	//	addresses = addresses[:len(addresses)-1]
	//	buf.WriteString(fmt.Sprintf("To: %s\n", addresses))
	//
	//}
	//if len(m.CC) > 0 {
	//	addresses := ""
	//	for _, addr := range m.CC {
	//		addresses += addr.Address + ","
	//	}
	//	addresses = addresses[:len(addresses)-1]
	//	buf.WriteString(fmt.Sprintf("Cc: %s\n", addresses))
	//
	//}
	//
	//if len(m.BCC) > 0 {
	//	addresses := ""
	//	for _, addr := range m.BCC {
	//		addresses += addr.Address + ","
	//	}
	//	addresses = addresses[:len(addresses)-1]
	//	buf.WriteString(fmt.Sprintf("Bcc: %s\n", addresses))
	//
	//}

	buf.WriteString("MIME-Version: 1.0\r\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n\n")
	}

	buf.WriteString(m.Body)
	buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
	if withAttachments {
		for k, v := range m.Attachments {
			//buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()

}
