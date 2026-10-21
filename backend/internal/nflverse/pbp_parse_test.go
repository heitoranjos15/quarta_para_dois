package nflverse

import (
	"bytes"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePlays_Fields(t *testing.T) {
	data := writePlaysParquet(t, []map[string]any{{
		"game_id":                  "2026_06_DAL_GB",
		"fixed_drive":              5,
		"fixed_drive_result":       "Turnover",
		"drive_quarter_start":      1,
		"drive_quarter_end":        2,
		"drive_time_of_possession": "2:35",
		"field_goal_result":        "made",
		"td_team":                  "DAL",
		"epa":                      0.42,
	}})

	plays, err := ParsePlays(bytes.NewReader(data))
	require.NoError(t, err)
	require.Len(t, plays, 1)

	p := plays[0]
	assert.Equal(t, "2026_06_DAL_GB", p.GameID)
	assert.Equal(t, 5, p.FixedDrive)
	assert.Equal(t, "Turnover", p.FixedDriveResult)
	assert.Equal(t, 1, p.DriveQuarterStart)
	assert.Equal(t, 2, p.DriveQuarterEnd)
	assert.Equal(t, "2:35", p.DriveTimeOfPossession)
	assert.Equal(t, "made", p.FieldGoalResult)
	assert.Equal(t, "DAL", p.TDTeam)
	require.NotNil(t, p.EPA)
	assert.Equal(t, 0.42, *p.EPA)
}

func TestParsePlays_NullEPA(t *testing.T) {
	data := writePlaysParquet(t, []map[string]any{{"game_id": "g", "desc": nil, "qtr": nil, "epa": nil}})

	plays, err := ParsePlays(bytes.NewReader(data))
	require.NoError(t, err)
	require.Len(t, plays, 1)

	assert.Nil(t, plays[0].EPA)
	assert.Equal(t, "", plays[0].Description)
	assert.Equal(t, 0, plays[0].Quarter)
}

func TestParsePlays_IntColumnTypes(t *testing.T) {
	rows := []map[string]any{{"game_id": "g", "qtr": 3, "fixed_drive": 7}}
	for name, typ := range map[string]arrow.DataType{
		"int32":   arrow.PrimitiveTypes.Int32,
		"int64":   arrow.PrimitiveTypes.Int64,
		"float64": arrow.PrimitiveTypes.Float64,
	} {
		t.Run(name, func(t *testing.T) {
			plays, err := ParsePlays(bytes.NewReader(writePlaysParquetTyped(t, rows, typ)))
			require.NoError(t, err)
			require.Len(t, plays, 1)
			assert.Equal(t, 3, plays[0].Quarter)
			assert.Equal(t, 7, plays[0].FixedDrive)
		})
	}
}

func TestParsePlays_MultipleBatches(t *testing.T) {
	const n = 10000
	rows := make([]map[string]any, n)
	for i := range rows {
		rows[i] = map[string]any{"game_id": "g", "play_id": i}
	}

	plays, err := ParsePlays(bytes.NewReader(writePlaysParquet(t, rows)))
	require.NoError(t, err)
	require.Len(t, plays, n)
	for i, p := range plays {
		if p.PlayID != i {
			t.Fatalf("play %d out of order: got play_id %d", i, p.PlayID)
		}
	}
}
