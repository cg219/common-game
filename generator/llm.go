package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GenerateParams struct {
    Model string `json:"model"`
    Prompt string `json:"prompt"`
    Suffix string `json:"suffix"`
    Images []string `json:"images"`
    Format string `json:"format"`
    Stream bool `json:"stream"`
    System string `json:"system"`
}

type GenerateResponse struct {
    Response string `json:"response"`
}

func GenerateFromLLM(url string) (*GenerateResponse, error) {
    // body := []byte(`{"model":"llama3.1","prompt":"Describe the sky as you would to a blind person"}`)
    params := &GenerateParams{
        Model: "llama3.1",
        // Prompt: "Create a list of 10000 words from the dictionary and urban dictionary. From these 10000, create 15 groups of 4 words that relate to a detailed category. Return this in a JSON format where that is an array with a 'category' key for the category and a 'words' key that is an array of the 4 associated words.",
        Prompt: "Create a list of 10000 unique English words that can also include proper nouns. Return this in a JSON format where that is an array of the words generated in the previous step and place it under the key 'words'.",
        System: "You are an expert of the English Language.",
        Format: "json",
        Stream: false,
    }

    body, err := json.Marshal(params)

    if err != nil {
        panic(err)
    }
    
    // prompt :=
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    data := &GenerateResponse{}

    decoder := json.NewDecoder(resp.Body)

    err = decoder.Decode(&data)

    if err != nil {
        panic(err)
    }

    return data, nil
}

func GenerateWords(url string) ([]string, error) {
    params := &GenerateParams{
        Model: "llama3.1",
        Prompt: "Create a list of 100 unique English words that can also include proper nouns. Return in JSON format under the key 'words'.",
        System: "You are an expert of the English Language.",
        Format: "json",
        Stream: false,
    }

    body, err := json.Marshal(params)

    if err != nil {
        panic(err)
    }
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    data := &GenerateResponse{}

    decoder := json.NewDecoder(resp.Body)

    err = decoder.Decode(&data)

    if err != nil {
        panic(err)
    }

    fmt.Println(data.Response)

    return make([]string, 1), err

}
