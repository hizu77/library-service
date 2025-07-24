package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/infra/model/outbox"
)

type (
	KindHandler   = func(ctx context.Context, data []byte) error
	GlobalHandler = func(kind outbox.Kind) (KindHandler, error)
)

func Handler() GlobalHandler {
	return func(kind outbox.Kind) (KindHandler, error) {
		switch kind {
		case outbox.KindAuthor:
			return authorHandler(), nil
		case outbox.KindBook:
			return bookHandler(), nil
		default:
			return nil, outbox.ErrUnknownKind
		}
	}
}

func authorHandler() KindHandler {
	// TODO KAFKA CONSUMER
	return func(ctx context.Context, data []byte) error {
		author := entity.Author{}
		err := json.Unmarshal(data, &author)
		if err != nil {
			return err
		}

		fmt.Println("Send book" + string(data))

		return nil
	}
}

func bookHandler() KindHandler {
	//TODO KAFKA CONSUMER
	return func(ctx context.Context, data []byte) error {
		book := entity.Book{}
		err := json.Unmarshal(data, &book)
		if err != nil {
			return err
		}

		fmt.Println("Send book" + string(data))

		return nil
	}
}
