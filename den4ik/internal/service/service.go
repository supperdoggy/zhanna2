package service

import (
	"github.com/supperdoggy/superSecretDevelopement/den4ik/internal/db"
	"github.com/supperdoggy/superSecretDevelopement/structs"
	den4ikdata "github.com/supperdoggy/superSecretDevelopement/structs/request/den4ik"
	"gopkg.in/mgo.v2"
	"log"
	"math/rand"
	"time"
)

type Service struct {
	DB *db.DbStruct
}

func (s Service) formCardDeck() [36]structs.Card {
	cards := structs.Cards
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return cards
}

func (s Service) createNewSession(id int) structs.Session {
	arr := s.formCardDeck()
	slc := arr[:]
	session := structs.Session{
		ID:           id,
		Cards:        slc,
		CreationTime: time.Now(),
	}
	return session
}

func (s Service) GetCard(req den4ikdata.GetCardReq) (resp den4ikdata.GetCardResp, err error) {
	// first of all we look for the session
	log.Println("getting session", req.SessionID)
	session, err := s.DB.GetSession(req.SessionID)
	if err != nil && err != mgo.ErrNotFound {
		resp.Err = err.Error()
		return resp, err
	} else if err == mgo.ErrNotFound { // if we cant find session, we create it
		log.Println("creating session", req.SessionID)
		session = s.createNewSession(req.SessionID)
		// putting session into db
		log.Println("inserting session", req.SessionID)
		if err = s.DB.InsertGameSession(session); err != nil {
			resp.Err = err.Error()
			resp.Card = structs.Card{}
			return
		}
		resp.SessionIsNew = true
		return
	}
	// if we found session

	// check if session is older than 1 day
	// if it is we create new session
	if !time.Now().Before(session.CreationTime.AddDate(0, 0, 1)) {
		log.Println("session is older than 1 day...creating new session", req.SessionID)
		if err = s.DB.UpdateSession(req.SessionID, s.createNewSession(req.SessionID)); err != nil {
			resp.Err = err.Error()
			return
		}
		resp.SessionIsNew = true
		return
	}

	// check if we have cards left in structs
	if len(session.Cards) == 0 {
		resp.SessionEnd = true
		if err = s.DB.UpdateSession(req.SessionID, s.createNewSession(req.SessionID)); err != nil {
			resp.Err = err.Error()
			return
		}
		return
	}

	resp.Card = session.Cards[0]
	// if we got only card left, which we already took, we create new session, but do not send session_is_new = true
	if len(session.Cards) == 1 {
		log.Println("only 1 card left", req.SessionID)
		session.Cards = []structs.Card{}
		if err = s.DB.UpdateSession(req.SessionID, session); err != nil {
			resp.Card = structs.Card{}
			resp.Err = err.Error()
			return
		}
		return
	}
	// removing first card from deck
	session.Cards = session.Cards[1:]
	log.Println("updating session", req.SessionID)
	if err = s.DB.UpdateSession(req.SessionID, session); err != nil {
		resp.Err = err.Error()
		resp.Card = structs.Card{}
		return
	}
	return

}

func (s Service) ResetSession(req den4ikdata.ResetSessionReq) (resp den4ikdata.ResetSessionResp, err error) {
	err = s.DB.DeleteSession(req.SessionID)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.OK = true
	return
}
