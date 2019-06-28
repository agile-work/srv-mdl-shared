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
	Filter     map[string]MetadataFilter `json:"filter,omitempty"`
	Pagination map[string]int            `json:"pagination,omitempty"`
	Order      []map[string]string       `json:"order,omitempty"`
	Columns    string                    `json:"columns,omitempty"`
	Group      string                    `json:"group,omitempty"`
}

// MetadataFilter defines filter metadata on the response
type MetadataFilter struct {
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// Load gets metadata from request
func (m *Metadata) Load(req *http.Request) error {
	metadataStr := req.URL.Query().Get("metadata")
	if metadataStr != "" {
		if err := json.Unmarshal([]byte(metadataStr), m); err != nil {
			return customerror.New(http.StatusBadRequest, "metadata load unmarshal", err.Error())
		}
	}
	return nil
}

// GenerateDBOptions convert the metadata in a db options
func (m *Metadata) GenerateDBOptions() *db.Options {
	opt := &db.Options{}
	if m.Columns != "" {
		opt.Columns = strings.Split(m.Columns, ",")
	}
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
	for column, filter := range m.Filter {
		var condition builder.Builder
		switch filter.Operator {
		case "=":
			condition = builder.Equal(column, filter.Value)
			break
		case "!=":
			condition = builder.NotEqual(column, filter.Value)
			break
		case ">":
			condition = builder.GreaterThen(column, filter.Value)
			break
		case ">=":
			condition = builder.GreaterOrEqual(column, filter.Value)
			break
		case "<":
			condition = builder.LowerThen(column, filter.Value)
			break
		case "<=":
			condition = builder.LowerOrEqual(column, filter.Value)
			break
		}
		opt.AddCondition(condition)
	}
	return opt
}
