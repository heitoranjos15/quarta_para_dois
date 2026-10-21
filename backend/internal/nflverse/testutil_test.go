package nflverse

import (
	"bytes"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/apache/arrow/go/v16/parquet/pqarrow"
	"github.com/stretchr/testify/require"
)

// Columns of nflverse's play-by-play file that the fixtures can fill. Numbers
// are float64 in the real 2026 file (qtr, down, fixed_drive, ...).
var (
	fixtureStringCols = []string{
		"game_id", "play_type", "desc", "posteam", "defteam", "home_team", "away_team",
		"fixed_drive_result", "drive_time_of_possession", "field_goal_result",
		"extra_point_result", "two_point_conv_result", "td_team",
	}
	fixtureNumberCols = []string{
		"play_id", "qtr", "down", "ydstogo", "yardline_100", "yards_gained",
		"fixed_drive", "drive_quarter_start", "drive_quarter_end", "interception",
		"fumble_lost", "touchdown",
	}
)

// writePlaysParquet builds a parquet file with nflverse's column types: integer-like
// columns as float64, strings as utf8, everything nullable, plus one unused column.
func writePlaysParquet(t *testing.T, rows []map[string]any) []byte {
	t.Helper()
	return writePlaysParquetTyped(t, rows, arrow.PrimitiveTypes.Float64)
}

// writePlaysParquetTyped is writePlaysParquet with the integer-like columns stored as numType.
func writePlaysParquetTyped(t *testing.T, rows []map[string]any, numType arrow.DataType) []byte {
	t.Helper()

	var fields []arrow.Field
	for _, n := range fixtureStringCols {
		fields = append(fields, arrow.Field{Name: n, Type: arrow.BinaryTypes.String, Nullable: true})
	}
	for _, n := range fixtureNumberCols {
		fields = append(fields, arrow.Field{Name: n, Type: numType, Nullable: true})
	}
	fields = append(fields,
		arrow.Field{Name: "epa", Type: arrow.PrimitiveTypes.Float64, Nullable: true},
		arrow.Field{Name: "unused_col", Type: arrow.BinaryTypes.String, Nullable: true},
	)
	schema := arrow.NewSchema(fields, nil)

	b := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer b.Release()

	for _, row := range rows {
		for i, f := range fields {
			v, present := row[f.Name]
			if !present || v == nil {
				b.Field(i).AppendNull()
				continue
			}
			appendFixtureValue(t, b.Field(i), v)
		}
	}

	rec := b.NewRecord()
	defer rec.Release()
	tbl := array.NewTableFromRecords(schema, []arrow.Record{rec})
	defer tbl.Release()

	var buf bytes.Buffer
	require.NoError(t, pqarrow.WriteTable(tbl, &buf, 4096, nil, pqarrow.DefaultWriterProps()))
	return buf.Bytes()
}

func appendFixtureValue(t *testing.T, b array.Builder, v any) {
	t.Helper()
	switch bb := b.(type) {
	case *array.StringBuilder:
		bb.Append(v.(string))
	case *array.Float64Builder:
		bb.Append(toFloat(v))
	case *array.Float32Builder:
		bb.Append(float32(toFloat(v)))
	case *array.Int32Builder:
		bb.Append(int32(toFloat(v)))
	case *array.Int64Builder:
		bb.Append(int64(toFloat(v)))
	default:
		t.Fatalf("unsupported builder %T", b)
	}
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case float64:
		return n
	}
	panic("fixture number must be int or float64")
}
