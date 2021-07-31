package connection

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var notified bool

func TrialZoneMsg(url string) []string {
	if !notified && strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")) != "" && strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID")) != "" {
		notified = true
		http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))),
			"application/json",
			strings.NewReader(fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "disable_notification": false}`, strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID")), url)))
	}
	return []string{
		"Your free tunnel has started! Visit it:",
		"  " + url,
	}
}
