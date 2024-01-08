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
	ID       int32              `json:"id"`
	Provider string             `json:"provider"`
	Location string             `json:"location"`
	Service  string             `json:"service"`
	Status   IngestionJobStatus `json:"status"`
	ErrorMsg string             `json:"error_msg"`
}

func (s Scheduler) RunIngestionJobScheduler() {
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
