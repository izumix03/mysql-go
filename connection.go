package mysqlgo

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Connection is MySQL connection wrapper
type Connection struct {
	DB *sql.DB
	tx *sql.Tx
}

// WithConnection provides connection to function.(Not transaction)
func WithConnection(ds string, fn func(con *Connection) error) error {
	con, err := createConnection(ds)
	if err != nil {
		return err
	}
	defer con.Close()
	return fn(con)
}

// Query for select and map rows
func (con *Connection) Query(sql string,
	rowMapperProvider func(rows *sql.Rows) func(rows *sql.Rows) error,
	args ...interface{}) error {

	rows, err := con.DB.Query(sql, args...)
	if err != nil {
		log.Printf("Failed to execute query. Sql => %s. Args => %+v", sql, args)
		return err
	}

	rowMapper := rowMapperProvider(rows)
	for rows.Next() {
		err = rowMapper(rows)
		if err != nil {
			return err
		}
	}
	return nil
}

// Exec executes and logs
func (con *Connection) Exec(s string, values ...interface{}) error {
	result, err := con.Execute(s, values...)
	if err != nil {
		log.Println(`Failed to executes...`)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Print(`Failed to get execute result...`)
	} else {
		log.Printf("Affected row number is %d", affected)
	}
	return nil
}

// Execute switches DB or tx and and executes
func (con *Connection) Execute(query string, args ...interface{}) (sql.Result, error) {
	if con.tx != nil {
		return con.tx.Exec(query, args...)
	}
	return con.DB.Exec(query, args...)
}

// Close Connection.
func (con *Connection) Close() {
	err := con.DB.Close()
	if err != nil {
		log.Println(`Failed to Close connection.`)
	}
}

// createConnection returns Connection if failed then panic
func createConnection(datasource string) (*Connection, error) {
	log.Println(datasource)
	db, err := sql.Open(`mysql`, datasource)
	if err != nil {
		log.Printf("Failed to connect DB. Destination => %s", datasource)
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping. Destination => %s", datasource)
		return nil, err
	}
	return setUpConnection(db), nil
}

// Arrange Connection
func setUpConnection(db *sql.DB) (con *Connection) {
	con = &Connection{DB: db}
	con.DB.SetMaxIdleConns(3)
	return
}
