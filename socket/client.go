package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

// 客户端连接详情
type wsClients struct {
	Conn       *websocket.Conn `json:"conn"`
	RemoteAddr string          `json:"remote_addr"`
	Uid        float64         `json:"uid"`
	Username   string          `json:"username"`
	RoomId     string          `json:"room_id"`
	AvatarId   string          `json:"avatar_id"`
}

// client 和 serve 的消息体
type msg struct {
	Status int             `json:"status"`
	Data   interface{}     `json:"data"`
	Conn   *websocket.Conn `json:"conn"`
}

type ServeInterface interface {
	RunWs(ctx *gin.Context)
	GetOnlineUserCount() int
	GetOnlineRoomUserCount(roomId int) int
}

// 变量定义初始化
var (
	wsUpgrader = websocket.Upgrader{}
	clientMsg  = msg{}
	mutex      = sync.Mutex{}
	rooms      = make(map[int][]wsClients)
	enterRooms = make(chan wsClients)
	sMsg       = make(chan msg)
	offline    = make(chan *websocket.Conn)
	chNotify   = make(chan int, 1)
)

//  定义消息类型
const msgTypeOnline = 1        // 上线
const msgTypeOffline = 2       // 离线
const msgTupeSend = 3          // 消息发送
const msgTypeGetOnlineUser = 4 // 获取用户列表
const msgTypePrivateChat = 5   // 私聊

const roomCount = 6 // 房间总数

type GoServe struct {
	ServeInterface
}

func (goServe *GoServe) RunWs(ctx *gin.Context) {
	Run(ctx)
}

func (goServe *GoServe) GetOnlineUserCount() int {
	return GetOnlineUserCount()
}

func (goServe *GoServe) GetOnlineRoomUserCount(roomId int) int {
	return GetOnlineRoomUserCount(roomId)
}

func Run(ctx *gin.Context) {
	// 允许跨域
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	go readLoop(conn)
	go writeLoop()

	select {}
}

func readLoop(conn *websocket.Conn) {
	defer func() {
		// 捕获read抛出的panic
		if err := recover(); err != nil {
			log.Println("readLoop发生错误", err)
		}
	}()

	for {
		// 读取消息
		mt, message, err := conn.ReadMessage()
		log.Println(mt, string(message))
		// 离线通知
		if err != nil {
			offline <- conn
			log.Println("用户已离线", err)
			break
		}

		serveMsgStr := message

		// 处理心跳响应,heartbeat为与客户段约定的值
		if string(serveMsgStr) == "heartbeat" {
			// 写入消息
			err = conn.WriteMessage(mt, []byte(`{"status":0,"data":"heartbeat ok"}`))
			if err != nil {
				break
			}
			continue
		}

		_ = json.Unmarshal(message, &clientMsg)
		log.Println("来自客户端的消息", clientMsg, conn.RemoteAddr())
		if clientMsg.Data != nil {
			if clientMsg.Status == msgTypeOnline { //  进入房间，建立链接
				roomId, _ := getRoomId()

				enterRooms <- wsClients{
					Conn:       conn,
					RemoteAddr: conn.RemoteAddr().String(),
					Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
					Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
					RoomId:     roomId,
					AvatarId:   clientMsg.Data.(map[string]interface{})["avatar_id"].(string),
				}
			}
			_, serveMsg := formatServeMsgStr(clientMsg.Status, conn)
			sMsg <- serveMsg
		}
	}
}

func writeLoop() {
	defer func() {
		// 捕获write抛出的panic
		if err := recover(); err != nil {
			log.Println("writeLoop发送错误", err)
		}
	}()

	for {
		select {
		case r := <-enterRooms:
			handleConnClients(r.Conn)
		case cl := <-sMsg:
			serveMsgStr, _ := json.Marshal(cl)
			switch cl.Status {
			case msgTypeOnline, msgTupeSend:
				notify(cl.Conn, string(serveMsgStr))
			case msgTypeGetOnlineUser:
				chNotify <- 1
				cl.Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				<-chNotify
			case msgTypePrivateChat:
				chNotify <- 1
				toC := findToUserConnClient()
				if toC != nil {
					toC.(wsClients).Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				}
				<-chNotify
			}
		case o := <-offline:
			disconnenct(o)
		}
	}
}

func handleConnClients(conn *websocket.Conn) {
	roomId, roomIdInt := getRoomId()

	for ckey, wcl := range rooms[roomIdInt] {
		if wcl.Uid == clientMsg.Data.(map[string]interface{})["uid"].(float64) {
			mutex.Lock()
			// 通知当前影虎下线
			wcl.Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
			rooms[roomIdInt] = append(rooms[roomIdInt][:ckey], rooms[roomIdInt][ckey+1:]...)
			wcl.Conn.Close()
			mutex.Unlock()
		}
	}

	mutex.Lock()
	rooms[roomIdInt] = append(rooms[roomIdInt], wsClients{
		Conn:       conn,
		RemoteAddr: conn.RemoteAddr().String(),
		Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
		Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
		RoomId:     roomId,
		AvatarId:   clientMsg.Data.(map[string]interface{})["avatar_id"].(string),
	})
	mutex.Unlock()
}

// 获取私聊的用户连接
func findToUserConnClient() interface{} {
	_, roomIdInt := getRoomId()
	toUserUid := clientMsg.Data.(map[string]interface{})["to_uid"].(string)

	for _, conn := range rooms[roomIdInt] {
		stringUid := strconv.FormatFloat(conn.Uid, "f", -1, 64)
		if stringUid == toUserUid {
			return conn
		}
	}
	return nil
}

// 统一消息发放
func notify(conn *websocket.Conn, msg string) {
	// 利用channel阻塞 避免并发去对同一个连接发送消息出现panic: concurrent write to websocket connection这样的异常
	chNotify <- 1
	_, roomIdInt := getRoomId()
	for _, con := range rooms[roomIdInt] {
		if con.RemoteAddr != conn.RemoteAddr().String() {
			con.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
	<-chNotify
}

// 离线通知
