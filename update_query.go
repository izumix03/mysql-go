package mysqlgo

import (
	"fmt"
	"strings"
)

type updateQuery struct {
	table           string
	fields          []string
	whereConditions []string
	args            []interface{}
}

func (q *updateQuery) appendField(field string) {
	q.fields = append(q.fields, field)
}

func (q *updateQuery) appendArg(arg interface{}) {
	q.args = append(q.args, arg)
}

func (q *updateQuery) appendWhere(whereCondition string, args ...interface{}) *updateQuery {
	q.whereConditions = append(q.whereConditions, whereCondition)
	q.args = append(q.args, args...)
	return q
}

func (q *updateQuery) build() (query string, args []interface{}) {
	query = fmt.Sprintf("UPDATE %s", q.table)
	query = spaceJoin(query, `SET`)

	var setters []string
	for _, field := range q.fields {
		setters = append(setters, join(field, `?`, `=`))
	}
	query = spaceJoin(query, strings.Join(setters, `,`))

	if len(q.whereConditions) != 0 {
		query = spaceJoin(query, `WHERE`)
		query = spaceJoin(query, strings.Join(q.whereConditions, ` and `))
	}

	args = q.args
	return
}
