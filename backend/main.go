package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// var templates = template.Must(template.ParseFiles(
// 	"edit.html",
// 	"view.html",
// 	"index.html",
// 	"newGame.html",
// 	"Game.html",
// ))

type Game struct {
	Cards   *[]Card
	Players map[string]int
	Code    int
	Started bool
}

type Card struct {
	Name string
}

type NewGameRequest struct {
	GameCode int
}

type StartRequest struct {
	GameCode int
}

type NewCardRequest struct {
	GameCode int
	CardName string
}

var games map[int]Game = make(map[int]Game)

// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "index", nil)
// }

func newGameHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("NEW: game request")
	var n NewGameRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &n); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if games[n.GameCode].Cards == nil {
		games[n.GameCode] = Game{Cards: new([]Card), Code: n.GameCode}
	}

	json.NewEncoder(w).Encode(games[n.GameCode])
}

// func playHandler(w http.ResponseWriter, r *http.Request) {
// 	code, err := getCode(r)
// 	if err != nil || !r.URL.Query().Has("name") {
// 		http.Redirect(w, r, "/", http.StatusBadRequest)
// 		return
// 	}
// 	if games[code].Players == nil || games[code].Cards == nil {
// 		http.Redirect(w, r, fmt.Sprint("/game/?code=%i", code), http.StatusBadRequest)
// 		return

// 	}

// 	name := r.URL.Query().Get("name")
// 	moderator := r.URL.Query().Has("moderator")
// 	taken := make([]int, 0, len(*games[code].Cards))
// 	for _, v := range games[code].Players {
// 		taken = append(taken, v)
// 	}
// 	available := make([]int, 0, len(*games[code].Cards))
// 	for i := range *games[code].Cards {
// 		if !slices.Contains(taken, i) {
// 			available = append(available, i)
// 		}
// 	}

// 	if !moderator {
// 		_, exists := games[code].Players[name]
// 		if len(available) == 0 && !exists {
// 			http.Redirect(w, r, fmt.Sprintf("/game/?code=%d", code), http.StatusFound)
// 			return
// 		}
// 		if !exists {
// 			games[code].Players[name] = available[rand.Intn(len(available))]
// 		}
// 	}

// 	renderTemplate(w, "Game", struct {
// 		Player    string
// 		Game      Game
// 		Moderator bool
// 	}{Player: name, Game: games[code], Moderator: moderator})

// }

func startHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NEW: game start request")
	var n StartRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &n); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	game := games[n.GameCode]
	game.Started = true
	game.Players = make(map[string]int)
	games[n.GameCode] = game
	w.WriteHeader(http.StatusOK)
}

func cardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NEW: card request")
	var n NewCardRequest

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &n); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	card := Card{n.CardName}

	if card.Name == "" || games[n.GameCode].Cards == nil {
		http.Error(w, "No Content", http.StatusNoContent)
		return
	}

	*games[n.GameCode].Cards = append(*games[n.GameCode].Cards, card)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(*(games[n.GameCode].Cards))

}

// func getCode(r *http.Request) (int, error) {
// 	var code int
// 	if !r.URL.Query().Has("code") {
// 		return 0, errors.New("no Code")
// 	}
// 	code, err := strconv.Atoi(r.URL.Query().Get("code"))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return code, nil
// }

func main() {
	http.HandleFunc("/api/game/", newGameHandler)
	http.HandleFunc("/api/card/", cardHandler)
	// http.HandleFunc("/api/playing/", playHandler)
	http.HandleFunc("/api/startGame/", startHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
