package app

import (
	"errors"

	"github.com/opensourceways/xihe-finetune/domain"
	"github.com/opensourceways/xihe-finetune/domain/finetune"
	"github.com/opensourceways/xihe-finetune/domain/watch"
)

type FinetuneCreateCmd struct {
	User domain.Account
	Id   string

	domain.FinetuneConfig
}

func (cmd *FinetuneCreateCmd) Validate() error {
	b := cmd.User != nil &&
		cmd.Name != nil &&
		cmd.Param != nil

	if !b {
		return errors.New("invalid cmd of creating finetune")
	}

	return nil
}

type JobInfoDTO struct {
	JobId string `json:"job_id"`
}

type LogURLDTO struct {
	URL string `json:"url"`
}

type FinetuneService interface {
	Create(cmd *FinetuneCreateCmd) (JobInfoDTO, string, error)
	Delete(jobId string) error
	Terminate(jobId string) error
	GetLogDownloadURL(jobId string) (LogURLDTO, error)
}

func NewFinetuneService(
	ts finetune.Finetune,
	ws watch.WatchService,
) FinetuneService {
	return &finetuneService{
		ts: ts,
		ws: ws,
	}
}

type finetuneService struct {
	ts finetune.Finetune
	ws watch.WatchService
}

func (s *finetuneService) Create(cmd *FinetuneCreateCmd) (JobInfoDTO, string, error) {
	dto := JobInfoDTO{}
	code := ""

	f := func(info *watch.FinetuneInfo) error {
		v, err := s.ts.Create(cmd.User, &cmd.FinetuneConfig)
		if err != nil {
			// TODO check parameter error
			return err
		}

		dto.JobId = v.JobId

		*info = watch.FinetuneInfo{
			User:    cmd.User,
			Id:      cmd.Id,
			JobInfo: v,
		}

		return nil
	}

	err := s.ws.ApplyWatch(f)

	return dto, code, err
}

func (s *finetuneService) Delete(jobId string) error {
	return s.ts.Delete(jobId)
}

func (s *finetuneService) Terminate(jobId string) error {
	return s.ts.Terminate(jobId)
}

func (s *finetuneService) GetLogDownloadURL(jobId string) (r LogURLDTO, err error) {
	r.URL, err = s.ts.GetLogDownloadURL(jobId)

	return
}
