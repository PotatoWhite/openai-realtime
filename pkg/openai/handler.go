package openai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"openai-realtime/pkg/openai/events"
	"os"
)

// ReceiveServerEvent 서버 이벤트 수신 (go routine)
func (c *Client) ReceiveServerEvent(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		log.Debug("Closing Server Event Receiver")
	}()

	for {

		select {
		case <-ctx.Done():
			log.Info("Context done. Closing Server Event Receiver")
			return
		default:
			messageChan := make(chan []byte, 10)
			errorChan := make(chan error)

			go func() {
				_, message, err := c.readEvent()
				if err != nil {
					errorChan <- err
					return
				}
				messageChan <- message
			}()

			select {
			case <-ctx.Done():
				log.Info("Context done. Closing Server Event Receiver")
				return
			case err := <-errorChan:
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					cancel()
					log.Warn("Connection closed:", err)
				} else {
					log.Error("Read error:", err)
				}
				return
			case message := <-messageChan:
				var event events.ServerEvent
				if err := json.Unmarshal(message, &event); err != nil {
					log.Error("Error unmarshalling server events:", err)
					continue
				}

				if err := c.handleServerEvent(ctx, event, message); err != nil {
					cancel()
					return
				}
			}
		}
	}
}

// 서버 이벤트 핸들링 함수
func (c *Client) handleServerEvent(ctx context.Context, event events.ServerEvent, message []byte) error {
	//logEventAsJSON("[RECV]", event, message)
	switch event.Type {
	case "error":
		if event.Error != nil {
			log.Error("Error:", event.Error.Message)
			return fmt.Errorf("Error: %s", event.Error.Message)
		}
	case "session.created":
		var sessionCreated events.SessionCreated
		if err := json.Unmarshal(message, &sessionCreated); err != nil {
			log.Error("Error unmarshalling session created events:", err)
			return err
		}
		c.status = StatusReady
	case "session.updated":
		var sessionUpdated events.SessionUpdated
		if err := json.Unmarshal(message, &sessionUpdated); err != nil {
			log.Error("Error unmarshalling session updated events:", err)
			return err
		}
	case "conversation.item.created":
		var conversationItemCreatedEvent events.ConversationItemCreated
		if err := json.Unmarshal(message, &conversationItemCreatedEvent); err != nil {
			log.Error("Error unmarshalling conversation item created events:", err)
			return err
		}
		c.status = StatusReady
	case "response.text.delta":
		var responseTextDelta events.ResponseTextDelta
		if err := json.Unmarshal(message, &responseTextDelta); err != nil {
			log.Error("Error unmarshalling text delta events:", err)
			return err
		}
		fmt.Print(responseTextDelta.Delta)
	case "response.text.done":
		var responseTextDone events.ResponseTextDone
		if err := json.Unmarshal(message, &responseTextDone); err != nil {
			log.Error("Error unmarshalling text done events:", err)
			return err
		}
		fmt.Printf("\n\n")
	case "response.audio.delta":
		fmt.Print("-")
		var responseAudioDelta events.ResponseAudioDelta
		if err := json.Unmarshal(message, &responseAudioDelta); err != nil {
			log.Error("Error unmarshalling audio delta events:", err)
			return err
		}
		decoded, err := base64.StdEncoding.DecodeString(responseAudioDelta.Delta)
		if err != nil {
			log.Error("Error decoding PCM data:", err)
			return err
		}

		c.AudioOutputChan <- decoded
	case "response.audio.done":
		var audioDoneEvent events.ResponseAudioDone
		if err := json.Unmarshal(message, &audioDoneEvent); err != nil {
			log.Error("Error unmarshalling audio done events:", err)
			return err
		}

		c.status = StatusReady
	case "response.audio_transcript.delta":
		var audioTranscriptDeltaEvent events.ResponseAudioTranscriptDelta
		if err := json.Unmarshal(message, &audioTranscriptDeltaEvent); err != nil {
			log.Error("Error unmarshalling audio transcript delta events:", err)
			return err
		}
		fmt.Print(audioTranscriptDeltaEvent.Delta)
	case "response.audio_transcript.done":
		var audioTranscriptDoneEvent events.ResponseAudioTranscriptDone
		if err := json.Unmarshal(message, &audioTranscriptDoneEvent); err != nil {
			log.Error("Error unmarshalling audio transcript done events:", err)
			return err
		}
		fmt.Print("\n\n")
	case "response.content_part.done":
		var contentPartDone events.ResponseContentPartDone
		if err := json.Unmarshal(message, &contentPartDone); err != nil {
			log.Error("Error unmarshalling content part done events:", err)
			return err
		}
	case "response.output_item.done":
		var outputItemDone events.ResponseOutputItemDone
		if err := json.Unmarshal(message, &outputItemDone); err != nil {
			log.Error("Error unmarshalling output item done events:", err)
			return err
		}
	case "response.done":
		var responseDone events.ResponseDone
		if err := json.Unmarshal(message, &responseDone); err != nil {
			log.Error("Error unmarshalling response done events:", err)
			return err
		}

		if responseDone.Response.Status == "failed" {
			log.Error("Response failed:", responseDone.Response.StatusDetails.Error.Message)
			return fmt.Errorf("Response failed: %s", responseDone.Response.StatusDetails.Error.Message)
		}

	case "input_audio_buffer.speech_started":
		var speechStartedEvent events.InputAudioBufferSpeechStarted
		if err := json.Unmarshal(message, &speechStartedEvent); err != nil {
			log.Error("Error unmarshalling speech started events:", err)
			return err
		}
	case "input_audio_buffer.speech_stopped":
		var speechStoppedEvent events.InputAudioBufferSpeechStopped
		if err := json.Unmarshal(message, &speechStoppedEvent); err != nil {
			log.Error("Error unmarshalling speech stopped events:", err)
			return err
		}
	case "input_audio_buffer.committed":
		var audioBufferCommittedEvent events.InputAudioBufferCommitted
		if err := json.Unmarshal(message, &audioBufferCommittedEvent); err != nil {
			log.Error("Error unmarshalling input audio buffer committed events:", err)
			return err
		}
		//c.status = StatusProcessing

	case "response.created":
		var responseCreatedEvent events.ResponseCreated
		if err := json.Unmarshal(message, &responseCreatedEvent); err != nil {
			log.Error("Error unmarshalling response created events:", err)
			return err
		}

	case "rate_limits.updated":
		var rateLimitsUpdatedEvent events.RateLimitsUpdated
		if err := json.Unmarshal(message, &rateLimitsUpdatedEvent); err != nil {
			log.Error("Error unmarshalling rate limits updated events:", err)
			return err
		}

		logEventAsJSON("[RECV]", event, message)

	case "response.output_item.added":
		var responseOutputItemAddedEvent events.ResponseOutputItemAdded
		if err := json.Unmarshal(message, &responseOutputItemAddedEvent); err != nil {
			log.Error("Error unmarshalling response output item added events:", err)
			return err
		}

	case "response.content_part.added":
		var contentPartAddedEvent events.ResponseContentPartAdded
		if err := json.Unmarshal(message, &contentPartAddedEvent); err != nil {
			log.Error("Error unmarshalling response content part added events:", err)
			return err
		}
	case "conversation.item.input_audio_transcription.completed":
		var conversationItemInputAudioTranscriptionCompleted events.ConversationItemInputAudioTranscriptionCompleted
		if err := json.Unmarshal(message, &conversationItemInputAudioTranscriptionCompleted); err != nil {
			log.Error("Error unmarshalling conversation item input audio transcription completed events:", err)
			return err
		}
	default:
		stringMessage := string(message)
		log.Error("Unknown events:", stringMessage)
	}

	return nil
}

// 외부네어 변환 할 예정 ffmpeg -f s16le -ar 24000 -ac 1 -i input.pcm output.wav
func AppendPCMDataToFile(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}
