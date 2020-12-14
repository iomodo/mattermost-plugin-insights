package insights

import "github.com/mattermost/mattermost-plugin-insights/server/store"

func (s *ServiceImpl) GetSentiment(channelID string) float64 {
	return s.getAverageSentiment(channelID, 30, 0)
}

func (s *ServiceImpl) GetSentiments(channelID string, days int) []float64 {
	res := []float64{}
	for day := 0; day < days; day++ {
		score := s.getAverageSentiment(channelID, day+1, day)
		res = append(res, score)
	}
	return res
}

func (s *ServiceImpl) getAverageSentiment(channelID string, from, to int) float64 {
	messages, err := s.store.GetPosts(store.PostOptions{
		ChannelID: channelID,
		From:      from,
		To:        to,
		Limit:     1000,
	})
	if err != nil {
		s.pluginAPI.Log.Debug("error getting post messages", err)
	}
	if len(messages) == 0 {
		return 0
	}
	scoreSum := float64(0)
	for _, message := range messages {
		scoreSum += s.sentimentAnalyzer.GetSentiment(message.Message)
	}
	return scoreSum / float64(len(messages))
}
