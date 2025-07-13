package usecase_test

import (
	"context"
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/usecase/book"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func newBookUseCase(t *testing.T) (*book.UseCaseImpl, *MockAuthorRepository, *MockBookRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockBookRepository := NewMockBookRepository(ctrl)
	mockAuthorRepository := NewMockAuthorRepository(ctrl)
	logger := zap.NewNop()
	usecase := book.NewUseCase(logger, mockAuthorRepository, mockBookRepository)

	return usecase, mockAuthorRepository, mockBookRepository
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
		mock    func(authorRepo *MockAuthorRepository, bookRepo *MockBookRepository)
		want    []entity.Book
		wantErr error
	}{
		{
			name: "author not found",
			args: args{
				ctx: context.Background(),
				id:  TestUUID,
			},
			mock: func(authorRepo *MockAuthorRepository, bookRepo *MockBookRepository) {
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
			mock: func(authorRepo *MockAuthorRepository, bookRepo *MockBookRepository) {
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
			mock: func(authorRepo *MockAuthorRepository, bookRepo *MockBookRepository) {
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
		mock    func(bookRepo *MockBookRepository)
		want    entity.Book
		wantErr error
	}{
		{
			name: "book not found",
			args: args{
				ctx: context.Background(),
				id:  TestUUID,
			},
			mock: func(bookRepo *MockBookRepository) {
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
			mock: func(bookRepo *MockBookRepository) {
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
		mock    func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository)
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
			mock: func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository) {
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
			mock: func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository) {
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
			mock: func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository) {
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
		mock    func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository)
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
			mock: func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository) {
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
			mock: func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository) {
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
			mock: func(bookRepo *MockBookRepository, authorRepo *MockAuthorRepository) {
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
