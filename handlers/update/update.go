package query

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	rc "github.com/grokify/go-ringcentral/client"
	//"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/html/htmlutil"
	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/text/markdown"
	"github.com/grokify/gotilla/type/strutil"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/botblox"
)

func NewIntent() botblox.Intent {
	return botblox.Intent{
		Type: botblox.MatchRegexpCapture,
		Regexps: []*regexp.Regexp{
			regexp.MustCompile(`^(?i)^\s*(?:update skills|update) (?P<query_any>.+)`),
		},
		HandleIntent: handleIntent,
	}
}

func handleIntent(bot *botblox.Bot, slots map[string]string, glipPostEventInfo *botblox.GlipPostEventInfo) (*hum.ResponseInfo, error) {
	glipPost := buildPost(bot, slots, glipPostEventInfo)
	return bot.SendGlipPost(glipPostEventInfo, glipPost)
}

func buildPost(bot *botblox.Bot, slots map[string]string, glipPostEventInfo *botblox.GlipPostEventInfo) rc.GlipCreatePost {
	log.Info("BUILD_POST_EXPERTFINDER_QUERY")
	reqBody := rc.GlipCreatePost{}

	qry, ok := slots["query_any"]
	qry = strings.TrimSpace(qry)

	if !ok || len(qry) == 0 {
		reqBody.Text = "No query found."
		return reqBody
	}

	bot.BotConfig.AlgoliaIndex = strings.TrimSpace(bot.BotConfig.AlgoliaIndex)
	if len(bot.BotConfig.AlgoliaIndex) == 0 {
		reqBody.Text = "Configuration error: No Index."
		return reqBody
	}

	client, err := botblox.GetAlgoliaApiClient(bot.BotConfig)
	if err != nil {
		log.Info(fmt.Sprintf("ALGOLIA_CLIENT_ERROR [%v]", err))
		reqBody.Text = fmt.Sprintf("Configuration error: Bad Client [%v].", err)
		return reqBody
	}
	index := client.InitIndex(bot.BotConfig.AlgoliaIndex)

	res, err := index.Search(qry, nil)
	if err != nil {
		log.Info(fmt.Sprintf("ALGOLIA_SEARCH_ERROR [%v]", err))
		reqBody.Text = fmt.Sprintf("ALGOLIA_SEARCH_ERROR [%v].", err)
		return reqBody
	}

	for _, hit := range res.Hits {
		fmt.Println(hit["email"])

		attachment := rc.GlipMessageAttachmentInfoRequest{
			Type:   "Card",
			Fields: []rc.GlipMessageAttachmentFieldsInfo{},
		}

		if hit["name"] != nil {
			attachment.Fields = append(attachment.Fields,
				rc.GlipMessageAttachmentFieldsInfo{
					Title: "Name",
					Value: hit["name"].(string),
					Style: "Short"})
		}
		if hit["title"] != nil {
			attachment.Fields = append(attachment.Fields,
				rc.GlipMessageAttachmentFieldsInfo{
					Title: "Title",
					Value: hit["title"].(string),
					Style: "Short"})
		}
		if hit["email"] != nil {
			attachment.Fields = append(attachment.Fields,
				rc.GlipMessageAttachmentFieldsInfo{
					Title: "Email",
					Value: hit["email"].(string),
					Style: "Short"})
		}
		if hit["phone"] != nil {
			attachment.Fields = append(attachment.Fields,
				rc.GlipMessageAttachmentFieldsInfo{
					Title: "Phone",
					Value: hit["phone"].(string),
					Style: "Short"})
		}

		if hit["skills"] != nil {
			skills := strings.Join(strutil.InterfaceToSliceString(hit["skills"]), ", ")
			skills = markdown.BoldText(skills, qry)
			attachment.Fields = append(attachment.Fields,
				rc.GlipMessageAttachmentFieldsInfo{
					Title: "Skills",
					Value: skills})
		}

		if len(attachment.Fields) > 0 {
			mod := math.Mod(float64(len(reqBody.Attachments)), 2)
			if mod == 0 {
				attachment.Color = htmlutil.RingCentralOrangeHex
			} else {
				attachment.Color = htmlutil.RingCentralBlueHex
			}
			reqBody.Attachments = append(reqBody.Attachments, attachment)
		}
	}

	if len(reqBody.Attachments) > 0 {
		suffix := ""
		if len(reqBody.Attachments) > 1 {
			suffix = "s"
		}
		reqBody.Text = fmt.Sprintf("Displaying %v of %v matching user%s for **%s**.", len(reqBody.Attachments), res.NbHits, suffix, qry)
	} else {
		reqBody.Text = "No users found."
	}
	return reqBody
}
