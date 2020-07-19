package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-insights/server/bot"
	"github.com/mattermost/mattermost-plugin-insights/server/insights"
)

type InsightsHandler struct {
	insights  insights.Service
	pluginAPI *pluginapi.Client
	log       bot.Logger
}

type Item struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

func NewInsightsHandler(insights insights.Service, router *mux.Router, api *pluginapi.Client, log bot.Logger) *InsightsHandler {
	handler := &InsightsHandler{
		insights:  insights,
		pluginAPI: api,
		log:       log,
	}
	insightsRouter := router.PathPrefix("/insights").Subrouter()
	insightsRouter.HandleFunc("/posts", handler.getPostData).Methods(http.MethodGet)
	insightsRouter.HandleFunc("/teams", handler.getTeams).Methods(http.MethodGet)
	insightsRouter.HandleFunc("/channels", handler.getChannels).Methods(http.MethodGet)
	return handler
}

func (h *InsightsHandler) getTeams(w http.ResponseWriter, r *http.Request) {
	println("in getTeams", r.URL.String())
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

func (h *InsightsHandler) getChannels(w http.ResponseWriter, r *http.Request) {
	println("in getChannels", r.URL.String())
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

func (h *InsightsHandler) getPostData(w http.ResponseWriter, r *http.Request) {
	println("in getPostData", r.URL.String())
	query := r.URL.Query()
	teamID := query.Get("team_id")
	channelID := query.Get("channel_id")

	rows := h.insights.GetPostCounts(teamID, channelID, "daily", "month", false)
	println(fmt.Sprintf("rows - %v", rows))
	b, _ := json.Marshal(rows)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
