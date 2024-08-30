package pg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgtype"
)

func ResultValueToString(value interface{}) (s string, err error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case time.Duration:
		return v.String(), nil
	case time.Time:
		return v.String(), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int8:
		return fmt.Sprintf("%d", v), nil
	case int16:
		return fmt.Sprintf("%d", v), nil
	case int32:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case uint:
		return fmt.Sprintf("%d", v), nil
	case uint8:
		return fmt.Sprintf("%d", v), nil
	case uint16:
		return fmt.Sprintf("%d", v), nil
	case uint32:
		return fmt.Sprintf("%d", v), nil
	case uint64:
		return fmt.Sprintf("%d", v), nil
	case []byte:
		return fmt.Sprintf("%x", v), nil
	case pgtype.Float4Array:
		var l []string
		for _, e := range v.Elements {
			l = append(l, fmt.Sprintf("%f", e.Float))
		}
		return fmt.Sprintf("[%s]", strings.Join(l, ",")), nil
	case pgtype.Line:
		return fmt.Sprintf("%f %f %f", v.A, v.B, v.C), nil
	case pgtype.Interval:
		var retVal time.Duration
		if err = v.AssignTo(&retVal); err != nil {
			return
		} else {
			return retVal.String(), nil
		}
	case pgtype.Numeric:
		var retVal float64
		if err = v.AssignTo(&retVal); err != nil {
			return
		} else {
			return strconv.FormatFloat(retVal, 'g', -1, 64), nil
		}
	case *pgtype.JSON:
		var retVal string
		if err = v.AssignTo(&retVal); err != nil {
			return
		} else {
			return retVal, nil
		}
	case *pgtype.JSONB:
		var retVal string
		if err = v.AssignTo(&retVal); err != nil {
			return
		} else {
			return retVal, nil
		}
	case nil:
		return "nil", nil
	default:
		return fmt.Sprintf("unknown datatype %v (%T)", value, value), nil
	}
}

type Result map[string]string
type Results []Result

func NewResultFromByteArrayArray(cols []string, values []interface{}) (ofr Result, err error) {
	ofr = make(Result)
	if len(cols) != len(values) {
		return ofr, fmt.Errorf("number of cols different then number of values")
	}
	for i, col := range cols {
		fmt.Printf("datatype: %T\n", values[i])
		if ofr[col], err = ResultValueToString(values[i]); err != nil {
			return
		}
	}
	return
}

func (ofr Result) String() (s string) {
	var results []string
	for key, value := range ofr {
		results = append(results, fmt.Sprintf("%s: %s",
			FormattedString(key),
			FormattedString(value)))
	}
	return fmt.Sprintf("{ %s }", strings.Join(results, ", "))
}

func (ofr Result) Columns() (result []string) {
	for key := range ofr {
		result = append(result, key)
	}
	return result
}

func FormattedString(s string) string {
	return fmt.Sprintf("'%s'", strings.Replace(s, "'", "\\'", -1))
}

func (ofr Result) Compare(other Result) (err error) {
	if len(ofr) != len(other) {
		return fmt.Errorf("number of columns different between row %v and compared row %v",
			ofr.Columns(), other.Columns())
	}
	for key, value := range ofr {
		otherValue, exists := other[key]
		if !exists {
			return fmt.Errorf("column row (%s) not in compared row", FormattedString(key))
		}
		if matched, err := regexp.MatchString(otherValue, value); err != nil {
			if value != otherValue {
				return fmt.Errorf("comparedrow is not an re, and column %s differs between row (%s), and comparedrow (%s)",
					FormattedString(key),
					FormattedString(value),
					FormattedString(otherValue))
			}
		} else if !matched {
			return fmt.Errorf("column %s value (%s) does not match with regular expression (%s)",
				FormattedString(key),
				FormattedString(value),
				FormattedString(otherValue))
		}
	}
	return nil
}

func (results Results) String() (s string) {
	var arr []string
	if len(results) == 0 {
		return "[ ]"
	}
	for _, result := range results {
		arr = append(arr, result.String())
	}
	return fmt.Sprintf("[ %s ]", strings.Join(arr, ", "))
}

func (results Results) Compare(other Results) (err error) {
	if len(results) != len(other) {
		return fmt.Errorf("different result (%s) then expected (%s)", results.String(),
			other.String())
	}
	for i, result := range results {
		err = result.Compare(other[i])
		if err != nil {
			return fmt.Errorf("different %d'th result: %s", i, err.Error())
		}
	}
	return nil
}
