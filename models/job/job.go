package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/agile-work/srv-mdl-core/models/dataset"
	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Job defines the struct of this object
type Job struct {
	ID          string                  `json:"id" sql:"id"`
	Code        string                  `json:"code" sql:"code" updatable:"false" validate:"required"`
	Name        translation.Translation `json:"name" sql:"name" field:"jsonb" validate:"required"`
	Description translation.Translation `json:"description" sql:"description" field:"jsonb" validate:"required"`
	JobType     string                  `json:"job_type" sql:"job_type"`
	ExecTimeout int                     `json:"exec_timeout" sql:"exec_timeout"`
	Params      []Param                 `json:"parameters" sql:"parameters" field:"jsonb"`
	Active      bool                    `json:"active" sql:"active"`
	CreatedBy   string                  `json:"created_by" sql:"created_by"`
	CreatedAt   time.Time               `json:"created_at" sql:"created_at"`
	UpdatedBy   string                  `json:"updated_by" sql:"updated_by"`
	UpdatedAt   time.Time               `json:"updated_at" sql:"updated_at"`
}

// Create persists the struct creating a new object in the database
func (j *Job) Create(trs *db.Transaction, columns ...string) error {
	id, err := db.InsertStructTx(trs.Tx, constants.TableCoreConfigLanguages, j, columns...)
	if err != nil {
		customerror.New(http.StatusInternalServerError, "language create", err.Error())
	}
	j.ID = id
	return nil
}

// Load defines only one object from the database
func (j *Job) Load() error {
	if err := db.SelectStruct(constants.TableCoreJobs, j, &db.Options{
		Conditions: builder.Equal("code", j.Code),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "job load", err.Error())
	}
	return nil
}

// Update updates object data in the database
func (j *Job) Update(trs *db.Transaction, columns []string, translations map[string]string) error {
	opt := &db.Options{Conditions: builder.Equal("code", j.Code)}

	if len(columns) > 0 {
		if err := db.UpdateStructTx(trs.Tx, constants.TableCoreJobs, j, opt, strings.Join(columns, ",")); err != nil {
			return customerror.New(http.StatusInternalServerError, "job update", err.Error())
		}
	}

	if len(translations) > 0 {
		statement := builder.Update(constants.TableCoreJobs)
		for col, val := range translations {
			statement.JSON(col, translation.FieldsRequestLanguageCode)
			jsonVal, _ := json.Marshal(val)
			statement.Values(jsonVal)
		}
		statement.Where(opt.Conditions)
		if _, err := trs.Query(statement); err != nil {
			return customerror.New(http.StatusInternalServerError, "job update", err.Error())
		}
	}

	return nil
}

// Delete deletes object from the database
func (j *Job) Delete(trs *db.Transaction) error {
	if err := db.DeleteStructTx(trs.Tx, constants.TableCoreJobs, &db.Options{
		Conditions: builder.Equal("code", j.Code),
	}); err != nil {
		return customerror.New(http.StatusInternalServerError, "job delete", err.Error())
	}

	ds := &dataset.Dataset{
		Code: fmt.Sprintf("ds_job_%s", j.Code),
	}

	return ds.Delete(trs)
}

// Jobs defines the array struct of this object
type Jobs []Job

// LoadAll defines all instances from the object
func (t *Jobs) LoadAll(opt *db.Options) error {
	if err := db.SelectStruct(constants.TableCoreJobs, t, opt); err != nil {
		return customerror.New(http.StatusInternalServerError, "jobs load", err.Error())
	}
	return nil
}

// Instance defines the struct of this object
type Instance struct {
	ID                     string                 `json:"id" sql:"id"`
	JobCode                string                 `json:"job_code" sql:"job_code"`
	ServiceID              string                 `json:"service_id" sql:"service_id"`
	ExecTimeout            int                    `json:"exec_timeout" sql:"exec_timeout"`
	Params                 []Param                `json:"parameters" sql:"parameters" field:"jsonb"`
	Results                map[string]interface{} `json:"results" sql:"results" field:"jsonb"`
	WorkflowStepInstanceID string                 `json:"bpm_step_instance_id" sql:"bpm_step_instance_id"`
	WorkflowStepActionCode string                 `json:"bpm_step_action_code" sql:"bpm_step_action_code"`
	Status                 string                 `json:"status" sql:"status"`
	StartAt                time.Time              `json:"start_at" sql:"start_at"`
	FinishAt               time.Time              `json:"finish_at" sql:"finish_at"`
	CreatedBy              string                 `json:"created_by" sql:"created_by"`
	CreatedAt              time.Time              `json:"created_at" sql:"created_at"`
	UpdatedBy              string                 `json:"updated_by" sql:"updated_by"`
	UpdatedAt              time.Time              `json:"updated_at" sql:"updated_at"`
}

// Create create a new job instance
func (i *Instance) Create(trs *db.Transaction, owner string, code string, params map[string]interface{}) (string, error) {
	job := Job{
		Code: code,
	}

	if err := job.Load(); err != nil {
		return "", err
	}

	if err := i.fillParameters(job.Params, params); err != nil {
		return "", err
	}

	date := time.Now()
	i.JobCode = job.Code
	i.ExecTimeout = job.ExecTimeout
	i.Status = constants.JobStatusCreating
	i.CreatedBy = owner
	i.CreatedAt = date
	i.UpdatedBy = owner
	i.UpdatedAt = date

	return db.InsertStructTx(trs.Tx, constants.TableCoreJobInstances, i)
}

// fillParameters fill parameters with values
func (i *Instance) fillParameters(params []Param, values map[string]interface{}) error {
	if len(params) != len(values) {
		return errors.New("the number of parameters can not be different from the number of values")
	}

	for _, param := range params {
		if value, ok := values[param.Key]; ok {
			param.Value = value.(string)
			i.Params = append(i.Params, param)
		} else {
			return errors.New("parameter invalid")
		}
	}

	return nil
}

// CreateFromJSON create a new job instance based on a json file
func (i *Instance) CreateFromJSON(trs *db.Transaction, owner, code, definitionPath string, timeout int, params map[string]interface{}) error {
	date := time.Now()
	i.JobCode = code
	i.ExecTimeout = timeout
	i.Status = constants.JobStatusCreating
	i.CreatedBy = owner
	i.CreatedAt = date
	i.UpdatedBy = owner
	i.UpdatedAt = date

	id, err := db.InsertStructTx(trs.Tx, constants.TableCoreJobInstances, i)
	if err != nil {
		return err
	}

	if err := importJSONTasks(trs, id, code, definitionPath); err != nil {
		return err
	}
	return nil
}

// Param defines the struct of this object
type Param struct {
	Type      string `json:"type"`
	Reference string `json:"ref"`
	Field     string `json:"field"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}
