package controller

import (
	"github.com/opensourceways/xihe-finetune/app"
	"github.com/opensourceways/xihe-finetune/domain"
)

type FinetuneCreateRequest struct {
	User            string            `json:"user"`
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	Task            string            `json:"task"`
	Model           string            `json:"model"`
	Hyperparameters map[string]string `json:"hyperparameter"`
}

func (req *FinetuneCreateRequest) toCmd() (cmd app.FinetuneCreateCmd, err error) {
	if cmd.User, err = domain.NewAccount(req.User); err != nil {
		return
	}

	cmd.Id = req.Id

	if cmd.Name, err = domain.NewFinetuneName(req.Name); err != nil {
		return
	}

	cmd.Param, err = domain.NewFinetuneParameter(req.Model, req.Task, req.Hyperparameters)

	return
}
