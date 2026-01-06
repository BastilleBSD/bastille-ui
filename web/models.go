package web

type PageData struct {
	Title  string
	Output string
	Error  string
	Jails []Jails
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