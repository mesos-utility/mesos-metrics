package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type MasterConfig struct {
	Enable   bool   `json:"enable"`
	Apiurl   string `json:"apiurl"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SlaveConfig struct {
	Enable   bool   `json:"enable"`
	Apiurl   string `json:"apiurl"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GlobalConfig struct {
	Debug    bool          `json:"debug"`
	Interval int           `json:"interval"`
	Http     *HttpConfig   `json:"http"`
	Master   *MasterConfig `json:"master"`
	Slave    *SlaveConfig  `json:"slave"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g:ParseConfig, ok, ", cfg)
}
