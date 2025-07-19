package book

import (
	"context"
	"testing"

	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/usecase/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	TestUUID = "4d4a8cd8-501b-4bd4-8589-6be8dcca7c09"
)

func newBookUseCase(t *testing.T) (*UseCaseImpl, *mock.MockAuthorRepository, *mock.MockBookRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockBookRepository := mock.NewMockBookRepository(ctrl)
	mockAuthorRepository := mock.NewMockAuthorRepository(ctrl)
	mockTransactor := mock.NewMockTransactor(ctrl)
	mockTransactor.EXPECT().WithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, function func(ctx context.Context) error) error {
			return function(ctx)
		})
	logger := zap.NewNop()
	uc := NewUseCase(logger, mockAuthorRepository, mockBookRepository, mockTransactor)

	return uc, mockAuthorRepository, mockBookRepository
}

func TestGetAuthorBooks(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(authorRepo *mock.MockAuthorRepository, bookRepo *mock.MockBookRepository)
		want    []entity.Book
		wantErr error
	}{
		{
			name: "author not found",
			args: args{
				ctx: context.Background(),
				id:  TestUUID,
			},
			mock: func(authorRepo *mock.MockAuthorRepository, bookRepo *mock.MockBookRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).Return(entity.Author{}, entity.ErrAuthorNotFound)
				bookRepo.EXPECT().GetBooksByAuthorID(gomock.Any(), TestUUID).Times(0)
			},
			wantErr: entity.ErrAuthorNotFound,
			want:    nil,
		},
		{
			name: "author zero books",
			args: args{
				ctx: context.Background(),
				id:  TestUUID,
			},
			mock: func(authorRepo *mock.MockAuthorRepository, bookRepo *mock.MockBookRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).Return(entity.Author{
					ID:   TestUUID,
					Name: "misha",
				}, nil)
				bookRepo.EXPECT().GetBooksByAuthorID(gomock.Any(), TestUUID).Return([]entity.Book{}, nil)
			},
			wantErr: nil,
			want:    []entity.Book{},
		},
		{
			name: "author more than zero books",
			args: args{
				ctx: context.Background(),
				id:  TestUUID,
			},
			mock: func(authorRepo *mock.MockAuthorRepository, bookRepo *mock.MockBookRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).Return(entity.Author{
					ID:   TestUUID,
					Name: "misha",
				}, nil)
				bookRepo.EXPECT().GetBooksByAuthorID(gomock.Any(), TestUUID).Return([]entity.Book{
					{
						ID:         TestUUID,
						Name:       "misha book",
						AuthorsIDs: []string{TestUUID},
					},
				}, nil)
			},
			wantErr: nil,
			want: []entity.Book{
				{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			uc, authorRepo, bookRepo := newBookUseCase(st)
			tt.mock(authorRepo, bookRepo)

			got, err := uc.GetAuthorBooks(tt.args.ctx, tt.args.id)

			require.Equal(st, tt.want, got)
			require.ErrorIs(st, err, tt.wantErr)
		})
	}
}

func TestGetBookInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(bookRepo *mock.MockBookRepository)
		want    entity.Book
		wantErr error
	}{
		{
			name: "book not found",
			args: args{
				ctx: context.Background(),
				id:  TestUUID,
			},
			mock: func(bookRepo *mock.MockBookRepository) {
				bookRepo.EXPECT().GetBook(gomock.Any(), TestUUID).Return(entity.Book{}, entity.ErrBookNotFound)
			},
			wantErr: entity.ErrBookNotFound,
			want:    entity.Book{},
		},
		{
			name: "book success",
			args: args{
				ctx: context.Background(),
				id:  TestUUID,
			},
			mock: func(bookRepo *mock.MockBookRepository) {
				bookRepo.EXPECT().GetBook(gomock.Any(), TestUUID).Return(entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				}, nil)
			},
			wantErr: nil,
			want: entity.Book{
				ID:         TestUUID,
				Name:       "misha book",
				AuthorsIDs: []string{TestUUID},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			uc, _, bookRepo := newBookUseCase(st)
			tt.mock(bookRepo)

			got, err := uc.GetBookInfo(tt.args.ctx, tt.args.id)

			require.Equal(st, tt.want, got)
			require.ErrorIs(st, err, tt.wantErr)
		})
	}
}

