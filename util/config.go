package util

import (
	"github.com/spf13/viper"
)

const ConfigName = "config"
const ConfigType = "yaml"

var Configuration Config

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`

	Logger struct {
		Dir        string `mapstructure:"dir"`
		FileName   string `mapstructure:"file_name"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		Compress   bool   `mapstructure:"compress"`
		LocalTime  bool   `mapstructure:"local_time"`
	} `mapstructure:"logger"`
	App struct {
		AppID              string `mapstructure:"app_id"`
		ConfigID           string `mapstructure:"config_id"`
		Secret             string `mapstructure:"secret"`
		HostURLCallback    string `mapstructure:"host_url_callback"`
		HostClientCallback string `mapstructure:"host_client_callback"`
	} `mapstructure:"app"`
	MongoDB struct {
		Host           string `mapstructure:"host"`
		Port           int    `mapstructure:"port"`
		Username       string `mapstructure:"username"`
		Password       string `mapstructure:"password"`
		Database       string `mapstructure:"database"`
		CredentialDB   string `mapstructure:"credential_database"`
		MaxPoolSize    uint64 `mapstructure:"max_pool_size"`
		ConnectTimeout int    `mapstructure:"connect_timeout"`
		AuthSource     string `mapstructure:"auth_source"`
		Collections    struct {
			Token            string `mapstructure:"token"`
			ClientCredential string `mapstructure:"client_credential"`
		} `mapstructure:"collection"`
	} `mapstructure:"mongodb"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	var config Config
	err = viper.Unmarshal(&config)
	Configuration = config
	return
}
