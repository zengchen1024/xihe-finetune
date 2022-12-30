package finetuneimpl

// respone
type response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

// token
type tokenRequest struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

type tokenResp struct {
	Duration int64  `json:"duration"`
	Token    string `json:"token"`
	Msg      string `json:"msg"`
}

// create
type createRequest struct {
	User       string           `json:"user"`
	Name       string           `json:"task_name"`
	Task       string           `json:"task_type"`
	Model      string           `json:"foundation_model"`
	Parameters []hyperparameter `json:"parameters,omitempty"`
}

type hyperparameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type createResp struct {
	response

	JobId string `json:"job_id"`
}

// get
type getResp struct {
	Status int        `json:"status"`
	Msg    string     `json:"msg"`
	Data   detailData `json:"data"`
}

type detailData struct {
	Phase      string `json:"phase"`
	TaskType   string `json:"task_type"`
	TaskName   string `json:"task_name"`
	Framework  string `json:"framework"`
	CreatedAt  string `json:"created_at"`
	EngineName string `json:"engine_name"`
	Runtime    int    `json:"runtime"`
}

// log
type logResp struct {
	response

	OBSURL string `json:"obs_url"`
}
