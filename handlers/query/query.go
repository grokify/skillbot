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

	"github.com/grokify/chatblox"
)

/*

Who can help me book a flight with Delta?
help with Delta
who can help me with detla
I need a delta expert



*/

const Fields = "name,title,skills"

var fieldsMap = stringToMap(Fields, ",")

func NewIntent() chatblox.Intent {
	return chatblox.Intent{
		Type: chatblox.MatchRegexpCapture,
		Regexps: []*regexp.Regexp{
			regexp.MustCompile(`^(?i)^\s*(?:I need an?|Find(?: me)? an?|find) (?P<query_any>.+)(?:\s+expert)\s*[\?.]?\s*$`),
			regexp.MustCompile(`^(?i)^\s*(?:who can help me book a flight with|who can help(?: me)? with) (?P<query_any>.+?)\s*[\?.]?\s*$`),
			regexp.MustCompile(`^(?i)^\s*(?:find|query|search|help with|skill) (?P<query_any>.+?)\s*[\?.]?\s*$`),
		},
		HandleIntent: handleIntent,
	}
}

func stringToMap(s, delim string) map[string]int {
	fields := map[string]int{}
	arr := strings.Split(s, delim)
	for _, item := range arr {
		item := strings.TrimSpace(item)
		if len(item) > 0 {
			fields[item] = 1
		}
	}
	return fields
}

func handleIntent(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) (*hum.ResponseInfo, error) {
	glipPost := buildPost(bot, slots, glipPostEventInfo)
	return bot.SendGlipPost(glipPostEventInfo, glipPost)
}

func buildPost(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) rc.GlipCreatePost {
	log.Info("BUILD_POST_EXPERTFINDER_QUERY")
	reqBody := rc.GlipCreatePost{}

	qry, ok := slots["query_any"]
	qry = strings.TrimSpace(qry)

	if !ok || len(qry) == 0 {
		reqBody.Text = "I don't understand what you are asking for. Please try asking 'search <your skill>'"
		return reqBody
	}

	bot.BotConfig.AlgoliaIndex = strings.TrimSpace(bot.BotConfig.AlgoliaIndex)
	if len(bot.BotConfig.AlgoliaIndex) == 0 {
		reqBody.Text = "I have a configuration problem. Please contact support."
		return reqBody
	}

	client, err := chatblox.GetAlgoliaApiClient(bot.BotConfig)
	if err != nil {
		log.Info(fmt.Sprintf("ALGOLIA_CLIENT_ERROR [%v]", err))
		reqBody.Text = "I had a problem opening my database. Please try again."
		return reqBody
	}
	index := client.InitIndex(bot.BotConfig.AlgoliaIndex)

	res, err := index.Search(qry, nil)
	if err != nil {
		log.Info(fmt.Sprintf("ALGOLIA_SEARCH_ERROR [%v]", err))
		reqBody.Text = "I had a problem searching my directory. Please try again."
		return reqBody
	}

	for _, hit := range res.Hits {
		fmt.Println(hit["email"])

		attachment := rc.GlipMessageAttachmentInfoRequest{
			Type:   "Card",
			Fields: []rc.GlipMessageAttachmentFieldsInfo{},
		}

		avatarUri := ""
		if hit["avatar"] != nil {
			avatarUri = strings.TrimSpace(hit["avatar"].(string))
		}

		if hit["name"] != nil {
			if 1 == 1 && len(avatarUri) > 0 {
				attachment.Author = rc.GlipMessageAttachmentAuthorInfo{
					Name:    hit["name"].(string),
					IconUri: hit["avatar"].(string)}
			} else {
				attachment.Fields = append(attachment.Fields,
					rc.GlipMessageAttachmentFieldsInfo{
						Title: "Name",
						Value: markdown.BoldText(hit["name"].(string), ""),
						Style: "Short"})
			}
		}
		if hit["title"] != nil {
			attachment.Fields = append(attachment.Fields,
				rc.GlipMessageAttachmentFieldsInfo{
					Title: "Title",
					Value: markdown.BoldText(hit["title"].(string), qry),
					Style: "Short"})
		}
		if _, ok := fieldsMap["email"]; ok {
			if hit["email"] != nil {
				attachment.Fields = append(attachment.Fields,
					rc.GlipMessageAttachmentFieldsInfo{
						Title: "Email",
						Value: hit["email"].(string),
						Style: "Short"})
			}
		}
		if _, ok := fieldsMap["phone"]; ok {
			if hit["phone"] != nil {
				attachment.Fields = append(attachment.Fields,
					rc.GlipMessageAttachmentFieldsInfo{
						Title: "Phone",
						Value: hit["phone"].(string),
						Style: "Short"})
			}
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

	// Here is a person who has skills matching "Delta".

	if len(reqBody.Attachments) > 0 {
		text := ""
		if len(reqBody.Attachments) > 1 {
			text = fmt.Sprintf("Here are some people who have skills matching **%v**", qry)
		} else if len(reqBody.Attachments) == 1 {
			text = fmt.Sprintf("Here is a person who has skills matching **\"%v\"**", qry)
		} else {
			text = fmt.Sprintf("I could not find anyone with skills matching **%v**", qry)
		}
		reqBody.Text = text
		//reqBody.Text = fmt.Sprintf("Displaying %v of %v matching user%s for **%s**.", len(reqBody.Attachments), res.NbHits, suffix, qry)
	} else {
		reqBody.Text = "No users found."
	}
	return reqBody
}
