package contract

import (
	"github.com/ArtisanCloud/go-libs/object"
	"time"
)

type MiddlewareInterface interface {
	GetName() string
	SetName(name string)

	RetryDecider(conditions *object.HashMap) bool
	Retries() int
	Delay() time.Duration
}
