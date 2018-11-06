package main

import (
	"fmt"
	"net/http"
	"os"

	ru "github.com/grokify/go-ringcentral/clientutil"
	cfg "github.com/grokify/gotilla/config"
	ro "github.com/grokify/oauth2more/ringcentral"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/chatblox"
	"github.com/grokify/chatblox/glip"

	"github.com/grokify/skillbot/handlers/alert_demo/alert"
	"github.com/grokify/skillbot/handlers/alert_demo/alert_team"
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
		err := cfg.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
		if err != nil {
			log.Warn(err)
		}
		engine = os.Getenv("CHATBLOX_ENGINE")
	}
	return engine
}

func NewGlipHandler() glip.RcOAuthManager {
	creds := ro.ApplicationCredentials{
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
	mux.HandleFunc("/oauth2callback", http.HandlerFunc(glipHandler.HandleOAuthNetHttp))

	serverURL := os.Getenv("RINGCENTRAL_SERVER_URL")
	httpClient, err := ro.NewHttpClientEnvFlexStatic("")
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
		chatblox.ServeAwsLambda(intentRouter)
	case "nethttp":
		log.Info("Starting Engine [nethttp]")
		chatblox.ServeNetHttp(intentRouter, mux)
	default:
		log.Fatal(fmt.Sprintf("E_NO_HTTP_ENGINE: [%v]", engine))
	}
}
