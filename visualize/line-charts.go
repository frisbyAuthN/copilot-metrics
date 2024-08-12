package copilotvisualize

import (
	copilotclient "copilot-metrics/client"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"fmt"
)

func BuildActiveUsersLineChart(copilotData []copilotclient.CopilotMetricsBody) *charts.Line {
	activeUsersLineData := make([]opts.LineData, 0)
	activeChatUsersLineData := make([]opts.LineData, 0)
	for i := range copilotData {
		activeUsersLineData = append(activeUsersLineData, opts.LineData{Value: copilotData[i].TotalActiveUsers})
		activeChatUsersLineData = append(activeChatUsersLineData, opts.LineData{Value: copilotData[i].TotalActiveChatUsers})
	}

	dateRange := make([]string, 0)
	for i := range copilotData {
		dateRange = append(dateRange, copilotData[i].Day)
	}
	chartSubtitle := fmt.Sprintf("Last %d days", len(copilotData))

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts((opts.Initialization{Theme: types.ThemeWesteros})),
		charts.WithTitleOpts(opts.Title{
			Title:    "Copilot Active Users",
			Subtitle: chartSubtitle,
		}))
	line.SetXAxis(dateRange).AddSeries("Active Copilot Users", activeUsersLineData).
		AddSeries("Active Copilot Chat Users", activeChatUsersLineData).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{
			Smooth:     opts.Bool(false),
			Symbol:     "rect",
			SymbolSize: 8}))

	return line
}

func BuildSuggestionsLineChart(copilotData []copilotclient.CopilotMetricsBody) *charts.Line {
	totalSuggestions := make([]opts.LineData, 0)
	totalAcceptances := make([]opts.LineData, 0)
	totalAcceptancePercentage := make([]opts.LineData, 0)
	for i := range copilotData {
		totalSuggestions = append(totalSuggestions, opts.LineData{Value: copilotData[i].TotalSuggestionsCount})
		totalAcceptances = append(totalAcceptances, opts.LineData{Value: copilotData[i].TotalAcceptancesCount})
		percentage := 0.0
		if copilotData[i].TotalSuggestionsCount != 0 {
			percentage = 100 * float64(copilotData[i].TotalAcceptancesCount) / float64(copilotData[i].TotalSuggestionsCount)
		}
		totalAcceptancePercentage = append(totalAcceptancePercentage, opts.LineData{Value: percentage})
	}

	dateRange := make([]string, 0)
	for i := range copilotData {
		dateRange = append(dateRange, copilotData[i].Day)
	}
	chartSubtitle := fmt.Sprintf("Last %d days", len(copilotData))

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts((opts.Initialization{Theme: types.ThemeWesteros})),
		charts.WithTitleOpts(opts.Title{
			Title:    "Copilot Suggestions",
			Subtitle: chartSubtitle,
		}),
	)

	line.ExtendYAxis(opts.YAxis{
		Name:  "Percentage",
		Show:  opts.Bool(true),
		Scale: opts.Bool(true),
		Type:  "value",
	})

	line.SetXAxis(dateRange).
		AddSeries("Suggestions", totalSuggestions).
		AddSeries("Suggestions Accepted", totalAcceptances).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{
			Smooth:     opts.Bool(false),
			ShowSymbol: opts.Bool(true),
			Symbol:     "rect",
			SymbolSize: 8,
		}))

	line.AddSeries("Acceptance Percentage", totalAcceptancePercentage,
		charts.WithLineChartOpts(opts.LineChart{
			YAxisIndex: 1,
			Smooth:     opts.Bool(false),
			SymbolSize: 8,
		}))

	return line
}

