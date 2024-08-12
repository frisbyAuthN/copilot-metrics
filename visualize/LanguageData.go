package copilotvisualize

type LanguageData struct {
	Day             string
	LanguageDetails []LanguageDetails
}

type LanguageDetails struct {
	Language             string
	SuggestionsCount     int
	AcceptancesCount     int
	AcceptancePercentage float64
}
