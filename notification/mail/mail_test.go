package mail

import (
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"net/mail"
	"os"
	"testing"
	"time"
)

func Test_SendFile(t *testing.T) {
	sender := NewSender(nil)
	m := NewMessage("Test", "test for PowerLib Mail Sender."+time.Now().String())
	//m.To = []string{"dev@artisan-cloud.com"}
	m.To = []mail.Address{mail.Address{"", "383819640@qq.com"}}
	m.CC = []mail.Address{mail.Address{"", "tech@artisan-cloud.com"}}
	m.BCC = []mail.Address{mail.Address{"", "wechat@artisan-cloud.com"}}
	m.AttachFile(os.Getenv("test_attach_file"))
	err := sender.Send(m)
	if err != nil {
		t.Error(err)
	}
}

func Test_SendData(t *testing.T) {
	sender := NewSender(nil)
	m := NewMessage("Test", "test for PowerLib Mail Sender.")
	m.To = []mail.Address{mail.Address{"", "383819640@qq.com"}}
	m.CC = []mail.Address{mail.Address{"", "tech@artisan-cloud.com"}}
	m.BCC = []mail.Address{mail.Address{"", "wechat@artisan-cloud.com"}}

	data := &object.HashMap{
		"test":  1,
		"test2": "hello",
	}

	strData, err := object.JsonEncode(data)
	if err != nil {
		t.Error(err)
	}

	err = m.AttachData("dataFile.xlsx", []byte(strData))
	if err != nil {
		t.Error(err)
	}

	err = sender.Send(m)
	if err != nil {
		t.Error(err)
	}
}
