package entry

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Vector []float32

type Entry struct {
	Time      time.Time
	Content   string
	Embedding Vector
}

// string [1,2,3] -> floats [1, 2, 3]
func DeserializeEmbeddings(input string) (Vector, error) {
	var result Vector

	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// floats [1, 2, 3] -> string vector('[1,2,3]')
func SerializeEmbeddingsWithVectorPrefix(embeddings Vector) string {
	builder := strings.Builder{}
	builder.WriteString("vector(")
	builder.WriteString(SerializeEmbeddings(embeddings))
	builder.WriteString(")")

	return builder.String()
}

// floats [1, 2, 3] -> string '[1,2,3]'
func SerializeEmbeddings(embeddings Vector) string {
	strNums := make([]string, len(embeddings))
	for i, num := range embeddings {
		strNums[i] = strconv.FormatFloat(float64(num), 'f', -1, 32)
	}

	builder := strings.Builder{}
	builder.WriteString("'[")
	builder.WriteString(strings.Join(strNums, ","))
	builder.WriteString("]'")

	return builder.String()
}
