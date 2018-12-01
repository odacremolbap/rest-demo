package clauses

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Query is a placeholder for SQL clauses
type Query struct {
	Where         string
	WhereParams   []interface{}
	Pagination    string
	OrderByClause string
}

// FilterItem is a placeholder for SQL where clause items
type FilterItem struct {
	Field      string
	Comparison string
	Value      interface{}
}

// AllowedWhere keeps the allowed fields to build queries
// Type will be checked at validation
// TODO this info might be extracted using reflection from
// the model type, or be generated
type AllowedWhere struct {
	URLField string
	DBField  string
	Type     string
}

// OrderItem is a placeholder for SQL orderby clause items
type OrderItem struct {
	Field string
	Sort  string
}

// PaginationClause will return the pagination sql clause
// - page will start at 1 and page size will default to 50
// - returned value doesn't include trailing spaces
// - arguments are int typed, no risk for injection
func PaginationClause(page, pageSize int) (string, error) {
	if page < 1 {
		return "", errors.New("pagination starts at page 1")
	}
	if pageSize < 1 {
		return "", errors.New("pagination page size needs to be 1 or greater")
	}

	offset := (page - 1) * pageSize
	return fmt.Sprintf("offset %d limit %d", offset, pageSize), nil
}

// WhereClause will return sql where items clause
// - returned value doesn't include trailing spaces
func WhereClause(filters []FilterItem) (string, []interface{}, error) {

	l := len(filters)
	where := strings.Builder{}
	values := []interface{}{}
	for i, f := range filters {
		if len(f.Field) == 0 {
			return "", nil, errors.New("missing 'field' at the filter clause")
		}
		if f.Comparison != ">" &&
			f.Comparison != "<" &&
			f.Comparison != "=" {
			return "", nil,
				fmt.Errorf("%s is not one of the supported compare clauses", f.Comparison)
		}
		if f.Value == nil {
			return "", nil, errors.New("missing 'value' at the filter clause")
		}

		where.WriteString(fmt.Sprintf("%s %s $%d", f.Field, f.Comparison, i+1))
		if i != l-1 {
			where.WriteString(" and ")
		}
		values = append(values, f.Value)
	}
	return where.String(), values, nil
}

// TODO create order by generic (currently only OrderByClause at the URL helper)
