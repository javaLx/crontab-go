package master

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

// 全局变量
var (
	WorkerManager *WorkerMgr
)

// 工作管理
type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// 初始化etcd相关信息
func InitWorkerMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	// 1 创建配置对象
	config = clientv3.Config{
		Endpoints:   Config.EtcdEndpoints,
		DialTimeout: time.Duration(Config.EtcdDialTimeout) * time.Millisecond,
	}
	// 2 创建连接获得客户端
	if client, err = clientv3.New(config); err != nil {
		return
	}
	// 3 创建kv对象
	kv = clientv3.NewKV(client)
	// 4 创建租约对象
	lease = clientv3.NewLease(client)
	// 创建WorkerMgr对象
	WorkerManager = &WorkerMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}
