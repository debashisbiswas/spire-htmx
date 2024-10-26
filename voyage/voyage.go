package voyage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type VoyageClient struct {
	apiKey string
}

func NewClient(apiKey string) VoyageClient {
	return VoyageClient{apiKey: apiKey}
}

type voyageEmbeddingResponse struct {
	Object string // always "list"
	Data   []struct {
		Object    string // always "embedding"
		Embedding []float32
		Index     int
	}
	Model string
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	}
}

func parseEmbeddingResponse(input []byte) (voyageEmbeddingResponse, error) {
	parsed := voyageEmbeddingResponse{}
	json.Unmarshal(input, &parsed)
	return parsed, nil
}

func (vc VoyageClient) GetEmbedding(input string) ([]float32, error) {
	requestBody := struct {
		model string
		input string
	}{
		model: "voyage-3-lite",
		input: input,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		"POST",
		"https://api.voyageai.com/v1/embeddings",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+vc.apiKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	result, err := parseEmbeddingResponse(body)
	if err != nil {
		return nil, err
	}

	if len(result.Data) != 1 {
		return nil, errors.New(
			fmt.Sprintf("expected 1 result from API, but got %d",
			len(result.Data)),
		)
	}

	return result.Data[0].Embedding, err
}
