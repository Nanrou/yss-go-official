package orm

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// -------- 配置文件的定义 ------------

// 每个sql的data
type sqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// 白名单
type trustMeta struct {
	Ip      []string
	Network []string
}

// config yaml的定义
type config struct {
	Local  sqlConfig
	Prod   sqlConfig
	Mssql  sqlConfig
	Secret string
	Trust  trustMeta
}

// 根据运行环境选择sql的配置
func (c *config) sqlMeta() sqlConfig {
	if runtime.GOOS != "linux" {
		return c.Local
	} else {
		return c.Prod
	}
}

func (c *config) GetTrustIps() []string {
	return c.Trust.Ip
}

func (c *config) GetTrustNetwork() []string {
	return c.Trust.Network
}

// 获取生产config
func GetConfig() *config {
	s, _ := os.Getwd()
	s = filepath.Join(s, "config.yaml")
	content, err := ioutil.ReadFile(s)
	if err != nil {
		log.Fatal(err)
	}
	c := config{}
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

// -------- 表结构的定义 ------------

// wechat_profile表的定义
type WeChatProfile struct {
	ID             int
	IDCardNumber   string
	DefaultAccount string
}

// user_data表的定义
type UserData struct {
	ID           int
	IDCardNumber string
	Account      string
	Name         string
	Phone        string
	// IsPaid bool
	// LastUpdate timestamp
}

// -------- 存储过程结果的定义 ------------

type AccountCheckColumns struct {
	Validated string
	Address   string
	Name      string
}

type BillListColumns struct {
	Date          string
	Charge        string
	CurrentMeter  string
	PreviousMeter string
	IsPaid        bool
	Yszbh         string
}

type BillDetailColumns struct {
	address          string
	name             string
	charge           string
	currentMeter     string
	meterReadingDate string
	isPaid           bool
	previousMeter    string
	waterCharge      string
	waterProperty    string
	// ... 各种其他费用
}
