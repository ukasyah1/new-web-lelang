package dbutil

import "strings"

// QualifiedTable returns a safely normalized Oracle schema-qualified table name.
func QualifiedTable(schema, table string) string {
	schema = strings.ToUpper(strings.TrimSpace(schema))
	if schema == "" {
		return table
	}
	return schema + "." + table
}
