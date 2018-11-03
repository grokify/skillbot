package alertteam

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	rc "github.com/grokify/go-ringcentral/client"
	"github.com/grokify/go-ringcentral/clientutil/glipgroups"
	hum "github.com/grokify/gotilla/net/httputilmore"
	//log "github.com/sirupsen/logrus"

	"github.com/grokify/chatblox"
	"github.com/grokify/go-glip/examples"
)

type Factory struct {
	HTTPClient *http.Client
	APIClient  *rc.APIClient
	ServerURL  string
}

func (f *Factory) GetGlipTeams() (glipgroups.GroupsSet, error) {
	return glipgroups.NewGroupsSetApiRequest(
		f.HTTPClient,
		f.ServerURL,
		"Team")
}

func (f *Factory) NewIntent() chatblox.Intent {
	return chatblox.Intent{
		Type: chatblox.MatchRegexpCapture,
		Regexps: []*regexp.Regexp{
			regexp.MustCompile(`^(?i)^\s*(?:alert team|alert) (?P<query_any>.+?)\s*[\?.]?\s*$`),
		},
		HandleIntent: f.HandleIntent,
	}
}

func (f *Factory) HandleIntent(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) (*hum.ResponseInfo, error) {
	qry, ok := slots["query_any"]
	if !ok {
		return bot.SendGlipPost(
			glipPostEventInfo,
			rc.GlipCreatePost{
				Text: "I'm sorry but I didn't understand you.\n\nPlease say `alert <my teamm>` to post alert into \"my team\"."})
	}
	qry = strings.TrimSpace(qry)

	teamsSet, err := f.GetGlipTeams()
	if err != nil {
		return bot.SendGlipPost(
			glipPostEventInfo,
			rc.GlipCreatePost{
				Text: "I'm sorry but I ran into a problem."})
	}
	groups := teamsSet.FindGroupsByNameLower(qry)

	if len(groups) == 0 {
		reqBody := rc.GlipCreatePost{
			Text: "I'm sorry but I couldn't find any groups for: " + qry}
		return bot.SendGlipPost(glipPostEventInfo, reqBody)
	}

	glipPosts := []rc.GlipCreatePost{
		examples.GetExamplePostAlertWarning(),
		examples.GetExamplePostAlertSOS()}

	attachments := []rc.GlipMessageAttachmentInfoRequest{}

	for _, group := range groups {
		for _, reqBody := range glipPosts {
			_, resp, err := bot.RingCentralClient.GlipApi.CreatePost(
				context.Background(), group.ID, reqBody,
			)
			if err != nil || resp.StatusCode >= 300 {
				reqBody := rc.GlipCreatePost{
					Text: "I'm sorry but I ran into a problem."}
				return bot.SendGlipPost(glipPostEventInfo, reqBody)
			}
		}
		attachment := rc.GlipMessageAttachmentInfoRequest{
			Type: "Card",
			Fields: []rc.GlipMessageAttachmentFieldsInfo{
				{
					Title: "Name",
					Value: group.Name,
					Style: "Short"},
				{
					Title: "group ID",
					Value: group.ID,
					Style: "Short"},
			},
		}
		attachments = append(attachments, attachment)
	}

	reqBody := rc.GlipCreatePost{
		Text:        "I successfully sent messages for the following groups:",
		Attachments: attachments}

	return bot.SendGlipPost(glipPostEventInfo, reqBody)
}
