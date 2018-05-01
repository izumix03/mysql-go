package mysqlgo

import (
	"database/sql"
	"log"
)

type Updater struct {
	con   *Connection
	query *updateQuery
}

func (con *Connection) UpdateTable(table string) *Updater {
	sq := &updateQuery{
		table: table,
	}
	return &Updater{
		con:   con,
		query: sq,
	}
}

func (u *Updater) Set(attrs ...interface{}) *Updater {
	if len(attrs) == 0 || len(attrs)%2 == 1 {
		log.Printf("invalid args... %+v", attrs)
		return u
	}
	for i, att := range attrs {
		if i%2 == 0 {
			field, ok := att.(string)
			if !ok {
				log.Printf(`Invalid args when setting update field.`)
			}
			u.query.appendField(field)
		} else {
			u.query.appendArg(att)
		}
	}
	return u
}

func (u *Updater) Where(query string, args ...interface{}) *Updater {
	u.query.appendWhere(query, args...)
	return u
}

func (u *Updater) Execute() (sql.Result, error) {
	sql, args := u.query.build()
	log.Printf("Update query => %s", sql)
	return u.con.Execute(sql, args...)
}
