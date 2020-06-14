package api

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (h *Handler) getTeams(w http.ResponseWriter, r *http.Request) {
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
