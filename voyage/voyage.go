package voyage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"spire/entry"
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
		Embedding entry.Vector
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

func (vc VoyageClient) GetEmbedding(input string) (entry.Vector, error) {
	requestBody := struct {
		Model string `json:"model"`
		Input string `json:"input"`
	}{
		Model: "voyage-3-lite",
		Input: input,
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

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf(
			"server returned status %d with data %s",
			resp.StatusCode,
			string(body),
		))
	}

	result, err := parseEmbeddingResponse(body)
	if err != nil {
		return nil, err
	}

	if len(result.Data) != 1 {
		return nil, errors.New(fmt.Sprintf(
			"expected 1 result from API, but got %d",
			len(result.Data),
		))
	}

	return result.Data[0].Embedding, err
}
