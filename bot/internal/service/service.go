package service

import (
	"bytes"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/communication"
	"github.com/supperdoggy/superSecretDevelopement/bot/internal/db"
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
			{File: telebot.FromReader(bytes.NewReader(logo.Data)), Caption: "logo"},
			{File: telebot.FromReader(bytes.NewReader(rules.Data)), Caption: "rules"},
		}, nil
	}
	// todo session end
	if resp.SessionEnd {}

	resp.Card.Suit = s.adjustSuit(resp.Card.Suit)
	// todo add localization for every value we can get to add it as a caption
	pic, err := s.GetAndFormPicMessage(resp.Card.Value+"_"+resp.Card.Suit, "lol test")
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
		File: telebot.FromReader(bytes.NewReader(p.Data)),
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
