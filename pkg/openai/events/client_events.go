package events

// 클라이언트 이벤트 구조체
type ClientEvent struct {
	EventID string         `json:"event_id"`
	Type    string         `json:"type"`
	Session *SessionUpdate `json:"session,omitempty"`
	Audio   *string        `json:"audio,omitempty"`
	Item    *Item          `json:"item,omitempty"`
}

func (e ClientEvent) GetType() string {
	return e.Type
}

type Item struct {
	Content []Content `json:"content"`
	Type    string    `json:"type"`
	Role    string    `json:"role"`
}

func (e Item) GetType() string {
	return e.Type
}

type SessionUpdate struct {
	Model                   string                   `json:"model"`
	Modalities              []string                 `json:"modalities"`
	Instructions            string                   `json:"instructions"`
	Voice                   string                   `json:"voice"`
	InputAudioFormat        string                   `json:"input_audio_format"`
	OutputAudioFormat       string                   `json:"output_audio_format"`
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`
	TurnDetection           *TurnDetection           `json:"turn_detection,omitempty"`
	Tools                   []Tool                   `json:"tools"`
	ToolChoice              string                   `json:"tool_choice"`
	Temperature             float64                  `json:"temperature"`
	MaxResponseOutputTokens int                      `json:"max_response_output_tokens,omitempty"`
}

func (e SessionUpdate) GetType() string {
	return "session.update"
}
