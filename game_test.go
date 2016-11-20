package mlbgameday

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestLineScore(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/linescore.xml")
	if err != nil {
		t.Fatal("Could not read data file")
	}

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/linescore.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.LineScore()
	if err != nil {
		t.Fatalf("Game.LineScore returned error: %v", err)
	}

	want := &LineScore{
		Game{
			gid, "2016/09/05 2:10", "ET", "PM",
			"KC", "MIN", "Kansas City", "Minnesota", "Royals", "Twins",
			"Y", "In Progress", 7, 0,
			5, 4, 10, 9, 0, 1, 0,
		},
		"N", "N",
		Review{0, 1, 0, 1},
		[]LineScoreInning{
			{1, 1, 0},
			{2, 0, 2},
			{3, 2, 0},
			{4, 1, 0},
			{5, 0, 3},
			{6, 0, 0},
			{7, 0, 0},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Game.LineScore returned %v, want %v", got, want)
	}
}

func TestLineScoreErrorHTTP404(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/linescore.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "Not Found", 404)
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.LineScore()
	if got != nil {
		t.Errorf("Game.LineScore returned %v, want nil", got)
	}

	testError(t, "Game.LineScore", "HTTP 404", err)
}

func TestLineScoreErrorEOF(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/linescore.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "")
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.LineScore()
	if got != nil {
		t.Errorf("Game.LineScore returned %v, want nil", got)
	}

	testError(t, "Game.LineScore", "EOF", err)
}

func TestPlayers(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/players.xml")
	if err != nil {
		t.Fatalf("Could not read data file")
	}

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/players.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.Players()
	if err != nil {
		t.Fatalf("Game.Players returned error: %v", err)
	}

	want := &Players{
		[]Roster{
			{
				"KC",
				[]Player{
					{
						521692, "Salvador", "Perez", 13, "R", "R", "C", "A",
						0.254, 20, 58, 0, 0, 0,
					},
					{
						572044, "Brooks", "Pounders", 62, "R", "R", "P", "A",
						0.000, 0, 0, 1, 1, 10.29,
					},
				},
				[]Coach{{124681, "Ned", "Yost", 3, "manager"}},
			},
			{
				"MIN",
				[]Player{
					{
						542953, "Buddy", "Boshers", 62, "L", "L", "P", "A",
						0.000, 0, 0, 2, 0, 5.33,
					},
					{
						621439, "Byron", "Buxton", 25, "R", "R", "CF", "A",
						0.221, 4, 25, 0, 0, 0,
					},
				},
				[]Coach{{119236, "Paul", "Molitor", 4, "manager"}},
			},
		},
		[]Umpire{{427019, "Ted", "Barrett", "home", "Ted Barrett"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Game.Players returned %v, want %v", got, want)
	}
}

func TestPlayersErrorHTTP404(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/players.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "Not Found", 404)
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.Players()
	if got != nil {
		t.Errorf("Game.Players returned %v, want nil", got)
	}

	testError(t, "Game.Players", "HTTP 404", err)
}

func TestPlayersErrorEOF(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/players.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "")
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.Players()
	if got != nil {
		t.Errorf("Game.Players returned %v, want nil", got)
	}

	testError(t, "Game.Players", "EOF", err)
}

