package socket

import (
	"encoding/json"
	"goAdmin/controller"
	"goAdmin/util"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	ImgUrl     string          `json:"img_url"`
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
const msgTypeHeart = 0         // 心跳
const msgTypeOnline = 1        // 上线
const msgTypeOffline = 2       // 离线
const msgTupeSend = 3          // 群聊消息发送
const msgTypeGetOnlineUser = 4 // 获取用户列表
const msgTypePrivateChat = 5   // 私聊

const roomCount = 4 // 房间总数

var key = []byte("xViwLY6Auz1wd9sU")
var iv = []byte("0BiT8j0ZelK4CjHa")

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
		// 离线通知
		if err != nil {
			offline <- conn
			log.Println("用户已离线", err)
			break
		}

		// 解密消息
		serveMsgStr, err := util.AesDecrypt(message, key, iv)
		if err != nil {
			conn.Close() // 关闭ws连接
			log.Println("解密失败", err)
		}

		// 处理心跳响应,heartbeat为与客户端约定的值
		if string(serveMsgStr) == "heartbeat" {
			// 写入消息
			mes := []byte(`{"status":0,"data":"heartbeat ok"}`)
			enMsg, _ := util.AesEncrypt(mes, key, iv)
			err = conn.WriteMessage(mt, enMsg)
			if err != nil {
				break
			}
			continue
		}

		_ = json.Unmarshal(serveMsgStr, &clientMsg)
		log.Println("来自客户端的消息", clientMsg, conn.RemoteAddr())

		if clientMsg.Data != nil {
			roomId, _ := getRoomId()
			if clientMsg.Status == msgTypeOnline { //  进入房间，建立链接
				enterRooms <- wsClients{
					Conn:       conn,
					RemoteAddr: conn.RemoteAddr().String(),
					Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
					Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
					RoomId:     roomId,
					ImgUrl:     clientMsg.Data.(map[string]interface{})["img_url"].(string),
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
			enMsg, err := util.AesEncrypt(serveMsgStr, key, iv)
			if err != nil {
				cl.Conn.Close() // 关闭ws连接
				log.Println(err, "加密失败")
			}
			switch cl.Status {
			case msgTypeOnline, msgTupeSend:
				notify(cl.Conn, string(enMsg))
			case msgTypeGetOnlineUser:
				chNotify <- 1
				cl.Conn.WriteMessage(websocket.TextMessage, enMsg)
				<-chNotify
			case msgTypePrivateChat:
				chNotify <- 1
				toC := findToUserConnClient()
				if toC != nil {
					toC.(wsClients).Conn.WriteMessage(websocket.TextMessage, enMsg)
				}
				<-chNotify
			}
		case o := <-offline:
			disconnect(o)
		}
	}
}

// 处理建立连接的用户
func handleConnClients(conn *websocket.Conn) {
	roomId, roomIdInt := getRoomId()
	for ckey, wcl := range rooms[roomIdInt] {
		if wcl.Uid == clientMsg.Data.(map[string]interface{})["uid"].(float64) {
			mutex.Lock()
			// 通知当前用户下线
			// wcl.Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
			rooms[roomIdInt] = append(rooms[roomIdInt][:ckey], rooms[roomIdInt][ckey+1:]...)
			// wcl.Conn.Close()
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
		ImgUrl:     clientMsg.Data.(map[string]interface{})["img_url"].(string),
	})
	mutex.Unlock()
}

// 获取私聊的用户连接
func findToUserConnClient() interface{} {
	_, roomIdInt := getRoomId()
	toUserUid := clientMsg.Data.(map[string]interface{})["to_uid"].(string)

	for _, conn := range rooms[roomIdInt] {
		stringUid := strconv.FormatFloat(conn.Uid, 'f', -1, 64)
		if stringUid == toUserUid {
			return conn
		}
	}
	return nil
}

// 统一消息发放
func notify(conn *websocket.Conn, message string) {
	// 利用channel阻塞 避免并发去对同一个连接发送消息出现panic: concurrent write to websocket connection这样的异常
	chNotify <- 1
	roomId, roomIdInt := getRoomId()

	for _, con := range rooms[roomIdInt] {
		if con.RemoteAddr != conn.RemoteAddr().String() {
			con.Conn.WriteMessage(websocket.TextMessage, []byte(message))
			// 只要收到发的信息就推送chatIndex和房间id
			if clientMsg.Status == msgTupeSend || clientMsg.Status == msgTypePrivateChat {
				infos := map[string]interface{}{
					"data": map[string]interface{}{
						"room_id":   roomId,
						"uid":       clientMsg.Data.(map[string]interface{})["uid"].(float64),
						"chatIndex": clientMsg.Data.(map[string]interface{})["chatIndex"],
						"code":      200,
						"msg":       "消息发送成功",
					},
				}
				jsonStr, _ := json.Marshal(infos)
				enMsg, _ := util.AesEncrypt(jsonStr, key, iv)
				err := conn.WriteMessage(websocket.TextMessage, []byte(enMsg))
				if err != nil {
					log.Println(err, "err")
				}
			}
		}
	}

	<-chNotify
}

// 离线通知
func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomId()
	for index, con := range rooms[roomIdInt] {
		if con.RemoteAddr == conn.RemoteAddr().String() {
			data := map[string]interface{}{
				"username": con.Username,
				"uid":      con.Uid,
				"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
			}

			jsonStrServeMsg := msg{
				Status: msgTypeOffline,
				Data:   data,
			}
			serveMsgStr, _ := json.Marshal(jsonStrServeMsg)
			enMsg, _ := util.AesEncrypt(serveMsgStr, key, iv)
			disMsg := string(enMsg)
			mutex.Lock()
			rooms[roomIdInt] = append(rooms[roomIdInt][:index], rooms[roomIdInt][index+1:]...)
			con.Conn.Close()
			mutex.Unlock()
			notify(conn, disMsg)
		}
	}
}

