package finetuneimpl

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/opensourceways/community-robot-lib/utils"
)

func newClient(cfg *Config) client {
	return client{
		cfg: *cfg,
		hc:  utils.NewHttpClient(3),
	}
}

type client struct {
	cfg Config
	hc  utils.HttpClient
}

func (cli *client) tokenURL() string {
	return cli.cfg.Endpoint + "/foundation-model/token"
}

func (cli *client) createURL() string {
	return cli.cfg.Endpoint + "/v1/foundation-model/finetune"
}

func (cli *client) jobURL(jobId string) string {
	return cli.createURL() + "/" + jobId
}

func (cli *client) logURL(jobId string) string {
	return cli.jobURL(jobId) + "/log"
}

func (cli *client) token() (t string, err error) {
	payload, err := utils.JsonMarshal(&tokenRequest{
		User:     cli.cfg.Username,
		Password: cli.cfg.Password,
	})
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, cli.tokenURL(), bytes.NewBuffer(payload),
	)
	if err != nil {
		return
	}

	resp := new(tokenResp)
	if err = cli.forwardTo(req, "", resp); err != nil {
		return
	}

	if resp.Msg != "" {
		err = errors.New(resp.Msg)
	} else {
		t = resp.Token
	}

	return
}

func (cli *client) createJob(options *createRequest) (jobId string, err error) {
	payload, err := utils.JsonMarshal(options)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, cli.createURL(), bytes.NewBuffer(payload),
	)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp := new(createResp)
	if err = cli.forwardTo(req, token, resp); err != nil {
		return
	}

	if resp.Status == 201 {
		jobId = resp.JobId
	} else {
		err = errors.New(resp.Msg)
	}

	return
}

func (cli *client) getJob(jobId string) (info detailData, err error) {
	req, err := http.NewRequest(http.MethodGet, cli.jobURL(jobId), nil)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	res := new(getResp)
	if err = cli.forwardTo(req, token, res); err != nil {
		return
	}

	if res.Status == 200 {
		info = res.Data
	} else {
		err = errors.New(res.Msg)
	}

	return
}

func (cli *client) deleteJob(jobId string) (err error) {
	req, err := http.NewRequest(http.MethodDelete, cli.jobURL(jobId), nil)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp := new(response)
	if err = cli.forwardTo(req, token, resp); err != nil {
		return
	}

	if resp.Status != 204 {
		err = errors.New(resp.Msg)
	}

	return
}

func (cli *client) terminateJob(jobId string) (err error) {
	req, err := http.NewRequest(http.MethodPut, cli.jobURL(jobId), nil)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp := new(response)
	if err = cli.forwardTo(req, token, resp); err != nil {
		return
	}

	if resp.Status != 202 {
		err = errors.New(resp.Msg)
	}

	return
}

func (cli *client) getLogURL(jobId string) (log string, err error) {
	req, err := http.NewRequest(http.MethodGet, cli.logURL(jobId), nil)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp := new(logResp)
	if err = cli.forwardTo(req, token, resp); err != nil {
		return
	}

	if resp.Status == 200 {
		log = resp.OBSURL
	} else {
		err = errors.New(resp.Msg)
	}

	return
}

func (cli *client) forwardTo(req *http.Request, token string, jsonResp interface{}) (err error) {
	if token != "" {
		req.Header.Set("Authorization", "JWT "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	_, err = cli.hc.ForwardTo(req, jsonResp)

	return
}
