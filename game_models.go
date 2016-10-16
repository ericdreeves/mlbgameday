package mlbgameday

// Game represents an MLB game.
type Game struct {
	// GID is a unique ID for each game.
	GID            string `xml:"gameday_link,attr"`
	TimeDate       string `xml:"time_date,attr"`
	TimeZone       string `xml:"time_zone,attr"`
	AMPM           string `xml:"ampm,attr"`
	AwayNameAbbrev string `xml:"away_name_abbrev,attr"`
	HomeNameAbbrev string `xml:"home_name_abbrev,attr"`
	AwayTeamCity   string `xml:"away_team_city,attr"`
	HomeTeamCity   string `xml:"home_team_city,attr"`
	AwayTeamName   string `xml:"away_team_name,attr"`
	HomeTeamName   string `xml:"home_team_name,attr"`
	TopInning      string `xml:"top_inning,attr"`
	Status         string `xml:"status,attr"`
	Inning         int    `xml:"inning,attr"`
	Outs           int    `xml:"outs,attr"`
	AwayTeamRuns   int    `xml:"away_team_runs,attr"`
	HomeTeamRuns   int    `xml:"home_team_runs,attr"`
	AwayHitsRuns   int    `xml:"away_team_hits,attr"`
	HomeHitsRuns   int    `xml:"home_team_hits,attr"`
	AwayTeamErrors int    `xml:"away_team_errors,attr"`
	HomeTeamErrors int    `xml:"home_team_errors,attr"`

	// BaseState indicates which bases have runners.
	// 0: Empty, 1: 1B, 2: 2B, 3: 3B
	// 4: 1B,2B, 5: 1B,3B, 6: 2B,3B
	// 7: 1B,2B,3B
	BaseState int `xml:"runner_on_base_status,attr"`
}

// LineScore represents the line score of a game, including runs by inning
// and review information.
type LineScore struct {
	Game
	NoHitter    string            `xml:"is_no_hitter,attr"`
	PerfectGame string            `xml:"is_perfect_game,attr"`
	Review      Review            `xml:"review"`
	Innings     []LineScoreInning `xml:"linescore"`
}

// Review represents data about reviews used during a game as part of the
// line score.
type Review struct {
	ChallengesAwayUsed      int `xml:"challenges_away_used,attr"`
	ChallengesAwayRemaining int `xml:"challenges_away_remaining,attr"`
	ChallengesHomeUsed      int `xml:"challenges_home_used,attr"`
	ChallengesHomeRemaining int `xml:"challenges_home_remaining,attr"`
}

// LineScoreInning represents runs per inning as part of the line score.
type LineScoreInning struct {
	Inning   int `xml:"inning,attr"`
	HomeRuns int `xml:"home_inning_runs,attr,omitempty"`
	AwayRuns int `xml:"away_inning_runs,attr,omitempty"`
}

// Players represents a list of team rosters and umpires.
type Players struct {
	Teams   []Roster `xml:"team"`
	Umpires []Umpire `xml:"umpires>umpire"`
}

// Roster represents the list of players on a team.
type Roster struct {
	TeamID  string   `xml:"id,attr"`
	Players []Player `xml:"player"`
	Coaches []Coach  `xml:"coach"`
}

// Player represents a baseball player and some basic stats.
type Player struct {
	ID       int     `xml:"id,attr"`
	First    string  `xml:"first,attr"`
	Last     string  `xml:"last,attr"`
	Num      int     `xml:"num,attr"`
	RL       string  `xml:"rl,attr"`
	Bats     string  `xml:"bats,attr"`
	Position string  `xml:"position,attr"`
	Status   string  `xml:"status,attr"`
	Avg      float32 `xml:"avg,attr"`
	HR       int     `xml:"hr,attr"`
	RBI      int     `xml:"rbi,attr"`
	Wins     int     `xml:"wins,attr"`
	Losses   int     `xml:"losses,attr"`
	ERA      float32 `xml:"era,attr"`
}

// Coach represents a member of a team's coaching staff.
type Coach struct {
	ID       int    `xml:"id,attr"`
	First    string `xml:"first,attr"`
	Last     string `xml:"last,attr"`
	Num      int    `xml:"num,attr"`
	Position string `xml:"position,attr"`
}

