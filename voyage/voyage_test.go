package voyage

import (
	"os"
	"testing"
)

func TestParseEmbeddingResponse(t *testing.T) {
	data, err := os.ReadFile("./testfiles/example_embedding_response.json")
	if err != nil {
		t.Fatal(err)
	}

	result, err := parseEmbeddingResponse(data)
	if err != nil {
		t.Fatal(err)
	}

	// Test root-level properties

	if result.Object != "list" {
		t.Errorf(`root-level object key should be "list", but got "%s"`, result.Object)
	}

	if result.Model != "voyage-3-lite" {
		t.Errorf(`root-level model key should be "voyage-3-lite", but got "%s"`, result.Model)
	}

	if result.Usage.TotalTokens != 4 {
		t.Errorf("total tokens should be 4, but got %d", result.Usage.TotalTokens)
	}

	// Test properties of data field

	if len(result.Data) != 1 {
		t.Fatalf("result data length should be 1, but got %d", len(result.Data))
	}

	dataField := result.Data[0]

	if dataField.Object != "embedding" {
		t.Errorf(`data's object key should be "embedding", but got "%s"`, dataField.Object)
	}

	if len(dataField.Embedding) != 512 {
		t.Errorf("data's embedding length should be 512, but got %d", len(dataField.Embedding))
	}
}
