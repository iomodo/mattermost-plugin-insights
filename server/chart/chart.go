package chart

import (
	"bytes"

	"github.com/mattermost/mattermost-server/v5/model"
	ch "github.com/wcharczuk/go-chart"
)

func CreateBarChart(title string, rows model.AnalyticsRows) *bytes.Buffer {
	values := []ch.Value{}
	for _, row := range rows {
		values = append(values, ch.Value{Value: row.Value, Label: row.Name})
	}

	graph := ch.BarChart{
		Title: title,
		Background: ch.Style{
			Padding: ch.Box{
				Top: 40,
			},
		},
		Height:   1024,
		BarWidth: 60,
		Bars:     values,
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(ch.PNG, buffer)
	return buffer
}
