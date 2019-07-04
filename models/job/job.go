package job

import (
	"errors"
	"net/http"
	"time"

	"github.com/agile-work/srv-mdl-shared/models/customerror"
	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/constants"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// Job defines the struct of this object
type Job struct {
	ID          string                  `json:"id" sql:"id" pk:"true"`
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

// ViewJobInstance defines the struct of this object
type ViewJobInstance struct {
	ID          string    `json:"id" sql:"id" pk:"true"`
	JobID       string    `json:"job_id" sql:"job_id"`
	Code        string    `json:"code" sql:"code"`
	Name        string    `json:"name" sql:"name"`
	Description string    `json:"description" sql:"description"`
	JobType     string    `json:"job_type" sql:"job_type"`
	ExecTimeout int       `json:"exec_timeout" sql:"exec_timeout"`
	Params      []Param   `json:"parameters" sql:"parameters" field:"jsonb"`
	StartAt     time.Time `json:"start_at" sql:"start_at"`
	FinishAt    time.Time `json:"finish_at" sql:"finish_at"`
	Status      string    `json:"status" sql:"status"`
	CreatedBy   string    `json:"created_by" sql:"created_by"`
	CreatedAt   time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy   string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at" sql:"updated_at"`
}

// JobTask defines the struct of this object
type JobTask struct {
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

// ViewJobTaskInstance defines the struct of this object
type ViewJobTaskInstance struct {
	ID               string    `json:"id" sql:"id" pk:"true"`
	JobID            string    `json:"job_id" sql:"job_id"`
	JobInstanceID    string    `json:"job_instance_id" sql:"job_instance_id"`
	TaskID           string    `json:"task_id" sql:"task_id"`
	Code             string    `json:"code" sql:"code"`
	Name             string    `json:"name" sql:"name"`
	Description      string    `json:"description" sql:"description"`
	TaskSequence     int       `json:"task_sequence" sql:"task_sequence"`
	ExecTimeout      int       `json:"exec_timeout" sql:"exec_timeout"`
	Params           []Param   `json:"parameters" sql:"parameters" field:"jsonb"`
	ParentID         string    `json:"parent_id" sql:"parent_id" fk:"true"`
	ExecAction       string    `json:"exec_action" sql:"exec_action"`
	ExecAddress      string    `json:"exec_address" sql:"exec_address"`
	ExecPayload      string    `json:"exec_payload" sql:"exec_payload"`
	ExecResponse     string    `json:"exec_response" sql:"exec_response"`
	ActionOnFail     string    `json:"action_on_fail" sql:"action_on_fail"`
	MaxRetryAttempts int       `json:"max_retry_attempts" sql:"max_retry_attempts"`
	RollbackAction   string    `json:"rollback_action" sql:"rollback_action"`
	RollbackAddress  string    `json:"rollback_address" sql:"rollback_address"`
	RollbackPayload  string    `json:"rollback_payload" sql:"rollback_payload"`
	RollbackResponse string    `json:"rollback_response" sql:"rollback_response"`
	StartAt          time.Time `json:"start_at" sql:"start_at"`
	FinishAt         time.Time `json:"finish_at" sql:"finish_at"`
	Status           string    `json:"status" sql:"status"`
	CreatedBy        string    `json:"created_by" sql:"created_by"`
	CreatedAt        time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy        string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt        time.Time `json:"updated_at" sql:"updated_at"`
}

// JobFollowers defines the struct of this object
type JobFollowers struct {
	ID           string    `json:"id" sql:"id" pk:"true"`
	JobID        string    `json:"job_id" sql:"job_id" fk:"true"`
	Name         string    `json:"name" sql:"name"`
	LanguageCode string    `json:"language_code" sql:"language_code"`
	FollowerID   string    `json:"follower_id" sql:"follower_id"`
	FollowerType string    `json:"follower_type" sql:"follower_type"`
	Active       bool      `json:"active" sql:"active"`
	CreatedBy    string    `json:"created_by" sql:"created_by"`
	CreatedAt    time.Time `json:"created_at" sql:"created_at"`
	// UpdatedBy    string    `json:"updated_by" sql:"updated_by"`
	// UpdatedAt    time.Time `json:"updated_at" sql:"updated_at"`
}

// ViewFollowerAvailable defines the struct of this object
type ViewFollowerAvailable struct {
	ID                    string    `json:"id" sql:"id" pk:"true"`
	Name                  string    `json:"name" sql:"name"`
	LanguageCode          string    `json:"language_code" sql:"language_code"`
	FollowerAvailableType string    `json:"ug_type" sql:"ug_type"`
	Active                bool      `json:"active" sql:"active"`
	CreatedBy             string    `json:"created_by" sql:"created_by"`
	CreatedAt             time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy             string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt             time.Time `json:"updated_at" sql:"updated_at"`
}

// Param defines the struct of this object
type Param struct {
	Type      string `json:"type"`
	Reference string `json:"ref"`
	Field     string `json:"field"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// JobInstance defines the struct of this object
type JobInstance struct {
	ID          string    `json:"id" sql:"id" pk:"true"`
	JobCode     string    `json:"job_code" sql:"job_code" fk:"true"`
	ServiceID   string    `json:"service_id" sql:"service_id" fk:"true"`
	ExecTimeout int       `json:"exec_timeout" sql:"exec_timeout"`
	Params      []Param   `json:"parameters" sql:"parameters" field:"jsonb"`
	Status      string    `json:"status" sql:"status"`
	StartAt     time.Time `json:"start_at" sql:"start_at"`
	FinishAt    time.Time `json:"finish_at" sql:"finish_at"`
	CreatedBy   string    `json:"created_by" sql:"created_by"`
	CreatedAt   time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy   string    `json:"updated_by" sql:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at" sql:"updated_at"`
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

// Create create a new job instance
func (j *JobInstance) Create(trs *db.Transaction, owner string, code string, params map[string]interface{}) (string, error) {
	job := Job{
		Code: code,
	}

	if err := job.Load(); err != nil {
		return "", err
	}

	if err := j.fillParameters(job.Params, params); err != nil {
		return "", err
	}

	date := time.Now()
	j.JobCode = job.Code
	j.ExecTimeout = job.ExecTimeout
	j.Status = constants.JobStatusCreating
	j.CreatedBy = owner
	j.CreatedAt = date
	j.UpdatedBy = owner
	j.UpdatedAt = date

	return db.InsertStructTx(trs.Tx, constants.TableCoreJobInstances, j)
}

// fillParameters fill parameters with values
func (j *JobInstance) fillParameters(params []Param, values map[string]interface{}) error {
	if len(params) != len(values) {
		return errors.New("the number of parameters can not be different from the number of values")
	}

	for _, param := range params {
		if value, ok := values[param.Key]; ok {
			param.Value = value.(string)
			j.Params = append(j.Params, param)
		} else {
			return errors.New("parameter invalid")
		}
	}

	return nil
}
