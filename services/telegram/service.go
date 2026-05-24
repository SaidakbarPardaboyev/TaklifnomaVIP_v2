package telegram

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"saidakbar.origin/caching"
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/plugins"
	"saidakbar.origin/repository"
)

const btnGenerateCode = "🔑 Kod olish"

type Service interface {
	StartPolling()
}

type telegramService struct {
	bot         *tgbotapi.BotAPI
	accountRepo repository.AccountRepository
	otpCache    caching.OtpCache
}

func NewService(token string, accountRepo repository.AccountRepository, otpCache caching.OtpCache) (Service, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("telegram: %w", err)
	}

	if _, err = bot.Request(tgbotapi.SetMyCommandsConfig{
		Commands: []tgbotapi.BotCommand{
			{Command: "start", Description: "Botni ishga tushirish"},
			{Command: "help", Description: "Yordam"},
		},
	}); err != nil {
		log.Printf("telegram: set commands: %v", err)
	}

	return &telegramService{
		bot:         bot,
		accountRepo: accountRepo,
		otpCache:    otpCache,
	}, nil
}

func (s *telegramService) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := update.Message
		chatID := msg.Chat.ID

		// Contact share
		if msg.Contact != nil {
			phone := normalizePhone(msg.Contact.PhoneNumber)
			fullName := strings.TrimSpace(msg.Contact.FirstName + " " + msg.Contact.LastName)

			code, err := s.handleContact(chatID, phone, fullName)
			if err != nil {
				log.Printf("telegram: handle contact: %v", err)
				s.send(chatID, "Something went wrong. Please try again.")
				continue
			}
			s.sendCodeWithKeyboard(chatID, code)
			continue
		}

		// Commands
		if msg.IsCommand() {
			switch msg.Command() {
			case "start":
				s.handleStart(chatID, msg)
			case "help":
				s.handleHelp(chatID)
			}
			continue
		}

		// Reply keyboard buttons
		if msg.Text == btnGenerateCode {
			s.handleGenerateCode(chatID)
			continue
		}
	}
}

func (s *telegramService) handleStart(chatID int64, msg *tgbotapi.Message) {
	acc, err := s.accountRepo.GetByChatID(chatID)
	if err != nil {
		log.Printf("telegram: get account by chat id: %v", err)
		s.send(chatID, "Something went wrong. Please try again.")
		return
	}

	if acc == nil {
		// Unknown user — ask to share contact
		fullName := strings.TrimSpace(msg.Chat.FirstName + " " + msg.Chat.LastName)
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonContact("📱 Kontaktni ulashish"),
			),
		)
		keyboard.OneTimeKeyboard = true
		keyboard.ResizeKeyboard = true
		reply := tgbotapi.NewMessage(chatID, fmt.Sprintf(`🇺🇿
Salom %s 👋
Rasmiy botimizga xush kelibsiz

⬇️ Kontaktingizni yuboring (tugmani bosib)`, fullName))
		reply.ReplyMarkup = keyboard
		if _, err = s.bot.Send(reply); err != nil {
			log.Printf("telegram: send start: %v", err)
		}
		return
	}

	// Known user — warn if active code exists, else generate new one
	if existing, getErr := s.otpCache.GetOTP(acc.Phone); getErr == nil && existing != "" {
		s.sendWithKeyboard(chatID, "⚠️ Sizda allaqachon faol tasdiqlash kodi mavjud.\nIltimos, ilovaga kodni kiriting.")
		return
	}

	code, err := s.generateOrResendOTP(acc.Phone)
	if err != nil {
		log.Printf("telegram: generate otp: %v", err)
		s.send(chatID, "Something went wrong. Please try again.")
		return
	}
	s.sendCodeWithKeyboard(chatID, code)
}

func (s *telegramService) handleHelp(chatID int64) {
	s.sendWithKeyboard(chatID, `ℹ️ Yordam

/start — Botni ishga tushirish yoki yangi kod olish
[ 🔑 Kod olish ] — Yangi tasdiqlash kodi olish

Muammo yuzaga kelsa, @Saidakbar_Pardaboyev bilan bog'laning.`)
}

func (s *telegramService) handleGenerateCode(chatID int64) {
	acc, err := s.accountRepo.GetByChatID(chatID)
	if err != nil {
		log.Printf("telegram: get account by chat id: %v", err)
		s.send(chatID, "Something went wrong. Please try again.")
		return
	}
	if acc == nil {
		s.send(chatID, "Avval /start buyrug'ini yuboring.")
		return
	}

	if existing, getErr := s.otpCache.GetOTP(acc.Phone); getErr == nil && existing != "" {
		s.sendWithKeyboard(chatID, "⚠️ Sizda allaqachon faol tasdiqlash kodi mavjud.\nIltimos, ilovaga kodni kiriting.")
		return
	}

	code, err := s.generateOrResendOTP(acc.Phone)
	if err != nil {
		log.Printf("telegram: generate otp: %v", err)
		s.send(chatID, "Something went wrong. Please try again.")
		return
	}
	s.sendCodeWithKeyboard(chatID, code)
}

func (s *telegramService) handleContact(chatID int64, phone, fullName string) (string, error) {
	existing, err := s.accountRepo.GetByPhone(phone)
	if err != nil {
		return "", err
	}

	now := plugins.GetNow()

	if existing == nil {
		acc := &mysql_entity.AccountModel{
			BaseEntity: mysql_entity.BaseEntity{
				ID:        uuid.New().String(),
				CreatedAt: now,
				UpdatedAt: now,
				IsDeleted: false,
			},
			FullName: fullName,
			Phone:    phone,
			ChatID:   chatID,
			IsActive: false,
		}
		if err = s.accountRepo.Create(acc); err != nil {
			return "", err
		}
	} else {
		existing.FullName = fullName
		existing.ChatID = chatID
		existing.UpdatedAt = now
		if err = s.accountRepo.Update(existing); err != nil {
			return "", err
		}
	}

	return s.generateOrResendOTP(phone)
}

func (s *telegramService) generateOrResendOTP(phone string) (string, error) {
	if stored, err := s.otpCache.GetOTP(phone); err == nil && stored != "" {
		return stored, nil
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	if err := s.otpCache.SetOTP(phone, code); err != nil {
		return "", err
	}
	return code, nil
}

func (s *telegramService) mainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(btnGenerateCode),
		),
	)
	keyboard.ResizeKeyboard = true
	return keyboard
}

func (s *telegramService) sendCodeWithKeyboard(chatID int64, code string) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🔒 Code:\n<code>%s</code>", code))
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = s.mainKeyboard()
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("telegram: send code: %v", err)
	}
}

func (s *telegramService) sendWithKeyboard(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = s.mainKeyboard()
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("telegram: send: %v", err)
	}
}

func (s *telegramService) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("telegram: send: %v", err)
	}
}

func (s *telegramService) sendHTML(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("telegram: send: %v", err)
	}
}

func normalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}
	return phone
}
