package copilotvisualize

import (
	copilotclient "copilot-metrics/client"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func BuildLinesSuggestedBar(copilotData []copilotclient.CopilotMetricsBody) *charts.Bar {

	totalLinesSuggested := 0.0
	totalLinesAccepted := 0.0
	days := float64(len(copilotData))
	for i := range copilotData {
		totalLinesSuggested += float64(copilotData[i].TotalLinesSuggested)
		totalLinesAccepted += float64(copilotData[i].TotalLinesAccepted)
	}
	averageLinesSuggested := totalLinesSuggested / days
	averageLinesSuggested = truncate(averageLinesSuggested)
	averageLinesAccepted := totalLinesAccepted / days
	averageLinesAccepted = truncate(averageLinesAccepted)

	averageLinesSuggestedBarData := make([]opts.BarData, 0)
	averageLinesSuggestedBarData = append(averageLinesSuggestedBarData, opts.BarData{Value: averageLinesSuggested})
	averageLinesAcceptedBarData := make([]opts.BarData, 0)
	averageLinesAcceptedBarData = append(averageLinesAcceptedBarData, opts.BarData{Value: averageLinesAccepted})

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Average Daily Copilot Lines of Code"}),
		charts.WithLegendOpts(opts.Legend{Left: "right", Orient: "vertical"}))

	bar.SetXAxis([]string{"Average Daily Copilot Lines"}).
		AddSeries("Lines Suggested", averageLinesSuggestedBarData).
		AddSeries("Lines Accepted", averageLinesAcceptedBarData).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "top"}))

	return bar
}

func truncate(input float64) float64 {
	return float64(int(input))
}
