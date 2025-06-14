package config

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type JiraConfig struct {
	Url           string `yaml:"url"`
	ThreadCount   int    `yaml:"thread_count"`
	IssueInOneReq int    `yaml:"issue_in_one_request"`
	MinSleep      int    `yaml:"min_sleep"`
	MaxSleep      int    `yaml:"max_sleep"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Env       string       `yaml:"env"`
	LogFile   string       `yaml:"log_file"`
	DBCfg     DBConfig     `yaml:"database"`
	JiraCfg   JiraConfig   `yaml:"jira-connector"`
	ServerCfg ServerConfig `yaml:"server"`
}
