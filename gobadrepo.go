package gobadrepo

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type StructFieldsInfo struct {
	FieldName  string
	FieldType  string
	FieldValue string
}

type StructInfo struct {
	StructName       string
	StructFieldsInfo []StructFieldsInfo
}

var SQL_DATA_TYPES = map[string]string{
	"int":    "INTEGER",
	"string": "TEXT",
	"float":  "REAL",
}

var ID_AUTO_INCREMENT = map[string]string{
	"sqlite3":  "INTEGER PRIMARY KEY AUTOINCREMENT",
	"postgres": "SERIAL PRIMARY KEY",
	"mysql":    "INT AUTO_INCREMENT PRIMARY KEY",
}

func ExtractTableStructInfo(tableStruct any, fillSQLTypes bool, driverName string) StructInfo {
	objType := reflect.TypeOf(tableStruct)
	objValue := reflect.ValueOf(tableStruct)

	if objType.Kind() != reflect.Struct {
		panic("input is not a struct")
	}

	var structInfo StructInfo
	structName := strings.ToLower(objType.Name())
	if !strings.HasSuffix(structName, "s") {
		structName = structName + "s"
	}
	structInfo.StructName = structName

	idFieldFound := false
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		if field.Type.Kind() == reflect.Struct {
			panic("only flat structs are accepted")
		}

		fieldValue := objValue.Field(i)
		fieldType := fieldValue.Type().String()
		fieldValueStr := fmt.Sprintf("%v", fieldValue.Interface())
		dbTag := field.Tag.Get("db")

		if dbTag == "" {
			panic("db tag must be provided for all struct fields ex: Id int `db:\"id\"`")
		}

		if fillSQLTypes {
			mappedDataType := false
			for goDT, sqliteDT := range SQL_DATA_TYPES {
				if strings.HasPrefix(goDT, fieldType) {
					mappedDataType = true
					fieldType = sqliteDT
					break
				}
			}

			if !mappedDataType {
				fieldType = "TEXT"
			}

			if field.Name == "Id" {
				fieldType = ID_AUTO_INCREMENT[driverName]
				idFieldFound = true
			}
		}

		fieldInfo := StructFieldsInfo{
			FieldName:  dbTag,
			FieldType:  fieldType,
			FieldValue: fieldValueStr,
		}

		structInfo.StructFieldsInfo = append(structInfo.StructFieldsInfo, fieldInfo)
	}

	if !idFieldFound && fillSQLTypes {
		panic("table struct must have 'Id int' field")
	}

	return structInfo
}

func generateCreateTableSQL(info StructInfo, driverName string) string {
	var columns []string

	for _, field := range info.StructFieldsInfo {
		col := fmt.Sprintf("%s %s", field.FieldName, field.FieldType)
		columns = append(columns, col)
	}

	columnDefinitions := strings.Join(columns, ",\n\t")

	var createTableSQL string
	switch driverName {
	case "postgres":
		createTableSQL = fmt.Sprintf(`
DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_class c 
      WHERE c.relkind = 'S' 
      AND c.relname = '%[1]s_id_seq'
   ) THEN
      CREATE SEQUENCE %[1]s_id_seq;
   END IF;
END
$do$;
CREATE TABLE IF NOT EXISTS %[1]s (
	%[2]s
);
ALTER SEQUENCE %[1]s_id_seq OWNED BY %[1]s.id;
`, info.StructName, columnDefinitions)

	case "mysql":
		createTableSQL = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %[1]s (
	%[2]s
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`, info.StructName, columnDefinitions)
	case "sqlite3":
		createTableSQL = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %[1]s (
	%[2]s
);`, info.StructName, columnDefinitions)
	default:
		fmt.Println("Unsupported driver")
		return ""
	}

	return createTableSQL
}

func generateInsertQuery(structInfo []StructInfo, driverName string) (string, []interface{}) {
	var columns []string
	var values []string
	var args []interface{}

	for _, field := range structInfo[0].StructFieldsInfo {
		if field.FieldName == "id" {
			continue
		}
		columns = append(columns, field.FieldName)
	}

	for _, si := range structInfo {
		var tempValues []string
		for _, field := range si.StructFieldsInfo {
			if field.FieldName == "id" {
				continue
			}
			args = append(args, field.FieldValue)

			if driverName == "postgres" {
				tempValues = append(tempValues, fmt.Sprintf("$%d", len(args)))
			} else {
				tempValues = append(tempValues, "?")
			}
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(tempValues, ", ")))
	}

	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(values, ", ")

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", structInfo[0].StructName, columnsStr, valuesStr)

	return insertQuery, args
}

func generateFindQuery(structInfo StructInfo, limit, skip int, driverName string) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	placeholder := func(i int) string {
		if driverName == "postgres" {
			return fmt.Sprintf("$%d", i)
		}
		return "?"
	}

	i := 1
	for _, field := range structInfo.StructFieldsInfo {
		if field.FieldName == "id" && field.FieldValue == "0" {
			continue
		}
		if field.FieldValue != "" {
			condition := fmt.Sprintf("%s = %s", field.FieldName, placeholder(i))
			conditions = append(conditions, condition)
			args = append(args, field.FieldValue)
			i++
		}
	}

	conditionsStr := strings.Join(conditions, " AND ")
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", structInfo.StructName, conditionsStr)

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
		if skip > 0 {
			query += fmt.Sprintf(" OFFSET %d", skip)
		}
	}

	query += ";"

	return query, args
}

