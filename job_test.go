package shared

import (
	"testing"

	"github.com/agile-work/srv-mdl-shared/models"
	"github.com/agile-work/srv-shared/sql-builder/db"
	"github.com/stretchr/testify/suite"
)

type JobTestSuite struct {
	suite.Suite
}

func (suite *JobTestSuite) SetupTest() {
	db.Connect(
		"cryo.cdnm8viilrat.us-east-2.rds-preview.amazonaws.com",
		5432,
		"cryoadmin",
		"x3FhcrWDxnxCq9p",
		"cryo",
		false,
	)
}

func (suite *JobTestSuite) Test00001CreateJobInstance() {
	data := map[string]interface{}{
		"schema_name": "Contrato",
		"schema_id":   "2394234-234237426937-23497234",
	}

	id, err := CreateJobInstance("307e481c-69c5-11e9-96a0-06ea2c43bb20", "job_system_create_schema", data)

	if id != "" || err != nil {

	}
}

// func (suite *JobTestSuite) Test00002LoadAvaibleJobInstances() {
// 	data := map[string]interface{}{
// 		"schema_name": "Contrato",
// 		"schema_id":   "2394234-234237426937-23497234",
// 	}

// 	id, err := LoadAvaibleJobInstances("307e481c-69c5-11e9-96a0-06ea2c43bb20", "job_system_create_schema", data)

// 	if id != "" || err != nil {

// 	}
// }

func (suite *JobTestSuite) Test00003SetAvaibleJobInstancesToInQueue() {
	jobInstances := []models.JobInstance{}
	LoadAvaibleJobInstances(&jobInstances)
	SetAvaibleJobInstancesToInQueue(&jobInstances)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestJobSuite(t *testing.T) {
	suite.Run(t, new(JobTestSuite))
}
