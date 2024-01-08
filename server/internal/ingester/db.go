package ingester

import (
	"fmt"
	"strings"
)

func (s Scheduler) MakeJob(provider, service, region string) (int64, error) {
	invalid := true
	for _, p := range validProviders {
		if p == strings.ToLower(provider) {
			invalid = false
			break
		}
	}
	if invalid {
		return 0, fmt.Errorf("invalid provider")
	}
	if region == "" {
		region = "all"
	}
	q := `
		INSERT INTO ingestion_jobs (provider, location, service, status)
		VALUES (?, ?, ?, ?)
	`

	result, err := s.db.Exec(q, provider, region, service, IngestionJobCreated)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, err
}

func (s Scheduler) GetJobs(status, provider, service, location string) ([]IngestionJob, error) {
	q := "SELECT id, provider, location, service, status FROM ingestion_jobs"
	var args []any
	if status != "" || provider != "" || service != "" || location != "" {
		q += " WHERE "
		var conditions []string
		if status != "" {
			conditions = append(conditions, "status = ?")
			args = append(args, status)
		}
		if provider != "" {
			conditions = append(conditions, "provider = ?")
			args = append(args, provider)
		}
		if service != "" {
			conditions = append(conditions, "service = ?")
			args = append(args, service)
		}
		if location != "" {
			conditions = append(conditions, "location = ?")
			args = append(args, location)
		}

		q += strings.Join(conditions, " AND ")
	}

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	var jobs []IngestionJob
	for rows.Next() {
		var job IngestionJob
		err := rows.Scan(&job.ID, &job.Provider, &job.Location, &job.Service, &job.Status)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s Scheduler) GetJobById(id int32) (*IngestionJob, error) {
	q := "SELECT id, provider, location, service, status FROM ingestion_jobs WHERE id = ? LIMIT 1"

	rows, err := s.db.Query(q, id)
	if err != nil {
		return nil, err
	}

	var job IngestionJob
	_ = rows.Next()
	err = rows.Scan(&job.ID, &job.Provider, &job.Location, &job.Service, &job.Status)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (s Scheduler) getCreatedJobs() ([]IngestionJob, error) {
	q := `
		SELECT id, provider, location, service, status
		FROM ingestion_jobs
		WHERE status IN (?)
	`
	rows, err := s.db.Query(q, IngestionJobCreated)
	if err != nil {
		return nil, err
	}

	var jobs []IngestionJob
	for rows.Next() {
		var job IngestionJob
		err := rows.Scan(&job.ID, &job.Provider, &job.Location, &job.Service, &job.Status)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (s Scheduler) updateJobStatus(id int32, status IngestionJobStatus, errorMsg string) error {
	q := `UPDATE ingestion_jobs SET status = ?, error_msg = ? WHERE id = ?`

	_, err := s.db.Exec(q, status, errorMsg, id)
	if err != nil {
		return err
	}
	return err
}
