package store

import (
	"database/sql"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-insights/server/bot"
)

// Store .
type Store struct {
	db     *sql.DB
	driver string
	log    bot.Logger
}

// NewStore .
func NewStore(pluginAPI *pluginapi.Client, log bot.Logger) *Store {
	db, err := pluginAPI.Store.GetReplicaDB()
	if err != nil {
		log.Errorf("error while getting DB replica", err)
		return nil
	}
	driverName := pluginAPI.Store.DriverName()
	return &Store{
		db:     db,
		driver: driverName,
		log:    log,
	}
}

func (s *Store) Close() {
	if err := s.db.Close(); err != nil {
		s.log.Errorf("error while closing DB", err)
	}
}
