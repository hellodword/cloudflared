package tunnel

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

var m sync.Map

func Notify(url string) {
	_, loaded := m.LoadOrStore(url, 0)
	if !loaded && strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")) != "" && strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID")) != "" {
		http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))),
			"application/json",
			strings.NewReader(fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "disable_notification": false}`, strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID")), url)))
	}
}
