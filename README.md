# mlbgameday

mlbgameday is a Go client for the Major League Baseball (MLB) GameDay API.

The MLB Gameday API provides data for Major League Baseball games in several
different formats: e.g., XML, JSON, and p-list. The raw game data can be viewed
at <http://gd2.mlb.com/components/game/mlb/>.

## MLBAM Data Copyright

MLB Advanced Media, L.P. ("MLBAM") provides a copyright for the MLB Gameday API
data. It can be read in full [here](http://gdx.mlb.com/components/copyright.txt).

By using the mlbclient package to access any such data you accept the terms of
the copyright.

## Features

Access MLB Gameday data for Major League Baseball games, past and present:

### By Game Day:
 * Scoreboard
 * Games In Progress

### By Game:
 * Line Score
 * At-bats: outcome, PitchF/X data, etc.
 * Filter at-bats by custom criteria.
 * Rosters and umpires
 * Notifications

## API

The mlbgameday client mostly follows the structure and XML schemas of the MLB
Gameday API.

	MLB Gameday API:
	---------------
	components/game/mlb/year_YYYY/month_MM/day_DD/
	|__ miniscoreboard.xml
	|__ gid_YYYY_MM_DD_ATMmlb_HTMmlb_1/
	    |__ linescore.xml
	    |__ players.xml
		|__ inning/
		|__ notifications/

	mlbgameday:
	----------
	GamedayService(year_YYYY/month_MM/day_DD)
	|__ Scoreboard()
	|__ GameService(gid_YYYY_MM_DD_ATMmlb_HTMmlb_1)
	    |__ LineScore()
	    |__ Players()
		|__ AtBats()
		|__ Notifications()

	YYYY: Year
	  MM: Month
	  DD: Day
	 ATM: AwayTeam
	 HTM: HomeTeam

Please note that the mlbgameday API is subject to change.

## Usage

```go
import "github.com/ericdreeves/mlbgameday"
```

After creating a new MLB Gameday client, use the Gameday method to access a
particular game day.

The Gameday type provides the scoreboard for a chosen day. Use the Game
method to access a game by its Game ID.

## Examples

### Create a new MLB Gameday client:

```go
client := mlbgameday.NewClient(nil)
```

### Access data for a game day:

```go
loc, _ := time.LoadLocation("America/New_York")
date := time.Date(2016, 9, 5, 0, 0, 0, 0, loc)
gameday := client.Gameday(date)

scoreboard, err := gameday.Scoreboard()
if err != nil {
	panic(err)
}

games := scoreboard.Games
```

### Access data for a specific game:

```go
game, err := gameday.Game(games[0].GID)
if nil != err {
	panic(err)
}

lineScore, err := game.LineScore()
if nil != err {
	panic(err)
}

fmt.Printf("Inning: %v\nOuts: %v\nBaseState: %v\n",
	lineScore.Inning, lineScore.Outs, lineScore.BaseState)
for _, inn := range lineScore.Innings {
	fmt.Printf("  Inning: %v | Away %v | Home %v\n",
		inn.Inning, inn.AwayRuns, inn.HomeRuns)
}
```

## Documentation

The MLB Gameday API data can be viewed at <http://gd2.mlb.com/components/game/mlb/>.

The [Unofficial MLB Gameday API Reference](https://github.com/brianmpalma/gameday-api-docs)
was last updated in 2014 however it is still a fantastic reference.

[FanGraphs](http://www.fangraphs.com/library/misc/pitch-fx/) provides a primer
on PITCHF/x.

## License

The mlbgameday project is released under the terms of the 
[MIT License](https://en.wikipedia.org/wiki/MIT_License).
