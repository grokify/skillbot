package help

import (
	// rc "github.com/grokify/go-ringcentral/client"
	rc "github.com/grokify/go-ringcentral-client/office/v1/client"
	hum "github.com/grokify/mogo/net/http/httputilmore"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/chatblox"
)

const (
	InstructionsMD = "Please try asking me as follows:\n\n* Find me a <your skill> expert.\n* Who can help me with <your skill>?"
)

func NewIntent() chatblox.Intent {
	return chatblox.Intent{
		Type:         chatblox.MatchStringLowerCase,
		Strings:      []string{"help", "help me", "help me!", "help me."},
		HandleIntent: handleIntent,
	}
}

func handleIntent(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) (*hum.ResponseInfo, error) {
	glipPost := buildPost(bot, slots, glipPostEventInfo)
	return bot.SendGlipPost(glipPostEventInfo, glipPost)
}

func buildPost(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) rc.GlipCreatePost {
	log.Info("BUILD_POST_EXPERTFINDER_QUERY")
	reqBody := rc.GlipCreatePost{}

	reqBody.Text = "I can help you find an expert.\n\n" + InstructionsMD
	return reqBody
}
