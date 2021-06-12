package service

import (
	"database/sql"
	"github.com/SemmiDev/fiber-go-blog/internal/app/model"
	"github.com/SemmiDev/fiber-go-blog/internal/app/repository"
	"github.com/SemmiDev/fiber-go-blog/internal/constant"
	"github.com/SemmiDev/fiber-go-blog/internal/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AccountService interface {
	Create(c *fiber.Ctx, req model.AccountCreateRequest) (*model.AccountResponse, error)
	List(c *fiber.Ctx, req model.AccountListRequest) ([]*model.AccountResponse, error)
	Get(c *fiber.Ctx, req model.AccountGetRequest) (*model.AccountResponse, error)
	Update(c *fiber.Ctx, req model.AccountUpdateRequest) (*model.AccountResponse, error)
	UpdatePassword(c *fiber.Ctx, req model.AccountPasswordUpdateRequest) (*model.AccountResponse, error)
	Delete(c *fiber.Ctx, req model.AccountDeleteRequest) error
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService{accountRepository}
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func (s *accountService) Create(c *fiber.Ctx, req model.AccountCreateRequest) (*model.AccountResponse, error) {
	_, err := s.accountRepository.GetByEmail(c.Context(), req.Email)

	if err != nil && err != sql.ErrNoRows {
		logger.Log().Err(err).Msg("failed to get account by email")
		return nil, constant.ErrServer
	} else if err == nil {
		return nil, constant.ErrEmailRegistered
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Err(err).Msg("failed to generate from password")
		return nil, constant.ErrServer
	}

	account := &model.Account{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(password),
		CreatedAt: time.Now(),
	}

	err = s.accountRepository.Create(c.Context(), account)
	if err != nil {
		logger.Log().Err(err).Msg("failed to create account")
		return nil, constant.ErrServer
	}

	return model.NewAccountResponse(account), nil
}

func (s *accountService) List(c *fiber.Ctx, req model.AccountListRequest) ([]*model.AccountResponse, error) {
	accounts, err := s.accountRepository.List(c.Context(), req.Limit, req.Offset, req.Name)
	if err != nil {
		logger.Log().Err(err).Msg("failed to list accounts")
		return nil, constant.ErrServer
	}

	return model.NewAccountListResponse(accounts), nil
}

func (s *accountService) Get(c *fiber.Ctx, req model.AccountGetRequest) (*model.AccountResponse, error) {
	account, err := s.accountRepository.Get(c.Context(), req.ID)
	if err != nil {
		return nil, s.switchErrAccountNotFoundOrErrServer(err)
	}

	return model.NewAccountResponse(account), nil
}

func (s *accountService) Update(c *fiber.Ctx, req model.AccountUpdateRequest) (*model.AccountResponse, error) {
	account, err := s.accountRepository.GetByEmail(c.Context(), req.Email)
	if err != nil && err != sql.ErrNoRows {
		logger.Log().Err(err).Msg("failed to get account by email")
		return nil, constant.ErrServer
	} else if err == nil && account.ID != req.ID {
		return nil, constant.ErrEmailRegistered
	}

	account, err = s.accountRepository.Get(c.Context(), req.ID)
	if err != nil {
		return nil, s.switchErrAccountNotFoundOrErrServer(err)
	}

	account.Name = req.Name
	account.Email = req.Email
	account.UpdatedAt.Time = time.Now()

	err = s.accountRepository.Update(c.Context(), account)
	if err != nil {
		return nil, s.switchErrAccountNotFoundOrErrServer(err)
	}

	return model.NewAccountResponse(account), nil
}

func (s *accountService) UpdatePassword(c *fiber.Ctx, req model.AccountPasswordUpdateRequest) (*model.AccountResponse, error) {
	account, err := s.accountRepository.Get(c.Context(), req.ID)
	if err != nil {
		return nil, s.switchErrAccountNotFoundOrErrServer(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.OldPassword))
	if err != nil {
		return nil, constant.ErrWrongPassword
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Err(err).Msg("failed to generate from password")
		return nil, constant.ErrServer
	}

	account.Password = string(password)
	account.UpdatedAt.Time = time.Now()

	err = s.accountRepository.Update(c.Context(), account)
	if err != nil {
		return nil, s.switchErrAccountNotFoundOrErrServer(err)
	}

	return model.NewAccountResponse(account), nil
}

func (s *accountService) Delete(c *fiber.Ctx, req model.AccountDeleteRequest) error {
	err := s.accountRepository.Delete(c.Context(), req.ID)
	if err != nil {
		return s.switchErrAccountNotFoundOrErrServer(err)
	}

	return nil
}

func (s *accountService) switchErrAccountNotFoundOrErrServer(err error) error {
	switch err {
	case sql.ErrNoRows:
		return constant.ErrAccountNotFound
	default:
		logger.Log().Err(err).Msg("failed to execute operation account repository")
		return constant.ErrServer
	}
}