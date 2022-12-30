package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const pathSpliter = "/"

var (
	reName      = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	reDirectory = regexp.MustCompile("^[a-zA-Z0-9_/-]+$")
	reFilePath  = regexp.MustCompile("^[a-zA-Z0-9_/.-]+$")

	FinetuneStatusFailed      = finetuneStatus("Failed")
	FinetuneStatusPending     = finetuneStatus("Pending")
	FinetuneStatusRunning     = finetuneStatus("Running")
	FinetuneStatusCreating    = finetuneStatus("Creating")
	FinetuneStatusAbnormal    = finetuneStatus("Abnormal")
	FinetuneStatusCompleted   = finetuneStatus("Completed")
	FinetuneStatusTerminated  = finetuneStatus("Terminated")
	FinetuneStatusTerminating = finetuneStatus("Terminating")

	finetuneDoneStatus = map[string]bool{
		"Failed":     true,
		"Abnormal":   true,
		"Completed":  true,
		"Terminated": true,
	}
)

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" || strings.ToLower(v) == "root" || !reName.MatchString(v) {
		return nil, errors.New("invalid user name")
	}

	return dpAccount(v), nil
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

// FinetuneName
type FinetuneName interface {
	FinetuneName() string
}

func NewFinetuneName(v string) (FinetuneName, error) {
	max := config.MaxFinetuneNameLength
	min := config.MinFinetuneNameLength

	if n := len(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !reName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return finetuneName(v), nil
}

type finetuneName string

func (r finetuneName) FinetuneName() string {
	return string(r)
}

// FinetuneStatus
type FinetuneStatus interface {
	FinetuneStatus() string
	IsDone() bool
	IsSuccess() bool
}

type finetuneStatus string

func (s finetuneStatus) FinetuneStatus() string {
	return string(s)
}

func (s finetuneStatus) IsDone() bool {
	return finetuneDoneStatus[string(s)]
}

func (s finetuneStatus) IsSuccess() bool {
	return string(s) == FinetuneStatusCompleted.FinetuneStatus()
}
