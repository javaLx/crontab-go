package master

import (
	"context"
	"crontab-go/common"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"math"
)

var (
	JobMgr *JobManager
)

// 任务管理器
type JobManager struct {
}

// 初始化任务管理器
func InitJobManager() error {
	JobMgr = &JobManager{}
	return nil
}

// 保存任务
func (j *JobManager) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey      string
		jobValue    string
		putResponse *clientv3.PutResponse
	)
	// etcd保存的key
	jobKey = common.JOB_SAVE_DIR + job.Name
	//序列化json
	if jobValue, err = job.MarshalJob(); err != nil {
		return
	}
	// 保存到etcd中
	if putResponse, err = WorkerManager.kv.Put(context.TODO(), jobKey,
		jobValue, clientv3.WithPrevKV()); err != nil {
		return
	}
	// 如果是更新 返回旧值
	if putResponse.PrevKv != nil {
		// 反序列化job
		if oldJob, err = common.UnmarshalJob(putResponse.PrevKv.Value); err != nil {
			err = nil
			return
		}
	}
	return
}

// 删除任务
func (j *JobManager) DelJob(jobName string) (oldJob *common.Job, err error) {
	var (
		jobKey         string
		deleteResponse *clientv3.DeleteResponse
	)
	// etcd保存的key
	jobKey = common.JOB_SAVE_DIR + jobName
	// 删除key
	if deleteResponse, err = WorkerManager.kv.Delete(context.TODO(),
		jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	// 获取被删除之前到信息
	if len(deleteResponse.PrevKvs) != 0 {
		// 反序列化job
		if oldJob, err = common.UnmarshalJob(deleteResponse.PrevKvs[0].Value); err != nil {
			err = nil
			return
		}
	}
	return
}

// 任务列表分页
func (j *JobManager) JobList(page, size int64) (pageInfo common.PageInfo, err error) {
	var (
		getResponse *clientv3.GetResponse
		keyValue    *mvccpb.KeyValue
		key         int
		startIndex  int64
		endIndex    int64
		job         *common.Job
		jobList     []*common.Job
	)

	// 获取所以的任务
	if getResponse, err = WorkerManager.kv.Get(context.TODO(),
		common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	// 分页对象
	pageInfo.Count = getResponse.Count
	pageInfo.CurrPage = page
	pageInfo.TotalPage = int64(math.Ceil(float64(pageInfo.Count) / float64(size)))

	// 计算开始和结束数据位置
	{
		startIndex = (page - 1) * size
		if startIndex < 0 {
			startIndex = 0
		}
		endIndex = (startIndex + size) - 1
		pageInfo.Size = size
	}

	// 处理数据
	for key, keyValue = range getResponse.Kvs {
		// 超出页码
		if pageInfo.CurrPage > pageInfo.TotalPage {
			break
		}
		//小于开始位置
		if int64(key) < startIndex {
			continue
		}
		// 大于结束位置
		if int64(key) > endIndex {
			break
		}
		// 反序列化
		if job, err = common.UnmarshalJob(keyValue.Value); err != nil {
			continue
		}
		jobList = append(jobList, job)
	}
	pageInfo.Data = jobList
	return
}

// 杀死任务
func (j *JobManager) KillJob(jobName string) (err error) {
	// 更新一下key=/cron/killer/任务名
	var (
		killerKey      string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
	)
	// 通知worker杀死对应任务
	killerKey = common.JOB_KILLER_DIR + jobName

	// 让worker监听到一次put操作, 创建一个租约让其稍后自动过期即可
	if leaseGrantResp, err = WorkerManager.lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	// 租约ID
	leaseId = leaseGrantResp.ID

	// 设置killer标记
	if _, err = WorkerManager.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}
	return
}
