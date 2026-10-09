package nflverse

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/apache/arrow/go/v16/parquet/file"
	"github.com/apache/arrow/go/v16/parquet/pqarrow"
	"golang.org/x/sync/singleflight"
)

const (
	pbpBaseURL      = "https://github.com/nflverse/nflverse-data/releases/download/pbp/play_by_play_%d.parquet"
	schedulesBaseURL = "https://github.com/nflverse/nflverse-data/releases/download/schedules/schedules_%d.parquet"
	teamsBaseURL    = "https://github.com/nflverse/nflverse-data/releases/download/team/team_desc.parquet"
	rostersBaseURL  = "https://github.com/nflverse/nflverse-data/releases/download/rosters/rosters_%d.parquet"
	
	cacheTTL = time.Hour
)

type Client struct {
	cache       Cache
	httpClient  *http.Client
	githubToken string
	sfGroup     singleflight.Group
}

func NewClient(cache Cache, githubToken string) *Client {
	return &Client{
		cache:       cache,
		githubToken: githubToken,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *Client) GetPBP(season int) ([]Play, error) {
	cacheKey := fmt.Sprintf("pbp:%d", season)
	var plays []Play
	if c.cache.Get(cacheKey, &plays) {
		return plays, nil
	}

	result, err, _ := c.sfGroup.Do(cacheKey, func() (any, error) {
		plays, err := c.downloadAndParsePBP(season)
		if err != nil {
			return nil, err
		}
		c.cache.Set(cacheKey, plays, cacheTTL)
		return plays, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]Play), nil
}

func (c *Client) GetSchedules(season int) ([]Game, error) {
	cacheKey := fmt.Sprintf("schedules:%d", season)
	var games []Game
	if c.cache.Get(cacheKey, &games) {
		return games, nil
	}

	result, err, _ := c.sfGroup.Do(cacheKey, func() (any, error) {
		games, err := c.downloadAndParseSchedules(season)
		if err != nil {
			return nil, err
		}
		c.cache.Set(cacheKey, games, cacheTTL)
		return games, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]Game), nil
}

func (c *Client) GetTeams() ([]Team, error) {
	cacheKey := "teams"
	var teams []Team
	if c.cache.Get(cacheKey, &teams) {
		return teams, nil
	}

	result, err, _ := c.sfGroup.Do(cacheKey, func() (any, error) {
		teams, err := c.downloadAndParseTeams()
		if err != nil {
			return nil, err
		}
		c.cache.Set(cacheKey, teams, cacheTTL)
		return teams, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]Team), nil
}

func (c *Client) GetRosters(season int) ([]Roster, error) {
	cacheKey := fmt.Sprintf("rosters:%d", season)
	var rosters []Roster
	if c.cache.Get(cacheKey, &rosters) {
		return rosters, nil
	}

	result, err, _ := c.sfGroup.Do(cacheKey, func() (any, error) {
		rosters, err := c.downloadAndParseRosters(season)
		if err != nil {
			return nil, err
		}
		c.cache.Set(cacheKey, rosters, cacheTTL)
		return rosters, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]Roster), nil
}

func (c *Client) downloadAndParsePBP(season int) ([]Play, error) {
	url := fmt.Sprintf(pbpBaseURL, season)
	return c.downloadAndParsePBPURL(url)
}

func (c *Client) downloadAndParseSchedules(season int) ([]Game, error) {
	url := fmt.Sprintf(schedulesBaseURL, season)
	return c.downloadAndParseSchedulesURL(url)
}

func (c *Client) downloadAndParseTeams() ([]Team, error) {
	url := teamsBaseURL
	return c.downloadAndParseTeamsURL(url)
}

func (c *Client) downloadAndParseRosters(season int) ([]Roster, error) {
	url := fmt.Sprintf(rostersBaseURL, season)
	return c.downloadAndParseRostersURL(url)
}

func (c *Client) downloadAndParsePBPURL(url string) ([]Play, error) {
	result, err := c.downloadAndParse(url, func(r *file.Reader) (any, error) {
		return c.parsePBP(r)
	})
	if err != nil {
		return nil, err
	}
	return result.([]Play), nil
}

func (c *Client) downloadAndParseSchedulesURL(url string) ([]Game, error) {
	result, err := c.downloadAndParse(url, func(r *file.Reader) (any, error) {
		return c.parseSchedules(r)
	})
	if err != nil {
		return nil, err
	}
	return result.([]Game), nil
}

func (c *Client) downloadAndParseTeamsURL(url string) ([]Team, error) {
	result, err := c.downloadAndParse(url, func(r *file.Reader) (any, error) {
		return c.parseTeams(r)
	})
	if err != nil {
		return nil, err
	}
	return result.([]Team), nil
}

func (c *Client) downloadAndParseRostersURL(url string) ([]Roster, error) {
	result, err := c.downloadAndParse(url, func(r *file.Reader) (any, error) {
		return c.parseRosters(r)
	})
	if err != nil {
		return nil, err
	}
	return result.([]Roster), nil
}

func (c *Client) downloadAndParse(url string, parseFn func(*file.Reader) (any, error)) (any, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.githubToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.githubToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited by GitHub")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pqReader, err := file.NewParquetReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer pqReader.Close()

	return parseFn(pqReader)
}

func (c *Client) parsePBP(r *file.Reader) ([]Play, error) {
	props := pqarrow.ArrowReadProperties{
		BatchSize: 8192,
	}
	
	arrowReader, err := pqarrow.NewFileReader(r, props, memory.DefaultAllocator)
	if err != nil {
		return nil, err
	}

	var plays []Play
	
	recordReader, err := arrowReader.GetRecordReader(context.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer recordReader.Release()

	for recordReader.Next() {
		record := recordReader.Record()
		
		numRecords := record.NumRows()
		for i := 0; i < int(numRecords); i++ {
			play := Play{
				PlayID:             getInt32(record, "play_id", i),
				GameID:             getString(record, "game_id", i),
				Quarter:            getInt32(record, "qtr", i),
				Down:               getInt32(record, "down", i),
				YardsToGo:          getInt32(record, "ydstogo", i),
				PlayType:           getString(record, "play_type", i),
				Description:        getString(record, "desc", i),
				EPA:                getFloat64(record, "epa", i),
				PasserPlayerName:   getString(record, "passer_player_name", i),
				ReceiverPlayerName: getString(record, "receiver_player_name", i),
				RusherPlayerName:   getString(record, "rusher_player_name", i),
				TacklerPlayerName:  getString(record, "tackler_player_name", i),
				Posteam:            getString(record, "posteam", i),
				Defteam:            getString(record, "defteam", i),
				Yardline100:        getInt32(record, "yardline_100", i),
				YardsGained:        getInt32(record, "yards_gained", i),
				PassLocation:       getString(record, "pass_location", i),
				RunLocation:        getString(record, "run_location", i),
				RunGap:             getString(record, "run_gap", i),
				PassAttempt:        getInt32(record, "pass_attempt", i),
				RushAttempt:        getInt32(record, "rush_attempt", i),
				CompletePass:       getInt32(record, "complete_pass", i),
				Interception:       getInt32(record, "interception", i),
				Touchdown:          getInt32(record, "touchdown", i),
				FumbleLost:         getInt32(record, "fumble_lost", i),
				Sack:               getInt32(record, "sack", i),
				QBHit:              getInt32(record, "qb_hit", i),
				QBPressure:         getInt32(record, "qb_pressure", i),
				RusherPlayerId:     getString(record, "rusher_player_id", i),
				PasserPlayerId:     getString(record, "passer_player_id", i),
				ReceiverPlayerId:   getString(record, "receiver_player_id", i),
				Season:             getInt32(record, "season", i),
				Week:               getInt32(record, "week", i),
				HomeTeam:           getString(record, "home_team", i),
				AwayTeam:           getString(record, "away_team", i),
				HomeScore:          getInt32(record, "home_score", i),
				AwayScore:          getInt32(record, "away_score", i),
				GameSecondsRemaining: getInt32(record, "game_seconds_remaining", i),
				Drive:              getInt32(record, "drive", i),
				Down12:             getInt32(record, "down12", i),
				Down34:             getInt32(record, "down34", i),
				GoalToGo:           getInt32(record, "goal_to_go", i),
				Time:               getString(record, "time", i),
				SideOfField:        getString(record, "side_of_field", i),
				ThirdDownConverted: getInt32(record, "third_down_converted", i),
				FourthDownConverted: getInt32(record, "fourth_down_converted", i),
			}
			plays = append(plays, play)
		}
	}

	return plays, nil
}

func (c *Client) parseSchedules(r *file.Reader) ([]Game, error) {
	props := pqarrow.ArrowReadProperties{
		BatchSize: 8192,
	}
	
	arrowReader, err := pqarrow.NewFileReader(r, props, memory.DefaultAllocator)
	if err != nil {
		return nil, err
	}

	var games []Game
	
	recordReader, err := arrowReader.GetRecordReader(context.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer recordReader.Release()

	for recordReader.Next() {
		record := recordReader.Record()
		
		numRecords := record.NumRows()
		for i := 0; i < int(numRecords); i++ {
			game := Game{
				GameID:   getString(record, "game_id", i),
				Season:   getInt32(record, "season", i),
				Week:     getInt32(record, "week", i),
				GameType: getString(record, "game_type", i),
				GameDate: getString(record, "gameday", i),
				Weekday:  getString(record, "weekday", i),
				GameTime: getString(record, "gametime", i),
				HomeTeam: getString(record, "home_team", i),
				AwayTeam: getString(record, "away_team", i),
				HomeScore: getInt32(record, "home_score", i),
				AwayScore: getInt32(record, "away_score", i),
				Location: getString(record, "location", i),
				Stadium:  getString(record, "stadium", i),
				Surface:  getString(record, "surface", i),
				Roof:     getString(record, "roof", i),
				Temp:     getInt32(record, "temp", i),
				Wind:     getInt32(record, "wind", i),
				HomeCoach: getString(record, "home_coach", i),
				AwayCoach: getString(record, "away_coach", i),
				Referee:  getString(record, "referee", i),
				StadiumId: getString(record, "stadium_id", i),
			}
			games = append(games, game)
		}
	}

	return games, nil
}

func (c *Client) parseTeams(r *file.Reader) ([]Team, error) {
	props := pqarrow.ArrowReadProperties{
		BatchSize: 8192,
	}
	
	arrowReader, err := pqarrow.NewFileReader(r, props, memory.DefaultAllocator)
	if err != nil {
		return nil, err
	}

	var teams []Team
	
	recordReader, err := arrowReader.GetRecordReader(context.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer recordReader.Release()

	for recordReader.Next() {
		record := recordReader.Record()
		
		numRecords := record.NumRows()
		for i := 0; i < int(numRecords); i++ {
			team := Team{
				TeamAbbr:      getString(record, "team_abbr", i),
				TeamName:      getString(record, "team_name", i),
				TeamId:        getString(record, "team_id", i),
				TeamNick:      getString(record, "team_nick", i),
				TeamConf:      getString(record, "team_conf", i),
				TeamDiv:       getString(record, "team_div", i),
				TeamColor:     getString(record, "team_color", i),
				TeamColor2:    getString(record, "team_color2", i),
				TeamColor3:    getString(record, "team_color3", i),
				TeamColor4:    getString(record, "team_color4", i),
				TeamLogo:      getString(record, "team_logo", i),
				TeamLogoEspn:  getString(record, "team_logo_espn", i),
				TeamWordmark:  getString(record, "team_wordmark", i),
			}
			teams = append(teams, team)
		}
	}

	return teams, nil
}

func (c *Client) parseRosters(r *file.Reader) ([]Roster, error) {
	props := pqarrow.ArrowReadProperties{
		BatchSize: 8192,
	}
	
	arrowReader, err := pqarrow.NewFileReader(r, props, memory.DefaultAllocator)
	if err != nil {
		return nil, err
	}

	var rosters []Roster
	
	recordReader, err := arrowReader.GetRecordReader(context.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer recordReader.Release()

	for recordReader.Next() {
		record := recordReader.Record()
		
		numRecords := record.NumRows()
		for i := 0; i < int(numRecords); i++ {
			roster := Roster{
				Season:             getInt32(record, "season", i),
				Week:               getInt32(record, "week", i),
				Team:               getString(record, "team", i),
				PlayerId:           getString(record, "player_id", i),
				PlayerName:         getString(record, "player_name", i),
				Position:           getString(record, "position", i),
				DepthTeam:          getString(record, "depth_team", i),
				DepthPos:           getString(record, "depth_pos", i),
				Status:             getString(record, "status", i),
				GsisId:             getString(record, "gsis_id", i),
				EspnId:             getString(record, "espn_id", i),
				SportradarId:       getString(record, "sportradar_id", i),
				YahooId:            getString(record, "yahoo_id", i),
				RotowireId:         getString(record, "rotowire_id", i),
				PffId:              getString(record, "pff_id", i),
				PfrId:              getString(record, "pfr_id", i),
				FantasyDataId:      getString(record, "fantasy_data_id", i),
				StatusDescriptionAbbr: getString(record, "status_description_abbr", i),
				StatusShortDescription: getString(record, "status_short_description", i),
				HeadshotUrl:        getString(record, "headshot_url", i),
				NgSId:              getString(record, "ngs_id", i),
				College:            getString(record, "college", i),
				CollegeConf:        getString(record, "college_conf", i),
				DraftYear:          getInt32(record, "draft_year", i),
				DraftRound:         getInt32(record, "draft_round", i),
				DraftPick:          getInt32(record, "draft_pick", i),
				Height:             getString(record, "height", i),
				Weight:             getInt32(record, "weight", i),
				BirthDate:          getString(record, "birth_date", i),
				Age:                getInt32(record, "age", i),
				Experience:         getInt32(record, "experience", i),
			}
			rosters = append(rosters, roster)
		}
	}

	return rosters, nil
}

func getColumnIndex(record arrow.Record, name string) int {
	schema := record.Schema()
	for i, field := range schema.Fields() {
		if field.Name == name {
			return i
		}
	}
	return -1
}

func getInt32(record arrow.Record, name string, row int) int {
	idx := getColumnIndex(record, name)
	if idx == -1 {
		return 0
	}
	col := record.Column(idx)
	if col == nil {
		return 0
	}
	
	switch arr := col.(type) {
	case *array.Int32:
		return int(arr.Value(row))
	case *array.Int64:
		return int(arr.Value(row))
	case *array.Float64:
		return int(arr.Value(row))
	case *array.Float32:
		return int(arr.Value(row))
	default:
		return 0
	}
}

func getString(record arrow.Record, name string, row int) string {
	idx := getColumnIndex(record, name)
	if idx == -1 {
		return ""
	}
	col := record.Column(idx)
	if col == nil {
		return ""
	}
	
	switch arr := col.(type) {
	case *array.String:
		return arr.Value(row)
	case *array.Binary:
		return string(arr.Value(row))
	default:
		return ""
	}
}

func getFloat64(record arrow.Record, name string, row int) float64 {
	idx := getColumnIndex(record, name)
	if idx == -1 {
		return 0
	}
	col := record.Column(idx)
	if col == nil {
		return 0
	}
	
	switch arr := col.(type) {
	case *array.Float64:
		return arr.Value(row)
	case *array.Float32:
		return float64(arr.Value(row))
	case *array.Int32:
		return float64(arr.Value(row))
	case *array.Int64:
		return float64(arr.Value(row))
	default:
		return 0
	}
}