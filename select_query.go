package mysqlgo

import (
	"fmt"
	"strconv"
	"strings"
)

type selectQuery struct {
	table           *mainTable
	fields          []string
	joinQuery       string
	whereConditions []string
	orderBys        []string
	limitRowCnt     int
	args            []interface{}
}

type mainTable struct {
	name  string
	alias string
}

// selectFields sets field that will be selected
func (q *selectQuery) selectFields(fields ...string) *selectQuery {
	q.fields = fields
	return q
}

// from sets target mainTable
func (q *selectQuery) from(table string, aliasOption ...string) *selectQuery {
	var alias string
	if len(aliasOption) == 1 {
		alias = aliasOption[0]
	}

	if q.table == nil {
		q.table = &mainTable{
			name:  table,
			alias: alias,
		}
	} else {
		q.table.name = table
		if present(alias) {
			q.table.alias = alias
		}
	}
	return q
}

// where sets where conditions and args
func (q *selectQuery) where(whereCondition string, args ...interface{}) *selectQuery {
	q.whereConditions = append(q.whereConditions, whereCondition)
	q.args = append(q.args, args...)
	return q
}

// joins joins table... Cannot check strictly... Take care
func (q *selectQuery) joins(joinQuery string) *selectQuery {
	q.joinQuery = joinQuery
	return q
}

// orderBy appends order by conditions
func (q *selectQuery) orderBy(queries ...string) *selectQuery {
	q.orderBys = append(q.orderBys, queries...)
	return q
}

// limit sets limitRowCnt
func (q *selectQuery) limit(limitCnt int) *selectQuery {
	q.limitRowCnt = limitCnt
	return q
}

// alias returns main table alias if present
func (q *selectQuery) aliasPrefix() string {
	if blank(q.table.alias) {
		return ``
	}
	return q.table.alias + `.`
}

// fieldList builds field list query for sql
// ex. `id`, `name`, ...
// ex. t.`id`, t.`name`, ...
func (q *selectQuery) fieldsList() string {
	if len(q.fields) == 0 {
		return q.aliasPrefix() + `*`
	}
	return prefixJoin(q.aliasPrefix(), q.fields, `,`)
}

// tableName returns table name with alias
// ex. `segment_tasks`
// ex. `segment_tasks` s
func (q *selectQuery) tableName() string {
	if blank(q.table.alias) {
		return backQuoted(q.table.name)
	}
	return spaceJoin(backQuoted(q.table.name), q.table.alias)
}

// build builds query string with struct
func (q *selectQuery) build() (query string, args []interface{}) {
	query = fmt.Sprintf("SELECT")
	query = spaceJoin(query, q.fieldsList())
	query = spaceJoin(query, `FROM`)
	query = spaceJoin(query, q.tableName())

	if present(q.joinQuery) {
		query = spaceJoins(query, `join`, q.joinQuery)
	}

	if len(q.whereConditions) != 0 {
		query = spaceJoin(query, `WHERE`)
		query = spaceJoin(query, strings.Join(q.whereConditions, ` and `))
	}
	if len(q.orderBys) != 0 {
		query = spaceJoin(query, `ORDER BY `+strings.Join(q.orderBys, `,`))
	}
	if q.limitRowCnt != 0 {
		query = spaceJoin(query, `LIMIT `+strconv.Itoa(q.limitRowCnt))
	}
	args = q.args
	return
}
