package author

import (
	"context"
	"testing"

	"github.com/hizu77/library-service/internal/usecase/mock"

	"github.com/hizu77/library-service/internal/entity"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	TestUUID = "4d4a8cd8-501b-4bd4-8589-6be8dcca7c09"
)

func newAuthorUseCase(t *testing.T) (*UseCaseImpl, *mock.MockAuthorRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockAuthorRepository := mock.NewMockAuthorRepository(ctrl)
	mockOutbox := mock.NewMockOutboxRepository(ctrl)
	mockTransactor := mock.NewMockTransactor(ctrl)
	mockOutbox.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockTransactor.EXPECT().WithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, function func(ctx context.Context) error) error {
			return function(ctx)
		})
	logger := zap.NewNop()
	uc := NewUseCase(logger, mockAuthorRepository, mockTransactor, mockOutbox)

	return uc, mockAuthorRepository
}

func TestGetAuthorInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		uuid string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(repo *mock.MockAuthorRepository)
		want    entity.Author
		wantErr error
	}{
		{
			name: "author not found",
			args: args{
				ctx:  context.Background(),
				uuid: TestUUID,
			},
			mock: func(repo *mock.MockAuthorRepository) {
				repo.EXPECT().GetAuthor(gomock.Any(), TestUUID).
					Return(entity.Author{}, entity.ErrAuthorNotFound)
			},
			want:    entity.Author{},
			wantErr: entity.ErrAuthorNotFound,
		},
		{
			name: "author success",
			args: args{
				ctx:  context.Background(),
				uuid: TestUUID,
			},
			mock: func(repo *mock.MockAuthorRepository) {
				repo.EXPECT().GetAuthor(gomock.Any(), TestUUID).
					Return(entity.Author{
						ID:   TestUUID,
						Name: "misha",
					}, nil)
			},
			want: entity.Author{
				ID:   TestUUID,
				Name: "misha",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			uc, repo := newAuthorUseCase(st)
			tt.mock(repo)

			res, err := uc.GetAuthorInfo(tt.args.ctx, tt.args.uuid)

			require.Equal(st, tt.want, res)
			require.ErrorIs(st, err, tt.wantErr)
		})
	}
}

func TestRegisterAuthor(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req entity.Author
	}

	tests := []struct {
		name    string
		args    args
		mock    func(repo *mock.MockAuthorRepository)
		want    entity.Author
		wantErr error
	}{
		{
			name: "author success",
			args: args{
				ctx: context.Background(),
				req: entity.Author{
					Name: "misha",
				},
			},
			mock: func(repo *mock.MockAuthorRepository) {
				repo.EXPECT().AddAuthor(gomock.Any(), gomock.Any()).
					Return(entity.Author{
						ID:   TestUUID,
						Name: "misha",
					}, nil)
			},
			want: entity.Author{
				ID:   TestUUID,
				Name: "misha",
			},
			wantErr: nil,
		},
		{
			name: "author already exists",
			args: args{
				ctx: context.Background(),
				req: entity.Author{
					ID:   TestUUID,
					Name: "misha",
				},
			},
			mock: func(repo *mock.MockAuthorRepository) {
				repo.EXPECT().AddAuthor(gomock.Any(), gomock.Any()).
					Return(entity.Author{}, entity.ErrAuthorAlreadyExists)
			},
			want:    entity.Author{},
			wantErr: entity.ErrAuthorAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			uc, repo := newAuthorUseCase(st)
			tt.mock(repo)

			res, err := uc.RegisterAuthor(tt.args.ctx, tt.args.req)

			require.Equal(st, tt.want, res)
			require.ErrorIs(st, err, tt.wantErr)
		})
	}
}

func TestChangeAuthorInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req entity.Author
	}

	tests := []struct {
		name    string
		args    args
		mock    func(repo *mock.MockAuthorRepository)
		want    entity.Author
		wantErr error
	}{
		{
			name: "author success",
			args: args{
				ctx: context.Background(),
				req: entity.Author{
					ID:   TestUUID,
					Name: "misha",
				},
			},
			mock: func(repo *mock.MockAuthorRepository) {
				repo.EXPECT().UpdateAuthor(gomock.Any(), entity.Author{
					ID:   TestUUID,
					Name: "misha",
				}).Return(entity.Author{
					ID:   TestUUID,
					Name: "misha",
				}, nil)
			},
			want: entity.Author{
				ID:   TestUUID,
				Name: "misha",
			},
			wantErr: nil,
		},
		{
			name: "author does not exist",
			args: args{
				ctx: context.Background(),
				req: entity.Author{
					ID:   TestUUID,
					Name: "misha",
				},
			},
			mock: func(repo *mock.MockAuthorRepository) {
				repo.EXPECT().UpdateAuthor(gomock.Any(), entity.Author{
					ID:   TestUUID,
					Name: "misha",
				}).Return(entity.Author{}, entity.ErrAuthorNotFound)
			},
			want:    entity.Author{},
			wantErr: entity.ErrAuthorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			uc, repo := newAuthorUseCase(st)
			tt.mock(repo)

			res, err := uc.ChangeAuthorInfo(tt.args.ctx, tt.args.req)

			require.Equal(st, tt.want, res)
			require.ErrorIs(st, err, tt.wantErr)
		})
	}
}
