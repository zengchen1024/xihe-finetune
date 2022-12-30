package finetune

import "github.com/opensourceways/xihe-finetune/domain"

type Finetune interface {
	Create(domain.Account, *domain.Finetune) (domain.JobInfo, error)
	Delete(string) error
	Terminate(string) error

	// GetLogDownloadURL returns the log url which can be used
	// to download the log of running training.
	GetLogDownloadURL(string) (string, error)

	GetDetail(string) (domain.JobDetail, error)
}
