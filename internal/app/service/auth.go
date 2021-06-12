package service

import (
	"context"
	"database/sql"
	"github.com/SemmiDev/fiber-go-blog/internal/app/model"
	"github.com/SemmiDev/fiber-go-blog/internal/app/repository"
	"github.com/SemmiDev/fiber-go-blog/internal/auth"
	"github.com/SemmiDev/fiber-go-blog/internal/constant"
	"github.com/SemmiDev/fiber-go-blog/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req model.AuthRequest) (*model.AuthResponse, error)
}

func NewAuthService(accountRepository repository.AccountRepository) AuthService {
	return &authService{accountRepository}
}

type authService struct {
	accountRepository repository.AccountRepository
}

func (s *authService) Login(ctx context.Context, req model.AuthRequest) (*model.AuthResponse, error) {
	account, err := s.accountRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, constant.ErrEmailNotRegistered
		default:
			logger.Log().Err(err).Msg("failed to get account by email")
			return nil, constant.ErrServer
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password))
	if err != nil {
		return nil, constant.ErrWrongPassword
	}

	accessToken, err := auth.CreateToken(account.ID)
	if err != nil {
		logger.Log().Err(err).Msg("failed to generate token")
		return nil, constant.ErrServer
	}

	return &model.AuthResponse{Token: accessToken}, nil
}