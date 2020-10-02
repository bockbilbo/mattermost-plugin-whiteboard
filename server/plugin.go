package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botID string
}

// OnActivate register the plugin command
func (p *Plugin) OnActivate() error {

	if err := p.API.RegisterCommand(createWhiteboardCommand()); err != nil {
		return err
	}
	whiteboardBot := &model.Bot{
		Username:    "whiteboard",
		DisplayName: "Whiteboard",
		Description: "A bot account created by the Whiteboard Plugin",
	}
	options := []plugin.EnsureBotOption{
		plugin.ProfileImagePath("assets/icon.png"),
	}

	botID, ensureBotError := p.Helpers.EnsureBot(whiteboardBot, options...)
	if ensureBotError != nil {
		return errors.Wrap(ensureBotError, "failed to ensure whiteboard bot user.")
	}

	p.botID = botID

	return nil
}

func (p *Plugin) deleteEphemeralPost(userID, postID string) {
	p.API.DeleteEphemeralPost(userID, postID)
}

func getUserName(user *model.User) string {
	var username = user.Username
	if len(user.LastName) != 0 {
		username = user.FirstName + " " + user.LastName
	}
	return username
}

func (p *Plugin) startWhiteboard(user *model.User, channel *model.Channel, whiteboardID string) (string, error) {
	config := p.getConfiguration()

	whiteboardURL := strings.TrimRight(strings.TrimSpace(config.WhiteboardURL), "/")

	if len(whiteboardURL) == 0 {
		switch config.WhiteboardTool {
		case "cracker0dks":
			whiteboardURL = "https://cloud13.de/testwhiteboard/"

		case "excalidraw":
			whiteboardURL = "https://excalidraw.com"

		default:
			whiteboardURL = "https://wbo.ophir.dev"
		}
	}

	urlPath := ""

	switch config.WhiteboardTool {
	case "cracker0dks":
		urlPath = "?whiteboardid="

	case "excalidraw":
		urlPath = "#room="

	default:
		urlPath = "boards/"
	}

	if config.WhiteboardTool == "excalidraw" {
		whiteboardID = generateExcalidrawID()
	} else {
		if whiteboardID == "" {
			whiteboardID = generateEnglishTitleName()
		}
	}

	whiteboardLink := fmt.Sprintf("%s/%s%s", whiteboardURL, urlPath, whiteboardID)

	slackWhiteboardTopic := "New Whiteboard"

	slackAttachment := model.SlackAttachment{
		Fallback: "Whiteboard shared at [" + whiteboardID + "](" + whiteboardLink + ").\n\n[Join Whiteboard](" + whiteboardLink + ")",
		Title:    slackWhiteboardTopic,
		Text:     "[" + whiteboardID + "](" + whiteboardURL + ")\n\n",
	}

	post := &model.Post{
		UserId:    p.botID,
		ChannelId: channel.Id,
		Type:      "custom_whiteboard",
		Props: map[string]interface{}{
			"attachments":                 []*model.SlackAttachment{&slackAttachment},
			"whiteboard_id":               whiteboardID,
			"whiteboard_link":             whiteboardLink,
			"whiteboard_creator_username": getUserName(user),
		},
	}

	if _, err := p.API.CreatePost(post); err != nil {
		return "", err
	}

	return "ok", nil
}
