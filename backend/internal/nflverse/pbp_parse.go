package nflverse

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/apache/arrow/go/v16/parquet"
	"github.com/apache/arrow/go/v16/parquet/file"
	"github.com/apache/arrow/go/v16/parquet/pqarrow"
)

// playColumns maps an nflverse column name to the index of the Play field it fills.
var playColumns = func() map[string]int {
	t := reflect.TypeOf(Play{})
	m := make(map[string]int, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		if name := t.Field(i).Tag.Get("parquet"); name != "" {
			m[name] = i
		}
	}
	return m
}()

// ParsePlays reads a play-by-play parquet file, decoding only the columns Play uses.
// Null strings become "", null numbers become 0 and a null epa stays nil.
func ParsePlays(r parquet.ReaderAtSeeker) ([]Play, error) {
	pf, err := file.NewParquetReader(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = pf.Close() }()

	arrowReader, err := pqarrow.NewFileReader(pf, pqarrow.ArrowReadProperties{BatchSize: 8192}, memory.DefaultAllocator)
	if err != nil {
		return nil, err
	}
	schema, err := arrowReader.Schema()
	if err != nil {
		return nil, err
	}

	// The schema is flat, so a field's index is its leaf column index.
	var cols []int
	for i, f := range schema.Fields() {
		if _, ok := playColumns[f.Name]; ok {
			cols = append(cols, i)
		}
	}

	recordReader, err := arrowReader.GetRecordReader(context.Background(), cols, nil)
	if err != nil {
		return nil, err
	}
	defer recordReader.Release()

	var plays []Play
	for recordReader.Next() {
		batch := recordReader.Record()
		setters, err := batchSetters(batch)
		if err != nil {
			return nil, err
		}
		rows := int(batch.NumRows())
		for row := 0; row < rows; row++ {
			plays = append(plays, Play{})
			p := reflect.ValueOf(&plays[len(plays)-1]).Elem()
			for _, set := range setters {
				set(p, row)
			}
		}
	}
	if err := recordReader.Err(); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return plays, nil
}

type playSetter func(p reflect.Value, row int)

// batchSetters resolves each column of the batch once, so the row loop does no lookups.
func batchSetters(batch arrow.Record) ([]playSetter, error) {
	var setters []playSetter
	for i, f := range batch.Schema().Fields() {
		fieldIdx, ok := playColumns[f.Name]
		if !ok {
			continue
		}
		col := batch.Column(i)
		fieldType := reflect.TypeOf(Play{}).Field(fieldIdx).Type

		switch fieldType.Kind() {
		case reflect.String:
			c, ok := col.(*array.String)
			if !ok {
				return nil, fmt.Errorf("column %s: want string, got %s", f.Name, col.DataType())
			}
			setters = append(setters, func(p reflect.Value, row int) {
				if !c.IsNull(row) {
					p.Field(fieldIdx).SetString(c.Value(row))
				}
			})
		case reflect.Int:
			get, err := intGetter(col)
			if err != nil {
				return nil, fmt.Errorf("column %s: %w", f.Name, err)
			}
			setters = append(setters, func(p reflect.Value, row int) {
				if !col.IsNull(row) {
					p.Field(fieldIdx).SetInt(get(row))
				}
			})
		case reflect.Pointer: // *float64
			get, err := floatGetter(col)
			if err != nil {
				return nil, fmt.Errorf("column %s: %w", f.Name, err)
			}
			setters = append(setters, func(p reflect.Value, row int) {
				if !col.IsNull(row) {
					v := get(row)
					p.Field(fieldIdx).Set(reflect.ValueOf(&v))
				}
			})
		}
	}
	return setters, nil
}

func intGetter(col arrow.Array) (func(int) int64, error) {
	switch c := col.(type) {
	case *array.Int32:
		return func(i int) int64 { return int64(c.Value(i)) }, nil
	case *array.Int64:
		return func(i int) int64 { return c.Value(i) }, nil
	case *array.Float32:
		return func(i int) int64 { return int64(c.Value(i)) }, nil
	case *array.Float64:
		return func(i int) int64 { return int64(c.Value(i)) }, nil
	}
	return nil, fmt.Errorf("want number, got %s", col.DataType())
}

func floatGetter(col arrow.Array) (func(int) float64, error) {
	switch c := col.(type) {
	case *array.Float64:
		return c.Value, nil
	case *array.Float32:
		return func(i int) float64 { return float64(c.Value(i)) }, nil
	case *array.Int32:
		return func(i int) float64 { return float64(c.Value(i)) }, nil
	case *array.Int64:
		return func(i int) float64 { return float64(c.Value(i)) }, nil
	}
	return nil, fmt.Errorf("want number, got %s", col.DataType())
}
