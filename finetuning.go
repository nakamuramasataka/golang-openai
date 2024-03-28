package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FineTuningJobRequest struct {
	Model           string                 `json:"model"`
	TrainingFile    string                 `json:"training_file"`
	Hyperparameters map[string]interface{} `json:"hyperparameters,omitempty"`
	Suffix          *string                `json:"suffix,omitempty"`
	ValidationFile  *string                `json:"validation_file,omitempty"`
}

// FineTuningJobResponse はファインチューニングジョブのレスポンスを表します。
type FineTuningJobResponse struct {
	ID              string                  `json:"id"`
	CreatedAt       int64                   `json:"created_at"`
	Error           *map[string]interface{} `json:"error"`
	FineTunedModel  *string                 `json:"fine_tuned_model"`
	FinishedAt      *int64                  `json:"finished_at"`
	Hyperparameters map[string]interface{}  `json:"hyperparameters"`
	Model           string                  `json:"model"`
	Object          string                  `json:"object"`
	OrganizationID  string                  `json:"organization_id"`
	ResultFiles     []string                `json:"result_files"`
	Status          string                  `json:"status"`
	TrainedTokens   *int64                  `json:"trained_tokens"`
	TrainingFile    string                  `json:"training_file"`
	ValidationFile  *string                 `json:"validation_file"`
}

// CreateFineTuningJob は新しいファインチューニングジョブを作成します。
func (c *Client) CreateFineTuningJob(req *FineTuningJobRequest) (*FineTuningJobResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/fine_tuning/jobs", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request to fine tuning jobs API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fine tuning jobs API returned non-200 status code: %d", resp.StatusCode)
	}

	var fineTuningJobResponse FineTuningJobResponse
	err = json.NewDecoder(resp.Body).Decode(&fineTuningJobResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &fineTuningJobResponse, nil
}
