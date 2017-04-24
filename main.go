package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/go-gomail/gomail"
	"os"
	"time"
	"flag"
	"log"
)
var oldIp = ""
func main() {
	interval := flag.Int64("i",30,"检查间隔,单位:秒")
	flag.Parse()
	log.Println("已经设置检查间隔:",*interval)
	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
		if err:=recover();err!=nil{
			log.Println(err) // 这里的err其实就是panic传入的内容
		}
	}()
	ticker := time.NewTicker(time.Second * time.Duration(*interval))
	for _ = range ticker.C{
		ip := getPublicIp()
		if len(ip) > 0 && oldIp != ip{
			sendEmail(ip)
		}
		oldIp = ip
	}
}
func sendEmail(ip string){
	hostName,err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	subject := "[" + time.Now().Format("2006-01-02 15:04:05") + "]" + hostName + "IP地址已变更为:"
	content := fmt.Sprintf(`
		<html>
		<body>
		<h3>
		%s
		</h3>
		</body>
		</html>
		`,ip)
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
		log.Println(err)
	}
}
func getPublicIp() string{
	resp, err := http.Get("http://pv.sohu.com/cityjson?ie=utf-8")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	js := string(body)
	ip := strings.Split(js, `=`)[1]
	ipJson := strings.TrimSpace(ip)

	if strings.HasSuffix(ipJson, `;`) {
		ipJson = ipJson[0 : len(ipJson)-1]
	}

	IP,err := simplejson.NewJson([]byte(ipJson))
	if err != nil {
		return ""
	}
	ipAddress,err := IP.Get("cip").String()
	if err != nil {
		return ""
	}
	return ipAddress
}