package insights

import "github.com/mattermost/mattermost-server/v5/model"

// Service is the insights/service interface.
type Service interface {
	GetPosts() int
	GetPostCounts(teamID, channelID, frequency, span string, botsOnly bool) model.AnalyticsRows
	GetSentiment(channelID string) float64
	GetSentiments(channelID string, days int) []float64
}
