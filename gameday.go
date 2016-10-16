package mlbgameday

import (
	"encoding/xml"
	"fmt"
	"time"
)

// GamedayService is an interface for accessing game information from the MLB
// Gameday API.
type GamedayService interface {
	Scoreboard() (*Scoreboard, error)
	GamesInProgress() ([]Game, error)
	Game(string) (GameService, error)
}

// GamedayServiceOp communicates with the MLB Gameday API to access a game.
type GamedayServiceOp struct {
	// The MLB Gameday API client used to communicate with the MLB Gameday API.
	client *Client

	// The date with which this service is associated.
	date time.Time

	// The path to the MLB Gameday API at which resources for this day reside.
	path string
}

// NewGamedayService returns a new GamedayServiceOp for the provided date.
// Communication with the MLB Gameday API occurs through the provided Client.
// The provided date will be converted to ET to be consistent with the
// Gameday API.
func NewGamedayService(client *Client, date time.Time) GamedayService {
	l, _ := time.LoadLocation("America/New_York")

	return &GamedayServiceOp{
		client: client,
		date:   date.In(l),
		path:   pathFromDate(date),
	}
}

// Scoreboard represents a list of MLB games and their statuses.
type Scoreboard struct {
	Games []Game `xml:"game"`
}

// Scoreboard lists all the games scheduled on the provided date.
func (s *GamedayServiceOp) Scoreboard() (*Scoreboard, error) {
	sb, err := s.getScoreboard()
	if err != nil {
		return nil, err
	}

	return sb, nil
}

// GamesInProgress lists all the games for the current day that are in progress.
func (s *GamedayServiceOp) GamesInProgress() ([]Game, error) {
	sb, err := s.getScoreboard()
	if err != nil {
		return nil, err
	}

	var games []Game
	for _, game := range sb.Games {
		if game.Status == "In Progress" {
			games = append(games, game)
		}
	}

	return games, nil
}

// Game returns a new GameServiceOp linked to the game referenced by the
// provided game ID (gid).
func (s *GamedayServiceOp) Game(gid string) (GameService, error) {
	svc, err := NewGameService(s.client, gid)
	if err != nil {
		return nil, err
	}
	return svc, err
}

// getScoreboard retrieves all the games on the provided date.
func (s *GamedayServiceOp) getScoreboard() (*Scoreboard, error) {
	data, err := s.client.get(fmt.Sprintf("%s/miniscoreboard.xml", s.path))
	if err != nil {
		return nil, err
	}

	sb := new(Scoreboard)
	err = xml.Unmarshal(data, &sb)
	if err != nil {
		return nil, err
	}

	return sb, nil
}

// pathFromDate returns the URL path to all Gameday data for a given date.
func pathFromDate(date time.Time) string {
	return date.Format("components/game/mlb/year_2006/month_01/day_02/")
}
