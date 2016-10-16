package mlbgameday

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// GameService is an interface for retrieving information for a single game
// from the MLB Gameday API.
type GameService interface {
	LineScore() (*LineScore, error)
	Players() (*Players, error)
	AtBats() (*AtBats, error)
	CurrentAtBat() (*AtBat, error)
	FilterAtBats(func(*AtBat) bool) (*AtBats, error)
	Notifications() (*Notifications, error)
}

// GameServiceOp communicates with the MLB Gameday API to retrieve information
// for a single game.
type GameServiceOp struct {
	// The MLB Gameday API client used to communicate with the MLB Gameday API.
	client *Client

	// The unique ID associated with the game.
	gid string

	// The path to the MLB Gameday API at which resources for this game reside.
	path string
}

// NewGameService returns a new GameServiceOp that is linked to the
// provided game ID. Communication with the MLB Gameday API occurs through
// the provided Client.
func NewGameService(client *Client, gid string) (GameService, error) {
	path, err := pathFromGID(gid, "")
	if err != nil {
		return nil, err
	}

	return &GameServiceOp{client: client, gid: gid, path: path}, err
}

// LineScore retrieves the line score for this game.
func (s *GameServiceOp) LineScore() (*LineScore, error) {
	data, err := s.client.get(s.path + "linescore.xml")
	if err != nil {
		return nil, err
	}

	ls := new(LineScore)
	if err = xml.Unmarshal(data, &ls); err != nil {
		return nil, err
	}

	return ls, nil
}

// Players returns the team rosters for this game.
func (s *GameServiceOp) Players() (*Players, error) {
	data, err := s.client.get(s.path + "players.xml")
	if err != nil {
		return nil, err
	}

	p := new(Players)
	if err = xml.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	return p, nil
}

// AtBats returns all at-bats for this game, including any that are in
// progress.
func (s *GameServiceOp) AtBats() (*AtBats, error) {
	data, err := s.client.get(s.path + "inning/inning_all.xml")
	if err != nil {
		return nil, err
	}

	g := new(AtBats)
	if err = xml.Unmarshal(data, &g); err != nil {
		return nil, err
	}

	return g, nil
}

// CurrentAtBat returns the most current at-bat during a game in progress.
// Calling CurrentAtBat for a game that has ended will return the game's final
// at-bat.
func (s *GameServiceOp) CurrentAtBat() (*AtBat, error) {
	g, err := s.AtBats()
	if err != nil {
		return nil, err
	}

	i, top, err := lastHalfInning(g)
	if err != nil {
		return nil, err
	}
	inning := g.Innings[i]
	abs := inning.Bottom
	if top == "Y" {
		abs = inning.Top
	}
	ab := abs[len(abs)-1]

	return &ab, nil
}

// FilterAtBats returns all at-bats for which the evaluation in the provided
// function returns true.
func (s *GameServiceOp) FilterAtBats(f func(*AtBat) bool) (*AtBats, error) {
	all, err := s.AtBats()
	if err != nil {
		return nil, err
	}

	abs := AtBats{}
	for _, inning := range all.Innings {
		inn := AtBatInning{Number: inning.Number}
		for _, ab := range inning.Top {
			if f(&ab) {
				inn.Top = append(inn.Top, ab)
			}
		}
		for _, ab := range inning.Bottom {
			if f(&ab) {
				inn.Bottom = append(inn.Bottom, ab)
			}
		}
		if len(inn.Top) > 0 || len(inn.Bottom) > 0 {
			abs.Innings = append(abs.Innings, inn)
		}
	}

	return &abs, nil
}

// Notifications returns all notifications for this game: general game
// notifications as well as notifications by team/inning.
func (s *GameServiceOp) Notifications() (*Notifications, error) {
	data, err := s.client.get(s.path + "notifications/notifications_full.xml")
	if err != nil {
		return nil, err
	}

	n := new(Notifications)
	if err = xml.Unmarshal(data, &n); err != nil {
		return nil, err
	}

	return n, nil
}

// pathFromGID returns the URL path to a resource provided by the MLB
// Gameday API based on the provided gid.
func pathFromGID(gid string, path string) (string, error) {
	toks := strings.Split((string)(gid), "_")
	if len(toks) != 6 {
		return "", errors.New("Could not derive date from id " + gid)
	}

	gDate := fmt.Sprintf("year_%s/month_%s/day_%s", toks[0], toks[1], toks[2])
	return fmt.Sprintf("components/game/mlb/%s/gid_%s/%s", gDate, gid, path), nil
}

// lastHalfInning returns the most recent half inning, i.e. the inning number
// and whether it is the top or bottom of the inning.
func lastHalfInning(g *AtBats) (int, string, error) {
	n := len(g.Innings) - 1
	for n >= 0 {
		i := g.Innings[n]
		if len(i.Bottom) > 0 {
			return n, "N", nil
		} else if len(i.Top) > 0 {
			return n, "Y", nil
		}

		n--
	}

	return 0, "", errors.New("Game has not started")
}
