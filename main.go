package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/go-gomail/gomail"
)

func main() {
	resp, _ := http.Get("http://pv.sohu.com/cityjson?ie=utf-8")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	js := string(body)
	ip := strings.Split(js, `=`)[1]
	ipJson := strings.TrimSpace(ip)

	if strings.HasSuffix(ipJson, `;`) {
		ipJson = ipJson[0 : len(ipJson)-1]
	}

	IP,_ := simplejson.NewJson([]byte(ipJson))
	ipAddress,_ := IP.Get("cip").String()
	fmt.Println(ipAddress)
	subject := "家里的IP地址已变更为:"
	content := fmt.Sprintf(`
		<html>
		<body>
		<h3>
		%s
		</h3>
		</body>
		</html>
		`,ipAddress)
	fmt.Println("send email")
	msg := gomail.NewMessage()
	msg.SetHeader("From", "454679718@163.com")
	msg.SetHeader("To", "454679718@qq.com")
	//msg.SetAddressHeader("Cc", "dan@example.com", "Dan")
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	// Send the email to Bob, Cora and Dan
	mailer := gomail.NewDialer("smtp.163.com",25, "454679718@163.com", "brady1988")
	if err := mailer.DialAndSend(msg); err != nil {
		panic(err)
	}
}