package account

import (
	"time"

	"github.com/google/uuid"
	"saidakbar.origin/caching"
	"saidakbar.origin/caching/models"
	"saidakbar.origin/core"
	"saidakbar.origin/dto"
	"saidakbar.origin/repository"
)

type accountService struct {
	accountRepo repository.AccountRepository
	otpCache    caching.OtpCache
	tokenCache  caching.TokenCache
}

func NewService(
	accountRepo repository.AccountRepository,
	otpCache caching.OtpCache,
	tokenCache caching.TokenCache,
) Service {
	return &accountService{
		accountRepo: accountRepo,
		otpCache:    otpCache,
		tokenCache:  tokenCache,
	}
}

func (s *accountService) VerifyCode(model VerifyCodeModel) (result *dto.VerifyCodeResult, err error) {
	result = new(dto.VerifyCodeResult)

	// get account
	acc, getErr := s.accountRepo.GetByPhone(model.Phone)
	if getErr != nil {
		err = getErr
		return
	}
	if acc == nil {
		result.UserNotFound = true
		return
	}

	// validate OTP
	{
		stored, getErr := s.otpCache.GetOTP(model.Phone)
		if getErr != nil || stored != model.Code {
			result.InvalidCode = true
			return
		}
		_ = s.otpCache.DeleteOTP(model.Phone)
	}

	// activate account on first login
	if !acc.IsActive {
		if err = s.accountRepo.Activate(acc.ID); err != nil {
			return
		}
	}

	// issue bearer token
	{
		token := uuid.New().String()
		cachedAccount := &models.Account{
			ID:        acc.ID,
			FullName:  acc.FullName,
			Phone:     acc.Phone,
			ChatID:    acc.ChatID,
			IsActive:  acc.IsActive,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
			DeletedAt: acc.DeletedAt,
			IsDeleted: acc.IsDeleted,
			TokenType: 0,
		}
		if err = s.tokenCache.SetAccount(token, cachedAccount, core.TokenTTL); err != nil {
			return
		}
		result.Token = token
	}

	return
}

func (s *accountService) UpdateAccount(model UpdateAccountModel) (*UpdateAccountResult, error) {
	acc, err := s.accountRepo.GetByID(model.ID)
	if err != nil {
		return nil, err
	}

	acc.FullName = model.FullName
	acc.UpdatedAt = time.Now()
	if err = s.accountRepo.Update(acc); err != nil {
		return nil, err
	}

	return &UpdateAccountResult{Account: acc}, nil
}
