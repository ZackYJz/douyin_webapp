package settings

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

// NewSetting 初始化配置
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New() //创建新的配置
	// viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	vp.SetConfigName("config") //配置文件名称(无扩展名)
	//将命令行参数传入的配置文件路径添加
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.AddConfigPath("configs/") //配置文件相对路径:可配置多个
	vp.SetConfigType("yaml")     //配置文件类型
	err := vp.ReadInConfig()     // 查找并读取配置文件
	if err != nil {
		return nil, err
	}

	s := &Setting{vp: vp}
	s.WatchSettingChange()
	return s, nil
}

//文件热更新的监听和更变处理
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		//当配置文件修改的钩子函数
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			log.Println("发现配置文件更变，执行热更新....")
			//让viper的变更重新更新到配置结构体中
			_ = s.ReloadAllSection()
		})
	}()
}
