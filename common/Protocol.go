package common

import "encoding/json"

// 定时任务
type Job struct {
	Name     string `json:"name"`     //  任务名
	Command  string `json:"command"`  // shell命令
	CronExpr string `json:"cronExpr"` // cron表达式
}

// HTTP接口应答
type ApiResponse struct {
	ErrorCode int         `json:"error_code"` // 错误码
	Msg       string      `json:"msg"`        // 消息
	Data      interface{} `json:"data"`       // 数据
}

// 分页信息
type PageInfo struct {
	Count     int64       `json:"count"`      // 总条数
	TotalPage int64       `json:"total_page"` //总页数
	CurrPage  int64       `json:"curr_page"`  // 当前页
	Size      int64       `json:"size"`       // 每页显示数量
	Data      interface{} `json:"data"`       // 数据
}

// 任务序列化
func (j *Job) MarshalJob() (jobJson string, err error) {
	var (
		context []byte
	)
	if context, err = json.Marshal(j); err != nil {
		return
	}
	jobJson = string(context)
	return
}

// 任务反序列化
func UnmarshalJob(jobJson []byte) (job *Job, err error) {
	job = &Job{}
	if err = json.Unmarshal(jobJson, job); err != nil {
		return
	}
	return
}
