package command

import (
	"fmt"
	"strings"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-insights/server/bot"
	"github.com/mattermost/mattermost-plugin-insights/server/insights"
	"github.com/mattermost/mattermost-server/v5/model"
)

const helpText = "###### Mattermost Workplace Insights Plugin - Slash Command Help\n" +
	"* `/insights posts` - Get post insights. \n" +
	""

// Command represents slash command of the plugin
type Command struct {
	args      *model.CommandArgs
	log       bot.Logger
	pluginAPI *pluginapi.Client
	poster    bot.Poster
	insight   insights.Service
}

// NewCommand creates new command
func NewCommand(args *model.CommandArgs, logger bot.Logger, api *pluginapi.Client, poster bot.Poster, insight insights.Service) *Command {
	return &Command{
		args:      args,
		log:       logger,
		pluginAPI: api,
		poster:    poster,
		insight:   insight,
	}
}

// Register is a function that allows the runner to register commands with the mattermost server.
type Register func(*model.Command) error

// RegisterCommands should be called by the plugin to register all necessary commands
func RegisterCommands(registerFunc Register) error {
	return registerFunc(getCommand())
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "insights",
		DisplayName:      "Insights Bot",
		Description:      "Get the workplace insights",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: ",
		AutoCompleteHint: "[command]",
	}
}

// Handle .
func (c *Command) Handle() error {
	split := strings.Fields(c.args.Command)
	command := split[0]
	parameters := []string{}
	cmd := ""
	if len(split) > 1 {
		cmd = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}
	println(parameters)

	if command != "/insights" {
		return nil
	}

	switch cmd {
	case "posts":
		c.posts()
	default:
		c.postCommandResponse(helpText)
	}

	return nil
}

func (c *Command) postCommandResponse(text string) {
	c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", text)
}

func (c *Command) posts() {
	rows := c.insight.GetPostCounts("", "", false, false)
	text := fmt.Sprintf("Number of rows are %d\n", len(rows))

	for _, row := range rows {
		text = text + row.Name + " - " + fmt.Sprintf("%f", row.Value) + "\n"
	}
	c.poster.Ephemeral(c.args.UserId, c.args.ChannelId, "%s", text)
}
