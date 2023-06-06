package module

type vpnConnect struct {
	target     string
	users      []string
	canGetUser bool
	isVul      bool
	timeout    int
	check      bool
	change     bool
	cookie     string
}
