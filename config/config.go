package config

import (
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bitly/go-nsq"
	"github.com/coreos/etcd/client"
	"github.com/lodastack/log"
	"strings"
)

var (
	mux        = new(sync.RWMutex)
	config     = new(Config)
	configPath = ""
)

type Config struct {
	Com  CommonConfig   `toml:"common"`
	Reg  RegistryConfig `toml:"registry"`
	Nsq  NsqConfig      `toml:"nsq"`
	Etcd EtcdConfig     `toml:"etcd"`
	Log  LogConfig      `toml:"log"`

	EtcdConfig client.Config `toml:"-"`
}
type EtcdConfig struct {
	Auth          bool          `toml:"auth"`
	Username      string        `toml:"username"`
	Password      string        `toml:"password"`
	Endpoints     []string      `toml:"endpoints"`
	HeaderTimeout time.Duration `toml:"timeout"`
}

type CommonConfig struct {
	Listen             string `toml:"listen"`
	InfluxdPort        int    `toml:"influxdPort"`
	TopicsPollInterval int    `toml:"topicsPollInterval"`
	HiddenMetricSuffix string `toml:"hiddenMetricSuffix"`
}

type LogConfig struct {
	Enable   bool   `toml:"enable"`
	Path     string `toml:"path"`
	Level    string `toml:"level"`
	FileNum  int    `toml:"file_num"`
	FileSize int    `toml:"file_size"`
}

type RegistryConfig struct {
	Link      string `toml:"link"`
	ExpireDur int    `toml:"expireDur"`
}

type NsqConfig struct {
	Enable              bool     `toml:"enable"`
	MaxAttempts         uint16   `toml:"maxAttempts"`
	MaxInFlight         int      `toml:"maxInFlight"`
	HeartbeatInterval   int      `toml:"heartbeatInterval"`
	ReadTimeout         int      `toml:"readTimeout"`
	LookupdPollInterval int      `toml:"lookupdPollInterval"`
	HandlerCount        int      `toml:"handlerCount"`
	Lookupds            []string `toml:"lookupds"`
	Chan                string   `toml:"chan"`
	TopicPrefix         string   `toml:"topicPrefix"`
}

func (this NsqConfig) GetNsqConfig() *nsq.Config {
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxAttempts = this.MaxAttempts
	nsqConfig.MaxInFlight = this.MaxInFlight
	nsqConfig.HeartbeatInterval = time.Duration(this.HeartbeatInterval) * time.Millisecond
	nsqConfig.ReadTimeout = time.Duration(this.ReadTimeout) * time.Millisecond
	nsqConfig.LookupdPollInterval = time.Duration(this.LookupdPollInterval) * time.Millisecond

	return nsqConfig
}

func Reload() {
	err := LoadConfig(configPath)
	if err != nil {
		os.Exit(1)
	}
}

func LoadConfig(path string) (err error) {
	mux.Lock()
	defer mux.Unlock()
	configPath = path
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Error while loading config %s.\n%s\n", path, err.Error())
		return
	}
	if _, err = toml.Decode(string(configFile), &config); err != nil {
		log.Errorf("Error while decode the config %s.\n%s\n", path, err.Error())
		return
	} else {
		return nil
	}
}

func GetConfig() *Config {
	mux.RLock()
	defer mux.RUnlock()
	return config
}
