package data

// BaseConfig .
// 基础文件配置
type BaseConfig struct {
	Server serverCfg `json:"server"`
	Cert   certCfg   `json:"cert"`
	MySQL  mysqlCfg  `json:"mysql"`
	WeChat wechat    `json:"wechat"`
}

type serverCfg struct {
	IP        string   `json:"ip"`
	Host      string   `json:"host"`
	Port      string   `json:"port"`
	Name      string   `json:"name"`
	Mode      string   `json:"mode"`
	Secret    string   `json:"secret"`
	Author    string   `json:"author"`
	Version   string   `json:"version"`
	WhiteList []string `json:"white-list"`
	Picture string `json:"picture"`
}
type certCfg struct {
	KeyPath  string `json:"key_path"`
	CertPath string `json:"cert_path"`
}

type mysqlCfg struct {
	Dsn     string `json:"dsn"`
	MaxConn int    `json:"max_conn"`
	MaxIdle int    `json:"max_idle"`
	ShowSQL bool   `json:"show_sql"`
}
type wechat struct {
	Name         string `json:"name"`
	Account      string `json:"account"`
	AppID        string `json:"appid"`
	Secret       string `json:"secret"`
	Aes          string `json:"aes"`
	State        string `json:"state"`
	Key          string `json:"key"`
	Token        string `json:"token"`
	MchID        string `json:"mch_id"`
	PayNotifyURL string `json:"notify_url"`
	CertKeyPath  string `json:"cert_key_path"`
	CertCertPath string `json:"cert_cert_path"`
	CertCaPath   string `json:"cert_ca_path"`
}
