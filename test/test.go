package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

func main() {
	var stime, etime, gotime int64

	cert, err := certificate.FromP12File("../public/source/cert_1.p12", "HelloHuoban88")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	notification := &apns2.Notification{}
	notification.DeviceToken = "7b5155050b74648d370854a884cce7b5281dd88709b0f06afc0517279bd8e2ff"
	notification.Topic = "com.huoban.cloudgrid"
	notification.Payload = []byte(`{"aps":{"alert":{"title":"Hello!","body":"Bob wants to play poker"},"badge":1,"sound":"default"}}`) // See Payload section below

	client := apns2.NewClient(cert).Production()

	stime = time.Now().UnixNano()
	fmt.Printf("时间戳（纳秒）：%v;\n", stime)
	res, err := client.Push(notification)

	etime = time.Now().UnixNano()
	fmt.Printf("时间戳（纳秒）：%v;\n", etime)
	gotime = etime - stime
	fmt.Printf("时间戳（纳秒）：%v;\n", gotime)

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)

}