func generateUpdateQuery(structQueryInfo, structDataInfo StructInfo, driverName string) (string, []interface{}) {
	var updates []string
	var conditions []string
	var args []interface{}

	placeholder := func(i int) string {
		if driverName == "postgres" {
			return fmt.Sprintf("$%d", i)
		}
		return "?"
	}

	argIndex := 1
	for _, field := range structDataInfo.StructFieldsInfo {
		if field.FieldName == "id" {
			continue
		}
		if field.FieldValue != "" {
			update := fmt.Sprintf("%s = %s", field.FieldName, placeholder(argIndex))
			updates = append(updates, update)
			args = append(args, field.FieldValue)
			argIndex++
		}
	}

	for _, field := range structQueryInfo.StructFieldsInfo {
		if field.FieldName == "id" && field.FieldValue == "0" {
			continue
		}
		if field.FieldValue != "" {
			condition := fmt.Sprintf("%s = %s", field.FieldName, placeholder(argIndex))
			conditions = append(conditions, condition)
			args = append(args, field.FieldValue)
			argIndex++
		}
	}

	updatesStr := strings.Join(updates, ", ")
	conditionsStr := strings.Join(conditions, " AND ")

	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s;", structDataInfo.StructName, updatesStr, conditionsStr)

	return updateQuery, args
}

func generateDeleteQuery(structQueryInfo StructInfo, driverName string) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	placeholder := func(i int) string {
		if driverName == "postgres" {
			return fmt.Sprintf("$%d", i)
		}
		return "?"
	}

	argIndex := 1
	for _, field := range structQueryInfo.StructFieldsInfo {
		if field.FieldName == "id" && field.FieldValue == "0" {
			continue
		}
		if field.FieldValue != "" {
			condition := fmt.Sprintf("%s = %s", field.FieldName, placeholder(argIndex))
			conditions = append(conditions, condition)
			args = append(args, field.FieldValue)
			argIndex++
		}
	}

	conditionsStr := strings.Join(conditions, " AND ")

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE %s;", structQueryInfo.StructName, conditionsStr)

	return deleteQuery, args
}

func generateDeleteAllRowsQuery(structQueryInfo StructInfo) string {
	return fmt.Sprintf("DELETE FROM %s;", structQueryInfo.StructName)
}

func CreateTable(db *sqlx.DB, tableStruct any) error {
	structInfo := ExtractTableStructInfo(tableStruct, true, db.DriverName())
	sqlQuery := generateCreateTableSQL(structInfo, db.DriverName())
	_, err := db.Exec(sqlQuery)
	return err
}

func InsertOne(db *sqlx.DB, tableStruct any) error {
	structInfo := ExtractTableStructInfo(tableStruct, false, db.DriverName())
	structInfos := []StructInfo{}
	structInfos = append(structInfos, structInfo)
	sqlQuery, args := generateInsertQuery(structInfos, db.DriverName())
	_, err := db.Exec(sqlQuery, args...)
	return err
}

func InsertMany(db *sqlx.DB, tablesStructs []interface{}) error {
	structInfos := []StructInfo{}
	for _, tableStruct := range tablesStructs {
		si := ExtractTableStructInfo(tableStruct, false, db.DriverName())
		structInfos = append(structInfos, si)
	}
	sqlQuery, args := generateInsertQuery(structInfos, db.DriverName())
	_, err := db.Exec(sqlQuery, args...)
	return err
}

func FindOne(db *sqlx.DB, tableStructQuery, tableStruct any) error {
	structInfo := ExtractTableStructInfo(tableStructQuery, false, db.DriverName())
	findQuery, args := generateFindQuery(structInfo, 1, 0, db.DriverName())
	err := db.Get(tableStruct, findQuery, args...)
	if err != nil {
		return err
	}
	return nil
}

func FindManyLimitSkip(db *sqlx.DB, tableStructQuery, tableStruct any, limit, skip int) error {
	structInfo := ExtractTableStructInfo(tableStructQuery, false, db.DriverName())
	findQuery, args := generateFindQuery(structInfo, limit, skip, db.DriverName())
	err := db.Select(tableStruct, findQuery, args...)
	if err != nil {
		return err
	}
	return nil
}

func FindMany(db *sqlx.DB, tableStructQuery, tableStruct any) error {
	return FindManyLimitSkip(db, tableStructQuery, tableStruct, 0, 0)
}

func UpdateMany(db *sqlx.DB, tableStructQuery, tableStructUpdate any) error {
	structQueryInfo := ExtractTableStructInfo(tableStructQuery, false, db.DriverName())
	structDataInfo := ExtractTableStructInfo(tableStructUpdate, false, db.DriverName())
	sqlQuery, args := generateUpdateQuery(structQueryInfo, structDataInfo, db.DriverName())
	_, err := db.Exec(sqlQuery, args...)
	return err
}

func DeleteMany(db *sqlx.DB, tableStructQuery any) error {
	structQueryInfo := ExtractTableStructInfo(tableStructQuery, false, db.DriverName())
	sqlQuery, args := generateDeleteQuery(structQueryInfo, db.DriverName())
	_, err := db.Exec(sqlQuery, args...)
	return err
}

func DeleteAllRows(db *sqlx.DB, tableStructQuery any) error {
	structQueryInfo := ExtractTableStructInfo(tableStructQuery, false, db.DriverName())
	sqlQuery := generateDeleteAllRowsQuery(structQueryInfo)
	_, err := db.Exec(sqlQuery)
	return err
}
