/**
 * 多协程推送
**/
package push

import (
	"encoding/json"
	"fmt"

	"github.com/qq65326/apns_push/lib"
	"github.com/qq65326/apns_push/server"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

type NotifyNoSign struct {
	PlatForm    string      `json:"platform"`
	ServiceConf ServiceConf `json:"service_conf"`
	ServiceData ServiceData `json:"service_data"`
}

type Notify struct {
	NotifyNoSign
	Sign string `json:"sign"`
}

type ServiceConf struct {
	Cert         string `json:"cert"`
	CertPassword string `json:"cert_password"`
	BundleID     string `json:"bundle_id"`
	DeviceToken  string `json:"device_token"`
	IsDev        bool   `json:"is_dev"`
}

type Payload struct {
	Title          string      `json:"title"`
	MobilePushData interface{} `json:"data"`
	Message        string      `json:"message"`
}

type ServiceData struct {
	Title    string  `json:"title"`
	Subtitle string  `json:"subtitle"`
	Body     string  `json:"body"`
	Badge    int     `json:"badge"`
	Sound    string  `json:"sound"`
	Payload  Payload `json:"payload"`
}

var (
	Client *apns2.Client
	// 推送协程数量
	ChanMax = 10000
	// 生产证书密码配置
	Cert_Production_Conf = map[string]string{
		"push.p12": "password",
	}
	// 证书存放路径
	Cert_Path = "/data/conf/apns_push_ssl/"

	// 开发客户端集合
	Development_Clients = make(map[string]*apns2.Client)
	// 生产客户端集合
	Production_Clients = make(map[string]*apns2.Client)
)

/**
 * 开启服务
 */
func StartWorker() {
	lib.MainWaitGroup.Add(1)
	// 控制协程数量10000
	ch := make(chan int, ChanMax)
	// 加载对应证书，生成对应的httpClient
	for key, value := range Cert_Production_Conf {
		cert, err := certificate.FromP12File(Cert_Path+key, value)
		if err != nil {
			fmt.Println("Error: load cert failed! file name:" + key)
			return
		}
		Production_Clients[key] = apns2.NewClient(cert).Production()
	}
	// 监听队列开启协程推送
	for {
		notification := <-server.Notification
		ch <- 1
		lib.PushWaitGroup.Add(1)
		go workerPush(ch, notification)
	}

}

/**
 * 关闭服务
 */
func Stop() {
	lib.PushWaitGroup.Wait()
	lib.MainWaitGroup.Done()
}

/**
 * 协程推送
 * @param ch 消息队列
 * @paramnotification 消息
 */
func workerPush(ch chan int, notification string) {

	var notifyNoSign NotifyNoSign
	error := json.Unmarshal([]byte(notification), &notifyNoSign)
	if error != nil {
		lib.PushWaitGroup.Done()
		return
	}

	if notifyNoSign.ServiceConf.IsDev {
		Client = Development_Clients[notifyNoSign.ServiceConf.Cert]
	} else {
		Client = Production_Clients[notifyNoSign.ServiceConf.Cert]
	}

	// 组装推送数据
	pushNotification := &apns2.Notification{}
	pushNotification.DeviceToken = notifyNoSign.ServiceConf.DeviceToken
	pushNotification.Topic = notifyNoSign.ServiceConf.BundleID

	PrePayload := payload.NewPayload()
	PrePayload.AlertTitle(notifyNoSign.ServiceData.Title)
	// PrePayload.AlertSubtitle(notifyNoSign.ServiceData.Subtitle)
	PrePayload.AlertBody(notifyNoSign.ServiceData.Body)

	PrePayload.Badge(notifyNoSign.ServiceData.Badge)
	PrePayload.Sound(notifyNoSign.ServiceData.Sound)

	payloadJson, err := json.Marshal(notifyNoSign.ServiceData.Payload) //验证签名的json
	if err != nil {
	}
	PrePayload.Custom("payload", payloadJson)

	pushNotification.Payload = PrePayload
	// 推送
	res, err := Client.Push(pushNotification)
	if err != nil {
		server.Notification <- notification
	} else {
		if res.StatusCode != 200 {
		}
	}
	defer outChan(ch)
	lib.PushWaitGroup.Done()
}

/**
 * 协程发送完毕，控制信号-1
 */
func outChan(ch chan int) {
	<-ch
}
