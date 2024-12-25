package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
)

var templates = template.Must(template.ParseFiles(
	"edit.html",
	"view.html",
	"index.html",
	"newGame.html",
	"Game.html",
))

type Game struct {
	Cards   *[]Card
	Players map[string]int
	Code    int
	Started bool
}

type Card struct {
	Name string
}

var games map[int]Game = make(map[int]Game)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	code, err := getCode(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	if games[code].Cards == nil {
		games[code] = Game{Cards: new([]Card), Code: code}
	}

	renderTemplate(w, "newGame", games[code])
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	code, err := getCode(r)
	if err != nil || !r.URL.Query().Has("name") {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	name := r.URL.Query().Get("name")
	moderator := r.URL.Query().Has("moderator")
	taken := make([]int, 0, len(*games[code].Cards))
	for _, v := range games[code].Players {
		taken = append(taken, v)
	}
	available := make([]int, 0, len(*games[code].Cards))
	for i := range *games[code].Cards {
		if !slices.Contains(taken, i) {
			available = append(available, i)
		}
	}

	if !moderator {
		_, exists := games[code].Players[name]
		if len(available) == 0 && !exists {
			http.Redirect(w, r, fmt.Sprintf("/game/?code=%d", code), http.StatusFound)
			return
		}
		if !exists {
			games[code].Players[name] = available[rand.Intn(len(available))]
		}
	}

	renderTemplate(w, "Game", struct {
		Player    string
		Game      Game
		Moderator bool
	}{Player: name, Game: games[code], Moderator: moderator})

}
func startHandler(w http.ResponseWriter, r *http.Request) {
	code, err := getCode(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	game := games[code]
	game.Started = true
	game.Players = make(map[string]int)
	games[code] = game
	http.Redirect(w, r, fmt.Sprintf("/game?code=%d", code), http.StatusFound)

}

func cardHandler(w http.ResponseWriter, r *http.Request) {
	code, err := getCode(r)
	if err != nil || !r.URL.Query().Has("card") {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	card := Card{r.URL.Query().Get("card")}

	if card.Name == "" || games[code].Cards == nil {
		http.Redirect(w, r, fmt.Sprintf("/game/?code=%d", code), http.StatusFound)
		return
	}
	*games[code].Cards = append(*games[code].Cards, card)
	http.Redirect(w, r, fmt.Sprintf("/game/?code=%d", code), http.StatusFound)

}

func getCode(r *http.Request) (int, error) {
	var code int
	if !r.URL.Query().Has("code") {
		return 0, errors.New("no Code")
	}
	code, err := strconv.Atoi(r.URL.Query().Get("code"))
	if err != nil {
		return 0, err
	}
	return code, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p any) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/game/", newGameHandler)
	http.HandleFunc("/addCard/", cardHandler)
	http.HandleFunc("/playing/", playHandler)
	http.HandleFunc("/startGame/", startHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
