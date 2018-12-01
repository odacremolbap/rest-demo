package clauses

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	pageQuery     = "page"
	pageSizeQuery = "page_size"
	// OrderByQuery at query string
	OrderByQuery = "order"

	// Ordering
	ascending  = "asc"
	descending = "desc"
)

// PaginationClauseFromRequest uses PaginationClause using values from a map
// as prefixed input parameters
func PaginationClauseFromRequest(values map[string]string) (string, error) {
	page := 1
	pageSize := 50
	var err error

	if p := values[pageQuery]; p != "" {
		page, err = strconv.Atoi(p)
		if err != nil {
			return "", errors.Wrapf(err, "error parsing pagination %s", pageQuery)
		}
	}
	if p := values[pageSizeQuery]; p != "" {
		pageSize, err = strconv.Atoi(p)
		if err != nil {
			return "", errors.Wrapf(err, "error parsing pagination %s", pageSizeQuery)
		}
		// page_size 0 means to list all
		if pageSize == 0 {
			return "", nil
		}
	}
	return PaginationClause(page, pageSize)
}

// WhereClauseFromRequest given a values map builds a filter
// looking for allowed filter fields
func WhereClauseFromRequest(values map[string]string, allowedWhere []AllowedWhere) (string, []interface{}, error) {
	// build the FilterItem array of items that conform
	// the where clause
	var fis []FilterItem
	for _, v := range allowedWhere {
		if value := values[v.URLField]; value != "" {
			if v.Type != "" {
				// There must be a better way of doing this
				var err error
				switch v.Type {
				case "integer":
					_, err = strconv.Atoi(value)
				case "boolean":
					_, err = strconv.ParseBool(value)
				}
				if err != nil {
					return "", nil, errors.Wrapf(err, "field %s value %s can't be converted to %s",
						v.URLField, value, v.Type)
				}
			}
			fi := FilterItem{
				Field: v.DBField,
				Value: value,
				// TODO add comparisons, for now assume equal
				Comparison: "=",
			}
			fis = append(fis, fi)
		}
	}

	return WhereClause(fis)
}

// OrderByClauseFromRequest given a values map builds an order by clause
// Order can be specified at requests as:
// - ?order=field1
// - ?order=field1:asc
// - ?order=field1,field2:desc
func OrderByClauseFromRequest(values map[string]string, allowedOrderBy []string) (string, error) {

	urlOrder := values[OrderByQuery]
	if urlOrder == "" {
		return "", nil
	}

	uo := strings.Split(urlOrder, ",")
	orderby := ""
	for _, value := range uo {
		v := strings.Split(value, ":")
		if len(v) > 1 &&
			v[1] != ascending &&
			v[1] != descending {
			return "", errors.Errorf("ordering clause %s has wrong modifier %s",
				v[0], v[1])
		}
		isAllowed := false
		for _, allowed := range allowedOrderBy {
			if v[0] == allowed {
				if len(orderby) == 0 {
					orderby = fmt.Sprintf("%s", v[0])
				} else {
					orderby = fmt.Sprintf("%s,%s", orderby, v[0])
				}
				if len(v) > 1 {
					orderby = fmt.Sprintf("%s %s", orderby, v[1])
				}
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return "", errors.Errorf("field %s is not allowed for sorting", v[0])
		}
	}

	return orderby, nil
}

// BuildQueryClauseFromRequest for objects
func BuildQueryClauseFromRequest(
	values map[string]string,
	allowedWhere []AllowedWhere,
	allowedOrderBy []string) (*Query, error) {

	where, whereParams, err := WhereClauseFromRequest(values, allowedWhere)
	if err != nil {
		wrap := errors.Wrap(err, "error parsing query filters")
		return nil, wrap
	}

	pag, err := PaginationClauseFromRequest(values)
	if err != nil {
		wrap := errors.Wrap(err, "error parsing pagination")
		return nil, wrap
	}

	order, err := OrderByClauseFromRequest(values, allowedOrderBy)
	if err != nil {
		wrap := errors.Wrap(err, "error parsing query order")
		return nil, wrap
	}

	q := &Query{
		Where:         where,
		WhereParams:   whereParams,
		Pagination:    pag,
		OrderByClause: order,
	}

	return q, nil
}
