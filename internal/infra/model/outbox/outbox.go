package outbox

import "errors"

var ErrUnknownKind = errors.New("unknown kind")

const (
	KindUndefined Kind = iota
	KindBook
	KindAuthor
)

type (
	Kind int
	Data struct {
		IdempotencyKey string
		RawData        []byte
		Kind           Kind
	}
)

func (k Kind) String() string {
	switch k {
	case KindBook:
		return "book"
	case KindAuthor:
		return "author"
	default:
		return "undefined"
	}
}