func TestCurrentAtBat(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/inning_all.xml")
	if err != nil {
		t.Fatal("Could not read data file")
	}

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/inning/inning_all.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.CurrentAtBat()
	if err != nil {
		t.Fatalf("Game.CurrentAtBat returned error: %v", err)
	}

	want := &AtBat{
		AtBatSummary{88, 500871, 518397, 1, 1, 3, "Grounded Into DP", 5, 11},
		[]Pitch{
			{
				"Ball",
				"B", "2016-09-05T21:34:13Z", 134.08, 199.09, 91.5, 85.8, 3.1,
				1.4, 8.41, 4.27, -0.448, 1.47, 1.261, 50.0, 5.548, -7.469,
				-133.898, -6.314, 15.536, 24.104, -24.206, 23.9, -29.4, 6.3,
				"SI", 0.887, 7, 62, 117.153, 1891.802,
			},
			{
				"Foul",
				"S", "2016-09-05T21:34:30Z", 100.57, 194.85, 92.8, 86.5, 3.1,
				1.4, 9.53, 1.2, 0.431, 1.627, 1.411, 50.0, 5.55, -5.981,
				-135.904, -5.034, 17.999, 27.035, -29.826, 23.9, -28.7, 7.5,
				"SI", 0.896, 9, 51, 97.433, 1937.682,
			},
			{
				"In play, out(s)",
				"X", "2016-09-05T21:35:11Z", 127.52, 179.03, 93.7, 88.0, 3.1,
				1.4, 9.37, 2.62, -0.276, 2.213, 1.298, 50.0, 5.639, -7.636,
				-137.19, -4.417, 18.233, 24.366, -27.009, 23.9, -31.4, 6.7,
				"SI", 0.902, 4, 34, 105.817, 2004.007,
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Game.CurrentAtBat returned %v, want %v", got, want)
	}
}

func TestCurrentAtBatTopInning(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/inning_top.xml")
	if err != nil {
		t.Fatal("Could not read data file")
	}

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/inning/inning_all.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.CurrentAtBat()
	if err != nil {
		t.Fatalf("Game.CurrentAtBat returned error: %v", err)
	}

	want := &AtBat{
		AtBatSummary{1, 502481, 621244, 0, 1, 0, "Single", 0, 0},
		[]Pitch{
			{
				"Called Strike",
				"S", "2016-09-05T18:10:36Z", 151.27, 180.81, 92.5, 85.4, 3.47,
				1.62, -0.98, 9.59, -0.899, 2.147, -1.76, 50.0, 5.492, 2.649,
				-135.539, -6.307, -1.831, 28.902, -14.218, 23.8, 4.7, 3.4,
				"FF", 0.914, 13, 72, 185.821, 1931.941,
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Game.CurrentAtBat returned %v, want %v", got, want)
	}
}

func TestCurrentAtBatErrorHTTP404(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/inning/inning_all.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "Not Found", 404)
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.CurrentAtBat()
	if got != nil {
		t.Errorf("Game.CurrentAtBat returned %v, want nil", got)
	}

	testError(t, "Game.CurrentAtBat", "HTTP 404", err)
}

func TestCurrentAtBatErrorEOF(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/inning/inning_all.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "")
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.CurrentAtBat()
	if got != nil {
		t.Errorf("Game.CurrentAtBat returned %v, want nil", got)
	}

	testError(t, "Game.CurrentAtBat", "EOF", err)
}

func TestCurrentAtBatErrorGameNotStarted(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/inning/inning_all.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		xml := "<game>" +
			"<inning num=\"1\" away_team=\"kca\" home_team=\"min\">" +
			"<top/><bottom/>" +
			"</inning>" +
			"</game>"
		fmt.Fprint(w, xml)
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.CurrentAtBat()
	if got != nil {
		t.Errorf("Game.CurrentAtBat returned %v, want nil", got)
	}

	testError(t, "Game.CurrentAtBat", "Game has not started", err)
}

func TestFilterAtBats(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/inning_all.xml")
	if err != nil {
		t.Fatal("Could not read data file")
	}

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/inning/inning_all.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("ByGID(%q) returned error: %q", gid, err)
	}

	var testCases = []struct {
		Batter  int
		Top     bool
		Innings [][]int
	}{
		{435559, false, [][]int{{7, 1}}},
		{502481, true, [][]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {8, 1}}},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			abs, err := game.FilterAtBats(func(ab *AtBat) bool {
				return ab.Batter == tc.Batter
			})
			if err != nil {
				t.Fatalf("Game.FilterAtBats returned error: %q", err)
			}

			var got [][]int
			for _, i := range abs.Innings {
				a := i.Bottom
				if tc.Top == true {
					a = i.Top
				}
				got = append(got, []int{i.Number, len(a)})
			}

			if !reflect.DeepEqual(got, tc.Innings) {
				t.Errorf("Got [[Inning, NumAtBats]] %v, want %v",
					got, tc.Innings)
			}
		})
	}
}

func TestFilterAtBatsErrorHTTP404(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/inning/inning_all.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "Not Found", 404)
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.FilterAtBats(func(ab *AtBat) bool {
		return true
	})
	if got != nil {
		t.Errorf("Game.FilterAtBats returned %v, want nil", got)
	}

	testError(t, "Game.FilterAtBats", "HTTP 404", err)
}

func TestNotifications(t *testing.T) {
	setup()
	defer teardown()

	data, err := ioutil.ReadFile("./mock/notifications_full.xml")
	if err != nil {
		t.Fatal("Could not read data file")
	}

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/notifications/notifications_full.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, string(data))
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.Notifications()
	if err != nil {
		t.Fatalf("Game.Notifications returned error: %v", err)
	}

	want := &Notifications{
		[]Notification{
			{
				Inning: 1,
				Top:    "Y",
				AtBat:  1,
				Outs:   0,
				Types:  []Type{{"lineups"}},
			},
		},
		[]Team{
			{
				118, "kca", []Notification{
					{
						7, "N", 69, 2,
						[]Player{{ID: 572044}, {ID: 543169}},
						[]Type{{"pitching change"}},
					},
				},
			},
			{
				142, "min", []Notification{
					{
						7, "N", 69, 2,
						[]Player{{ID: 435559}, {ID: 518542}},
						[]Type{{"pinch hitter"}},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Game.Notifications returned %v, want %v", got, want)
	}
}

func TestNotificationsErrorHTTP404(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/notifications/notifications_full.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "Not Found", 404)
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.Notifications()
	if got != nil {
		t.Errorf("Game.Notifications returned %v, want nil", got)
	}

	testError(t, "Game.Notifications", "HTTP 404", err)
}

func TestNotificationsErrorEOF(t *testing.T) {
	setup()
	defer teardown()

	gid := "2016_09_05_kcamlb_minmlb_1"
	path := "/components/game/mlb/year_2016/month_09/day_05/" +
		"gid_" + gid + "/notifications/notifications_full.xml"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "")
	})

	gameday := setupGameday()
	game, err := gameday.Game(gid)
	if err != nil {
		t.Fatalf("Games.ByGID returned error: %v", err)
	}

	got, err := game.Notifications()
	if got != nil {
		t.Errorf("Game.Notifications returned %v, want nil", got)
	}

	testError(t, "Game.Notifications", "EOF", err)
}
