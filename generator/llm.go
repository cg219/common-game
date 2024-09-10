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

    _ "github.com/tursodatabase/go-libsql"

	"github.com/cg219/common-game/internal/data"
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
    ddl, err := os.ReadFile("../configs/schema.sql")
    if err != nil {
        return err
    }

    dbc, err := sql.Open("libsql", "file:../database.db")
    if err != nil {
        return err
    }

    defer dbc.Close()

    if _, err := dbc.ExecContext(ctx, string(ddl)); err != nil {
        return err
    }

    q := data.New(dbc)
    fmt.Println(q)

    return nil

    // words := "The Breakfast Club, The Big Lebowski, The Shawshank Redemption, The Silence of the Lambs, The English Patient, The Sixth Sense, The Princess Bride, The Green Mile, The Godfather, The Good, The Bad and the Ugly, The Bourne Identity, The Dark Knight, The Lord of the Rings, The Fellowship of the Ring, The Two Towers, The Return of the King, The Curious Case of Benjamin Button, The Prestige, The Aviator, The Departed, The Pianist, The Lives of Others, The Hurt Locker, The Social Network, The Wolf of Wall Street, The Grand Budapest Hotel, The Imitation Game, The Theory of Everything, The Martian, The Revenant, The Shape of Water, The Irishman, The Favourite, The Lighthouse, The Souvenir, The Beach, The Perks of Being a Wallflower, The Spectre, The Huntsman: Winter's War, The Nice Guys, The Hateful Eight, The Girl with the Dragon Tattoo"
    // words := "Fargo, Flashback, Frailty, Frighteners, Galapagos, Gamer, Gigli, Gnomeo, Gone, Gorillas, Gravity, Gremlins, Hackers, Holes, Honey, Hoodwinked, Hostel, Hotshot, Houseboat, Howl, Hunger, Imitation, Insidious, Intruder, Jumper, Kicking, Killshot, Kindergarten, Ladykillers, Lastex, Liar, Luckytown, Madhouse, Magnolia, Malice, Marathon, Max, Misery, Mockingbird, Moneyball, Monster, Mortal, Mr., Mummy, Narcos, Nighthawks, Nomad, Norma, Oblivion, Ocean, Omen, Original, Outlaw, Pulp, Punchline, Quicksilver, Ransom, Razorback, Rebound, Redbox, Restless, Revenge, Ripley, Roadhouse, RoboCop, Rocky, Romancing, Rushmore, Sabotage, Salute, Sanctum, Savage, Sawdust, Shaft, Shelter, Shutter, Sideways, Slaughterhouse, Sleuth, Smokin', Snatch, Speed, Stakeout, Stalker, Stranger, Streetcar, Suburbia, Sunshine, Sweet, Switchback, Tank, Takedown, Tango, Taxman, Tempest, Terror, Thinner, Threshold, Thunderbolt, Torn, Traffic, Trance, Triangle, Troublemaker, Truth, Twist, Twister, Unforgiven, Vengeance, Vigilante, Violent, Virtuoso, Wakefield, Waking, Walkabout, Warlock, Waterboy, Waterloo, Webbed, Weekend, Werewolf, Whiplash, Whiteout, Wichita, Wilder, Windtalkers, Wishmaster, Witness"
    // words := "Tweets, Follows, Likes, Shares, Comments, Posts, Updates, Status, Profile, Bio, Handle, Tagged, Mentions, Hashtags, Filters, Effects, Emojis, GIFs, Videos, Clips, LiveStream, Reactions, Views, Engagement, Follower, Unfollow, Block, Report, Flag, Spam, Trolls, Cyberbullying, Online, Presence, Identity, Reputation, Branding, Marketing, Influencer, Sponsorship, Partnership, Collaboration, Community, Forum, Discussion, Polls, Quizzes, Games, Challenges, Viral, Trending, Meme, Joke, Humor, Satire, Irony, Sarcasm, Wit, Censorship, Freedom, Speech, Expression, Opinion, Perspective, Viewpoint, Insight, Analysis, Criticism, Review, Rating, Recommendation, Endorsement, Affiliate, Link, URL, Bookmark, Favorite, Shareable, Clickbait, Headline, Caption, Description, Keyword, SEO, Optimization, Algorithm, Feed, Timeline, ProfilePic, Avatar, Username, Password, Security, Safety, Privacy, Policy, Terms"

    words := "Tweet, gram, snap, chat, post, share, like, follow, unfollow, block, mute, hashtag, tag, caption, comment, reply, retweet, favorite, dislike, report, spam, troll, cyberbully, online, offline, profile, bio, username, password, login, logout, account, verification, badge, emoji, GIF, meme, viral, trending, algorithm, feed, timeline, notification, alert, buzz, hype, trend, fad, craze, phenomenon, movement, community, forum, group, chatroom, socialite, influencer, celebrity, fan, follower, supporter, hater, critic, troll, cyberstalker, harassment, digital footprint, internet fame, viral, trending, hashtag, livestream, vlog, blog, podcast, stream, channel, playlist, video, audio, music, art, culture, lifestyle, fashion, beauty, fitness, wellness, gaming, esports, streaming, multiplayer, game, player, gamer, avatar, character, roleplay, virtual, metaverse"

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
    prompt := fmt.Sprintf("Given this list of words: '%s'.\n\nCreate 50 unique groupings of 4 words from the previous list given and create a subject for them that relates them to one another. The words chosen should not be in the subject at all. Make sure the words chosen only come from the given list. Do not use other words outside the given list.\n\nHere is an example of the output:\n\n%s\n\nDo not use any of the subjects from the example. Do no put 'Types of' in the subject. Subjects shouldn't be more than 4 words long. Return only the JSON and nothing else.", words, example)

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

    fmt.Println(data.Response)

    if err := json.Unmarshal([]byte(data.Response), subjects); err != nil {
        panic(err)
    }

    // for _, v := range subjects.List {
    //     wd := &SubjectWordsJson{ Words: v.Words }
    //
    //     w, err := json.Marshal(wd); 
    //
    //     if err != nil {
    //         panic(err)
    //     }
    //     
    //     s := subjectsdb.SaveSubjectParams{ Subject: v.Subject, Words: string(w)}
    //     sq.SaveSubject(ctx, s)
    // }

    return err
}

func GenerateWords(url string) error {
    ctx := context.Background()
    ddl, err := os.ReadFile("../configs/schema.sql")
    if err != nil {
        return err
    }

    db, err := sql.Open("libsql", "file:../database.db")

    if err != nil {
        return err
    }

    if _, err := db.ExecContext(ctx, string(ddl)); err != nil {
        return err
    }

    q := data.New(db)
    fmt.Println(q)
    defer db.Close()

    params := &GenerateParams{
        Model: "cllama",
        Prompt: "Create a list of 75 words related to social media where the word is longer than 2 characters but less than 15 without repeating any words. The response should ONLY be a comma separated list without any pretext.",
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
                // s := strings.TrimSpace(tw[0])
                // q.SaveWord(ctx, s)
                word.WriteString(tw[1])
            } else {
                // s := tw[0]
                // q.SaveWord(ctx, strings.TrimSpace(s))
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


