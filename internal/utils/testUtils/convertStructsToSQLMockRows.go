package testUtils

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func ExtractColumnName(tag string) string {
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}

func ConvertStructsToSQLMockRows(s interface{}) *sqlmock.Rows {
	slice := reflect.ValueOf(s)
	if slice.Kind() != reflect.Slice {
		panic("ConvertStructsToSQLMockRows requires a slice")
	}

	var columnNames []string
	if slice.Len() > 0 {
		firstElem := slice.Index(0)
		if firstElem.Kind() == reflect.Ptr {
			firstElem = firstElem.Elem()
		}
		for i := 0; i < firstElem.NumField(); i++ {
			field := firstElem.Type().Field(i)
			gormTag := field.Tag.Get("gorm")
			columnName := ExtractColumnName(gormTag)
			if columnName == "" {
				columnName = field.Name
			}
			if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Interface {
				continue
			}
			columnNames = append(columnNames, columnName)
		}
	}

	rows := sqlmock.NewRows(columnNames)

	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		var columnValues []driver.Value
		for j := 0; j < elem.NumField(); j++ {
			field := elem.Type().Field(j)
			if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Interface {
				continue
			}

			fieldValue := elem.Field(j)
			var value driver.Value
			switch fieldValue.Interface().(type) {
			case time.Time, sql.NullString, uuid.UUID:
				value = fieldValue.Interface()
			default:
				if fieldValue.IsValid() && fieldValue.CanInterface() {
					value = fieldValue.Interface()
				} else {
					value = fieldValue.String()
				}
			}
			columnValues = append(columnValues, value)
		}
		rows.AddRow(columnValues...)
	}

	return rows
}
