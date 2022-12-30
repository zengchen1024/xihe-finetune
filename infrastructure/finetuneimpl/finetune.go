package finetuneimpl

import (
	"strings"

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

func NewFinetune(cfg *Config) finetune.Finetune {
	return finetuneImpl{
		cli: newClient(cfg),
	}
}

type finetuneImpl struct {
	cli client
}

func (impl finetuneImpl) Create(user domain.Account, t *domain.Finetune) (
	info domain.JobInfo, err error,
) {
	p := t.Param
	req := createRequest{
		User:  user.Account(),
		Name:  t.Id + "_" + t.Name.FinetuneName(),
		Task:  p.Task(),
		Model: p.Model(),
	}

	if hp := p.Hypeparameters(); len(hp) > 0 {
		i := 0
		items := make([]hyperparameter, len(hp))
		for k, v := range hp {
			items[i] = hyperparameter{
				Name:  k,
				Value: v,
			}

			i++
		}

		req.Parameters = items
	}

	info.JobId, err = impl.cli.createJob(&req)

	return
}

func (impl finetuneImpl) Delete(jobId string) error {
	return impl.cli.deleteJob(jobId)
}

func (impl finetuneImpl) GetDetail(jobId string) (r domain.JobDetail, err error) {
	v, err := impl.cli.getJob(jobId)
	if err != nil {
		return
	}

	if status, ok := statusMap[strings.ToLower(v.Phase)]; ok {
		r.Status = status
	} else {
		r.Status = domain.FinetuneStatusFailed
	}

	// convert millisecond to second
	r.Duration = v.Runtime / 1000

	return
}

func (impl finetuneImpl) Terminate(jobId string) error {
	return impl.cli.terminateJob(jobId)
}

func (impl finetuneImpl) GetLogDownloadURL(jobId string) (string, error) {
	return impl.cli.getLogURL(jobId)
}