func TestAddBook(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		book entity.Book
	}

	tests := []struct {
		name    string
		args    args
		mock    func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository)
		wantErr error
		want    entity.Book
	}{
		{
			name: "add book success",
			args: args{
				ctx: context.Background(),
				book: entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				},
			},
			mock: func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).Return(entity.Author{
					ID:   TestUUID,
					Name: "misha",
				}, nil)

				bookRepo.EXPECT().AddBook(gomock.Any(), gomock.Any()).
					Return(entity.Book{
						ID:         TestUUID,
						Name:       "misha book",
						AuthorsIDs: []string{TestUUID},
					}, nil)
			},
			wantErr: nil,
			want: entity.Book{
				ID:         TestUUID,
				Name:       "misha book",
				AuthorsIDs: []string{TestUUID},
			},
		},
		{
			name: "author not found",
			args: args{
				ctx: context.Background(),
				book: entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				},
			},
			mock: func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).
					Return(entity.Author{}, entity.ErrAuthorNotFound)

				bookRepo.EXPECT().AddBook(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: entity.ErrAuthorNotFound,
			want:    entity.Book{},
		},
		{
			name: "book already exists",
			args: args{
				ctx: context.Background(),
				book: entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				},
			},
			mock: func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).
					Return(entity.Author{
						ID:   TestUUID,
						Name: "misha",
					}, nil)

				bookRepo.EXPECT().AddBook(gomock.Any(), gomock.Any()).
					Return(entity.Book{}, entity.ErrBookAlreadyExists)
			},
			wantErr: entity.ErrBookAlreadyExists,
			want:    entity.Book{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			uc, authorRepo, bookRepo := newBookUseCase(st)
			tt.mock(bookRepo, authorRepo)

			got, err := uc.AddBook(tt.args.ctx, tt.args.book)

			require.Equal(st, tt.want, got)
			require.ErrorIs(st, err, tt.wantErr)
		})
	}
}

func TestUpdateBook(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		book entity.Book
	}

	tests := []struct {
		name    string
		args    args
		mock    func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository)
		wantErr error
		want    entity.Book
	}{
		{
			name: "update book success",
			args: args{
				ctx: context.Background(),
				book: entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				},
			},
			mock: func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).Return(entity.Author{
					ID:   TestUUID,
					Name: "misha",
				}, nil)

				bookRepo.EXPECT().UpdateBook(gomock.Any(), entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				}).Return(entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				}, nil)
			},

			wantErr: nil,
			want: entity.Book{
				ID:         TestUUID,
				Name:       "misha book",
				AuthorsIDs: []string{TestUUID},
			},
		},
		{
			name: "author not found",
			args: args{
				ctx: context.Background(),
				book: entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				},
			},
			mock: func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).
					Return(entity.Author{}, entity.ErrAuthorNotFound)

				bookRepo.EXPECT().UpdateBook(gomock.Any(), entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				}).Times(0)
			},
			wantErr: entity.ErrAuthorNotFound,
			want:    entity.Book{},
		},
		{
			name: "book not found",
			args: args{
				ctx: context.Background(),
				book: entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				},
			},
			mock: func(bookRepo *mock.MockBookRepository, authorRepo *mock.MockAuthorRepository) {
				authorRepo.EXPECT().GetAuthor(gomock.Any(), TestUUID).
					Return(entity.Author{
						ID:   TestUUID,
						Name: "misha",
					}, nil)

				bookRepo.EXPECT().UpdateBook(gomock.Any(), entity.Book{
					ID:         TestUUID,
					Name:       "misha book",
					AuthorsIDs: []string{TestUUID},
				}).Return(entity.Book{}, entity.ErrBookNotFound)
			},
			wantErr: entity.ErrBookNotFound,
			want:    entity.Book{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(st *testing.T) {
			st.Parallel()

			uc, authorRepo, bookRepo := newBookUseCase(st)
			tt.mock(bookRepo, authorRepo)

			got, err := uc.UpdateBook(tt.args.ctx, tt.args.book)

			require.Equal(st, tt.want, got)
			require.ErrorIs(st, err, tt.wantErr)
		})
	}
}
