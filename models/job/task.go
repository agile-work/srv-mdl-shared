package job

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Task defines the struct of this object
type Task struct {
	ID               string                  `json:"id" sql:"id" pk:"true"`
	Code             string                  `json:"code" sql:"code"`
	Name             translation.Translation `json:"name" sql:"name" field:"jsonb" validate:"required"`
	Description      translation.Translation `json:"description" sql:"description" field:"jsonb"`
	JobCode          string                  `json:"job_code" sql:"job_code" fk:"true"`
	TaskSequence     int                     `json:"task_sequence" sql:"task_sequence"`
	ExecTimeout      int                     `json:"exec_timeout" sql:"exec_timeout"`
	Params           []Param                 `json:"parameters" sql:"parameters" field:"jsonb"`
	ParentCode       string                  `json:"parent_code" sql:"parent_code" fk:"true"`
	ExecAction       string                  `json:"exec_action" sql:"exec_action"`
	ExecAddress      string                  `json:"exec_address" sql:"exec_address"`
	ExecPayload      string                  `json:"exec_payload" sql:"exec_payload"`
	ActionOnFail     string                  `json:"action_on_fail" sql:"action_on_fail"`
	MaxRetryAttempts int                     `json:"max_retry_attempts" sql:"max_retry_attempts"`
	RollbackAction   string                  `json:"rollback_action" sql:"rollback_action"`
	RollbackAddress  string                  `json:"rollback_address" sql:"rollback_address"`
	RollbackPayload  string                  `json:"rollback_payload" sql:"rollback_payload"`
	CreatedBy        string                  `json:"created_by" sql:"created_by"`
	CreatedAt        time.Time               `json:"created_at" sql:"created_at"`
	UpdatedBy        string                  `json:"updated_by" sql:"updated_by"`
	UpdatedAt        time.Time               `json:"updated_at" sql:"updated_at"`
}

// Create persists the struct creating a new object in the database
func (t *Task) Create(trs *db.Transaction, columns ...string) error {
	id, err := db.InsertStructTx(trs.Tx, constants.TableCoreJobTasks, t, columns...)
	if err != nil {
		customerror.New(http.StatusInternalServerError, "task create", err.Error())
	}
	t.ID = id
	return nil
}

// Load defines only one object from the database
func (t *Task) Load() error {
	if err := db.SelectStruct(constants.TableCoreJobTasks, t, &db.Options{
		Conditions: builder.And(
			builder.Equal("job_code", t.JobCode),
			builder.Equal("code", t.Code),
		),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "task load", err.Error())
	}
	return nil
}

// Update updates object data in the database
func (t *Task) Update(trs *db.Transaction, columns []string, translations map[string]string) error {
	opt := &db.Options{Conditions: builder.And(
		builder.Equal("job_code", t.JobCode),
		builder.Equal("code", t.Code),
	)}

	if len(columns) > 0 {
		if err := db.UpdateStructTx(trs.Tx, constants.TableCoreJobTasks, t, opt, strings.Join(columns, ",")); err != nil {
			return customerror.New(http.StatusInternalServerError, "task update", err.Error())
		}
	}

	if len(translations) > 0 {
		statement := builder.Update(constants.TableCoreJobTasks)
		for col, val := range translations {
			statement.JSON(col, translation.FieldsRequestLanguageCode)
			jsonVal, _ := json.Marshal(val)
			statement.Values(jsonVal)
		}
		statement.Where(opt.Conditions)
		if _, err := trs.Query(statement); err != nil {
			return customerror.New(http.StatusInternalServerError, "task update", err.Error())
		}
	}

	return nil
}

// Delete deletes object from the database
func (t *Task) Delete(trs *db.Transaction) error {
	if err := db.DeleteStructTx(trs.Tx, constants.TableCoreJobTasks, &db.Options{
		Conditions: builder.And(
			builder.Equal("job_code", t.JobCode),
			builder.Equal("code", t.Code),
		),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "task delete", err.Error())
	}
	return nil
}

// Tasks defines the array struct of this object
type Tasks []Task

// LoadAll defines all instances from the object
func (t *Tasks) LoadAll(opt *db.Options) error {
	if err := db.SelectStruct(constants.TableCoreJobTasks, t, opt); err != nil {
		return customerror.New(http.StatusInternalServerError, "tasks load", err.Error())
	}
	return nil
}