// Umpire represents a member of the game's umpire staff.
type Umpire struct {
	ID       int    `xml:"id,attr"`
	First    string `xml:"first,attr"`
	Last     string `xml:"last,attr"`
	Position string `xml:"position,attr"`
	Name     string `xml:"name,attr"`
}

// AtBats represents a list of at-bats by inning.
type AtBats struct {
	Innings []AtBatInning `xml:"inning"`
}

// AtBatInning represents the at-bats that occur in both halves of an inning.
type AtBatInning struct {
	Number int     `xml:"num,attr"`
	Top    []AtBat `xml:"top>atbat"`
	Bottom []AtBat `xml:"bottom>atbat"`
}

// AtBat represents a single at-bat.
type AtBat struct {
	Number       int     `xml:"num,attr"`
	Batter       int     `xml:"batter,attr"`
	Pitcher      int     `xml:"pitcher,attr"`
	Balls        int     `xml:"b,attr"`
	Strikes      int     `xml:"s,attr"`
	Outs         int     `xml:"o,attr"`
	Event        string  `xml:"event,attr"`
	HomeTeamRuns int     `xml:"home_team_runs,attr"`
	AwayTeamRuns int     `xml:"away_team_runs,attr"`
	Pitches      []Pitch `xml:"pitch"`
}

// Pitch represents the PitchF/X data for a single pitch.
type Pitch struct {
	Des            string  `xml:"des,attr"`
	Type           string  `xml:"type,attr"`
	TFSZulu        string  `xml:"tfs_zulu,attr"`
	X              float32 `xml:"x,attr"`
	Y              float32 `xml:"y,attr"`
	StartSpeed     float32 `xml:"start_speed,attr"`
	EndSpeed       float32 `xml:"end_speed,attr"`
	SZTop          float32 `xml:"sz_top,attr"`
	SZBot          float32 `xml:"sz_bot,attr"`
	PfxX           float32 `xml:"pfx_x,attr"`
	PfxZ           float32 `xml:"pfx_z,attr"`
	PX             float32 `xml:"px,attr"`
	PZ             float32 `xml:"pz,attr"`
	X0             float32 `xml:"x0,attr"`
	Y0             float32 `xml:"y0,attr"`
	Z0             float32 `xml:"z0,attr"`
	VX0            float32 `xml:"vx0,attr"`
	VY0            float32 `xml:"vy0,attr"`
	VZ0            float32 `xml:"vz0,attr"`
	AX             float32 `xml:"ax,attr"`
	AY             float32 `xml:"ay,attr"`
	AZ             float32 `xml:"az,attr"`
	BreakY         float32 `xml:"break_y,attr"`
	BreakAngle     float32 `xml:"break_angle,attr"`
	BreakLength    float32 `xml:"break_length,attr"`
	PitchType      string  `xml:"pitch_type,attr"`
	TypeConfidence float32 `xml:"type_confidence,attr"`
	Zone           float32 `xml:"zone,attr"`
	Nasty          float32 `xml:"nasty,attr"`
	SpinDir        float32 `xml:"spin_dir,attr"`
	SpinRate       float32 `xml:"spin_rate,attr"`
}

// Notifications represents both game and team notifications.
type Notifications struct {
	Game  []Notification `xml:"game>notification"`
	Teams []Team         `xml:"team"`
}

// Team represents all the notifications for a team.
type Team struct {
	ID            int            `xml:"id,attr"`
	Code          string         `xml:"code,attr"`
	Notifications []Notification `xml:"notification"`
}

// Notification represents a game notification from the MLB Gameday API.
type Notification struct {
	Inning  int      `xml:"inning,attr"`
	Top     string   `xml:"top,attr"`
	AtBat   int      `xml:"ab,attr"`
	Outs    int      `xml:"outs,attr"`
	Players []Player `xml:"player"`
	Types   []Type   `xml:"type"`
}

// Type represents the notification type.
type Type struct {
	Category string `xml:"category,attr"`
}
