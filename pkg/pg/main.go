package pg

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Initialize(logger *zap.SugaredLogger) {
	log = logger
}

var ErrInvalidOption = errors.New("invalid role option")

type Dsn map[string]string

//// identifier returns the object name ready to be used in a sql query as an object name (e.a. select * from %s)
//func identifier(objectName string) (escaped string) {
//	return fmt.Sprintf("\"%s\"", strings.Replace(objectName, "\"", "\"\"", -1))
//}
//
//// quotedSqlValue uses proper quoting for values in SQL queries
//func quotedSqlValue(objectName string) (escaped string) {
//	return fmt.Sprintf("'%s'", strings.Replace(objectName, "'", "''", -1))
//}

// connectStringValue uses proper quoting for connect string values
func connectStringValue(objectName string) (escaped string) {
	return fmt.Sprintf("'%s'", strings.Replace(objectName, "'", "\\'", -1))
}
