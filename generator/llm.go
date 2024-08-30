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

	"github.com/cg219/common-game/internal/subjectsdb"
	"github.com/cg219/common-game/internal/wordsdb"
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

type SubjectsResponse struct {
    List []struct {
        Subject string `json:"subject"`
        Words []string `json:"words"`
    } `json:"list"`
}

type SubjectWordsJson struct {
    Words []string `json:"words"`
}

func GenerateSubjects(url string) error {
    ctx := context.Background()
    ddl, err := os.ReadFile("../configs/words-schema.sql")
    if err != nil {
        return err
    }

    sddl, err := os.ReadFile("../configs/subjects-schema.sql")
    if err != nil {
        return err
    }

    wdb, err := sql.Open("sqlite", "../words.db")
    if err != nil {
        return err
    }

    defer wdb.Close()

    sdb, err := sql.Open("sqlite", "../subjects.db")
    if err != nil {
        return err
    }

    defer sdb.Close()

    if _, err := wdb.ExecContext(ctx, string(ddl)); err != nil {
        return err
    }

    if _, err := sdb.ExecContext(ctx, string(sddl)); err != nil {
        return err
    }

    q := wordsdb.New(wdb)
    sq := subjectsdb.New(sdb)

    fmt.Println(sq)

    words, err := q.GetWords(ctx)

    if err != nil {
        return err
    }

    example := `{ "list": [
        {
            "subject": "Venues",
            "words": ["arena", "stadium", "club", "festival"]
        },
        {
            "subject": "Nuts",
            "words": ["peanut", "cashew", "almond", "pecan"]
        },
        {
            "subject": "Cake Ingredients",
            "words": ["egg", "flour", "vanilla", "chocolate"]
        },
        {
            "subject": "Creative Expression",
            "words": ["painter", "musician", "author", "sculptor"]
        },
        {
            "subject": "Ways to Travel",
            "words": ["plane", "train", "boat", "bicycle"]
        },
        {
            "subject": "Sports Events",
            "words": ["gymnastics", "track", "basketball", "baseball"]
        },
        {
            "subject": "Environmental Awareness",
            "words": ["conservation", "ecosystem", "recycle", "reservoir"]
        },
        {
            "subject": "Culinary Delights",
            "words": ["cake", "croissant", "wellington", "crepe"]
        },
        {
            "subject": "Deals With The Past",
            "words": ["archaeologist", "historian", "museum", "relic"]
        },
        {
            "subject": "Scientific",
            "words": ["atom", "element", "experiment", "lab"]
        },
        {
            "subject": "Fashion and Style",
            "words": ["couture", "designer", "fashionista", "model"]
        },
        {
            "subject": "Gaming",
            "words": ["twitch", "console", "handheld", "controller"]
        },
        {
            "subject": "Outdoor Activities",
            "words": ["camping", "hiking", "skydiving", "walking"]
        },
        {
            "subject": "Educational",
            "words": ["classroom", "library", "professor", "university"]
        },
        {
            "subject": "Weather Conditions",
            "words": ["rainy", "sunny", "cloudy", "clear"]
        }]}`
    prompt := fmt.Sprintf("Given this list of words: '%s'.\n\nCreate 30 unique groupings of 4 words from the previous list given and create a subject for them that relates them to one another. The words chosen should not be in the subject. Make sure the words chosen only come from the given list. Do not use other words outside the given list.\n\nHere is an example of the output:\n\n%s\n\nDo not use any of the subjects from the example. Do no put 'Types of' in the subject. Return only the JSON and nothing else.", strings.Join(words, ", "), example)

    params := &GenerateParams{
        Model: "cllama",
        Prompt: prompt,
        System: "You are an expert of the English Language, who spends a lot of time consuming pop and hip hop culture through social media and entertainment. Do not hallucinate.",
        Format: "",
        Stream: false,
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

    data := &GenerateResponse{}

    if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
        panic(err)
    }

    subjects := &SubjectsResponse{}

    if err := json.Unmarshal([]byte(data.Response), subjects); err != nil {
        panic(err)
    }

    for _, v := range subjects.List {
        wd := &SubjectWordsJson{ Words: v.Words }

        w, err := json.Marshal(wd); 

        if err != nil {
            panic(err)
        }
        
        s := subjectsdb.SaveSubjectParams{ Subject: v.Subject, Words: string(w)}
        sq.SaveSubject(ctx, s)
    }

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

    q := wordsdb.New(db)
    fmt.Println(q)
    defer db.Close()

    params := &GenerateParams{
        Model: "cllama",
        Prompt: "Create a list of 30 common English words that start with the letter 'a', longer than 2 characters without repeating any. The response should ONLY be a comma separated list without any pretext.",
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


