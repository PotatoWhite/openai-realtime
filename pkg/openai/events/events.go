package events

type OpenAIEvent interface {
	GetType() string
}

type InputAudioTranscription struct {
	Model string `json:"model"`
}

type TurnDetection struct {
	Type              string  `json:"type"`
	Threshold         float64 `json:"threshold"`
	PrefixPaddingMs   int     `json:"prefix_padding_ms"`
	SilenceDurationMs int     `json:"silence_duration_ms"`
}

type Tool struct {
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type SessionCreated struct {
	ID                      string                   `json:"id"`
	Object                  string                   `json:"object"`
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
	MaxResponseOutputTokens interface{}              `json:"max_response_output_tokens,omitempty"`
}

type SessionUpdated struct {
	EventID string `json:"event_id"`
	Type    string `json:"type"`
	Session struct {
		ID                      string                  `json:"id"`
		Object                  string                  `json:"object"`
		Model                   string                  `json:"model"`
		Modalities              []string                `json:"modalities"`
		Instructions            string                  `json:"instructions"`
		Voice                   string                  `json:"voice"`
		InputAudioFormat        string                  `json:"input_audio_format"`
		OutputAudioFormat       string                  `json:"output_audio_format"`
		InputAudioTranscription InputAudioTranscription `json:"input_audio_transcription"`
		TurnDetection           interface{}             `json:"turn_detection"`
		Tools                   []interface{}           `json:"tools"`
		ToolChoice              string                  `json:"tool_choice"`
		Temperature             float64                 `json:"temperature"`
		MaxResponseOutputTokens int                     `json:"max_response_output_tokens"`
	} `json:"session"`
}
