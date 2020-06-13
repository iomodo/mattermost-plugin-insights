package store

import (
	"log"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

func (s *Store) GetPostCount() (int, error) {
	query :=
		`SELECT
			COUNT(Posts.Id) AS Value
		FROM
			Posts,
			Channels
		WHERE
			Posts.ChannelId = Channels.Id`

	var numPosts int
	err := s.db.QueryRow(query).Scan(&numPosts)
	s.log.Debugf("GetPostCount", numPosts)
	if err != nil {
		return 0, errors.Wrap(err, "failed to query database")
	}

	return numPosts, nil
}

type PostCountsOptions struct {
	TeamID    string
	ChannelID string
	BotsOnly  bool
	Start     int64
	End       int64
	Limit     int64
}

func (s *Store) GetPostCounts(options PostCountsOptions) (model.AnalyticsRows, error) {
	args := []interface{}{}
	query :=
		`SELECT
			DATE(FROM_UNIXTIME(Posts.CreateAt / 1000)) AS Name,
			COUNT(Posts.Id) AS Value
		FROM Posts`
	if options.BotsOnly {
		query += " INNER JOIN Bots ON Posts.UserId = Bots.Userid"
	}
	if len(options.TeamID) > 0 {
		query += " INNER JOIN Channels ON Posts.ChannelId = Channels.Id AND Channels.TeamId = ? AND"
		args = append(args, options.TeamID)
	}
	if options.Start != 0 && options.End != 0 {
		query += ` WHERE Posts.CreateAt <= ?
		AND Posts.CreateAt >= ? `
		args = append(args, options.Start)
		args = append(args, options.End)
	}

	query += ` GROUP BY DATE(FROM_UNIXTIME(Posts.CreateAt / 1000))
		ORDER BY Name DESC
		LIMIT ?`
	args = append(args, options.Limit)

	s.log.Debugf("query", query)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query database")
	}
	defer rows.Close()

	result := []*model.AnalyticsRow{}
	for rows.Next() {
		s.log.Debugf("next")

		var (
			Name  string
			Value float64
		)
		if err := rows.Scan(&Name, &Value); err != nil {
			log.Fatal(err)
		}
		s.log.Debugf("Name", Name)
		s.log.Debugf("Value", Value)

		result = append(result, &model.AnalyticsRow{Name: Name, Value: Value})
	}
	return result, nil
}

func MillisFromTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

func Yesterday() time.Time {
	return time.Now().AddDate(0, 0, -1)
}
