package watchimpl

import (
	"errors"
	"sync"
	"time"

	"github.com/opensourceways/xihe-grpc-protocol/grpc/client"
	pt "github.com/opensourceways/xihe-grpc-protocol/grpc/finetune"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-finetune/domain/finetune"
	"github.com/opensourceways/xihe-finetune/domain/watch"
)

type finetuneData = pt.FinetuneInfo

func NewWatcher(cfg *Config, fs finetune.Finetune) (*Watcher, error) {
	cli, err := client.NewFinetuneClient(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		cli:         cli,
		fs:          fs,
		interval:    time.Duration(cfg.Interval) * time.Second,
		stop:        make(chan struct{}),
		stopped:     make(chan struct{}),
		finetunes:   make(chan finetuneInfo, cfg.MaxWatchNum+1),
		maxWatchNum: cfg.MaxWatchNum,
	}, nil
}

type finetuneInfo struct {
	watch.FinetuneInfo

	result finetuneData

	done bool
}

func (t *finetuneInfo) toIndex() pt.FinetuneIndex {
	return pt.FinetuneIndex{
		Id:   t.Id,
		User: t.User.Account(),
	}
}

func (t *finetuneInfo) isDone() bool {
	return t.done
}

// Watcher
type Watcher struct {
	cli *client.FinetuneClient
	fs  finetune.Finetune

	interval time.Duration

	stop      chan struct{}
	stopped   chan struct{}
	finetunes chan finetuneInfo

	lock        sync.RWMutex
	currentNum  int
	maxWatchNum int
}

func (w *Watcher) ApplyWatch(f func(*watch.FinetuneInfo) error) (err error) {
	if !w.increase() {
		return errors.New("exceed max watch num")
	}

	info := new(watch.FinetuneInfo)

	if err = f(info); err != nil {
		w.decrease()
	} else {
		w.finetunes <- finetuneInfo{FinetuneInfo: *info}
	}

	return
}

func (w *Watcher) increase() (b bool) {
	w.lock.Lock()
	if w.currentNum+1 <= w.maxWatchNum {
		w.currentNum++
		b = true
	}
	w.lock.Unlock()

	return
}

func (w *Watcher) decrease() {
	w.lock.Lock()
	w.currentNum--
	w.lock.Unlock()
}

func (w *Watcher) Run() {
	start := time.Now()

	// add the tag
	w.finetunes <- finetuneInfo{}

	for {
		select {
		case info := <-w.finetunes:
			// use =="" stands for the case that the loop is done
			if info.User == nil {
				logrus.Debug("finish a loop")

				t := start.Add(w.interval)

				if n := time.Now(); t.After(n) {
					time.Sleep(t.Sub(n))
				}

				w.finetunes <- finetuneInfo{}

				start = time.Now()

			} else {
				changed := w.check(&info)
				logrus.Debugf("check finetune %s/%s", info.Id, info.JobId)

				if info.isDone() {
					index := info.toIndex()

					if err := w.cli.SetFinetuneInfo(&index, &info.result); err == nil {
						w.decrease()
					} else {
						logrus.Errorf("set finetune info failed, err:%s", err.Error())

						w.finetunes <- info
					}
				} else {
					if changed {
						index := info.toIndex()

						if err := w.cli.SetFinetuneInfo(&index, &info.result); err != nil {
							logrus.Errorf("set finetune info failed, err:%s", err.Error())
						}
					}

					w.finetunes <- info
				}
			}

		case <-w.stop:
			close(w.stopped)

			return
		}
	}
}

func (w *Watcher) Exit() {
	close(w.stop)

	<-w.stopped

	w.cli.Disconnect()
}

func (w *Watcher) check(info *finetuneInfo) (changed bool) {
	result := &info.result

	if !info.done {
		detail, err := w.fs.GetDetail(info.JobId)
		if err != nil {
			return
		}

		if result.Duration != detail.Duration {
			result.Duration = detail.Duration
			changed = true
		}

		if s := detail.Status.FinetuneStatus(); s != result.Status {
			result.Status = s
			changed = true
		}

		if detail.Status.IsDone() {
			info.done = true
		}
	}

	return
}
