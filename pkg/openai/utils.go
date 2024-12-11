package openai

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"openai-realtime/pkg/openai/events"
	"strings"
)

// 로그 클라이언트 이벤트
func logEventAsJSON(prefix string, event events.OpenAIEvent, message []byte) {
	eventMap := make(map[string]interface{})
	if err := json.Unmarshal(message, &eventMap); err != nil {
		log.Error("Error unmarshalling events for logging:", err)
		return
	}
	if delta, ok := eventMap["delta"].(string); ok {
		eventMap["delta"] = len(delta)
	}
	if audio, ok := eventMap["audio"].(string); ok {
		eventMap["audio"] = len(audio)
	}
	jsonEvent, err := json.Marshal(eventMap)
	if err != nil {
		log.Error("Error marshalling events for logging:", err)
		return
	}
	log.Info(prefix, event.GetType(), string(jsonEvent))
}

func generateEventID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return "event_" + strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

// getUrl WebSocket URL 생성
func (c *Client) getUrl() string {
	u := url.URL{
		Scheme:   "wss",
		Host:     c.host,
		Path:     c.path,
		RawQuery: "model=" + url.QueryEscape(c.model),
	}
	return u.String()
}
