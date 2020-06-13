package insights

import (
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/mattermost/mattermost-plugin-insights/server/api"
	"github.com/mattermost/mattermost-plugin-insights/server/bot"
	"github.com/mattermost/mattermost-plugin-insights/server/config"
	"github.com/mattermost/mattermost-plugin-insights/server/store"
)

// ServiceImpl holds the information needed by the InsightsService's methods to complete their functions.
type ServiceImpl struct {
	pluginAPI     *pluginapi.Client
	configService config.Service
	store         *store.Store
	poster        bot.Poster
	handler       *api.ChartsHandler
}

// NewService creates a new insights ServiceImpl.
func NewService(pluginAPI *pluginapi.Client, store *store.Store, poster bot.Poster, configService config.Service, handler *api.ChartsHandler) *ServiceImpl {
	return &ServiceImpl{
		pluginAPI:     pluginAPI,
		store:         store,
		poster:        poster,
		configService: configService,
		handler:       handler,
	}
}

func (s *ServiceImpl) GetPosts() int {
	count, err := s.store.GetPostCount()
	if err != nil {
		s.pluginAPI.Log.Debug("error in get posts", err)
	}
	return count
}

func (s *ServiceImpl) GetPostCounts(teamID, channelID string, botsOnly, yesterdayOnly bool) model.AnalyticsRows {
	rows, err := s.store.GetPostCounts(store.PostCountsOptions{
		TeamID:    teamID,
		ChannelID: channelID,
		BotsOnly:  botsOnly,
		Start:     0,
		End:       0,
		Limit:     50,
	})
	if err != nil {
		s.pluginAPI.Log.Debug("error getting post counts", err)
	}
	return rows
}
