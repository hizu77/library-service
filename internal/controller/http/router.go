package http

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/hizu77/library-service/internal/controller/http/v1"
)

func NewRouter(
	ctx context.Context,
	app fiber.Router,
	grpcAddr string,
) error {
	mux := runtime.NewServeMux()
	{
		err := v1.NewLibraryRoutes(ctx, mux, grpcAddr)
		if err != nil {
			return err
		}
	}

	app.Use(adaptor.HTTPHandler(mux))

	return nil
}