// 格式化传递给客户端的消息数据
func formatServeMsgStr(status int, conn *websocket.Conn) ([]byte, msg) {
	roomId, roomIdInt := getRoomId()

	data := map[string]interface{}{
		"username":    clientMsg.Data.(map[string]interface{})["username"].(string),
		"uid":         clientMsg.Data.(map[string]interface{})["uid"].(float64),
		"img_url":     clientMsg.Data.(map[string]interface{})["img_url"].(string),
		"room_id":     roomId,
		"remote_addr": conn.RemoteAddr().String(),
		"time":        time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
	}

	if status == msgTupeSend || status == msgTypePrivateChat {
		data["content"] = clientMsg.Data.(map[string]interface{})["content"].(string)

		toUidStr := clientMsg.Data.(map[string]interface{})["to_uid"].(string)
		toUid, _ := strconv.Atoi(toUidStr)

		// 保存消息
		stringUid := strconv.FormatFloat(data["uid"].(float64), 'f', -1, 64)
		intUid, _ := strconv.Atoi(stringUid)

		if _, ok := clientMsg.Data.(map[string]interface{})["image_url"].(string); ok {
			data["image_url"] = clientMsg.Data.(map[string]interface{})["image_url"].(string)
			// 存在图片
			controller.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"room_id":    data["room_id"],
				"content":    data["content"],
				"username":   data["username"],
				"img_url":    data["img_url"],   // 用户头像
				"image_url":  data["image_url"], // 传的图片
			})
		} else {
			controller.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data["content"],
				"username":   data["username"],
				"room_id":    data["room_id"],
				"img_url":    data["img_url"],
			})
		}
	}

	if status == msgTypeGetOnlineUser {
		ro := rooms[roomIdInt]
		data["count"] = len(ro)
		data["list"] = ro
	}

	jsonStrServeMsg := msg{
		Status: status,
		Data:   data,
		Conn:   conn,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

	return serveMsgStr, jsonStrServeMsg

}

func getRoomId() (string, int) {
	roomId := clientMsg.Data.(map[string]interface{})["room_id"].(string)
	roomIdInt, _ := strconv.Atoi(roomId)
	return roomId, roomIdInt
}

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i < roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}

func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
