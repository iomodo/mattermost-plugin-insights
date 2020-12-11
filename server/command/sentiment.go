package command

import (
	"fmt"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"
)

func createSentimentAutocompleteData() *model.AutocompleteData {
	sentiment := model.NewAutocompleteData("sentiment", "", "get channel sentiment")
	sentiment.AddNamedTextArgument("span", "Specify number of days", "", "", false)
	return sentiment
}

func (c *Command) handleSentiment(parameters []string) {
	span := 0
	var err error
	if len(parameters) > 0 && parameters[0] == "--span" {
		span, err = strconv.Atoi(parameters[1])
		if err != nil {
			c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", "please input integer")
		}
		parameters = parameters[2:]
	}
	if span == 0 {
		sentiment := c.insight.GetSentiment(c.args.ChannelId)
		text := fmt.Sprintf("Average Happiness score for the last month is %d", int(sentiment*100))
		c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", text)
		return
	}

}
