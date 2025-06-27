package v1

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	generated "github.com/hizu77/library-service/generated/api/library"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewLibraryRoutes(
	ctx context.Context,
	mux *runtime.ServeMux,
	grpcAddr string) error {
	{
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		err := generated.RegisterLibraryHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
