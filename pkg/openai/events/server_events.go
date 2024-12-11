package events

// 서버 이벤트 구조체
type ServerEvent struct {
	EventID string          `json:"event_id"`
	Type    string          `json:"type"`
	Session *SessionCreated `json:"session,omitempty"`
	Error   *EventError     `json:"error,omitempty"`
}

func (e ServerEvent) GetType() string {
	return e.Type
}

type ConversationItemInputAudioTranscriptionCompleted struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ItemID       string `json:"item_id"`
	ContentIndex int    `json:"content_index"`
	Transcript   string `json:"transcript"`
}

func (e ConversationItemInputAudioTranscriptionCompleted) GetType() string {
	return e.Type
}

type ResponseContentPartAdded struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Part         struct {
		Type       string `json:"type"`
		Transcript string `json:"transcript"`
	} `json:"part"`
}

func (e ResponseContentPartAdded) GetType() string {
	return e.Type
}

type ResponseCreated struct {
	Type     string `json:"type"`
	EventID  string `json:"event_id"`
	Response struct {
		Object        string        `json:"object"`
		ID            string        `json:"id"`
		Status        string        `json:"status"`
		StatusDetails interface{}   `json:"status_details"`
		Output        []interface{} `json:"output"`
		Usage         interface{}   `json:"usage"`
	} `json:"response"`
}

func (e ResponseCreated) GetType() string {
	return e.Type
}

type RateLimitsUpdated struct {
	Type       string `json:"type"`
	EventID    string `json:"event_id"`
	RateLimits []struct {
		Name         string  `json:"name"`
		Limit        int     `json:"limit"`
		Remaining    int     `json:"remaining"`
		ResetSeconds float64 `json:"reset_seconds"`
	} `json:"rate_limits"`
}

func (e RateLimitsUpdated) GetType() string {
	return e.Type
}

type ResponseOutputItemAdded struct {
	Type        string `json:"type"`
	EventID     string `json:"event_id"`
	ResponseID  string `json:"response_id"`
	OutputIndex int    `json:"output_index"`
	Item        struct {
		ID      string        `json:"id"`
		Object  string        `json:"object"`
		Type    string        `json:"type"`
		Status  string        `json:"status"`
		Role    string        `json:"role"`
		Content []interface{} `json:"content"`
	} `json:"item"`
}

func (e ResponseOutputItemAdded) GetType() string {
	return e.Type
}

type InputAudioBufferCommitted struct {
	Type           string  `json:"type"`
	EventID        string  `json:"event_id"`
	PreviousItemID *string `json:"previous_item_id"`
	ItemID         string  `json:"item_id"`
}

func (e InputAudioBufferCommitted) GetType() string {
	return e.Type
}

type ConversationItemCreated struct {
	Type           string  `json:"type"`
	EventID        string  `json:"event_id"`
	PreviousItemID *string `json:"previous_item_id"`
	Item           struct {
		ID      string    `json:"id"`
		Object  string    `json:"object"`
		Type    string    `json:"type"`
		Status  string    `json:"status"`
		Role    string    `json:"role"`
		Content []Content `json:"content"`
	} `json:"item"`
}

func (e ConversationItemCreated) GetType() string {
	return e.Type
}

type InputAudioBufferSpeechStarted struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	AudioStartMs int    `json:"audio_start_ms"`
	ItemID       string `json:"item_id"`
}

func (e InputAudioBufferSpeechStarted) GetType() string {
	return e.Type
}

type InputAudioBufferSpeechStopped struct {
	Type       string `json:"type"`
	EventID    string `json:"event_id"`
	AudioEndMs int    `json:"audio_end_ms"`
	ItemID     string `json:"item_id"`
}

func (e InputAudioBufferSpeechStopped) GetType() string {
	return e.Type
}

type ResponseAudioDone struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
}

func (e ResponseAudioDone) GetType() string {
	return e.Type
}

type ResponseAudioTranscriptDone struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Transcript   string `json:"transcript"`
}

func (e ResponseAudioTranscriptDone) GetType() string {
	return e.Type
}

type ResponseContentPartDone struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Part         struct {
		Type       string `json:"type"`
		Transcript string `json:"transcript"`
	} `json:"part"`
}

func (e ResponseContentPartDone) GetType() string {
	return e.Type
}

type ResponseOutputItemDone struct {
	Type        string `json:"type"`
	EventID     string `json:"event_id"`
	ResponseID  string `json:"response_id"`
	OutputIndex int    `json:"output_index"`
	Item        struct {
		ID      string    `json:"id"`
		Object  string    `json:"object"`
		Type    string    `json:"type"`
		Status  string    `json:"status"`
		Role    string    `json:"role"`
		Content []Content `json:"content"`
	} `json:"item"`
}

func (e ResponseOutputItemDone) GetType() string {
	return e.Type
}

type ResponseTextDelta struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

func (e ResponseTextDelta) GetType() string {
	return e.Type
}

type ResponseTextDone struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Text         string `json:"text"`
}

func (e ResponseTextDone) GetType() string {
	return e.Type
}

type ResponseDone struct {
	Type     string `json:"type"`
	EventID  string `json:"event_id"`
	Response struct {
		ID            string     `json:"id"`
		Object        string     `json:"object"`
		Output        []struct{} `json:"output"`
		Status        string     `json:"status"`
		StatusDetails struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
			Type string `json:"type"`
		} `json:"status_details"`
		Usage struct {
			InputTokenDetails struct {
				AudioTokens         int `json:"audio_tokens"`
				CachedTokens        int `json:"cached_tokens"`
				CachedTokensDetails struct {
					AudioTokens int `json:"audio_tokens"`
					TextTokens  int `json:"text_tokens"`
				} `json:"cached_tokens_details"`
				TextTokens int `json:"text_tokens"`
			} `json:"input_token_details"`
			InputTokens        int `json:"input_tokens"`
			OutputTokenDetails struct {
				AudioTokens int `json:"audio_tokens"`
				TextTokens  int `json:"text_tokens"`
			} `json:"output_token_details"`
			OutputTokens int `json:"output_tokens"`
			TotalTokens  int `json:"total_tokens"`
		} `json:"usage"`
	} `json:"response"`
}

func (e ResponseDone) GetType() string {
	return e.Type
}

type ResponseAudioTranscriptDelta struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

func (e ResponseAudioTranscriptDelta) GetType() string {
	return e.Type
}

type ResponseAudioDelta struct {
	Type         string `json:"type"`
	EventID      string `json:"event_id"`
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

func (e ResponseAudioDelta) GetType() string {
	return e.Type
}

type Content struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func (e Content) GetType() string {
	return e.Type
}
