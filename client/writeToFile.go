package copilotclient

import (
	"fmt"
	"os"
	"time"
)

func WriteToFile(body []byte) {
	textFileName := "copilot-metrics-" + time.Now().Format("2006-01-02") + ".json"
	err := os.WriteFile(textFileName, body, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(textFileName + " printed")
}
