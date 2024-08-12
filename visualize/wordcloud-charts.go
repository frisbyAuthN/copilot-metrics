package copilotvisualize

import (
	copilotclient "copilot-metrics/client"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func BuildLanguageWordCloud(copilotData []copilotclient.CopilotMetricsBody) *charts.WordCloud {
	languageCounts := make(map[string]int)
	for i := range copilotData {
		for j := range copilotData[i].Breakdown {
			if languageCounts[copilotData[i].Breakdown[j].Language] == 0 {
				languageCounts[copilotData[i].Breakdown[j].Language] = copilotData[i].Breakdown[j].SuggestionsCount
			} else {
				languageCounts[copilotData[i].Breakdown[j].Language] += copilotData[i].Breakdown[j].SuggestionsCount
			}
		}
	}

	languageCloudData := make([]opts.WordCloudData, 0)
	for key, val := range languageCounts {
		languageCloudData = append(languageCloudData, opts.WordCloudData{Name: key, Value: val})
	}
	languageCloudChart := charts.NewWordCloud()
	languageCloudChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Copilot Language Usage"}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}))
	languageCloudChart.AddSeries("Languages", languageCloudData).
		SetSeriesOptions(charts.WithWorldCloudChartOpts(opts.WordCloudChart{
			SizeRange:     []float32{10, 200},
			Shape:         "roundRect",
			RotationRange: []float32{0, 0}, // No rotation
		}))

	return languageCloudChart
}

func BuildEditorWordCloud(copilotData []copilotclient.CopilotMetricsBody) *charts.WordCloud {
	editorCounts := make(map[string]int)
	for i := range copilotData {
		for j := range copilotData[i].Breakdown {
			if editorCounts[copilotData[i].Breakdown[j].Editor] == 0 {
				editorCounts[copilotData[i].Breakdown[j].Editor] = copilotData[i].Breakdown[j].SuggestionsCount
			} else {
				editorCounts[copilotData[i].Breakdown[j].Editor] += copilotData[i].Breakdown[j].SuggestionsCount
			}
		}
	}

	editorCloudData := make([]opts.WordCloudData, 0)
	for key, val := range editorCounts {
		editorCloudData = append(editorCloudData, opts.WordCloudData{Name: key, Value: val})
	}
	editorCloudChart := charts.NewWordCloud()
	editorCloudChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Copilot Editor Usage"}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}))
	editorCloudChart.AddSeries("Editors", editorCloudData).
		SetSeriesOptions(charts.WithWorldCloudChartOpts(opts.WordCloudChart{
			SizeRange:     []float32{10, 100},
			Shape:         "roundRect",
			RotationRange: []float32{0, 0}, // No rotation
		}))

	return editorCloudChart
}
