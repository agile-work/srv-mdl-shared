package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Metadata defines metadata on the response
type Metadata struct {
	Filter     map[string]map[string]interface{} `json:"filter,omitempty"`
	Pagination map[string]int                    `json:"pagination,omitempty"`
	Order      []map[string]string               `json:"order,omitempty"`
	Columns    string                            `json:"columns,omitempty"`
	Group      string                            `json:"group,omitempty"`
}

// Load gets metadata from request
func (m *Metadata) Load(req *http.Request) error {
	metadataStr := req.URL.Query().Get("metadata")
	if err := json.Unmarshal([]byte(metadataStr), m); err != nil {
		return customerror.New(http.StatusBadRequest, "metadata load unmarshal", err.Error())
	}
	return nil
}

// GenerateDBOptions convert the metadata in a db options
func (m *Metadata) GenerateDBOptions() *db.Options {
	opt := &db.Options{}
	opt.Columns = strings.Split(m.Columns, ",")
	opt.Limit = m.Pagination["totalItens"]
	opt.Offset = (m.Pagination["selected"] * m.Pagination["totalItens"]) - m.Pagination["totalItens"]
	for _, row := range m.Order {
		for column, order := range row {
			if order == "asc" {
				opt.AddOrderBy(builder.Asc(column))
			} else {
				opt.AddOrderBy(builder.Desc(column))
			}
		}
	}
	for column, prop := range m.Filter {
		for operator, value := range prop {
			switch operator {
			case "=":
				opt.AddCondition(builder.Equal(column, value))
				break
			case "!=":
				opt.AddCondition(builder.NotEqual(column, value))
				break
			case ">":
				opt.AddCondition(builder.GreaterThen(column, value))
				break
			case ">=":
				opt.AddCondition(builder.GreaterOrEqual(column, value))
				break
			case "<":
				opt.AddCondition(builder.LowerThen(column, value))
				break
			case "<=":
				opt.AddCondition(builder.LowerOrEqual(column, value))
				break
			}
		}
	}
	return opt
}
