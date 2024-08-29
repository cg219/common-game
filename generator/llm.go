package generator

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	_ "modernc.org/sqlite"

	storage "github.com/cg219/common-game/internal/wordsdb"
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

func GenerateSubjects(url string) error {
    ctx := context.Background()
    ddl, err := os.ReadFile("../configs/words-schema.sql")
    if err != nil {
        return err
    }

    db, err := sql.Open("sqlite", "../words.db")

    if err != nil {
        return err
    }

    if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
        return err
    }

    q := storage.New(db)
    defer db.Close()

    words, err := q.GetWords(ctx)

    if err != nil {
        return err
    }

    example := `{ "list": [
        {
            "subject": "Music Festival",
            "words": ["concertina", "giggle", "gibbous", "glamour"]
        },
        {
            "subject": "Space Exploration",
            "words": ["galaxies", "geese", "gravity", "habitat"]
        },
        {
            "subject": "Foodie Culture",
            "words": ["biscuit", "cantaloupe", "chocolate", "honey"]
        },
        {
            "subject": "Artistic Expression",
            "words": ["canvas", "coloratura", "painter", "sculptor"]
        },
        {
            "subject": "Travel Adventure",
            "words": ["backpack", "beachball", "passport", "trekking"]
        },
        {
            "subject": "Sports and Fitness",
            "words": ["gymnast", "marathon", "runner", "weightlifter"]
        },
        {
            "subject": "Environmental Awareness",
            "words": ["conservation", "ecosystem", "recycle", "sustainability"]
        },
        {
            "subject": "Culinary Delights",
            "words": ["butterfly", "caterpillar", "honeycomb", "jam"]
        },
        {
            "subject": "Historical Significance",
            "words": ["archaeologist", "historian", "museum", "relic"]
        },
        {
            "subject": "Scientific Discovery",
            "words": ["atom", "element", "experiment", "lab"]
        },
        {
            "subject": "Fashion and Style",
            "words": ["couture", "designer", "fashionista", "model"]
        },
        {
            "subject": "Cultural Exchange",
            "words": ["diplomat", "embassy", "language", "translation"]
        },
        {
            "subject": "Outdoor Activities",
            "words": ["camping", "hiking", "outdoor", "survival"]
        },
        {
            "subject": "Educational Institutions",
            "words": ["classroom", "library", "professor", "university"]
        },
        {
            "subject": "Social Justice",
            "words": ["activist", "charity", "human rights", "protest"]
        }]}`
    prompt := fmt.Sprintf("Given this list of words: '%s'.\n\nCreate 30 unique groupings of 4 words from the previous list given and create a subject for them that relates them to one another. The words chosenshould not be in the subject. Make sure the words chosen only come from the given list. Do not use other words outside the list.\n\nHere is an example of the output:\n\n%s\n\nDo not use any of the subjects from the example. Return only the JSON and nothing else.", strings.Join(words, ", "), example)

    params := &GenerateParams{
        Model: "cllama",
        Prompt: prompt,
        System: "You are an expert of the English Language, who spends a lot of time consuming pop and hip hop culture through social media and entertainment.",
        Format: "",
        Stream: true,
    }

    body, err := json.Marshal(params)

    if err != nil {
        panic(err)
    }
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

    if err != nil {
        return err
    }

    defer resp.Body.Close()

    scanner := bufio.NewScanner(resp.Body)
    buffer := make([]byte, 0, 512000)
    scanner.Buffer(buffer, 512000)
    // ch := make(chan string, 32)

    // go func() {
    //     var word strings.Builder
    //     
    //     for w := range ch {
    //         if !strings.Contains(w, ",") {
    //             word.WriteString(w)
    //             continue
    //         }
    //
    //         tw := strings.Split(word.String(), ",")
    //         word.Reset()
    //
    //         if len(tw) > 1 {
    //             s := strings.TrimSpace(tw[0])
    //             q.SaveWord(ctx, s)
    //             word.WriteString(tw[1])
    //         } else {
    //             s := tw[0]
    //             q.SaveWord(ctx, strings.TrimSpace(s))
    //         }
    //     }
    // }()

    for scanner.Scan() {
        data := &GenerateResponse{}
        bts := scanner.Bytes()

        if err := json.Unmarshal(bts, data); err != nil {
            panic(err)
        }

        fmt.Print(data.Response)

        // ch <- data.Response 
    }

    // close(ch)
    return err
}

func GenerateWords(url string) error {
    ctx := context.Background()
    ddl, err := os.ReadFile("../configs/words-schema.sql")
    if err != nil {
        return err
    }

    db, err := sql.Open("sqlite", "../words.db")

    if err != nil {
        return err
    }

    if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
        return err
    }

    q := storage.New(db)
    defer db.Close()

    params := &GenerateParams{
        Model: "cllama",
        Prompt: "Create a list of 50 common English words that start with 'u', longer than 2 characters without repeating any. The response should ONLY be a comma separated list without any pretext.",
        System: "You are an expert of the English Language, who spends a lot of time consuming pop and hip hop culture through social media and entertainment.",
        Format: "",
        Stream: true,
    }

    body, err := json.Marshal(params)

    if err != nil {
        panic(err)
    }
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

    if err != nil {
        return err
    }

    defer resp.Body.Close()

    scanner := bufio.NewScanner(resp.Body)
    buffer := make([]byte, 0, 512000)
    scanner.Buffer(buffer, 512000)
    ch := make(chan string, 32)

    go func() {
        var word strings.Builder
        
        for w := range ch {
            if !strings.Contains(w, ",") {
                word.WriteString(w)
                continue
            }

            tw := strings.Split(word.String(), ",")
            word.Reset()

            if len(tw) > 1 {
                s := strings.TrimSpace(tw[0])
                q.SaveWord(ctx, s)
                word.WriteString(tw[1])
            } else {
                s := tw[0]
                q.SaveWord(ctx, strings.TrimSpace(s))
            }
        }
    }()

    for scanner.Scan() {
        data := &GenerateResponse{}
        bts := scanner.Bytes()

        if err := json.Unmarshal(bts, data); err != nil {
            panic(err)
        }

        fmt.Print(data.Response)

        ch <- data.Response 
    }

    close(ch)
    return err
}


