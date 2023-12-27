package data

type Instance struct {
	Id        string
	IpAddress string
	User      string
	Pass      string
	SshPort   int
	Tags      []string
}
