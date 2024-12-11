package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"openai-realtime/pkg/config"
	"openai-realtime/pkg/openai/events"
	"os"
	"time"
)

const (
	StatusConnected  = "connected"
	StatusClosed     = "closed"
	StatusWaiting    = "waiting"
	StatusReady      = "ready"
	StatusProcessing = "processing"

	reconnectInterval    = 5 * time.Second
	maxReconnectAttempts = 5
)

var log = func() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(config.LogLevel)
	return log
}()

// 클라이언트 설정 구조체
type Client struct {
	apiKey string
	conn   *websocket.Conn
	host   string
	path   string
	model  string

	status          string
	AudioOutputChan chan []byte
	ErrChan         chan error

	reconnectAttempts int
}

// NewClient 생성자 함수
func NewClient(ctx context.Context, host, path, model, apiKey string) (*Client, error) {
	client := Client{
		host:            host,
		model:           model,
		path:            path,
		apiKey:          apiKey,
		status:          StatusClosed,
		AudioOutputChan: make(chan []byte, 10),
		ErrChan:         make(chan error, 1),
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	return &client, nil
}

// Connect WebSocket 연결
func (c *Client) Connect(ctx context.Context) error {
	// API 키 확인
	if c.apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	// 인증 헤더 설정
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+c.apiKey)
	headers.Add("OpenAI-Beta", "realtime=v1")

	urlString := c.getUrl()

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	log.Info("Connecting to WebSocket server:", urlString)

	if conn, _, err := dialer.DialContext(ctx, urlString, headers); err != nil {
		log.Info("Dial error:", err)
		return err
	} else {
		c.conn = conn
		c.status = StatusConnected
	}

	log.Info("WebSocket connection established")

	return nil
}

// Close WebSocket 연결 종료
func (c *Client) Close() error {

	if c.conn == nil {
		return nil
	}

	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Error(fmt.Sprintf("Send close message error: %v", err))
	}

	err = c.conn.Close()
	if err != nil {
		log.Error(fmt.Sprintf("Connection close error: %v", err))
		return err
	}

	c.conn = nil
	c.status = StatusClosed
	log.Info("WebSocket connection closed")
	return nil
}

// readEvent WebSocket 메시지 읽기
func (c *Client) readEvent() (messageType int, p []byte, err error) {
	return c.conn.ReadMessage()
}

// sendEvent 클라이언트 이벤트 전송
func (c *Client) sendEvent(event events.ClientEvent, logging bool) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Error("Error marshalling events:", err)
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, eventJSON)
	if err != nil {
		log.Error("Error sending events:", err)
		return err
	}

	if logging {
		logEventAsJSON("[SEND]", event, eventJSON)
	}
	return nil
}
