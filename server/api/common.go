package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Item struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

func (h *Handler) getTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.pluginAPI.Team.List()
	if err != nil {
		HandleError(w, err)
		return
	}
	a := make([]Item, 0, len(teams))

	for _, team := range teams {
		a = append(a, Item{ID: team.Id, DisplayName: team.DisplayName})
	}
	b, _ := json.Marshal(a)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handler) getChannels(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	teamID := query.Get("team_id")
	page, _ := strconv.Atoi(query.Get("page"))
	perPage, _ := strconv.Atoi(query.Get("per_page"))

	channels, err := h.pluginAPI.Channel.ListPublicChannelsForTeam(teamID, page, perPage)
	if err != nil {
		HandleError(w, err)
		return
	}
	a := make([]Item, 0, len(channels))

	for _, channel := range channels {
		a = append(a, Item{ID: channel.Id, DisplayName: channel.DisplayName})
	}
	b, _ := json.Marshal(a)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
