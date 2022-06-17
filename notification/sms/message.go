package sms

import (
	"github.com/ArtisanCloud/PowerLibs/v2/notification/contract"
)

type Message struct {
	contract.MessageInterface

	To          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}
