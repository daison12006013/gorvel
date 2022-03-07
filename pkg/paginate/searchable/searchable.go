package searchable

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/daison12006013/gorvel/pkg/paginate"
)

type Input struct {
	Visible       bool
	Placeholder   interface{}
	Value         interface{}
	CanSearch     bool
	SearchColumn  []string
	SearchPattern string
}

type Header struct {
	Name  interface{}
	Input Input
}

type Table struct {
	Paginate    paginate.Paginate
	Headers     []Header
	Params      map[string]string
	OrderByCol  *string
	OrderBySort *string
}

func (st *Table) QueryCount(table string) sq.SelectBuilder {
	return sq.Select("count(*)").From(table)
}

func (st *Table) QuerySelect(table string) sq.SelectBuilder {
	builder := sq.Select("*").From(table).
		OrderBy(*st.OrderByCol + " " + *st.OrderBySort).
		Limit(uint64(st.Paginate.PerPage)).
		Offset(uint64(((st.Paginate.CurrentPage) - 1) * st.Paginate.PerPage))

	for _, header := range st.Headers {
		if !header.Input.CanSearch || header.Input.Value == "" {
			continue
		}

		var pred sq.Or
		for _, searchColumn := range header.Input.SearchColumn {
			switch header.Input.SearchPattern {
			case "-":
				pred = append(pred, sq.Eq{searchColumn: fmt.Sprintf("%v", header.Input.Value)})
			case "<-":
				pred = append(pred, sq.Like{searchColumn: "%" + fmt.Sprintf("%v", header.Input.Value)})
			case "->":
				pred = append(pred, sq.Like{searchColumn: fmt.Sprintf("%v", header.Input.Value) + "%"})
			case "<->":
				pred = append(pred, sq.Like{searchColumn: "%" + fmt.Sprintf("%v", header.Input.Value) + "%"})
			}
		}
		builder = builder.Where(pred)
	}
	return builder
}