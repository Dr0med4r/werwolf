package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

func newGameHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("NEW: game request")
	var n NewGameRequest

	// Parsing
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

	// Processing
	if games[n.GameCode].Cards == nil {
		games[n.GameCode] = Game{Cards: new([]Card), Code: n.GameCode}
	}

	// Response
	json.NewEncoder(w).Encode(games[n.GameCode])
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NEW: game start request")

	// Parsing
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

	// Processing
	if games[n.GameCode].Started {
		http.Error(w, "Already Reported", http.StatusAlreadyReported)
		return
	}

	game := games[n.GameCode]
	game.Started = true
	game.Players = make(map[string]int)
	games[n.GameCode] = game

	// Response
	w.WriteHeader(http.StatusOK)
}

func cardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NEW: card request")

	// Parsing
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

	// Processing
	if games[n.GameCode].Started {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if n.CardName == "" || games[n.GameCode].Cards == nil {
		http.Error(w, "No Content", http.StatusNoContent)
		return
	}

	card := Card{n.CardName}

	*games[n.GameCode].Cards = append(*games[n.GameCode].Cards, card)

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(*(games[n.GameCode].Cards))

}

func main() {
	http.HandleFunc("/api/game/", newGameHandler)
	http.HandleFunc("/api/card/", cardHandler)
	http.HandleFunc("/api/startGame/", startHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
