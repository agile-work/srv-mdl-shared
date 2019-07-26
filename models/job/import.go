package job

import (
	"io/ioutil"
	"time"

	"github.com/agile-work/srv-shared/constants"

	"github.com/agile-work/srv-shared/sql-builder/db"
	"github.com/tidwall/gjson"
)

func importJSONTasks(trs *db.Transaction, id, path string) error {
	jsonByte, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	tasks := gjson.GetBytes(jsonByte, "tasks")
	for _, task := range tasks.Array() {
		taskID := db.UUID()
		date := time.Now()
		instanceTask := &InstanceTask{
			ID:               taskID,
			JobInstanceID:    id,
			TaskCode:         taskID,
			TaskSequence:     int(task.Get("sequence").Int()),
			ExecTimeout:      60,
			ExecAction:       task.Get("exec_action").String(),
			ExecAddress:      task.Get("exec_address").String(),
			ExecPayload:      task.Get("exec_payload").String(),
			ActionOnFail:     constants.OnFailRetryAndCancel,
			MaxRetryAttempts: 2,
			Status:           constants.JobStatusCreated,
			CreatedBy:        "admin",
			CreatedAt:        date,
			UpdatedBy:        "admin",
			UpdatedAt:        date,
		}
		if err := instanceTask.Create(trs); err != nil {
			return err
		}
	}
	return nil
}
