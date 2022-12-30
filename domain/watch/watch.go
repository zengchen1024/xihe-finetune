package watch

import "github.com/opensourceways/xihe-finetune/domain"

type FinetuneInfo struct {
	User domain.Account
	Id   string

	domain.JobInfo
}

type WatchService interface {
	ApplyWatch(f func(*FinetuneInfo) error) error
}
