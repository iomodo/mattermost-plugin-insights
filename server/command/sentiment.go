package command

import (
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
)

func createSentimentAutocompleteData() *model.AutocompleteData {
	sentiment := model.NewAutocompleteData("sentiment", "", "get channel sentiment")
	sentiment.AddNamedTextArgument("period", "Specify number of days", "", "", false)
	return sentiment
}

func (c *Command) handleSentiment(parameters []string) {
	period := 0
	var err error
	if len(parameters) > 0 && parameters[0] == "--period" {
		period, err = strconv.Atoi(parameters[1])
		if err != nil {
			c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", "please input integer")
		}
		parameters = parameters[2:]
	}
	if period == 0 {
		sentiment := c.insight.GetSentiment(c.args.ChannelId)
		text := fmt.Sprintf("Average Happiness score for the last month is %d", int(sentiment*100))
		c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", text)
		return
	}
	sentiments := c.insight.GetSentiments(c.args.ChannelId, period)
	res := ""
	for _, sentiment := range sentiments {
		res = strconv.Itoa(int(sentiment*100)) + ", " + res
	}
	res = "[" + res[:len(res)-2] + "]"
	text := fmt.Sprintf("Average Happiness score for the last %d days is %v", period, res)
	c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", text)
}
