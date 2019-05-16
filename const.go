package shared

const (
	// ErrorParsingRequest unable to unmarshall json to struct
	ErrorParsingRequest string = "001-ErrorParsingRequest"
	// ErrorInsertingRecord unable to insert record on database
	ErrorInsertingRecord string = "002-ErrorInsertingRecord"
	// ErrorReturningData unable to return data
	ErrorReturningData string = "003-ErrorReturningData"
	// ErrorDeletingData unable to return data
	ErrorDeletingData string = "004-ErrorDeletingData"
	// ErrorLoadingData unable to load data
	ErrorLoadingData string = "005-ErrorLoadingData"
	// ErrorLogin unable to login user
	ErrorLogin string = "006-ErrorLoginUser"
	// ErrorJobExecution unable to execute job
	ErrorJobExecution string = "007-ErrorJobExecution"
	// JobSystemCreateSchema exec tasks to create a schema
	JobSystemCreateSchema string = "job_system_create_schema"
)
