package account

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"saidakbar.origin/caching"
	"saidakbar.origin/caching/models"
	"saidakbar.origin/core"
	"saidakbar.origin/dto"
	"saidakbar.origin/repository"

	"github.com/google/uuid"
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

	var accID, accFullName string
	var accIsActive bool

	// get account
	{
		acc, getErr := s.accountRepo.GetByPhone(model.Phone)
		if getErr != nil {
			err = getErr
			return
		}
		if acc == nil {
			result.UserNotFound = true
			return
		}
		accID = acc.ID
		accFullName = acc.FullName
		accIsActive = acc.IsActive
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
	{
		if !accIsActive {
			if err = s.accountRepo.Activate(accID); err != nil {
				return
			}
		}
	}

	// issue bearer token
	{
		token := uuid.New().String()
		cachedAccount := &models.Account{
			ID:        primitive.NewObjectID(),
			Username:  model.Phone,
			Name:      accFullName,
			TokenType: 0,
			ActiveOrganization: &models.Organization{
				ID:   accID,
				Name: accFullName,
			},
		}
		if err = s.tokenCache.SetAccount(token, cachedAccount, core.TokenTTL); err != nil {
			return
		}
		result.Token = token
	}

	return
}
