package mlbgameday

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func setupGameday() GamedayService {
	l, _ := time.LoadLocation("America/New_York")
	d := time.Date(2016, 9, 5, 0, 0, 0, 0, l)
	return client.Gameday(d)
}

func TestGame(t *testing.T) {
	gameday := setupGameday()

	_, err := gameday.Game("2016_09_05_kcamlb_minmlb_1")
	if nil != err {
		t.Fatalf("Gameday.Game returned error: %v", err)
	}
}

func TestGameErrorInvalidGID(t *testing.T) {
	gameday := setupGameday()

	gid := "invalid_gid"
	got, err := gameday.Game(gid)
	if got != nil {
		t.Errorf("Gameday.Game returned %v, want nil", got)
	}

	want := fmt.Sprintf("Could not derive date from id %v", gid)
	testError(t, "Gameday.Game", want, err)
}

func TestScoreboard(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/miniscoreboard.xml")
	if err != nil {
		t.Fatal("Could not read data file")
	}

	path := "/components/game/mlb/year_2016/month_09/day_05/miniscoreboard.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()

	got, err := gameday.Scoreboard()
	if err != nil {
		t.Fatalf("Gameday.Games returned error: %v", err)
	}

	want := &Scoreboard{
		[]Game{
			{
				"2016_09_05_tormlb_nyamlb_1",
				"2016/09/05 1:05", "ET", "PM",
				"TOR", "NYY", "Toronto", "NY Yankees", "Blue Jays", "Yankees",
				"Y", "In Progress", 1, 0,
				0, 0, 1, 0, 0, 0, 2,
			},
			{
				"2016_09_05_nynmlb_cinmlb_1",
				"2016/09/05 1:10", "ET", "PM",
				"NYM", "CIN", "NY Mets", "Cincinnati", "Mets", "Reds",
				"", "Preview", 0, 0,
				0, 0, 0, 0, 0, 0, 0,
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Gameday.Games returned %v, want %v", got, want)
	}
}

func TestScoreboardErrorDial(t *testing.T) {
	gameday := setupGameday()

	got, err := gameday.Scoreboard()
	if got != nil {
		t.Errorf("Gameday.Games returned %v, want nil", got)
	}

	if err == nil {
		t.Errorf("Gameday.Games did not return error")
	}
}

func TestScoreboardErrorEOF(t *testing.T) {
	setup()
	defer teardown()

	path := "/components/game/mlb/year_2016/month_09/day_05/miniscoreboard.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "")
	})

	gameday := setupGameday()

	got, err := gameday.Scoreboard()
	if got != nil {
		t.Errorf("Gameday.Games returned %v, want nil", got)
	}

	testError(t, "Gameday.Games", "EOF", err)
}

func TestGamesInProgress(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/miniscoreboard.xml")
	if err != nil {
		t.Fatal("Could not read data file")
	}

	// miniscoreboard.xml contains games from 2016-09-05 but mux serves a
	// URL that will match the current day.
	path := "/components/game/mlb/year_2016/month_09/day_05/miniscoreboard.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()

	got, err := gameday.GamesInProgress()
	if err != nil {
		t.Fatalf("Gameday.GamesInProgress returned error: %v", err)
	}

	want := []Game{
		{
			"2016_09_05_tormlb_nyamlb_1",
			"2016/09/05 1:05", "ET", "PM",
			"TOR", "NYY", "Toronto", "NY Yankees", "Blue Jays", "Yankees",
			"Y", "In Progress", 1, 0,
			0, 0, 1, 0, 0, 0, 2,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Gameday.GamesInProgress returned %v, want %v", got, want)
	}
}

func TestGamesInProgressErrorHTTP404(t *testing.T) {
	setup()
	defer teardown()

	path := fmt.Sprintf("/components/game/mlb/%v/miniscoreboard.xml",
		time.Now().Format("year_2006/month_01/day_02"))

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "Not Found", 404)
	})

	gameday := setupGameday()

	got, err := gameday.GamesInProgress()
	if got != nil {
		t.Errorf("Gameday.GamesInProgress returned %v, want nil", got)
	}

	testError(t, "Gameday.GamesInProgress", "HTTP 404", err)
}
