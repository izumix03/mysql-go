package mysqlgo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
)

// Searcher for searching
type Searcher struct {
	con   *Connection
	query *selectQuery
}

// Count counts records
func (con *Connection) Count(record Basic) (count int, err error) {
	sq := &selectQuery{}
	query, _ := sq.
		selectFields(`count(*)`).
		from(record.Table()).
		build()

	log.Println(query)

	err = con.Query(query, func(rows *sql.Rows) func(rows *sql.Rows) error {
		return func(*sql.Rows) error {
			return rows.Scan(&count)
		}
	})
	return
}

// Select returns Searcher struct that is for executing query.
func (con *Connection) Select(fields ...string) *Searcher {
	sq := &selectQuery{}
	sq = sq.selectFields(fields...)
	return &Searcher{
		con:   con,
		query: sq,
	}
}

// From returns Searcher struct that is for executing query.
func (con *Connection) From(name string, alias string) (searcher *Searcher) {
	sq := &selectQuery{}
	sq = sq.from(name, alias)
	searcher = &Searcher{
		con:   con,
		query: sq,
	}
	return
}

// Where returns Searcher struct that is for executing query.
func (con *Connection) Where(query string, args ...interface{}) *Searcher {
	sq := &selectQuery{}
	sq = sq.selectFields().where(query, args...)
	return &Searcher{
		con:   con,
		query: sq,
	}
}

// WhereIn returns Searcher struct that is for executing query. (appended conditions to query)
func (con *Connection) WhereIn(columnName string, args ...interface{}) *Searcher {
	sq := &selectQuery{}
	query := fmt.Sprintf("%s in (%s)", columnName, repeatJoin(`?`, `,`, len(args)))
	sq = sq.where(query, args)
	return &Searcher{
		con:   con,
		query: sq,
	}
}

// From append table to query
func (s *Searcher) From(table string) *Searcher {
	s.query.from(table)
	return s
}

// Joins append join query
func (s *Searcher) Joins(joinQuery string) *Searcher {
	s.query.joins(joinQuery)
	return s
}

// Where append conditions to query
func (s *Searcher) Where(query string, args ...interface{}) *Searcher {
	s.query.where(query, args...)
	return s
}

// WhereIn append conditions to query
func (s *Searcher) WhereIn(columnName string, args ...interface{}) *Searcher {
	query := fmt.Sprintf("%s in (%s)", columnName, repeatJoin(`?`, `,`, len(args)))
	s.query.where(query, args...)
	return s
}

// OrderBy appends sort column and direction
func (s *Searcher) OrderBy(queries ...string) *Searcher {
	s.query.orderBy(queries...)
	return s
}

// Limit sets upper limit
func (s *Searcher) Limit(limit int) *Searcher {
	s.query.limit(limit)
	return s
}

// FindMany sets value to entities(that is slice)
func (s *Searcher) FindMany(basicSlice interface{}) (err error) {
	basics := reflect.Indirect(reflect.ValueOf(basicSlice)).Interface()
	basicIF := newPtrFromSlice(basics)
	basic, ok := basicIF.Interface().(Basic)
	if !ok {
		log.Printf("Unsupported interface %+v", basics)
		panic(errors.New(`assertion error`))
	}

	query, args := s.query.selectFields(extractColumnNames(basic)...).
		from(basic.Table()).
		build()

	log.Println(query)

	err = s.con.Query(query, func(rows *sql.Rows) func(rows *sql.Rows) error {
		return provideRowMapper(rows, basicSlice, reflect.ValueOf(basic))
	}, args...)
	return
}

// Count returns count
func (s *Searcher) Count(tableName string) (count int, err error) {
	query, _ := s.query.
		selectFields(`count(*)`).
		from(tableName).
		build()

	log.Println(query)

	err = s.con.Query(query, func(rows *sql.Rows) func(rows *sql.Rows) error {
		return func(*sql.Rows) error {
			return rows.Scan(&count)
		}
	})
	return
}

// LoadMapList sets value to map.
func (s *Searcher) LoadMapList(records *[]map[string]interface{}) error {
	return s.LoadPrefixedMapList(records, ``)
}

// LoadMapList sets value to map.
func (s *Searcher) LoadMap() (map[string]interface{}, error) {
	var records []map[string]interface{}
	err := s.LoadPrefixedMapList(&records, ``)
	if err != nil {
		return nil, err
	}
	if len(records) < 1 {
		return nil, err
	}
	return records[0], err
}

// LoadPrefixedMapList sets value to map.
// Allow specify the key prefix of map
func (s *Searcher) LoadPrefixedMapList(records *[]map[string]interface{}, keyPrefix string) (err error) {
	query, args := s.query.build()

	log.Printf("query => %+v. args => %+v", query, args)

	err = s.con.Query(query, func(rows *sql.Rows) func(rows *sql.Rows) error {
		return provideMapRowMapper(rows, records, keyPrefix)
	}, args...)

	if err != nil {
		log.Println("Failed to find many map...")
	}
	return
}

// FindOne sets value to a entities. limitRowCnt 1 is set.
// Should be a Basic struct
func (s *Searcher) FindOne(basic Basic) (err error) {
	query, args := s.query.selectFields(extractColumnNames(basic)...).
		from(basic.Table()).
		limit(1).
		build()

	log.Printf("query => %+v", query)
	log.Printf("args => %+v", args)

	err = s.con.Query(query, func(rows *sql.Rows) func(rows *sql.Rows) error {
		return provideRowMapper(rows, basic, reflect.ValueOf(basic))
	}, args...)
	return
}
