package store

import (
	"fmt"
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

type FrequencyType string

const (
	FrequencyTypeDaily   FrequencyType = "daily"
	FrequencyTypeWeekly  FrequencyType = "weekly"
	FrequencyTypeMonthly FrequencyType = "monthly"
	FrequencyTypeYearly  FrequencyType = "yearly"
)

type PeriodType string

const (
	PeriodTypeWeek  PeriodType = "week"
	PeriodTypeMonth PeriodType = "month"
	PeriodTypeYear  PeriodType = "year"
	PeriodTypeAll   PeriodType = "all"
)

type PostCountsOptions struct {
	TeamID    string
	ChannelID string
	BotsOnly  bool
	Start     int64
	End       int64
	Limit     int64
	Frequency FrequencyType
	Period    PeriodType
}

func (s *Store) GetPostCounts(options PostCountsOptions) (model.AnalyticsRows, error) {
	args := []interface{}{}
	query :=
		`SELECT
			%s(FROM_UNIXTIME(Posts.CreateAt / 1000)) AS Name,
			COUNT(Posts.Id) AS Value
		FROM Posts`
	s.log.Errorf("error %v, %v, %v", s.driver, model.DATABASE_DRIVER_MYSQL, model.DATABASE_DRIVER_POSTGRES)
	if s.driver == model.DATABASE_DRIVER_POSTGRES {
		query = `SELECT
			TO_CHAR(%s(TO_TIMESTAMP(Posts.CreateAt / 1000)), 'YYYY-MM-DD') AS Name,
			COUNT(Posts.Id) AS Value
			FROM Posts`
	}

	freq := "DATE"
	if options.Frequency == FrequencyTypeDaily {
		freq = "DATE"
	} else if options.Frequency == FrequencyTypeWeekly {
		freq = "WEEK"
	} else if options.Frequency == FrequencyTypeMonthly {
		freq = "MONTH"
	} else if options.Frequency == FrequencyTypeYearly {
		freq = "YEAR"
	}
	query = fmt.Sprintf(query, freq)

	if options.BotsOnly {
		query += " INNER JOIN Bots ON Posts.UserId = Bots.Userid" //TODO fix
	}

	param := "WHERE"

	if len(options.TeamID) > 0 {
		query += " INNER JOIN Channels ON Posts.ChannelId = Channels.Id AND Channels.TeamId = $1 AND"
		args = append(args, options.TeamID)
		param = ""
		if len(options.ChannelID) > 0 {
			query += " Posts.ChannelId = $2 AND"
			args = append(args, options.ChannelID)
		}
	}
	if options.Period == PeriodTypeWeek {
		t := time.Now().Add(time.Hour * time.Duration(-24*7))
		args = append(args, t.Unix()*1000)
	} else if options.Period == PeriodTypeMonth {
		t := time.Now().Add(time.Hour * time.Duration(-24*30))
		args = append(args, t.Unix()*1000)
	} else if options.Period == PeriodTypeYear {
		t := time.Now().Add(time.Hour * time.Duration(-24*365))
		args = append(args, t.Unix()*1000)
	}
	query += fmt.Sprintf(` %s Posts.CreateAt >= $%d`, param, len(args))

	// if options.Start != 0 && options.End != 0 {
	// 	query += ` WHERE Posts.CreateAt <= ?
	// 	AND Posts.CreateAt >= ? `
	// 	args = append(args, options.Start)
	// 	args = append(args, options.End)
	// }

	args = append(args, options.Limit)
	if s.driver == model.DATABASE_DRIVER_MYSQL {
		query += fmt.Sprintf(` GROUP BY %s(FROM_UNIXTIME(Posts.CreateAt / 1000))
		ORDER BY Name DESC
		LIMIT $%d`, freq, len(args))
	} else if s.driver == model.DATABASE_DRIVER_POSTGRES {
		query += fmt.Sprintf(` GROUP BY TO_CHAR(%s(TO_TIMESTAMP(Posts.CreateAt / 1000)), 'YYYY-MM-DD')
		ORDER BY Name DESC
		LIMIT $%d`, freq, len(args))
	}

	s.log.Debugf("query", query)
	s.log.Debugf("args", args)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query database")
	}
	defer rows.Close()

	result := []*model.AnalyticsRow{}
	for rows.Next() {
		var (
			Name  string
			Value float64
		)
		if err := rows.Scan(&Name, &Value); err != nil {
			log.Fatal(err)
		}
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

type PostOptions struct {
	ChannelID string
	From      int // days before today
	To        int // days before today
	Limit     int // number of posts
}

type MessagesRow struct {
	Message string
}

func (s *Store) GetPosts(options PostOptions) ([]*MessagesRow, error) {
	args := []interface{}{}
	query :=
		`SELECT
			Posts.Message AS Message
		FROM Posts WHERE
		Posts.ChannelId = $1 AND Posts.CreateAt >= $2 AND Posts.CreateAt < $3 LIMIT $4`
	args = append(args, options.ChannelID)
	from := time.Now().Add(time.Hour * time.Duration(-24*options.From))
	args = append(args, from.Unix()*1000)
	to := time.Now().Add(time.Hour * time.Duration(-24*options.To))
	args = append(args, to.Unix()*1000)
	args = append(args, options.Limit)

	rows, err := s.db.Query(query, args...)
	if err != nil || rows == nil {
		return nil, errors.Wrap(err, "failed to query database")
	}
	defer rows.Close()

	result := []*MessagesRow{}
	for rows.Next() {
		var (
			Message string
		)
		if err := rows.Scan(&Message); err != nil {
			log.Fatal(err)
		}
		result = append(result, &MessagesRow{Message: Message})
	}
	return result, nil
}
