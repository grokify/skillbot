package main

import (
	"fmt"
	"net/http"
	"os"

	cfg "github.com/grokify/gotilla/config"
	ro "github.com/grokify/oauth2more/ringcentral"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/chatblox"
	"github.com/grokify/chatblox/glip"

	"github.com/grokify/skillbot/handlers/query"
)

func LoadEnv() string {
	// Check and load environment file if necessary
	engine := os.Getenv("BOTBLOX_ENGINE")
	if len(engine) == 0 {
		err := cfg.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
		if err != nil {
			log.Warn(err)
		}
		engine = os.Getenv("BOTBLOX_ENGINE")
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

	// Set intents
	intentRouter := chatblox.IntentRouter{
		Intents: []chatblox.Intent{
			query.NewIntent()}} // Default

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
