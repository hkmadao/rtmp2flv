package conf

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/core/config"
	_ "github.com/beego/beego/v2/core/config/yaml"
)

func init() {
	filePath := "./conf/conf.yml"
	err := config.InitGlobalInstance("yaml", filePath)
	if err != nil {
		logs.Error("read conf file [%s] error : %v", filePath, err)
		return
	}
}
