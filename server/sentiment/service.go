package sentiment

import (
	"github.com/cdipaolo/sentiment"
)

type Sentiment interface {
	GetSentiment(post string) float64
}

type service struct {
	model sentiment.Models
}

func NewSentimentAnalyzer() Sentiment {
	model, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}
	return &service{model: model}
}

func (s *service) GetSentiment(post string) float64 {
	analysis := s.model.SentimentAnalysis(post, sentiment.English)
	return float64(analysis.Score)
}
