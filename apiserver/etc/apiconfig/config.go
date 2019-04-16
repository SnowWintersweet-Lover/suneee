package apiconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// auto generate struct
// https://mholt.github.io/json-to-go/
// use mapstructure to replace json for '_' key words, e.g. rpc_port,big_data
type ConfigStruct struct {
	Server struct {
		Host struct {
			Address string `json:"address"`
			Port    int    `json:"port"`
		} `json:"host"`
		TLS struct {
			Enable     bool `json:"enable"`
			VerifyPeer bool `mapstructure:"verify_peer"`
		} `json:"tls"`
		Pki struct {
			Key  string `json:"key"`
			Cert string `json:"cert"`
			Ca   string `json:"ca"`
		} `json:"pki"`
	} `json:"server"`
	Log struct {
		Path string `json:"path"`
		Host struct {
			Address string `json:"address"`
			Port    int    `json:"port"`
		} `json:"host"`
	} `json:"log"`
	Minio struct {
		Enable    bool   `json:"enable"`
		Secure    bool   `json:"secure"`
		Address   string `json:"address"`
		Port      int    `json:"port"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
	} `json:"minio"`
	TokenAuth struct {
		Path string `json:"path"`
	} `mapstructure:"token_auth"`
}

var (
	defaultApiVersion = "v1"
	defaultFilePath   = "/etc/apiconfig/config.json"
	ViperConfig       *viper.Viper
	Config            *ConfigStruct
	serverPath        = os.Getenv("BAAS_PATH")
	serverType        = os.Getenv("BAAS_TYPE")
	serverTypeProd    = "production"
	serverTypeTest    = "development"
)

// 不再初始化时自动读取配置文件
func init() {
	if serverPath == "" {
		serverPath = "./"
		fmt.Println("BAAS_PATH env not set, use ./ as default")
	}
	//log.Sugar.Debugf("BAAS_PATH ===> %s", serverPath)
	InitConfig(filepath.Join(GetServerDir(), defaultFilePath))
}
func InitConfig(filePath string) {
	ViperConfig = viper.New()
	if filePath == "" {
		filePath = filepath.Join(GetServerDir(), defaultFilePath)
	}
	if !strings.HasSuffix(filePath, "config.json") {
		panic(fmt.Sprintf("配置文件路径必须以config.json为结尾 ===> %s", filePath))
	}
	ViperConfig.SetConfigFile(filePath)
	err := ViperConfig.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = ViperConfig.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
func GetApiDefaultVersion() string {
	return defaultApiVersion
}
func GetServerDir() string {
	return serverPath
}

func GetServerHostAddr() string {
	return Config.Server.Host.Address
}

func GetServerHostPort() int {
	return Config.Server.Host.Port
}

func GetServerTlsEnable() bool {
	return Config.Server.TLS.Enable
}

func GetServerPkiCa() string {
	return Config.Server.Pki.Ca
}

func GetServerPkiCert() string {
	return Config.Server.Pki.Cert
}

func GetServerPkiKey() string {
	return Config.Server.Pki.Key
}

func GetServerTlsVerifyPeer() bool {
	return Config.Server.TLS.VerifyPeer
}

func GetLogPath() string {
	return filepath.Join(GetServerDir(), Config.Log.Path)
}

func GetLogHostAddress() string {
	return Config.Log.Host.Address
}

func GetLogHostPort() int {
	return Config.Log.Host.Port
}

func ServerTypeIsProd() bool {
	if serverType == serverTypeProd {
		return true
	}
	return false
}

func GetMinioEnable() bool {
	return Config.Minio.Enable
}
func GetMinioSecure() bool {
	return Config.Minio.Secure
}
func GetMinioAddress() string {
	return Config.Minio.Address
}
func GetMinioPort() int {
	return Config.Minio.Port
}
func GetMinioAccessKey() string {
	return Config.Minio.AccessKey
}
func GetMinioSecretKey() string {
	return Config.Minio.SecretKey
}
func GetTokenAuthPath() string {
	return Config.TokenAuth.Path
}
