package service

import (
	"bytes"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/localization"
	"gopkg.in/tucnak/telebot.v2"
)

type Service struct {
	DB *db.DbStruct
}

func (s Service) GetCard(chatId int) ([]*telebot.Photo, error) {
	resp, err := communication.GetCard(chatId)
	if err != nil {
		return nil, err
	}
	if resp.SessionIsNew {
		logo, rules, err := s.DB.GetLogoAndRules()
		if err != nil {
			return nil, err
		}
		// turn pic type to telebot.Photo type
		return []*telebot.Photo{
			{File: telebot.FromReader(bytes.NewReader(logo.Data)), Caption: localization.GetLoc("den4ik_logo")},
			{File: telebot.FromReader(bytes.NewReader(rules.Data)), Caption: localization.GetLoc("den4ik_rules")},
		}, nil
	}
	// todo session end
	if resp.SessionEnd {

	}

	cardID := resp.Card.Value + "_" + s.adjustSuit(resp.Card.Suit)
	caption := localization.GetLoc(resp.Card.Value + "_card")
	pic, err := s.GetAndFormPicMessage(cardID, caption)
	if err != nil {
		return nil, err
	}
	return []*telebot.Photo{pic}, nil
}

func (s Service) GetAndFormPicMessage(id, caption string) (*telebot.Photo, error) {
	p, err := s.DB.GetPicFromDB(id)
	if err != nil {
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
