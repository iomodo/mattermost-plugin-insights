package command

import (
	"errors"

	"github.com/mattermost/mattermost-plugin-insights/server/chart"
	"github.com/mattermost/mattermost-plugin-insights/server/utils"
	"github.com/mattermost/mattermost-server/v5/model"
)

func createPostsAutocompleteData() *model.AutocompleteData {
	posts := model.NewAutocompleteData("posts", "", "Get insights for posts")
	posts.AddNamedDynamicListArgument("team", "Specify team", "teams_for_command", false)
	posts.AddNamedDynamicListArgument("channel", "Specify channel", "channels_for_commands", false)
	frequencyItems := []model.AutocompleteListItem{{
		Item:     "daily",
		Hint:     "",
		HelpText: "Aggregate per day",
	}, {
		Item:     "weekly",
		Hint:     "",
		HelpText: "Aggregate per week",
	}, {
		Item:     "monthly",
		Hint:     "",
		HelpText: "Aggregate per month",
	}, {
		Item:     "yearly",
		Hint:     "",
		HelpText: "Aggregate per year",
	},
	}
	posts.AddNamedStaticListArgument("frequency", "Specify frequency of the statistics", false, frequencyItems)
	spanItems := []model.AutocompleteListItem{{
		Item:     "week",
		Hint:     "",
		HelpText: "Last week statistics",
	}, {
		Item:     "month",
		Hint:     "",
		HelpText: "Last month statistics",
	}, {
		Item:     "3months",
		Hint:     "",
		HelpText: "Last 3 months statistics",
	}, {
		Item:     "6months",
		Hint:     "",
		HelpText: "Last 6 months statistics",
	}, {
		Item:     "year",
		Hint:     "",
		HelpText: "Last year statistics",
	},
	}
	posts.AddNamedStaticListArgument("span", "Timespan of the statistics", false, spanItems)
	return posts
}

func (c *Command) handlePosts(parameters []string) {
	team := ""
	channel := ""
	frequency := "daily"
	span := "month"
	if len(parameters) > 0 && parameters[0] == "--team" {
		team = parameters[1]
		parameters = parameters[2:]
	}
	if len(parameters) > 0 && parameters[0] == "--channel" {
		channel = parameters[1]
		parameters = parameters[2:]
	}
	if len(parameters) > 0 && parameters[0] == "--frequency" {
		frequency = parameters[1]
		parameters = parameters[2:]
	}
	if len(parameters) > 0 && parameters[0] == "--span" {
		span = parameters[1]
		parameters = parameters[2:]
	}

	team, _ = c.getTeamIDFromTeamName(team)
	channel, _ = c.getChannelIDFromChannelAndTeam(team, channel)

	rows := c.insight.GetPostCounts(team, channel, frequency, span, false)
	chart := chart.CreateBarChart("Posts per day", rows)
	id := utils.NewID()
	c.charts.AddChart(id, chart)

	siteURL := "http://localhost:8065"
	if c.pluginAPI.Configuration.GetConfig().ServiceSettings.SiteURL != nil {
		siteURL = *c.pluginAPI.Configuration.GetConfig().ServiceSettings.SiteURL
	}
	chartURL := siteURL + "/plugins/com.mattermost.plugin-insights/api/v1/charts?id=" + id

	text := "![" + "Posts per day" + "](" + chartURL + ")"
	c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", text)
}

func (c *Command) getTeamIDFromTeamName(teamName string) (string, error) {
	teams, err := c.pluginAPI.Team.List()
	if err != nil {
		return "", err
	}
	for _, team := range teams {
		if team.Name == teamName {
			return team.Id, nil
		}
	}
	return "", errors.New("team not found")
}

func (c *Command) getChannelIDFromChannelAndTeam(teamID, channelName string) (string, error) {
	perPage := 1000
	for page := 0; ; page++ {
		channels, err := c.pluginAPI.Channel.ListPublicChannelsForTeam(teamID, page, perPage)
		if err != nil {
			return "", err
		}
		if len(channels) == 0 {
			return "", errors.New("channel not found")
		}
		for _, channel := range channels {
			if channel.Name == channelName {
				return channel.Id, nil
			}
		}
	}
}
