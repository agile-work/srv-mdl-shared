package shared

import (
	"errors"
	"net/http"
	"time"

	"github.com/agile-work/srv-mdl-core/models"
	shared "github.com/agile-work/srv-shared"
	"github.com/agile-work/srv-shared/sql-builder/builder"
	"github.com/agile-work/srv-shared/sql-builder/db"
)

// CreateJobInstance create a new job instance
func CreateJobInstance(ownerID string, code string, params map[string]interface{}) (string, error) {
	jobTable := shared.TableCoreJobs
	conditions := builder.Equal("code", code)
	job := models.Job{}

	err := db.LoadStruct(jobTable, &job, conditions)
	if err != nil {
		return "", err
	}

	jobInstanceParams, err := fillParameters(job.Params, params)
	if err != nil {
		return "", err
	}

	jobInstanceTable := shared.TableCoreJobInstances
	date := time.Now()
	jobInstance := models.JobInstance{
		JobID:       job.ID,
		Code:        job.Code,
		ExecTimeout: job.ExecTimeout,
		Params:      jobInstanceParams,
		Status:      "creating",
		CreatedBy:   ownerID,
		CreatedAt:   date,
		UpdatedBy:   ownerID,
		UpdatedAt:   date,
	}

	return db.InsertStruct(jobInstanceTable, jobInstance)
}

// fillParameters fill parameters with values
func fillParameters(params []models.Param, values map[string]interface{}) ([]models.Param, error) {
	if len(params) != len(values) {
		return nil, errors.New("The number of parameters can not be different from the number of values")
	}

	result := []models.Param{}
	for _, param := range params {
		if value, ok := values[param.Key]; ok {
			param.Value = value.(string)
			result = append(result, param)
		} else {
			return nil, errors.New("Parameter invalid")
		}
	}

	return result, nil
}

// ExecJob execute a new job
func ExecJob(ownerID string, code string, params map[string]interface{}) (string, *Response) {
	response := &Response{
		Code: http.StatusOK,
	}

	id, err := CreateJobInstance(ownerID, code, params)
	if err != nil {
		response.Code = http.StatusInternalServerError
		response.Errors = append(response.Errors, NewResponseError(ErrorJobExecution, "Job execution", err.Error()))

		return "", response
	}

	return id, nil
}
