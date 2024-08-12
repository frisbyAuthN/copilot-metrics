package copilotclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func FetchCopilotMetrics() ([]byte, error) {
	client := http.Client{Timeout: time.Duration(1) * time.Second}

	enterpriseName := os.Getenv("METRICS_ENTERPRISE")
	copilotUrl := fmt.Sprintf("https://api.github.com/enterprises/%s/copilot/usage", enterpriseName)
	req, err := http.NewRequest("GET",
		copilotUrl,
		nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	token := os.Getenv("METRICS_TOKEN")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
