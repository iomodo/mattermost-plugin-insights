package insights

import "github.com/mattermost/mattermost-plugin-insights/server/store"

func (s *ServiceImpl) GetSentiment(channelID string) float64 {
	return s.getLastMonthAverageSentiment(channelID)
}

func (s *ServiceImpl) GetSentiments(days int) []float64 {
	return nil
}

func (s *ServiceImpl) getLastMonthAverageSentiment(channelID string) float64 {
	messages, err := s.store.GetPosts(store.PostOptions{
		ChannelID: channelID,
		From:      30,
		To:        0,
		Limit:     1000,
	})
	if err != nil {
		s.pluginAPI.Log.Debug("error getting post messages", err)
	}
	scoreSum := float64(0)
	for _, message := range messages {
		scoreSum += s.sentimentAnalyzer.GetSentiment(message.Message)
	}

	return scoreSum / float64(len(messages))
}
