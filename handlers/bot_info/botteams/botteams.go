package botteams

import (
	"net/http"

	rc "github.com/grokify/go-ringcentral/client"
	"github.com/grokify/go-ringcentral/clientutil/glipgroups"
	hum "github.com/grokify/gotilla/net/httputilmore"

	"github.com/grokify/chatblox"
)

type Factory struct {
	HTTPClient *http.Client
	APIClient  *rc.APIClient
	ServerURL  string
}

func (f *Factory) GetGlipTeams() (glipgroups.GroupsSet, error) {
	return glipgroups.NewGroupsSetApiRequest(f.HTTPClient, f.ServerURL, "Team")
}

func (f *Factory) NewIntent() chatblox.Intent {
	return chatblox.Intent{
		Type:         chatblox.MatchStringLowerCase,
		Strings:      []string{"bot teams"},
		HandleIntent: f.HandleIntent,
	}
}

func (f *Factory) HandleIntent(bot *chatblox.Bot, slots map[string]string, glipPostEventInfo *chatblox.GlipPostEventInfo) (*hum.ResponseInfo, error) {
	teamsSet, err := f.GetGlipTeams()
	if err != nil {
		return bot.SendGlipPost(glipPostEventInfo,
			rc.GlipCreatePost{Text: "I'm sorry but I ran into a problem."})
	}
	teams := teamsSet.GroupNamesSorted(true)

	text := ""
	if len(teams) > 0 {
		text = "I'm in the following teams:\n\n"
		for _, t := range teams {
			text += "* " + t
		}
	} else {
		text = "I'm not currently in any teams."
	}

	return bot.SendGlipPost(glipPostEventInfo,
		rc.GlipCreatePost{
			Text: text})
}
