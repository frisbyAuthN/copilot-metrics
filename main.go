package main

import (
	copilotclient "copilot-metrics/client"
	copilotvisualize "copilot-metrics/visualize"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-echarts/go-echarts/v2/components"
)

func main() {

	copilotData := make([]copilotclient.CopilotMetricsBody, 0)

	copilotData = aggregateDataFromFile(copilotData)
	copilotData = fetchCopilotDataFromGithub(copilotData)

	activeUsersChart := copilotvisualize.BuildActiveUsersLineChart(copilotData)
	suggestionsChart := copilotvisualize.BuildSuggestionsLineChart(copilotData)
	chatInfoChart := copilotvisualize.BuildChatLineChart(copilotData)
	linesBarChart := copilotvisualize.BuildLinesSuggestedBar(copilotData)
	languageCloudChart := copilotvisualize.BuildLanguageWordCloud(copilotData)
	editorCloudChart := copilotvisualize.BuildEditorWordCloud((copilotData))
	languageSuggestionsChart := copilotvisualize.BuildLanguageSpecificLineChart(copilotData)

	metricsReport := components.NewPage()
	metricsReport.SetPageTitle("Copilot Metrics Usage Report")
	metricsReport.AddCharts(
		suggestionsChart,
		activeUsersChart,
		chatInfoChart,
		linesBarChart,
		languageCloudChart,
		editorCloudChart,
		languageSuggestionsChart,
	)

	printMetricsReport(metricsReport)
}

func aggregateDataFromFile(copilotData []copilotclient.CopilotMetricsBody) []copilotclient.CopilotMetricsBody {
	files, err := os.ReadDir("older-data/")
	if err != nil {
		fmt.Println(err)
	}

	// ASSUMPTION ALERT - this code assumes the older data will be read first so new days will be appended in the correct time order
	for _, file := range files {
		bodyFromFile, err := os.ReadFile("older-data/" + file.Name())
		if err != nil {
			fmt.Println(err)
		}

		if filepath.Ext(file.Name()) == ".json" {
			copilotDataFromFile, err := unmarshalCopilotData(bodyFromFile)
			if err != nil {
				fmt.Println(err)
			}
			copilotData = appendUniqueDaysData(copilotData, copilotDataFromFile)
		}
	}
	return copilotData
}

func fetchCopilotDataFromGithub(copilotData []copilotclient.CopilotMetricsBody) []copilotclient.CopilotMetricsBody {
	body, err := copilotclient.FetchCopilotMetrics()
	if err != nil {
		fmt.Println(err)
	}
	copilotclient.WriteToFile(body)

	latestCopilotData, err := unmarshalCopilotData(body)
	if err != nil {
		fmt.Println(err)
	}

	copilotData = appendUniqueDaysData(copilotData, latestCopilotData)
	return copilotData
}

func appendUniqueDaysData(copilotData []copilotclient.CopilotMetricsBody, newCopilotData []copilotclient.CopilotMetricsBody) []copilotclient.CopilotMetricsBody {
	for i := range newCopilotData {
		if dataDoesNotContainDay(copilotData, newCopilotData[i].Day) {
			copilotData = append(copilotData, newCopilotData[i])
		}
	}
	return copilotData
}

func unmarshalCopilotData(body []byte) ([]copilotclient.CopilotMetricsBody, error) {
	copilotDaysOfData := 27
	copilotData := make([]copilotclient.CopilotMetricsBody, copilotDaysOfData)
	err := json.Unmarshal(body, &copilotData)
	if err != nil {
		fmt.Printf("Unable to unmarshal JSON due to %s\n", err)
		return nil, err
	}
	return copilotData, nil
}

func dataDoesNotContainDay(data []copilotclient.CopilotMetricsBody, day string) bool {
	for i := range data {
		if data[i].Day == day {
			return false
		}
	}
	return true
}

func printMetricsReport(metricsReport *components.Page) {
	reportName := "copilot-metrics-report-" + time.Now().Format("2006-01-02") + ".html"
	reportFile, err := os.Create(reportName)
	if err != nil {
		fmt.Println(err)
	}

	err = metricsReport.Render(io.MultiWriter(reportFile))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reportName + " printed")
}
