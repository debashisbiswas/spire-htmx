package entry

import (
	"slices"
	"testing"
)

func TestSerializeEmbeddings(t *testing.T) {
	actual := SerializeEmbeddings(Vector{1, 2, 3})
	expected := "'[1,2,3]'"
	if expected != actual {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestDeserializeEmbeddings(t *testing.T) {
	actual, err := DeserializeEmbeddings("[1,2,3]")

	if err != nil {
		t.Fatal(err)
	}

	expected := Vector{1, 2, 3}
	if !slices.Equal(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestSerializeEmbeddingsVector(t *testing.T) {
	actual := SerializeEmbeddingsWithVectorPrefix(Vector{1, 2, 3})
	expected := "vector('[1,2,3]')"
	if expected != actual {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
