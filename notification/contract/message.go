package contract

type MessageInterface interface {
	SetTo(to string) MessageInterface
	GetTo() string
	SetBody(body string) MessageInterface
	GetBody() string
	SetAttachments() MessageInterface
	GetAttachments() map[string][]byte
}
