package command

import (
	"github.com/mattermost/mattermost-plugin-insights/server/chart"
	"github.com/mattermost/mattermost-plugin-insights/server/utils"
	"github.com/mattermost/mattermost-server/v5/model"
)

func createPostsAutocompleteData() *model.AutocompleteData {
	posts := model.NewAutocompleteData("posts", "", "Get insights for posts")
	posts.AddNamedDynamicListArgument("team", "Specify team", "teams_for_command", false)
	posts.AddNamedTextArgument("channel", "Specify channel as ~some_channel", "[optional]", "", false)
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
