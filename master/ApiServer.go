package master

import (
	"crontab-go/common"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Server
type Server struct {
}

// ApiServer
func ApiServer() *Server {
	return &Server{}
}

// 路由注册
func Route(r *gin.Engine) {
	r.GET("/job/list", JobList)
	r.POST("/job/save", JobSave)
	r.POST("/job/remove", RemoveJob)
	r.POST("/job/kill", KillJob)
}

// 启动ApiServer
func (a *Server) Run() (err error) {
	r := gin.Default()
	// 注册路由
	Route(r)
	// 运行
	err = r.Run(":" + strconv.Itoa(Config.APIPort))
	return
}

// 任务列表
func JobList(ctx *gin.Context) {
	var (
		pageQuery int
		sizeQuery int
		page      int64
		size      int64
		pageInfo  common.PageInfo
		err       error
	)
	// 参数处理
	pageQuery, err = strconv.Atoi(ctx.Query("page"))
	sizeQuery, err = strconv.Atoi(ctx.Query("size"))
	page, size = int64(pageQuery), int64(sizeQuery)
	{
		// 默认页
		if page <= 0 {
			page = 1
		}
		// 默认显示条数
		if size <= 0 {
			size = Config.DefaultSize
		}
	}
	// 任务列表
	if pageInfo, err = JobMgr.JobList(page, size); err != nil {
		ctx.JSON(200, common.ApiResponse{
			ErrorCode: 500,
			Msg:       err.Error(),
		})
		return
	}
	// 返回数据
	ctx.JSON(200, common.ApiResponse{
		ErrorCode: 0,
		Msg:       "获取任务列表成功!",
		Data:      pageInfo,
	})
}

// 保存任务
func JobSave(ctx *gin.Context) {
	var (
		jobJson string
		job     *common.Job
		err     error
	)
	jobJson = ctx.PostForm("job")
	// 反序列化
	if job, err = common.UnmarshalJob([]byte(jobJson)); err != nil {
		ctx.JSON(200, common.ApiResponse{
			ErrorCode: 400,
			Msg:       "参数错误," + err.Error(),
		})
		return
	}
	// 保存失败
	if job, err = JobMgr.SaveJob(job); err != nil {
		ctx.JSON(200, common.ApiResponse{
			ErrorCode: 400,
			Msg:       err.Error(),
		})
		return
	}
	// 返回
	ctx.JSON(200, common.ApiResponse{
		ErrorCode: 0,
		Msg:       "保存任务成功!",
		Data:      job,
	})
}

// 杀死任务
func KillJob(ctx *gin.Context) {
	var (
		jobName string
		err     error
	)
	jobName = ctx.PostForm("jobName")
	// 杀死任务
	if err = JobMgr.KillJob(jobName); err != nil {
		ctx.JSON(200, common.ApiResponse{
			ErrorCode: 500,
			Msg:       err.Error(),
		})
		return
	}
	// 返回
	ctx.JSON(200, common.ApiResponse{
		ErrorCode: 0,
		Msg:       "杀死任务成功!",
	})
}

// 删除任务
func RemoveJob(ctx *gin.Context) {
	var (
		jobName string
		job     *common.Job
		err     error
	)
	jobName = ctx.PostForm("jobName")
	// 杀死任务
	if job, err = JobMgr.DelJob(jobName); err != nil {
		ctx.JSON(200, common.ApiResponse{
			ErrorCode: 500,
			Msg:       err.Error(),
		})
		return
	}
	// 返回
	ctx.JSON(200, common.ApiResponse{
		ErrorCode: 0,
		Msg:       "删除任务成功!",
		Data:      job,
	})
}
