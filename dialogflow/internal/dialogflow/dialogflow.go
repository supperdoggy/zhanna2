package dialogflow

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	dialogflowdata "github.com/supperdoggy/superSecretDevelopement/structs/request/dialogflow"
	cfg "github.com/supperdoggy/superSecretDevelopement/structs/services/dialogflow"
	"go.uber.org/zap"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"os"
)

type Dialogflow struct {
	ProjectID    string
	LanguageCode string
	Session      *dialogflow.SessionsClient
	Logger       *zap.Logger
}

var DF Dialogflow

func init() {
	logger, _ := zap.NewDevelopment()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", cfg.WorkCreds)

	sessionClient, err := dialogflow.NewSessionsClient(context.Background())
	if err != nil {
		panic(err.Error())
	}
	DF = Dialogflow{
		ProjectID:    "small-talk-qsespi",
		LanguageCode: "ru-RU",
		Session:      sessionClient,
		Logger:       logger,
	}
}

func (d *Dialogflow) DetectIntentText(req dialogflowdata.GetAnswerReq) dialogflowdata.GetAnswerResp {
	if req.ID == "" {
		return dialogflowdata.GetAnswerResp{Err: fmt.Sprintf("Received empty project (%s) or session (%s)", d.ProjectID, req.ID)}
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", d.ProjectID, req.ID)
	textInput := dialogflowpb.TextInput{Text: req.Text, LanguageCode: d.LanguageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := d.Session.DetectIntent(context.Background(), &request)
	if err != nil {
		d.Logger.Error("DetectIntent error", zap.Error(err), zap.Any("req", req))
		return dialogflowdata.GetAnswerResp{Err: err.Error()}
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	resp := dialogflowdata.GetAnswerResp{
		Answer: fulfillmentText,
	}
	return resp
}
