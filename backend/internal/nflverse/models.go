package nflverse

type Play struct {
	PlayID                int      `parquet:"play_id"`
	GameID                string   `parquet:"game_id"`
	Quarter               int      `parquet:"qtr"`
	Down                  int      `parquet:"down"`
	YardsToGo             int      `parquet:"ydstogo"`
	PlayType              string   `parquet:"play_type"`
	Description           string   `parquet:"desc"`
	EPA                   *float64 `parquet:"epa"` // nil when nflverse has no value
	PasserPlayerName      string   `parquet:"passer_player_name"`
	ReceiverPlayerName    string   `parquet:"receiver_player_name"`
	RusherPlayerName      string   `parquet:"rusher_player_name"`
	TacklerPlayerName     string   `parquet:"tackler_player_name"`
	Posteam               string   `parquet:"posteam"`
	Defteam               string   `parquet:"defteam"`
	Yardline100           int      `parquet:"yardline_100"`
	YardsGained           int      `parquet:"yards_gained"`
	PassLocation          string   `parquet:"pass_location"`
	RunLocation           string   `parquet:"run_location"`
	RunGap                string   `parquet:"run_gap"`
	PassAttempt           int      `parquet:"pass_attempt"`
	RushAttempt           int      `parquet:"rush_attempt"`
	CompletePass          int      `parquet:"complete_pass"`
	Interception          int      `parquet:"interception"`
	Touchdown             int      `parquet:"touchdown"`
	FumbleLost            int      `parquet:"fumble_lost"`
	Sack                  int      `parquet:"sack"`
	QBHit                 int      `parquet:"qb_hit"`
	QBPressure            int      `parquet:"qb_pressure"`
	RusherPlayerId        string   `parquet:"rusher_player_id"`
	PasserPlayerId        string   `parquet:"passer_player_id"`
	ReceiverPlayerId      string   `parquet:"receiver_player_id"`
	Season                int      `parquet:"season"`
	Week                  int      `parquet:"week"`
	HomeTeam              string   `parquet:"home_team"`
	AwayTeam              string   `parquet:"away_team"`
	HomeScore             int      `parquet:"home_score"`
	AwayScore             int      `parquet:"away_score"`
	GameSecondsRemaining  int      `parquet:"game_seconds_remaining"`
	Drive                 int      `parquet:"drive"`
	Down12                int      `parquet:"down12"`
	Down34                int      `parquet:"down34"`
	GoalToGo              int      `parquet:"goal_to_go"`
	Time                  string   `parquet:"time"`
	SideOfField           string   `parquet:"side_of_field"`
	ThirdDownConverted    int      `parquet:"third_down_converted"`
	FourthDownConverted   int      `parquet:"fourth_down_converted"`
	FixedDrive            int      `parquet:"fixed_drive"`
	FixedDriveResult      string   `parquet:"fixed_drive_result"`
	DriveQuarterStart     int      `parquet:"drive_quarter_start"`
	DriveQuarterEnd       int      `parquet:"drive_quarter_end"`
	DriveTimeOfPossession string   `parquet:"drive_time_of_possession"` // "m:ss"
	FieldGoalResult       string   `parquet:"field_goal_result"`
	ExtraPointResult      string   `parquet:"extra_point_result"`
	TwoPointConvResult    string   `parquet:"two_point_conv_result"`
	TDTeam                string   `parquet:"td_team"`
}
