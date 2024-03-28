package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type SpeechRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
}

type TranscriptionRequest struct {
	AudioBuffer            *bytes.Buffer
	FileName               string
	Model                  string
	Language               string
	Prompt                 string
	ResponseFormat         string
	Temperature            float64
	TimestampGranularities []string
}

type TranscriptionResponse struct {
	Text     string    `json:"text"`
	Language string    `json:"language,omitempty"`
	Duration float64   `json:"duration,omitempty"`
	Segments []Segment `json:"segments,omitempty"`
}

type Segment struct {
	ID               int     `json:"id"`
	Seek             float64 `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

type TranslationRequest struct {
	FileName       string
	AudioBuffer    *bytes.Buffer
	Model          string
	Prompt         string
	ResponseFormat string
	Temperature    float64
}

type TranslationResponse struct {
	Text string `json:"text"`
}

func (c *Client) CreateSpeech(req *SpeechRequest) ([]byte, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	endpoint := "https://api.openai.com/v1/audio/speech"
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request to audio/speech API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audio/speech API returned non-200 status code: %d", resp.StatusCode)
	}

	audioBuff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return audioBuff, nil
}

func (c *Client) CreateTranscription(req *TranscriptionRequest) (*TranscriptionResponse, error) {
	endpoint := "https://api.openai.com/v1/audio/transcriptions"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 音声データをマルチパートフォームデータに追加
	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}
	if _, err := part.Write(req.AudioBuffer.Bytes()); err != nil {
		return nil, fmt.Errorf("error writing audio buffer to form: %w", err)
	}

	// モデルと言語をフォームに追加（オプション）
	_ = writer.WriteField("model", req.Model)
	if req.Language != "" {
		_ = writer.WriteField("language", req.Language)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	httpReq, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request to transcription API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transcription API returned non-200 status code: %d", resp.StatusCode)
	}

	var transcriptionResponse TranscriptionResponse
	err = json.NewDecoder(resp.Body).Decode(&transcriptionResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &transcriptionResponse, nil
}

func (c *Client) CreateTranslation(req *TranslationRequest) (*TranslationResponse, error) {
	endpoint := "https://api.openai.com/v1/audio/translations"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}
	if _, err := part.Write(req.AudioBuffer.Bytes()); err != nil {
		return nil, fmt.Errorf("error writing audio buffer to form: %w", err)
	}

	_ = writer.WriteField("model", req.Model)
	if req.Prompt != "" {
		_ = writer.WriteField("prompt", req.Prompt)
	}
	if req.ResponseFormat != "" {
		_ = writer.WriteField("response_format", req.ResponseFormat)
	}
	if req.Temperature != 0 {
		_ = writer.WriteField("temperature", fmt.Sprintf("%f", req.Temperature))
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	httpReq, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request to translation API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("translation API returned non-200 status code: %d", resp.StatusCode)
	}

	var translationResponse TranslationResponse
	err = json.NewDecoder(resp.Body).Decode(&translationResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &translationResponse, nil
}
