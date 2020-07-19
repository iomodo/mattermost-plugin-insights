package chart

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	ch "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func CreateBarChart(title string, rows model.AnalyticsRows) *bytes.Buffer {
	println(fmt.Sprintf("rows = %v", rows))

	values := []ch.Value{}
	vx := []time.Time{}
	vy := []float64{}
	style := ch.Style{
		FillColor:   drawing.ColorFromHex("13c158"),
		StrokeColor: drawing.ColorFromHex("13c158"),
		StrokeWidth: 0,
	}

	for _, row := range rows {
		values = append(values, ch.Value{Style: style, Value: row.Value, Label: row.Name})
		parsed, _ := time.Parse(ch.DefaultDateFormat, row.Name)
		vx = append(vx, parsed)
		vy = append(vy, row.Value)
	}
	println(fmt.Sprintf("vx = %v", vx))
	println(fmt.Sprintf("vy = %v", vy))
	/*
		priceSeries := ch.TimeSeries{
			Name: "Posts",
			Style: ch.Style{
				Show:        true,
				StrokeColor: ch.GetDefaultColor(0),
			},
			XValues: vx,
			YValues: vy,
		}
		smaSeries := ch.SMASeries{
			Name: "SPY - SMA",
			Style: ch.Style{
				Show:            true,
				StrokeColor:     drawing.ColorRed,
				StrokeDashArray: []float64{5.0, 5.0},
			},
			InnerSeries: priceSeries,
		}

		graph := ch.Chart{
			XAxis: ch.XAxis{
				Name:         "Days",
				TickPosition: ch.TickPositionUnderTick,
			},
			YAxis: ch.YAxis{
				Name: "Posts",
				Range: &ch.ContinuousRange{
					Max: 300.0,
					Min: 0.0,
				},
			},
			Series: []ch.Series{
				priceSeries,
				smaSeries,
			},
		}
	*/
	graph2 := ch.BarChart{
		XAxis: ch.Style{
			FillColor: drawing.ColorFromHex("efefef"),
		},
		YAxis: ch.YAxis{
			Name: "Posts",
		},
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
	err := graph2.Render(ch.PNG, buffer)
	if err != nil {
		println("err", err.Error())
	}
	return buffer
}
