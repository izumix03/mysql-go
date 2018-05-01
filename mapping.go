package mysqlgo

import (
	"database/sql"
	"log"
	"reflect"
)

// provideRowMapper provides row mapper
func provideRowMapper(rows *sql.Rows, basicsIF interface{}, basic reflect.Value) func(*sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return func(*sql.Rows) (err error) {
			return
		}
	}

	structType := basic.Type().Elem()
	var indexes []int

	for _, column := range columns {
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)

			col := field.Tag.Get(`name`)
			if col == column {
				indexes = append(indexes, i)
			}
		}
	}

	isSlice := reflect.TypeOf(basicsIF).Elem().Kind() == reflect.Slice
	basicsIFPt := reflect.Indirect(reflect.ValueOf(basicsIF))

	return func(rs *sql.Rows) (err error) {
		if isSlice {
			e := reflect.New(basic.Type().Elem())
			err = rs.Scan(buildFieldList(e, indexes)...)
			basicsIFPt.Set(reflect.Append(basicsIFPt, e.Elem()))
		} else {
			err = rs.Scan(buildFieldList(basic, indexes)...)
		}
		return
	}
}

// buildFieldList builds fields list included result sets
func buildFieldList(any reflect.Value, indexes []int) (mappings []interface{}) {
	fields := any.Elem()
	for _, ind := range indexes {
		mappings = append(mappings, fields.Field(ind).Addr().Interface())
	}
	return
}

func provideMapRowMapper(rows *sql.Rows, records *[]map[string]interface{}, keyPrefix string) func(*sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return func(*sql.Rows) error {
			return err
		}
	}
	return func(rs *sql.Rows) error {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			log.Printf("Failed to scan rows are %+v", rows)
			return err
		}
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			entityVal := *val
			bytesVal, ok := entityVal.([]byte)

			var result interface{}
			if !ok {
				result = entityVal
			} else {
				result = string(bytesVal)
			}
			m[keyPrefix+colName] = result
		}
		*records = append(*records, m)
		return nil
	}
}
