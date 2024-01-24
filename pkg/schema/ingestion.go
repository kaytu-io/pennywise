package schema

type IngestionJobStatus string

const (
	IngestionJobCreated    IngestionJobStatus = "CREATED"
	IngestionJobInProgress IngestionJobStatus = "IN_PROGRESS"
	IngestionJobSucceeded  IngestionJobStatus = "SUCCEEDED"
	IngestionJobFailed     IngestionJobStatus = "FAILED"
)

type IngestionJob struct {
	ID       int32              `json:"id"`
	Provider string             `json:"provider"`
	Location string             `json:"location"`
	Service  string             `json:"service"`
	Status   IngestionJobStatus `json:"status"`
	ErrorMsg string             `json:"error_msg"`
}
