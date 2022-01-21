package znet

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connection map[uint32]ziface.IConnection //管理的链接集合
	connLock   sync.RWMutex                  //保护链接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connection: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源 map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入到connManager中
	connMgr.connection[conn.GetConnID()] = conn
	logrus.Infof("connID = %s connection add to ConnManager successful:conn num = %d", conn.GetConnID(), connMgr.Len())
}

//删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源 map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接信息
	delete(connMgr.connection, conn.GetConnID())
	logrus.Infof("connID = %s connection delete from ConnManager successful:conn num = %d", conn.GetConnID(), connMgr.Len())

}

//根据connID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源 map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	if conn, ok := connMgr.connection[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("Connection not Found")
	}

}

//得到当前链接数量
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connection)
}

//清除并终止所有的链接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源 map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connection {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connection, connID)
	}
	logrus.Info("Clear all connection success! conn num = ", connMgr.Len())
}
