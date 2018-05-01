package mysqlgo

import (
	"reflect"
)

// extractColumnNames returns column names
func extractColumnNames(basic Basic) (columnNames []string) {
	value := reflect.ValueOf(basic)
	structType := value.Type().Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		columnName := field.Tag.Get(`name`)
		if columnName == `` {
			continue
		}
		columnNames = append(columnNames, backQuoted(columnName))
	}
	return
}
