package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const whiteboardCommand = "whiteboard"

func startWhiteboardError(channelID string, detailedError string) (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			ChannelId:    channelID,
			Text:         "We could not share a whiteboard at this time.",
		}, &model.AppError{
			Message:       "We could not share a whiteboard at this time.",
			DetailedError: detailedError,
		}
}

func createWhiteboardCommand() *model.Command {
	return &model.Command{
		Trigger:          whiteboardCommand,
		AutoComplete:     true,
		AutoCompleteDesc: "Share a virtual whiteboard in the current channel with an optional whiteboard name (ignored when using Excalidraw)",
		AutoCompleteHint: "[whiteboard_name]",
	}
}

func getSanitizedName(whiteboardName string) string {
	// We only allow letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(whiteboardName, "")
}

//ExecuteCommand method
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]

	if command != "/"+whiteboardCommand {
		return &model.CommandResponse{}, nil
	}

	return p.executeStartWhiteboardCommand(c, args)
}

func (p *Plugin) executeStartWhiteboardCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	whiteboardID := getSanitizedName(strings.TrimSpace(strings.TrimPrefix(args.Command, "/"+whiteboardCommand)))

	user, appErr := p.API.GetUser(args.UserId)
	if appErr != nil {
		return startWhiteboardError(args.ChannelId, fmt.Sprintf("getUser() threw error: %s", appErr))
	}

	channel, appErr := p.API.GetChannel(args.ChannelId)
	if appErr != nil {
		return startWhiteboardError(args.ChannelId, fmt.Sprintf("getChannel() threw error: %s", appErr))
	}

	if _, err := p.startWhiteboard(user, channel, whiteboardID); err != nil {
		return startWhiteboardError(args.ChannelId, fmt.Sprintf("startWhiteboard() threw error: %s", appErr))
	}

	return &model.CommandResponse{}, nil
}
