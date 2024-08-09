package generator

import (
	"testing"
)

func TestGenerateWords(t *testing.T) {
    resp, err := GenerateWords("http://localhost:11434/api/generate")

    if err != nil {
        t.Fatalf("Error occured during generation: %v", err)
    }

    if resp == nil {
        t.Fatalf("No Response. Expected Response from generation.")
    }
}

// func TestGenerateFromLLM(t *testing.T) {
//     resp, err := GenerateFromLLM("http://localhost:11434/api/generate")
//
//     if err != nil {
//         t.Fatalf("Error occured during generation: %v", err)
//     }
//
//     if resp == nil {
//         t.Fatalf("No Response. Expected Response from generation.")
//     }
//
//     fmt.Println(resp.Response)
// }
