package master

import (
	"encoding/json"
	"io/ioutil"
)

// 全局变量
var (
	Config *Conf
)

// 配置
type Conf struct {
	APIPort               int      `json:"apiPort"`               //API接口服务端口
	EtcdEndpoints         []string `json:"etcdEndpoints"`         //	etcd的集群列表
	EtcdDialTimeout       int      `json:"etcdDialTimeout"`       //	etcd的连接超时
	MongodbURI            string   `json:"mongodbUri"`            //	mongodb地址
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"` //	mongodb连接超时时间

	// 分页相关
	DefaultSize int64 `json:"default_size"` // 默认每页显示数量
}

// 加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Conf
	)
	// 1 读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	// 2 json反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return err
	}
	// 全局赋值
	Config = &conf

	return
}
