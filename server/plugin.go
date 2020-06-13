package main

import (
	"fmt"
	"net/http"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-insights/server/bot"
	"github.com/mattermost/mattermost-plugin-insights/server/command"
	"github.com/mattermost/mattermost-plugin-insights/server/config"
	"github.com/mattermost/mattermost-plugin-insights/server/insights"
	"github.com/mattermost/mattermost-plugin-insights/server/store"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	config  *config.ServiceImpl
	bot     *bot.Bot
	insight insights.Service
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// OnActivate Called when this plugin is activated.
func (p *Plugin) OnActivate() error {
	pluginAPIClient := pluginapi.NewClient(p.API)
	p.config = config.NewConfigService(pluginAPIClient)
	pluginapi.ConfigureLogrus(logrus.New(), pluginAPIClient)

	botID, err := pluginAPIClient.Bot.EnsureBot(&model.Bot{
		Username:    "insights",
		DisplayName: "Insights Bot",
		Description: "A prototype demonstrating workplace insights in Mattermost.",
	},
		pluginapi.ProfileImagePath("assets/profile.png"),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to ensure insights bot")
	}

	err = p.config.UpdateConfiguration(func(c *config.Configuration) {
		c.BotUserID = botID
		c.AdminLogLevel = "debug"
	})
	if err != nil {
		return errors.Wrapf(err, "failed save bot to config")
	}

	p.bot = bot.New(pluginAPIClient, p.config.GetConfiguration().BotUserID, p.config)

	st := store.NewStore(pluginAPIClient, p.bot)
	p.insight = insights.NewService(pluginAPIClient, st, p.bot, p.config)

	if err := command.RegisterCommands(p.API.RegisterCommand); err != nil {
		return errors.Wrapf(err, "failed register commands")
	}

	p.API.LogDebug("Insights plugin Activated")
	return nil
}

// ExecuteCommand executes a given command and returns a command response.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	com := command.NewCommand(args, p.bot, pluginapi.NewClient(p.API), p.bot, p.insight)
	if err := com.Handle(); err != nil {
		return nil, model.NewAppError("insights.ExecuteCommand", "Unable to execute command.", nil, err.Error(), http.StatusInternalServerError)
	}

	return &model.CommandResponse{}, nil
}
