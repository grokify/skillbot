package main

import (
	"fmt"
	"net/http"
	"os"

	//ru "github.com/grokify/go-ringcentral/clientutil"
	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/goauth"
	ro "github.com/grokify/goauth/ringcentral"
	cfg "github.com/grokify/mogo/config"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/chatblox"
	"github.com/grokify/chatblox/glip"

	"github.com/grokify/skillbot/handlers/alert_demo/alert"
	alertteam "github.com/grokify/skillbot/handlers/alert_demo/alert_team"
	"github.com/grokify/skillbot/handlers/bot_info/botinfo"
	"github.com/grokify/skillbot/handlers/bot_info/botteams"
	"github.com/grokify/skillbot/handlers/help"
	"github.com/grokify/skillbot/handlers/query"
	"github.com/grokify/skillbot/handlers/sorry"
)

func LoadEnv() string {
	// Check and load environment file if necessary
	engine := os.Getenv("CHATBLOX_ENGINE")
	if len(engine) == 0 {
		_, err := cfg.LoadDotEnv([]string{os.Getenv("ENV_PATH"), "./.env"}, 1)
		if err != nil {
			log.Warn(err)
		}
		engine = os.Getenv("CHATBLOX_ENGINE")
	}
	return engine
}

func NewGlipHandler() glip.RcOAuthManager {
	creds := goauth.CredentialsOAuth2{
		ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
		ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET"),
		ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL"),
		RedirectURL:  os.Getenv("RINGCENTRAL_REDIRECT_URL")}

	return glip.NewRcOAuthManager(nil, creds)
}

func main() {
	engine := LoadEnv()

	glipHandler := NewGlipHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth2callback", http.HandlerFunc(glipHandler.HandleOAuthNetHTTP))

	serverURL := os.Getenv("RINGCENTRAL_SERVER_URL")
	httpClient, err := ro.NewHTTPClientEnvFlexStatic("")
	if err != nil {
		log.Fatal(err)
	}
	apiClient, err := ru.NewApiClientHttpClientBaseURL(httpClient, serverURL)
	if err != nil {
		log.Fatal(err)
	}

	alertteamFactory := alertteam.Factory{
		HTTPClient: httpClient,
		ServerURL:  serverURL,
		APIClient:  apiClient}

	botinfoFactory := botinfo.Factory{
		HTTPClient: httpClient,
		ServerURL:  serverURL,
		APIClient:  apiClient}

	botteamsFactory := botteams.Factory{
		HTTPClient: httpClient,
		ServerURL:  serverURL,
		APIClient:  apiClient}

	// Set intents
	intentRouter := chatblox.IntentRouter{
		Intents: []chatblox.Intent{
			botinfoFactory.NewIntent(),
			botteamsFactory.NewIntent(),
			alertteamFactory.NewIntent(),
			alert.NewIntent(),
			help.NewIntent(),
			query.NewIntent(),
			sorry.NewIntent()}}

	// Run engine
	switch engine {
	case "awslambda":
		log.Info("Starting Engine [awslambda]")
		chatblox.ServeAWSLambda(intentRouter)
	case "nethttp":
		log.Info("Starting Engine [nethttp]")
		chatblox.ServeNetHTTP(intentRouter, mux)
	default:
		log.Fatal(fmt.Sprintf("E_NO_HTTP_ENGINE: [%v]", engine))
	}
}
