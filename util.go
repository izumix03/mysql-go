package mysqlgo

import "reflect"

// concat returns joined string(a + b)
func concat(a string, b string) string {
	vals := make([]byte, 0, 10)
	vals = append(vals, a...)
	vals = append(vals, b...)
	return string(vals)
}

// join returns joined string(a + separator +  b)
func join(a string, b string, separator string) string {
	vals := make([]byte, 0, 10)
	vals = append(vals, a...)
	vals = append(vals, separator...)
	vals = append(vals, b...)
	return string(vals)
}

// spaceJoin returns joined string with space
func spaceJoin(a string, b string) string {
	return join(a, b, ` `)
}

// spaceJoins returns joined string with space from array
func spaceJoins(values ...string) string {
	if len(values) == 0 {
		return ``
	}
	vals := make([]byte, 0, 10)
	for index, val := range values {
		if index == 0 {
			vals = append(vals, val...)
		} else {
			vals = append(vals, ' ')
			vals = append(vals, val...)
		}
	}
	return string(vals)
}

// BackQuoted returns back quoted string
// ex. "aaa" -> "`aaa`"
func backQuoted(a string) string {
	return join("`", "`", a)
}

// repeatJoin returns joined string with separator the times
func repeatJoin(a string, separator string, count int) string {
	if count == 0 {
		return ``
	}
	vals := make([]byte, 0, 10)
	vals = append(vals, a...)
	for i := 1; i < count; i++ {
		vals = append(vals, separator...)
		vals = append(vals, a...)
	}
	return string(vals)
}

// newPtrFromSlice return struct instance from interface that is actually struct slice
func newPtrFromSlice(basics interface{}) reflect.Value {
	any := reflect.ValueOf(basics)
	return reflect.New(any.Type().Elem())
}

// blank returns string is blank or not
func blank(a string) bool {
	return len(a) == 0
}

// present returns string is not black or not
func present(a string) bool {
	return !blank(a)
}

// prefixJoin returns joined string ( prefix + element + separator + prefix + element...)
func prefixJoin(prefix string, array []string, separator string) (result string) {
	if len(array) == 0 {
		return
	}
	for index, val := range array {
		if index == 0 {
			result = val
		} else {
			result = join(result, concat(prefix, val), separator)
		}
	}
	return
}
