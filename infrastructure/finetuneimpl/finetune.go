package finetuneimpl

import (
	"github.com/opensourceways/xihe-finetune/domain"
	"github.com/opensourceways/xihe-finetune/domain/finetune"
)

var statusMap = map[string]domain.FinetuneStatus{
	"failed":      domain.FinetuneStatusFailed,
	"pending":     domain.FinetuneStatusPending,
	"running":     domain.FinetuneStatusRunning,
	"creating":    domain.FinetuneStatusCreating,
	"abnormal":    domain.FinetuneStatusAbnormal,
	"completed":   domain.FinetuneStatusCompleted,
	"terminated":  domain.FinetuneStatusTerminated,
	"terminating": domain.FinetuneStatusTerminating,
}

func NewFinetune(cfg *Config) (finetune.Finetune, error) {
	return finetuneImpl{}, nil
}

type finetuneImpl struct {
}

func (impl finetuneImpl) Create(user domain.Account, t *domain.FinetuneConfig) (
	info domain.JobInfo, err error,
) {
	return
}

func (impl finetuneImpl) Delete(jobId string) error {
	return nil
}

func (impl finetuneImpl) GetDetail(jobId string) (r domain.JobDetail, err error) {
	return
}

func (impl finetuneImpl) Terminate(jobId string) error {
	return nil
}

func (impl finetuneImpl) GetLogDownloadURL(jobId string) (string, error) {
	return "", nil
}
