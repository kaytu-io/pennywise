package ingester

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kaytu-io/pennywise/server/aws"
	"github.com/kaytu-io/pennywise/server/azurerm"
	"github.com/kaytu-io/pennywise/server/internal/backend"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type IngestionJobStatus string

const (
	IngestionJobCreated    IngestionJobStatus = "CREATED"
	IngestionJobInProgress IngestionJobStatus = "IN_PROGRESS"
	IngestionJobSucceeded  IngestionJobStatus = "SUCCEEDED"
	IngestionJobFailed     IngestionJobStatus = "FAILED"
)

var validProviders = []string{
	"aws",
	"azure",
}

type Scheduler struct {
	backend backend.Backend
	logger  *zap.Logger
	db      *sql.DB
}

func NewScheduler(b backend.Backend, logger *zap.Logger, db *sql.DB) Scheduler {
	return Scheduler{
		backend: b,
		logger:  logger,
		db:      db,
	}
}

type IngestionJob struct {
	ID       int32
	Provider string
	Location string
	Service  string
	Status   IngestionJobStatus
	ErrorMsg string
}

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
	if status != "" || provider != "" || service != "" || location != "" {
		q += " WHERE "
		var conditions []string
		if status != "" {
			conditions = append(conditions, fmt.Sprintf("status = '%s'", status))
		}
		if provider != "" {
			conditions = append(conditions, fmt.Sprintf("provider = '%s'", provider))
		}
		if service != "" {
			conditions = append(conditions, fmt.Sprintf("service = '%s'", service))
		}
		if location != "" {
			conditions = append(conditions, fmt.Sprintf("location = '%s'", location))
		}

		q += strings.Join(conditions, " AND ")
	}

	rows, err := s.db.Query(q)
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

func (s Scheduler) InsureScheduler() {
	go func() {
		err := s.runIngestionJobScheduler()
		if err != nil {
			s.logger.Error(err.Error())
			s.InsureScheduler()
		}
	}()
}

func (s Scheduler) runIngestionJobScheduler() error {
	s.logger.Info("IngestionJob scheduler started")

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for ; ; <-t.C {
		err := s.updateIngestion()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}
}

func (s Scheduler) updateIngestion() error {
	jobs, err := s.getCreatedJobs()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		go func() {
			s.logger.Info(fmt.Sprintf("ingestion job started: %v", job))
			err := s.updateJobStatus(job.ID, IngestionJobInProgress, "")
			if err != nil {
				s.logger.Error(fmt.Sprintf("updating job status failed for job %v with error: %s", job.ID, err.Error()))
			}
			err = s.runJob(job)
			if err != nil {
				err2 := s.updateJobStatus(job.ID, IngestionJobFailed, err.Error())
				if err2 != nil {
					s.logger.Error(fmt.Sprintf("ingestion job failed for job %v with error: %s", job.ID, err.Error()))
					s.logger.Error(fmt.Sprintf("updating job status failed for job %v with error: %s", job.ID, err2.Error()))
					return
				}
				s.logger.Error(fmt.Sprintf("ingestion job failed for job %v with error: %s", job.ID, err.Error()))
				return
			}
			s.updateJobStatus(job.ID, IngestionJobSucceeded, "")
			s.logger.Info(fmt.Sprintf("ingestion job succeeded: %v", job))
		}()
	}
	return nil
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

func (s Scheduler) runJob(job IngestionJob) error {
	if job.Location == "all" {
		job.Location = ""
	}
	if strings.ToLower(job.Provider) == "aws" {
		ingester, err := aws.NewIngester(job.Service, job.Location)
		if err != nil {
			return err
		}
		err = IngestPricing(context.Background(), s.backend, ingester)
		if err != nil {
			return err
		}
	} else if strings.ToLower(job.Provider) == "azure" || strings.ToLower(job.Provider) == "azurerm" {
		ingester, err := azurerm.NewIngester(job.Service, job.Location)
		if err != nil {
			return err
		}
		err = IngestPricing(context.Background(), s.backend, ingester)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid provider")
	}
	return nil
}

func (s Scheduler) updateJobStatus(id int32, status IngestionJobStatus, errorMsg string) error {
	q := `
		UPDATE ingestion_jobs SET status = ?, error_msg = ? WHERE id = ?
	`

	_, err := s.db.Exec(q, status, errorMsg, id)
	if err != nil {
		return err
	}
	return err
}
