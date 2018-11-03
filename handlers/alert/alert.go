package alert

import (
	rc "github.com/grokify/go-ringcentral/client"
	hum "github.com/grokify/gotilla/net/httputilmore"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/chatblox"
	"github.com/grokify/go-glip/examples"
)

func NewIntent() chatblox.Intent {
	return chatblox.Intent{
		Type:         chatblox.MatchStringLowerCase,
		Strings:      []string{"alert"},
		HandleIntent: handleIntent}
}

func handleIntent(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) (*hum.ResponseInfo, error) {
	glipPost := buildPost(bot, slots, glipPostEventInfo)
	return bot.SendGlipPost(glipPostEventInfo, glipPost)
}

func buildPost(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) rc.GlipCreatePost {
	log.Info("BUILD_POST_ALERT")
	return examples.GetExamplePostAlertWarning()
}
