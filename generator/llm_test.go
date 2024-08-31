package generator

import (
	"testing"
)

func ex(t *testing.T) {

}

// func TestGenerateWords(t *testing.T) {
//     err := GenerateWords("http://localhost:11434/api/generate")
//
//     if err != nil {
//         t.Fatalf("Error occured during generation: %v", err)
//     }
// }

// func TestGenerateSubjects(t *testing.T) {
//     err := GenerateSubjects("http://localhost:11434/api/generate")
//
//     if err != nil {
//         t.Fatalf("Error occured during generation: %v", err)
//     }
// }

// func TestGetWords(t *testing.T) {
//     ctx := context.Background()
//     ddl, err := os.ReadFile("../configs/words-schema.sql")
//     if err != nil {
//         t.Fatal(err)
//     }
//
//     db, err := sql.Open("sqlite", "../words.db")
//
//     if err != nil {
//         t.Fatal(err)
//     }
//
//     if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
//         t.Fatal(err)
//     }
//
//     q := storage.New(db)
//     defer db.Close()
//
//     res, err := q.GetWords(ctx)
//
//     if err != nil {
//         t.Fatal(err)
//     }
//
//     fmt.Print(res)
//     
// }

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
