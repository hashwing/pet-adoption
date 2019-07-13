package config

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/hashwing/log"
	"gopkg.in/yaml.v2"
)

// Config 配置信息
type Config struct {
	Port        int         `yaml:"listen_port"`
	WxConfig    WxConfig    `yaml:"wx_config"`
	MysqlConfig MysqlConfig `yaml:"mysql_config"`
	OssConfig   AliOss      `yaml:"oss"`
	Log         Log         `yaml:"log"`
}

type WxConfig struct {
	AppID  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type AliOss struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	InHost          string `yaml:"in_host"`
	Host            string `yaml:"host"`
	CallbackUrl     string `yaml:"callback_url"`
	UploadDir       string `yaml:"upload_dir"`
	ExpireTime      int64  `yaml:"expire_time"`
}

// Log log out setting
type Log struct {
	Level   int       `yaml:"level"`
	Outputs LogOutput `yaml:"outputs"`
}

// LogOutput log output type
type LogOutput struct {
	File    FileLog    `yaml:"file"`
	Console ConsoleLog `yaml:"console"`
}

// FileLog file log output
type FileLog struct {
	Enabled bool   `yaml:"enabled"`
	LogPath string `yaml:"log_path"`
}

// ConsoleLog console log output
type ConsoleLog struct {
	Enabled bool `yaml:"enabled"`
}

// Cfg 全局配置信息
var Cfg = &Config{}

// InitGlobal 初始化全局配置变量
func InitGlobal(file string) error {
	err := LoadConfig(file, Cfg)
	if err != nil {
		return err
	}
	if Cfg.Log.Outputs.File.LogPath == "" {
		Cfg.Log.Outputs.File.LogPath = "/var/log/pet-adoption/access.log"
	}
	logger, err := log.NewBeegoLog(Cfg.Log.Outputs.File.LogPath, Cfg.Log.Level, Cfg.Log.Outputs.Console.Enabled)
	if err != nil {
		log.Error(err)
	}
	log.SetHlogger(logger)
	return nil
}

// LoadConfig load config from config file
func LoadConfig(file string, settings interface{}) error {

	if file != "" {

		absConfPath, err := filepath.Abs(file)
		if err != nil {
			log.Debug(err)
			return err
		}

		if confData, err := ioutil.ReadFile(absConfPath); err != nil {
			log.Debug(err)
			return err
		} else {
			if err := yaml.Unmarshal(confData, settings); err != nil {
				log.Debug(err)
				return err
			}
		}
		return nil
	}

	return errors.New("file is nil")
}
