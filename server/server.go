/**
 * 接收消息入任务池
**/

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/qq65326/apns_push/gorilla/mux"
	lib "github.com/qq65326/apns_push/lib"
	util "github.com/qq65326/apns_push/util"
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

type Server struct {
	httpAddr   string
	httpServer *http.Server
	router     *mux.Router
	listener   *net.TCPListener
}

var (
	// 协程池 长度100000
	Notification = make(chan string, 100000)
	// 监听服务
	s *Server
)

// 签名加密盐
const Salt = "P7N25GlI"

/**
 * 主服务
 */
func MainServer(host string, port string) (s *Server) {
	addr := host
	if port != "" {
		addr += ":" + port
	}

	s = &Server{
		httpAddr: addr,
		router:   mux.NewRouter(),
	}

	s.router.HandleFunc("/notification/send", receiveNotification).Methods("post")

	s.httpServer = &http.Server{
		Addr:    s.httpAddr,
		Handler: s.router,
	}
	return s
}

/**
 *  接受消息服务开启
 */
func (s *Server) Start() {
	// 主流程main信号
	lib.MainWaitGroup.Add(1)
	go func() {
		laddr, err := net.ResolveTCPAddr("tcp", s.httpServer.Addr)
		if nil != err {
			fmt.Println("solve tcp addr error")
		}
		listener, err := net.ListenTCP("tcp", laddr)
		if nil != err {
			fmt.Println("listen tcp addr error")
		}
		s.listener = listener
		s.httpServer.Serve(listener)
	}()

}

/**
 * 接受消息服务关闭
 */
func (s *Server) Stop() {
	s.listener.Close()
	lib.MainWaitGroup.Done()
}

/**
 * 接受消息handler
 */
func receiveNotification(w http.ResponseWriter, r *http.Request) {
	con, _ := io.ReadAll(r.Body) //获取post的数据
	defer r.Body.Close()

	data := string(con) //获取到的json字符串

	var notification Notify

	err := json.Unmarshal([]byte(data), &notification)

	if err != nil {
		util.AssembleResult(w, 101, "data json cannot unmarshal")
		return
	}

	var notifyNoSign NotifyNoSign

	error := json.Unmarshal([]byte(data), &notifyNoSign)

	if error != nil {
		util.AssembleResult(w, 102, "data json cannot unmarshal")
		return
	}

	_, err := json.Marshal(notifyNoSign) //验证签名的json
	if err != nil {
		util.AssembleResult(w, 103, "data check sign json cannot marshal")
		return
	}

	Notification <- data
	util.AssembleResult(w, 200, "success")
	return

}
