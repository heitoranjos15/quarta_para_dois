package nflverse

type Play struct {
	PlayID            int     `parquet:"play_id"`
	GameID            string  `parquet:"game_id"`
	Quarter           int     `parquet:"qtr"`
	Down              int     `parquet:"down"`
	YardsToGo         int     `parquet:"ydstogo"`
	PlayType          string  `parquet:"play_type"`
	Description       string  `parquet:"desc"`
	EPA               float64 `parquet:"epa"`
	PasserPlayerName  string  `parquet:"passer_player_name"`
	ReceiverPlayerName string  `parquet:"receiver_player_name"`
	RusherPlayerName  string  `parquet:"rusher_player_name"`
	TacklerPlayerName string  `parquet:"tackler_player_name"`
	Posteam           string  `parquet:"posteam"`
	Defteam           string  `parquet:"defteam"`
	Yardline100       int     `parquet:"yardline_100"`
	YardsGained       int     `parquet:"yards_gained"`
	PassLocation      string  `parquet:"pass_location"`
	RunLocation       string  `parquet:"run_location"`
	RunGap            string  `parquet:"run_gap"`
	PassAttempt       int     `parquet:"pass_attempt"`
	RushAttempt       int     `parquet:"rush_attempt"`
	CompletePass      int     `parquet:"complete_pass"`
	Interception      int     `parquet:"interception"`
	Touchdown         int     `parquet:"touchdown"`
	FumbleLost        int     `parquet:"fumble_lost"`
	Sack              int     `parquet:"sack"`
	QBHit             int     `parquet:"qb_hit"`
	QBPressure        int     `parquet:"qb_pressure"`
	RusherPlayerId    string  `parquet:"rusher_player_id"`
	PasserPlayerId    string  `parquet:"passer_player_id"`
	ReceiverPlayerId  string  `parquet:"receiver_player_id"`
	Season            int     `parquet:"season"`
	Week              int     `parquet:"week"`
	HomeTeam          string  `parquet:"home_team"`
	AwayTeam          string  `parquet:"away_team"`
	HomeScore         int     `parquet:"home_score"`
	AwayScore         int     `parquet:"away_score"`
	GameSecondsRemaining int  `parquet:"game_seconds_remaining"`
	Drive             int     `parquet:"drive"`
	Down12            int     `parquet:"down12"`
	Down34            int     `parquet:"down34"`
	GoalToGo          int     `parquet:"goal_to_go"`
	Time              string  `parquet:"time"`
	SideOfField       string  `parquet:"side_of_field"`
	YardsToGoal       int     `parquet:"ydstogo"`
	ThirdDownConverted int    `parquet:"third_down_converted"`
	FourthDownConverted int   `parquet:"fourth_down_converted"`
}

type Game struct {
	GameID       string `parquet:"game_id"`
	Season       int    `parquet:"season"`
	Week         int    `parquet:"week"`
	GameType     string `parquet:"game_type"`
	GameDate     string `parquet:"gameday"`
	Weekday      string `parquet:"weekday"`
	GameTime     string `parquet:"gametime"`
	HomeTeam     string `parquet:"home_team"`
	AwayTeam     string `parquet:"away_team"`
	HomeScore    int    `parquet:"home_score"`
	AwayScore    int    `parquet:"away_score"`
	Location     string `parquet:"location"`
	Stadium      string `parquet:"stadium"`
	Surface      string `parquet:"surface"`
	Roof         string `parquet:"roof"`
	Temp         int    `parquet:"temp"`
	Wind         int    `parquet:"wind"`
	HomeCoach    string `parquet:"home_coach"`
	AwayCoach    string `parquet:"away_coach"`
	Referee      string `parquet:"referee"`
	StadiumId    string `parquet:"stadium_id"`
}

type Team struct {
	TeamAbbr    string `parquet:"team_abbr"`
	TeamName    string `parquet:"team_name"`
	TeamId      string `parquet:"team_id"`
	TeamNick    string `parquet:"team_nick"`
	TeamConf    string `parquet:"team_conf"`
	TeamDiv     string `parquet:"team_div"`
	TeamColor   string `parquet:"team_color"`
	TeamColor2  string `parquet:"team_color2"`
	TeamColor3  string `parquet:"team_color3"`
	TeamColor4  string `parquet:"team_color4"`
	TeamLogo    string `parquet:"team_logo"`
	TeamLogoEspn string `parquet:"team_logo_espn"`
	TeamWordmark string `parquet:"team_wordmark"`
}

type Roster struct {
	Season      int    `parquet:"season"`
	Week        int    `parquet:"week"`
	Team        string `parquet:"team"`
	PlayerId    string `parquet:"player_id"`
	PlayerName  string `parquet:"player_name"`
	Position    string `parquet:"position"`
	DepthTeam   string `parquet:"depth_team"`
	DepthPos    string `parquet:"depth_pos"`
	Status      string `parquet:"status"`
	GsisId      string `parquet:"gsis_id"`
	EspnId      string `parquet:"espn_id"`
	SportradarId string `parquet:"sportradar_id"`
	YahooId     string `parquet:"yahoo_id"`
	RotowireId  string `parquet:"rotowire_id"`
	PffId       string `parquet:"pff_id"`
	PfrId       string `parquet:"pfr_id"`
	FantasyDataId string `parquet:"fantasy_data_id"`
	StatusDescriptionAbbr string `parquet:"status_description_abbr"`
	StatusShortDescription string `parquet:"status_short_description"`
	HeadshotUrl string `parquet:"headshot_url"`
	NgSId       string `parquet:"ngs_id"`
	College     string `parquet:"college"`
	CollegeConf string `parquet:"college_conf"`
	DraftYear   int    `parquet:"draft_year"`
	DraftRound  int    `parquet:"draft_round"`
	DraftPick   int    `parquet:"draft_pick"`
	Height      string `parquet:"height"`
	Weight      int    `parquet:"weight"`
	BirthDate   string `parquet:"birth_date"`
	Age         int    `parquet:"age"`
	Experience  int    `parquet:"experience"`
}