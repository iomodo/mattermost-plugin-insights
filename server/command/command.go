package command

import (
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-insights/server/bot"
	"github.com/mattermost/mattermost-server/v5/model"
)

// Command represents slash command of the plugin
type Command struct {
	Args      *model.CommandArgs
	Log       bot.Logger
	pluginAPI *pluginapi.Client
	poster    bot.Poster
}

type logger interface {
	LogError(msg string, keyValuePairs ...interface{})
}

// NewCommand creates new command
func NewCommand(args *model.CommandArgs, logger bot.Logger, api *pluginapi.Client, poster bot.Poster) *Command {
	return &Command{
		Args:      args,
		Log:       logger,
		pluginAPI: api,
		poster:    poster,
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

func (c *Command) Handle() error {
	return nil
}
