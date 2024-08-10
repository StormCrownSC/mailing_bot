package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func CopyData(source interface{}, destination interface{}) {
	sourceValue := reflect.ValueOf(source)
	destValue := reflect.ValueOf(destination)
	if sourceValue.Kind() != reflect.Struct || destValue.Kind() != reflect.Ptr {
		return
	}

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldName := sourceValue.Type().Field(i).Name
		destFieldValue := destValue.Elem().FieldByName(sourceFieldName)

		if destFieldValue.IsValid() && destFieldValue.CanSet() {
			sourceFieldValue := sourceValue.Field(i)
			destFieldValue.Set(sourceFieldValue)
		}
	}
}

func SetPlaceholdersByNumber(quantity int, num int) string {
	var result strings.Builder

	result.WriteString("(")

	for i := num; i < quantity+num-1; i++ {
		result.WriteString(fmt.Sprintf("$%d, ", i))
	}

	result.WriteString(fmt.Sprintf("$%d", quantity+num-1))
	result.WriteString(")")

	return result.String()
}
