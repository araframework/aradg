package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/utils/conf"
	"log"
	"net"
	"time"
)

type Cluster struct {
	listener  net.Listener
	option    *conf.ClusterOption
	leader    bool
	startTime int64
}

type Data struct {
	Magic  uint16
	Leader bool
}

func NewCluster() *Cluster {
	c := &Cluster{}
	c.option = conf.LoadCluster()
	if c.option == nil {
		log.Fatal("Initialize conf failed.")
	}
	c.leader = false
	c.startTime = time.Now().UnixNano()
	return c
}

func (c *Cluster) Start() {
	go c.listen()
	time.Sleep(time.Second)
	go c.join()
}

func (c *Cluster) listen() {
	listener, err := net.Listen("tcp", c.option.Network.Interface)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go c.handleConnection(conn)
	}
}

func (c *Cluster) join() {
	for _, value := range c.option.Network.Join["tcp-ip"] {
		if value == c.option.Network.Interface {
			//continue
		}
		conn, err := net.Dial("tcp", value)
		if err != nil {
			log.Fatal(err)
		}

		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.LittleEndian, Data{consts.CmdMagic, true})
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(buf.Bytes())
	}
}

func (c *Cluster) Stop() {
	if c.listener != nil {
		c.listener.Close()
	}
}

func (c *Cluster) handleConnection(conn net.Conn) {
	//d := Data{}
	d := make([]byte, 8)
	i, err := conn.Read(d)
	fmt.Println(i, err)
	//err := binary.Read(conn, binary.LittleEndian, d)
	//fmt.Println("read done")
	//if err != nil {
	//	fmt.Println("binary.Read failed:", err)
	//}

	fmt.Printf("bin:% x\n", d)
	conn.Close()
}
