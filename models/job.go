package models

import "time"

// Job defines the struct of this object
type Job struct {
	ID            string    `json:"id" sql:"id" pk:"true"`
	Code          string    `json:"code" sql:"code"`
	Name          string    `json:"name" table:"core_translations" alias:"core_translations_name" sql:"value" on:"core_translations_name.structure_id = core_jobs.id and core_translations_name.structure_field = 'name'"`
	Description   string    `json:"description" table:"core_translations" alias:"core_translations_description" sql:"value" on:"core_translations_description.structure_id = core_jobs.id and core_translations_description.structure_field = 'description'"`
	JobType       string    `json:"job_type" sql:"job_type"`
	ExecTimeout   int       `json:"exec_timeout" sql:"exec_timeout"`
	Params        []Param   `json:"parameters" sql:"parameters" field:"jsonb"`
	Active        bool      `json:"active" sql:"active"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedByUser *User     `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_jobs.created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedByUser *User     `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_jobs.updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}

// JobTask defines the struct of this object
type JobTask struct {
	ID               string    `json:"id" sql:"id" pk:"true"`
	Code             string    `json:"code" sql:"code"`
	Name             string    `json:"name" table:"core_translations" alias:"core_translations_name" sql:"value" on:"core_translations_name.structure_id = core_job_tasks.id and core_translations_name.structure_field = 'name'"`
	Description      string    `json:"description" table:"core_translations" alias:"core_translations_description" sql:"value" on:"core_translations_description.structure_id = core_job_tasks.id and core_translations_description.structure_field = 'description'"`
	JobID            string    `json:"job_id" sql:"job_id" fk:"true"`
	TaskSequence     int       `json:"task_sequence" sql:"task_sequence"`
	ExecTimeout      int       `json:"exec_timeout" sql:"exec_timeout"`
	Params           []Param   `json:"parameters" sql:"parameters" field:"jsonb"`
	ParentID         string    `json:"parent_id" sql:"parent_id" fk:"true"`
	ExecAction       string    `json:"exec_action" sql:"exec_action"`
	ExecAddress      string    `json:"exec_address" sql:"exec_address"`
	ExecPayload      string    `json:"exec_payload" sql:"exec_payload"`
	ActionOnFail     string    `json:"action_on_fail" sql:"action_on_fail"`
	MaxRetryAttempts int       `json:"max_retry_attempts" sql:"max_retry_attempts"`
	RollbackAction   string    `json:"rollback_action" sql:"rollback_action"`
	RollbackAddress  string    `json:"rollback_address" sql:"rollback_address"`
	RollbackPayload  string    `json:"rollback_payload" sql:"rollback_payload"`
	CreatedBy        string    `json:"created_by" sql:"created_by"`
	CreatedByUser    *User     `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_job_tasks.created_by"`
	CreatedAt        time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy        string    `json:"updated_by" sql:"updated_by"`
	UpdatedByUser    *User     `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_job_tasks.updated_by"`
	UpdatedAt        time.Time `json:"updated_at" sql:"updated_at"`
}

// JobFollowers defines the struct of this object
type JobFollowers struct {
	ID            string    `json:"id" sql:"id" pk:"true"`
	JobID         string    `json:"job_id" sql:"job_id" fk:"true"`
	Name          string    `json:"name" sql:"name"`
	LanguageCode  string    `json:"language_code" sql:"language_code"`
	FollowerType  string    `json:"follower_type" sql:"follower_type"`
	Active        bool      `json:"active" sql:"active"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedByUser *User     `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_v_job_followers.created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedByUser *User     `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_v_job_followers.updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}

// ViewFollowerAvailable defines the struct of this object
type ViewFollowerAvailable struct {
	ID                    string    `json:"id" sql:"id" pk:"true"`
	Name                  string    `json:"name" sql:"name"`
	LanguageCode          string    `json:"language_code" sql:"language_code"`
	FollowerAvailableType string    `json:"ug_type" sql:"ug_type"`
	Active                bool      `json:"active" sql:"active"`
	CreatedBy             string    `json:"created_by" sql:"created_by"`
	CreatedByUser         *User     `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_v_users_and_groups.created_by"`
	CreatedAt             time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy             string    `json:"updated_by" sql:"updated_by"`
	UpdatedByUser         *User     `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_v_users_and_groups.updated_by"`
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
	ID            string    `json:"id" sql:"id" pk:"true"`
	JobID         string    `json:"job_id" sql:"job_id" fk:"true"`
	ServiceID     string    `json:"service_id" sql:"service_id" fk:"true"`
	Code          string    `json:"code" sql:"code"`
	ExecTimeout   int       `json:"exec_timeout" sql:"exec_timeout"`
	Params        []Param   `json:"parameters" sql:"parameters" field:"jsonb"`
	Status        string    `json:"status" sql:"status"`
	CreatedBy     string    `json:"created_by" sql:"created_by"`
	CreatedByUser *User     `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_jobs.created_by"`
	CreatedAt     time.Time `json:"created_at" sql:"created_at"`
	UpdatedBy     string    `json:"updated_by" sql:"updated_by"`
	UpdatedByUser *User     `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_jobs.updated_by"`
	UpdatedAt     time.Time `json:"updated_at" sql:"updated_at"`
}
