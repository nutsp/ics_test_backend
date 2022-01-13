package utils

import "strings"

func ConcatCondition(conditions []string) string {
	var condition string
	condition = strings.Join(conditions, " AND ")

	return condition
}