func BuildChatLineChart(copilotData []copilotclient.CopilotMetricsBody) *charts.Line {
	chatPromptAndResponse := make([]opts.LineData, 0)
	chatCodeAcceptances := make([]opts.LineData, 0)
	for i := range copilotData {
		chatPromptAndResponse = append(chatPromptAndResponse, opts.LineData{Value: copilotData[i].TotalChatTurns})
		chatCodeAcceptances = append(chatCodeAcceptances, opts.LineData{Value: copilotData[i].TotalChatAcceptances})
	}

	dateRange := make([]string, 0)
	for i := range copilotData {
		dateRange = append(dateRange, copilotData[i].Day)
	}
	chartSubtitle := fmt.Sprintf("Last %d days", len(copilotData))

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts((opts.Initialization{Theme: types.ThemeWesteros})),
		charts.WithTitleOpts(opts.Title{
			Title:    "Copilot Chat Interactions",
			Subtitle: chartSubtitle,
		}),
	)

	line.SetXAxis(dateRange).
		AddSeries("Prompt and Response Pairs", chatPromptAndResponse).
		AddSeries("Copilot Chat Code Acceptances", chatCodeAcceptances).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{
			Smooth:     opts.Bool(false),
			ShowSymbol: opts.Bool(true),
			Symbol:     "rect",
			SymbolSize: 8,
		}))

	return line
}

func BuildLanguageSpecificLineChart(copilotData []copilotclient.CopilotMetricsBody) *charts.Line {
	daysCount := len(copilotData)
	languageData := make([]LanguageData, daysCount)
	for i := range copilotData {
		languageData[i].Day = copilotData[i].Day
		for j := range copilotData[i].Breakdown {
			languageData[i].LanguageDetails = append(languageData[i].LanguageDetails, LanguageDetails{
				Language:         copilotData[i].Breakdown[j].Language,
				SuggestionsCount: copilotData[i].Breakdown[j].SuggestionsCount,
				AcceptancesCount: copilotData[i].Breakdown[j].AcceptancesCount,
			})
		}
	}

	javaData := extractSpecificLanguage(languageData, "java", "java-v18", daysCount)
	javaLineData := buildLanguagePercentageChart(javaData)

	javascriptData := extractSpecificLanguage(languageData, "javascript", "javascriptreact", daysCount)
	javascriptLineData := buildLanguagePercentageChart(javascriptData)

	typescriptData := extractSpecificLanguage(languageData, "typescript", "typescriptreact", daysCount)
	typescriptLineData := buildLanguagePercentageChart(typescriptData)

	pythonData := extractSpecificLanguage(languageData, "python", "python3", daysCount)
	pythonLineData := buildLanguagePercentageChart(pythonData)

	dateRange := make([]string, 0)
	for i := range copilotData {
		dateRange = append(dateRange, copilotData[i].Day)
	}
	chartSubtitle := fmt.Sprintf("Last %d days", len(copilotData))

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts((opts.Initialization{Theme: types.ThemeWesteros})),
		charts.WithTitleOpts(opts.Title{
			Title:    "Language-Specific Suggestion Acceptance Rate (Percentage)",
			Subtitle: chartSubtitle,
		}),
		charts.WithLegendOpts(opts.Legend{
			Orient: "vertical",
			Right:  "right",
		}),
	)

	line.SetXAxis(dateRange).
		AddSeries("Java", javaLineData).
		AddSeries("JavaScript", javascriptLineData).
		AddSeries("TypeScript", typescriptLineData).
		AddSeries("Python", pythonLineData).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{
			Smooth:     opts.Bool(false),
			ShowSymbol: opts.Bool(true),
			Symbol:     "rect",
			SymbolSize: 8,
		}))

	return line
}

func buildLanguagePercentageChart(data []LanguageDetails) []opts.LineData {
	for i := range data {
		if data[i].AcceptancesCount != 0 && data[i].SuggestionsCount != 0 {
			data[i].AcceptancePercentage = truncate(100.0 * (float64(data[i].AcceptancesCount) / float64(data[i].SuggestionsCount)))
		} else {
			data[i].AcceptancePercentage = 0
		}
	}

	lineData := make([]opts.LineData, 0)
	for i := range data {
		lineData = append(lineData, opts.LineData{Value: data[i].AcceptancePercentage})
	}
	return lineData
}

func extractSpecificLanguage(languageData []LanguageData, primaryName string, secondaryName string, daysCount int) []LanguageDetails {
	data := make([]LanguageDetails, daysCount)
	for i := range languageData {
		for j := range languageData[i].LanguageDetails {
			data[i].Language = primaryName
			if languageData[i].LanguageDetails[j].Language == primaryName || languageData[i].LanguageDetails[j].Language == secondaryName {
				data[i].SuggestionsCount += languageData[i].LanguageDetails[j].SuggestionsCount
				data[i].AcceptancesCount += languageData[i].LanguageDetails[j].AcceptancesCount
			}
		}
	}

	return data
}
