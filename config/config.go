package config

import "fmt"
import "github.com/spf13/viper"

// LoadConfig 加载配置
func LoadConfig() {
	// 设置配置文件的名称（不需要文件扩展名）
	viper.SetConfigName("conf")
	// 添加配置文件所在的路径，这里假设配置文件和main.go在同一目录
	viper.AddConfigPath("config/")
	viper.AddConfigPath("/config/")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("../../config/")
	viper.AddConfigPath("../../../config/")
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("读取配置文件出错：%v\n", err)
		panic(any(err))
	}
}
