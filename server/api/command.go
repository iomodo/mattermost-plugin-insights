package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (h *Handler) getTeamsForCommand(w http.ResponseWriter, r *http.Request) {
	teams, err := h.pluginAPI.Team.List()
	if err != nil {
		HandleError(w, err)
		return
	}
	a := make([]model.AutocompleteListItem, 0, len(teams))

	for _, team := range teams {
		a = append(a, model.AutocompleteListItem{HelpText: team.DisplayName + ": " + team.Description, Item: team.Name})
	}
	b, _ := json.Marshal(a)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handler) getChannelsForCommand(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	teamID := query.Get("team_id")
	page, _ := strconv.Atoi(query.Get("page"))
	perPage, _ := strconv.Atoi(query.Get("per_page"))
	if perPage > 10 || perPage == 0 {
		perPage = 10
	}

	channels, err := h.pluginAPI.Channel.ListPublicChannelsForTeam(teamID, page, perPage)
	if err != nil {
		HandleError(w, err)
		return
	}
	a := make([]model.AutocompleteListItem, 0, len(channels))

	for _, channel := range channels {
		a = append(a, model.AutocompleteListItem{HelpText: channel.DisplayName + ": " + channel.Purpose, Item: channel.Name})
	}
	b, _ := json.Marshal(a)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
