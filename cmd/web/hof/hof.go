package hof

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"tpf-aram-hof/cmd/database"
)

func HofPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	db := database.New()

	date, err := time.Parse("02/01/2006", r.FormValue("date"))
	if err != nil {
		date = time.Now()
	}

	player := database.Player{
		Nickname:    r.FormValue("nickname"),
		Champion:    r.FormValue("champion"),
		Description: r.FormValue("description"),
		GameMode:    r.FormValue("game_mode"),
		Frag:        r.FormValue("frag"),
		Date:        date,
	}

	db.CreatePlayer(player)

	component := HofSuccessMessage(player.Nickname)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HofPostHandler: %e", err)
	}
}

func HofGetHandler(w http.ResponseWriter, r *http.Request) {
	db := database.New()

	players, err := db.GetPlayers()
	if err != nil {
    log.Printf("Error getting players: %e", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	component := HofList(players)

	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HofGetHandler: %e", err)
	}
}

func HofDeleteHandler(w http.ResponseWriter, r *http.Request) {
  playerID := r.URL.Query().Get("playerId")
  id, err := strconv.Atoi(playerID)
  if err != nil {
    http.Error(w, "Bad Request", http.StatusBadRequest)
  }
  db := database.New()

  _, err = db.DeletePlayer(id)
  if err != nil {
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }

  http.Redirect(w, r, "/hof/players", http.StatusSeeOther)
}
