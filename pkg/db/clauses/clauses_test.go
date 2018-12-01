package clauses

import (
	"fmt"
	"testing"
)

func TestPaginationClause(t *testing.T) {
	var pagTests = []struct {
		page           int
		pageSize       int
		expectedClause string
		expectedErr    bool
	}{
		{1, 1, "offset 0 limit 1", false},
		{2, 5, "offset 5 limit 5", false},
		{3, 10, "offset 20 limit 10", false},
		{3, 0, "", true},
		{0, 0, "", true},
	}

	for _, pt := range pagTests {
		t.Run(fmt.Sprintf("pag %d, size %d", pt.page, pt.pageSize),
			func(t *testing.T) {
				out, err := PaginationClause(pt.page, pt.pageSize)
				if out != pt.expectedClause {
					t.Errorf("got %q, wanted %q", out, pt.expectedClause)
				}
				if (err != nil) != pt.expectedErr {
					t.Errorf("got %t, wanted %t", err != nil, pt.expectedErr)
				}
			})
	}
}

func TestOrderByClause(t *testing.T) {

	var sortTests = []struct {
		values         map[string]string
		allowedOrderBy []string
		expectedClause string
		expectedErr    bool
	}{
		{
			map[string]string{"order": "field1"},
			[]string{"field1"},
			"field1",
			false,
		},
		{
			map[string]string{"order": "field1:asc"},
			[]string{"field1"},
			"field1 asc",
			false,
		},
		{
			map[string]string{"order": "field1,field2"},
			[]string{"field1"},
			"",
			true,
		},
		{
			map[string]string{"order": "field1,field2"},
			[]string{"field1", "field2", "field3"},
			"field1,field2",
			false,
		},
		{
			map[string]string{"order": "field1,field3,field2"},
			[]string{"field1", "field2", "field3"},
			"field1,field3,field2",
			false,
		},
		{
			map[string]string{"order": "field1,field4,field2"},
			[]string{"field1", "field2", "field3"},
			"",
			true,
		},
		{
			map[string]string{"order": "field1:desc,field3:asc,field2:desc"},
			[]string{"field1", "field2", "field3"},
			"field1 desc,field3 asc,field2 desc",
			false,
		},
	}

	for i, st := range sortTests {
		t.Run(fmt.Sprintf("sort test %d, size %+v", i, st.values),
			func(t *testing.T) {
				out, err := OrderByClauseFromRequest(st.values, st.allowedOrderBy)
				if out != st.expectedClause {
					t.Errorf("got %q, wanted %q", out, st.expectedClause)
				}
				if (err != nil) != st.expectedErr {
					t.Logf("got error: %s", err)
					t.Errorf("got %t, wanted %t", err != nil, st.expectedErr)
				}
			})
	}
}
