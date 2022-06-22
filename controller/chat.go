package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"in-server/model"
	"in-server/service"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Node 本核心在于形成userid和Node的映射关系
type Node struct {
	Conn *websocket.Conn
	//并行转串行,
	DataQueue chan []byte
	GroupSets mapset.Set
}

//userid和Node映射关系表
var clientMap = make(map[uint64]*Node, 0)

var MsgDeal = service.Message{}

//读写锁
var rwLocker sync.RWMutex

func dispatch(data []byte) error {
	msg := model.Message{}
	//fmt.Println(string(data))
	err := json.Unmarshal(data, &msg)
	msg.CreateAt = time.Now()
	if err != nil {
		return errors.New("in dispatch json unmarshal error: " + err.Error())
	}
	msg.IsAd = tool.JudgeIsSpam(msg.Content)
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	switch msg.Type {
	case 0:
		err = SendMsg(msg.ToId, marshal, &msg)
		if err != nil {
			return errors.New("in dispatch sendMsg error, err = " + err.Error())
		}

	}
	fmt.Println(12345)
	go func() {
		err = MsgDeal.UpdateMessage(&msg)
		if err != nil {
			fmt.Println("in dispatch update message error, err = " + err.Error())
		}
	}()
	return nil
}

// Chat 聊天函数
func Chat(ctx *gin.Context) {
	//获取userid
	suserid := ctx.Query("userid")
	userid, _ := strconv.Atoi(suserid)
	// 获取websocket连接
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "websocket upgrade error",
		})
		return
	}
	// 创建节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 100),
		GroupSets: mapset.NewSet(),
	}

	// 如果之前有节点，则关闭之前链接
	if _, ok := clientMap[uint64(userid)]; ok {
		clientMap[uint64(userid)].Conn.Close()
		delete(clientMap, uint64(userid))
	}

	// 将节点添加到映射表
	rwLocker.Lock()
	clientMap[uint64(userid)] = node
	rwLocker.Unlock()

	//fmt.Println(clientMap)

	// 启动读协程
	go func() {
		err = recvproc(node, uint64(userid))
		if err != nil {
			fmt.Println("in Chat recvproc error, err = " + err.Error())
		}
	}()

	// 启动写协程
	go func() {
		err = sendproc(node)
		if err != nil {
			fmt.Println("in Chat sendproc error, err = " + err.Error())
		}
	}()
	marshal, err := json.Marshal("welcome!" + strconv.Itoa(userid))
	if err != nil {
		return
	}
	err = SendMsg(uint64(userid), marshal, nil)
	if err != nil {
		return
	}
}

//发送逻辑
func sendproc(node *Node) error {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return errors.New("in sendproc write message error, err = " + err.Error())
			}
		}
	}
}

//接收逻辑
func recvproc(node *Node, uid uint64) error {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			if err.Error() == "websocket: close 1001 (going away)" {
				fmt.Println("websocket连接断开")
				return errors.New("websocket连接断开")
			}
			return errors.New("in recvproc read message error: " + err.Error())
		}
		err = dispatch(data)
		if err != nil {
			return errors.New("in recvproc dispatch error: " + err.Error())
		}
		//todo对data进一步处理
		fmt.Printf("recv<=%s\n", string(data))
	}
}

func SendMsg(userid uint64, data []byte, message *model.Message) error {
	rwLocker.RLock()
	defer rwLocker.RUnlock()
	node, ok := clientMap[userid]
	if !ok {
		err := MsgDeal.UidNotOnline(message)
		if err != nil {
			return err
		}
		return nil
	}
	node.DataQueue <- data
	return nil
}
