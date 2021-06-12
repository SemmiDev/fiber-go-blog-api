package service

import (
	"database/sql"
	"github.com/SemmiDev/fiber-go-blog/internal/app/model"
	"github.com/SemmiDev/fiber-go-blog/internal/app/repository"
	"github.com/SemmiDev/fiber-go-blog/internal/auth"
	"github.com/SemmiDev/fiber-go-blog/internal/constant"
	"github.com/SemmiDev/fiber-go-blog/internal/logger"
	"github.com/SemmiDev/fiber-go-blog/internal/web"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

type PostService interface {
	Create(c *fiber.Ctx, req model.PostCreateRequest) (*model.PostResponse, error)
	List(c *fiber.Ctx, req model.PostListRequest) ([]*model.PostResponse, error)
	Get(c *fiber.Ctx, req model.PostGetRequest) (*model.PostResponse, error)
	Update(c *fiber.Ctx, req model.PostUpdateRequest) (*model.PostResponse, error)
	Delete(c *fiber.Ctx, req model.PostDeleteRequest) error
}

func NewPostService(postRepository repository.PostRepository) PostService {
	return &postService{postRepository}
}

type postService struct {
	postRepository repository.PostRepository
}

func (s *postService) Create(c *fiber.Ctx, req model.PostCreateRequest) (*model.PostResponse, error) {
	//Check if the auth token is valid and  get the user id from it
	accountID, _ := auth.ExtractTokenID(c)

	post := &model.Post{
		Title:     req.Title,
		Body:      req.Body,
		CreatedAt: time.Now(),
		AccountID: accountID,
	}

	err := s.postRepository.Create(c.Context(), post)

	if err != nil {
		logger.Log().Err(err).Msg("failed to create post")
		return nil, constant.ErrServer
	}

	var account *model.Account
	account, err = s.postRepository.GetAccount(c.Context(), accountID)

	return &model.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		UpdatedAt: nil,
		AccountID: accountID,
		Account: &model.AccountResponse{
			ID:        accountID,
			Name:      account.Name,
			Email:     account.Email,
			CreatedAt: account.CreatedAt,
			UpdatedAt: nil,
		},
	}, nil
}

func (s *postService) List(c *fiber.Ctx, req model.PostListRequest) ([]*model.PostResponse, error) {
	posts, err := s.postRepository.List(c.Context(), req.Limit, req.Offset, req.Title)
	if err != nil {
		logger.Log().Err(err).Msg("failed to list posts")
		return nil, constant.ErrServer
	}

	return model.NewPostListResponse(posts), nil
}

func (s *postService) Get(c *fiber.Ctx, req model.PostGetRequest) (*model.PostResponse, error) {
	post, err := s.postRepository.Get(c.Context(), req.ID)
	if err != nil {
		return nil, s.switchErrPostNotFoundOrErrServer(err)
	}

	return model.NewPostResponse(post), nil
}

func (s *postService) Update(c *fiber.Ctx, req model.PostUpdateRequest) (*model.PostResponse, error) {
	post, err := s.postRepository.Get(c.Context(), req.ID)
	if err != nil {
		return nil, s.switchErrPostNotFoundOrErrServer(err)
	}

	// authorize user who create post
	accountID, err := Authorize(c)
	if err != nil {
		return nil, constant.ErrUnauthorized
	}

	log.Println(accountID)
	log.Println(post.AccountID)

	if post.AccountID != accountID {
		return nil, constant.ErrUnauthorized
	}

	post.Title = req.Title
	post.Body = req.Body
	post.UpdatedAt.Time = time.Now()

	err = s.postRepository.Update(c.Context(), post)
	if err != nil {
		return nil, s.switchErrPostNotFoundOrErrServer(err)
	}

	return model.NewPostResponse(post), nil
}

func (s *postService) Delete(c *fiber.Ctx, req model.PostDeleteRequest) error {
	post, err := s.postRepository.Get(c.Context(), req.ID)
	if err != nil {
		return s.switchErrPostNotFoundOrErrServer(err)
	}

	// authorize user who create post
	accountID, err := Authorize(c)
	if err != nil {
		return constant.ErrUnauthorized
	}

	if post.AccountID != accountID {
		return constant.ErrUnauthorized
	}

	err = s.postRepository.Delete(c.Context(), req.ID)
	if err != nil {
		return s.switchErrPostNotFoundOrErrServer(err)
	}

	return nil
}

func (s *postService) switchErrPostNotFoundOrErrServer(err error) error {
	switch err {
	case sql.ErrNoRows:
		return constant.ErrPostNotFound
	default:
		logger.Log().Err(err).Msg("failed to execute operation post repository")
		return constant.ErrServer
	}
}

func Authorize(c *fiber.Ctx) (int64, error) {
	//Check if the auth token is valid and  get the user id from it
	id, err := auth.ExtractTokenID(c)
	if err != nil {
		return 0, web.MarshalError(c, fiber.StatusUnauthorized, err)
	}

	return id, nil
}
