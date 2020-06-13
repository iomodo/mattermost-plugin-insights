package api

import (
	"bytes"
	"errors"
	"net/http"
	"sync"

	pluginapi "github.com/mattermost/mattermost-plugin-api"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-plugin-insights/server/bot"
)

// ChartsHandler is the API handler.
type ChartsHandler struct {
	pluginAPI *pluginapi.Client
	poster    bot.Poster
	log       bot.Logger
	charts    map[string]*bytes.Buffer
	mux       sync.Mutex
}

// NewChartsHandler Creates a new Plugin API handler.
func NewChartsHandler(router *mux.Router, api *pluginapi.Client, poster bot.Poster, log bot.Logger) *ChartsHandler {
	handler := &ChartsHandler{
		pluginAPI: api,
		poster:    poster,
		log:       log,
		charts:    map[string]*bytes.Buffer{},
	}

	chartsRouter := router.PathPrefix("/charts").Subrouter()
	chartsRouter.HandleFunc("", handler.getCharts).Methods(http.MethodGet)

	return handler
}

func (h *ChartsHandler) AddChart(id string, chart *bytes.Buffer) {
	h.mux.Lock()
	h.charts[id] = chart
	h.mux.Unlock()
}

func (h *ChartsHandler) getCharts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	chartID := query["id"]

	h.mux.Lock()
	chart, ok := h.charts[chartID[0]]
	h.mux.Unlock()
	if !ok {
		HandleErrorWithCode(w, http.StatusOK, "can't find chart", errors.New("can't find chart"))
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(chart.Bytes())
}
