package settings

import "time"

type ServerSettingS struct {
	Name         string
	RunMode      string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	StartTime    string
	MachineId    int64
}

type AppSettingS struct {
	DefaultPageSize       int
	MaxPageSize           int
	DefaultContextTimeout time.Duration
	//UploadSavePath        string
	//UploadServerUrl       string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

type LogSettingS struct {
	Level             string
	Filename          string
	MaxSize           int
	MaxBackups        int
	MaxAge            int
	AccessLogFilename string
}

type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

type MySqlSettingS struct {
	Host         string
	UserName     string
	Password     string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type RedisSettingS struct {
	Host     string
	Password string
	Port     int
	Db       int
	PoolSize int
}

//将配置反序列化到 sections map中
var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
