package web

type ConfigStruct struct {
	User string `json:"user"`
	Password string `json:"password"`
	Host  string `json:"host"`
	Port     string `json:"port"`
	Nodes    []Node `json:"nodes"`
}

type Node struct {
	Name   string `json:"name"`
	Host     string `json:"host"`
	Port   string `json:"port"`
	Key    string `json:"key"`
}

type PageData struct {
	Title      string
	Output     string
	Error      string
	Jails      []Jails
	Config     *ConfigStruct
	Nodes      []Node
	ActiveNode *Node
}

type JailSettings struct {
	JID     string
	Name    string
	Boot    string
	Prio    string
	State   string
	Type    string
	IP      string
	Ports   string
	Release string
	Tags    string
}

type Jails struct {
	Jail JailSettings
}