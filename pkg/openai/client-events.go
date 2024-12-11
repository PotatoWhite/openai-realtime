package openai

import (
	"encoding/base64"
	"fmt"
	"openai-realtime/pkg/config"
	"openai-realtime/pkg/openai/events"
)

const (
	SessionUpdateEventType          = "session.update"
	ConversationItemCreateEventType = "conversation.item.create"
	InputAudioBufferAppendEventType = "input_audio_buffer.append"
	InputAudioBufferCommitEventType = "input_audio_buffer.commit"
)

// 로그 클라이언트 이벤트 (go routine)
func (c *Client) SessionUpdate(inputAudioTranscription events.InputAudioTranscription, turnDetection events.TurnDetection, tools []events.Tool) error {
	sessionUpdate := events.SessionUpdate{
		Modalities:              []string{"text", "audio"},
		Instructions:            config.SystemPrompt(),
		Voice:                   "alloy",
		InputAudioFormat:        "pcm16",
		OutputAudioFormat:       "pcm16",
		InputAudioTranscription: &inputAudioTranscription,
		TurnDetection:           &turnDetection,
		Tools:                   tools,
		ToolChoice:              "auto",
		Temperature:             1.2,
		MaxResponseOutputTokens: 1024,
	}

	return c.sendEvent(events.ClientEvent{
		EventID: generateEventID(),
		Type:    SessionUpdateEventType,
		Session: &sessionUpdate,
	}, true)
}

func (c *Client) ConversationItemCreate(content string, role string) error {
	item := events.Item{
		Content: []events.Content{
			{Text: content, Type: "input_text"},
		},
		Type: "message",
		Role: role,
	}

	return c.sendEvent(events.ClientEvent{
		EventID: generateEventID(),
		Type:    ConversationItemCreateEventType,
		Item:    &item,
	}, true)
}

func (c *Client) SendInputAudioBufferAppend(data []byte) error {
	fmt.Print(".")
	if len(data) == 0 {
		return fmt.Errorf("no streamAudio data to send")
	}

	if c.status != StatusReady {
		log.Warn("Cannot send audio data when not ready")
		return nil
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return c.sendEvent(events.ClientEvent{
		EventID: generateEventID(),
		Type:    InputAudioBufferAppendEventType,
		Audio:   &encoded,
	}, false)
}

func (c *Client) SendInputAudioBufferCommit() error {
	return c.sendEvent(events.ClientEvent{
		EventID: generateEventID(),
		Type:    InputAudioBufferCommitEventType,
	}, true)
}
