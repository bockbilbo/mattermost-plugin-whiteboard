package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const externalAPICacheTTL = 3600000

var externalAPICache []byte
var externalAPILastUpdate int64
var externalAPICacheMutex sync.Mutex

// StartWhiteboardRequest method
type StartWhiteboardRequest struct {
	ChannelID    string `json:"channel_id"`
	WhiteboardID int    `json:"whiteboard_id"`
}

// StartWhiteboardFromAction method
type StartWhiteboardFromAction struct {
	model.PostActionIntegrationRequest
	Context struct {
		WhiteboardID string `json:"whiteboard_id"`
	} `json:"context"`
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/api/v1/whiteboards":
		p.handleStartWhiteboard(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleStartWhiteboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")

	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		mlog.Debug("Unable to the user", mlog.Err(appErr))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req StartWhiteboardRequest
	var action StartWhiteboardFromAction

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		mlog.Debug("Unable to read request body", mlog.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err1 := json.NewDecoder(bytes.NewReader(bodyData)).Decode(&req)
	err2 := json.NewDecoder(bytes.NewReader(bodyData)).Decode(&action)
	if err1 != nil && err2 != nil {
		mlog.Debug("Unable to decode the request content as start whiteboard request or start whiteboard action")
		http.Error(w, "Unable to decode your request!", http.StatusBadRequest)
		return
	}

	channelID := req.ChannelID
	if channelID == "" {
		channelID = action.ChannelId
	}

	if _, err := p.API.GetChannelMember(channelID, userID); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	channel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var whiteboardID string
	if action.PostId != "" {
		whiteboardID, err = p.startWhiteboard(user, channel, action.Context.WhiteboardID)
		if err != nil {
			mlog.Error("Error starting a new whiteboard from ask response", mlog.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.deleteEphemeralPost(action.UserId, action.PostId)
	} else {
		whiteboardID, err = p.startWhiteboard(user, channel, "")
		if err != nil {
			mlog.Error("Error starting a new whiteboard", mlog.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	b, err := json.Marshal(map[string]string{"whiteboard_id": whiteboardID})
	if err != nil {
		mlog.Error("Error marshaling the WhiteboardID to json", mlog.Err(err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		mlog.Warn("Unable to write response body", mlog.String("handler", "handleStartWhiteboard"), mlog.Err(err))
	}
}
