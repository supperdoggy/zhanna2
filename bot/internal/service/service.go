package service

import (
	"bytes"
	"errors"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	"go.uber.org/zap"
	"gopkg.in/tucnak/telebot.v2"
)

type Service struct {
	db     db.IDbStruct
	logger *zap.Logger
}

type IService interface {
	GetCard(chatId int, langcode string) ([]*telebot.Photo, error)
	GetAndFormPicMessage(id, caption string) (*telebot.Photo, error)
	ResetDen4ik(id int, langcode string) (msg string, err error)
}

var (
	ErrSessionEnded = errors.New("session ended")
)

func NewService(logger *zap.Logger, db db.IDbStruct) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s Service) GetCard(chatId int, langcode string) ([]*telebot.Photo, error) {
	resp, err := communication.GetCard(chatId)
	if err != nil {
		s.logger.Error("error getting card", zap.Error(err), zap.Int("chat_id", chatId))
		return nil, err
	}
	if resp.SessionIsNew {
		logo, rules, err := s.db.GetLogoAndRules()
		if err != nil {
			s.logger.Error("error getting logo and rules", zap.Error(err), zap.Int("chat_id", chatId))
			return nil, err
		}
		// turn pic type to telebot.Photo type
		return []*telebot.Photo{
			{File: telebot.FromReader(bytes.NewReader(logo.Data))},
			{File: telebot.FromReader(bytes.NewReader(rules.Data))},
		}, nil
	}

	if resp.SessionEnd {
		return nil, ErrSessionEnded
	}

	cardID := resp.Card.Value + "_" + s.adjustSuit(resp.Card.Suit)
	caption := localization.GetLoc(resp.Card.Value + "_card", langcode)
	pic, err := s.GetAndFormPicMessage(cardID, caption)
	if err != nil {
		s.logger.Error("error getting and forming pic message",
			zap.Error(err),
			zap.Int("chat_id", chatId),
			zap.String("card_id", cardID),
			zap.String("caprion", caption))
		return nil, err
	}
	return []*telebot.Photo{pic}, nil
}

func (s Service) GetAndFormPicMessage(id, caption string) (*telebot.Photo, error) {
	p, err := s.db.GetPicFromDB(id)
	if err != nil {
		s.logger.Error("error getting card", zap.Error(err), zap.String("chat_id", id))
		return nil, err
	}
	photo := telebot.Photo{
		File:    telebot.FromReader(bytes.NewReader(p.Data)),
		Caption: caption,
	}
	return &photo, nil
}

// adjustSuit made because we have different names for suits in bot db and service db
func (s Service) adjustSuit(suit string) string {
	switch suit {
	case "diamonds":
		suit = "red"
	case "spades":
		suit = "purple"
	case "hearts":
		suit = "green"
	case "clubs":
		suit = "yellow"
	default:
		suit = ""
	}
	return suit
}

func (s Service) ResetDen4ik(id int, langcode string) (msg string, err error) {
	resp, err := communication.ResetDen4ik(id)
	if err != nil {
		s.logger.Error("error resetting den4ik", zap.Error(err), zap.Int("id", id))
		return "", err
	}
	if !resp.OK {
		s.logger.Error("got error in response", zap.String("error", resp.Err), zap.Int("id", id))
		return "", errors.New(resp.Err)
	}

	return localization.GetLoc("reset_ok", langcode), nil
}
