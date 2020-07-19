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
		query += " INNER JOIN Channels ON Posts.ChannelId = Channels.Id AND Channels.TeamId = ? AND"
		args = append(args, options.TeamID)
		param = ""
		if len(options.ChannelID) > 0 {
			query += " Posts.ChannelId = ? AND"
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
	query += fmt.Sprintf(` %s Posts.CreateAt >= ?`, param)

	// if options.Start != 0 && options.End != 0 {
	// 	query += ` WHERE Posts.CreateAt <= ?
	// 	AND Posts.CreateAt >= ? `
	// 	args = append(args, options.Start)
	// 	args = append(args, options.End)
	// }

	query += fmt.Sprintf(` GROUP BY %s(FROM_UNIXTIME(Posts.CreateAt / 1000))
		ORDER BY Name DESC
		LIMIT ?`, freq)
	args = append(args, options.Limit)

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
