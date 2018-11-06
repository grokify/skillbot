package botinfo

import (
	"context"
	"fmt"
	"net/http"

	rc "github.com/grokify/go-ringcentral/client"
	hum "github.com/grokify/gotilla/net/httputilmore"

	"github.com/grokify/chatblox"
)

type Factory struct {
	HTTPClient *http.Client
	APIClient  *rc.APIClient
	ServerURL  string
}

func (f *Factory) NewIntent() chatblox.Intent {
	return chatblox.Intent{
		Type:         chatblox.MatchStringLowerCase,
		Strings:      []string{"bot info"},
		HandleIntent: f.HandleIntent,
	}
}

func (f *Factory) HandleIntent(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) (*hum.ResponseInfo, error) {
	info, resp, err := f.APIClient.UserSettingsApi.LoadExtensionInfo(context.Background(), "~", "~")
	if err != nil {
		return bot.SendGlipPost(glipPostEventInfo,
			rc.GlipCreatePost{
				Text: fmt.Sprintf("I'm sorry, I received the following error: %s", err.Error())})
	} else if resp.StatusCode >= 300 {
		return bot.SendGlipPost(glipPostEventInfo,
			rc.GlipCreatePost{
				Text: fmt.Sprintf("I'm sorry, I received the following status code: [%v]", resp.StatusCode)})
	}
	return bot.SendGlipPost(glipPostEventInfo,
		rc.GlipCreatePost{
			Text: fmt.Sprintf("Here's my info:\n\n* botId: %v", info.Id)})
}
