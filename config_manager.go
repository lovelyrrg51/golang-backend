package main

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type App_config struct {
	User_db     string
	Password_db string
	IP_db       string
	Address     string
	Dbname      string
}

type Configuration_manager struct {
	mutex sync.Mutex
	app   App_config
	v     *viper.Viper
}

var instance *Configuration_manager
var once sync.Once

func (cm *Configuration_manager) Load(configuration_file string) bool {
	var tmp Configuration_manager
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.v.SetConfigFile(configuration_file)
	cm.v.SetConfigType("toml")
	err := cm.v.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error occurs while reading config file: %s \n", err)
		return false
	}
	err = cm.v.UnmarshalKey("app", &tmp.app)
	if err != nil {
		log.Fatalf("[app] part of config file is not valid: %s \n", err)
		return false
	}
	mi := cm.v.Get("app")
	if mi == nil {
		log.Fatalf("[app] part of config file is not valid\n")
		return false
	}
	cm.app = tmp.app
	return true

}
func (cm Configuration_manager) GetAppConfig() App_config {
	return cm.app
}
