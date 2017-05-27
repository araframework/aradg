package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/araframework/aradg/internal/consts/code"
	"github.com/araframework/aradg/internal/consts/status"
	"github.com/araframework/aradg/internal/utils/conf"
	"log"
	"net"
	"time"
)

type Node struct {
	// configurations
	option   *conf.ClusterOption
	listener net.Listener
	// self
	me Member
	// latest cluster info
	cluster *Cluster
	// connection to leader
	leaderConn net.Conn
}

// new Cluster instance
func NewNode() *Node {
	c := &Node{}
	c.option = conf.LoadCluster()
	if c.option == nil {
		log.Fatal("Initialize conf failed.")
	}

	ip := net.ParseIP(c.option.Network.Listen.Ip)
	port := c.option.Network.Listen.Port
	c.me = Member{status.New, time.Now().UnixNano(), ip, port}
	return c
}

// start this node
func (c *Node) Start() {
	go c.listen()
	time.Sleep(time.Second)
	go c.join()
}

// stop this cluster
func (c *Node) Stop() {
	if c.listener != nil {
		c.listener.Close()
	}

	if c.leaderConn != nil {
		c.leaderConn.Close()
	}
}

func (c *Node) listen() {
	laddr := c.option.Network.Listen.Ip + ":" + string(c.option.Network.Listen.Port)
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	c.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go c.handleConnection(conn)
	}
}

func (c *Node) join() {
	timeout := c.option.Network.Members.Timeout
	for _, addr := range c.option.Network.Members.Init {
		ip := net.ParseIP(addr.Ip)
		if ip == nil || bytes.Equal(ip, c.me.Ip) {
			// TODO skip self
			//continue
		}

		// remote address
		raddr := addr.Ip + ":" + string(addr.Port)
		duration := time.Millisecond * time.Duration(timeout)
		conn, err := net.DialTimeout("tcp", raddr, duration)
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		cmdJoin := newCmdJoin(c.me)
		fmt.Printf("cmdJoin:%v\n", cmdJoin)
		err = enc.Encode(cmdJoin)
		fmt.Printf("node:%v\n", c)
		fmt.Printf("me:%v\n", c.me)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		conn.Write(buf.Bytes())
	}
}

func (c *Node) handleConnection(conn net.Conn) {
	cmd, err := read(conn)
	if err != nil {
		log.Panicln("decode error:", err)
		// TODO handle read err
	}

	switch cmd.Code {
	case code.Join:
		fmt.Println(cmd)
	}

	// old node will be leader
	//if c.me.StartTime < cmd {
	//	c.me.Status = status.Leader
	//} else {
	//	c.me.Status = status.Member
	//}
	//
	//var buf bytes.Buffer
	//enc := gob.NewEncoder(&buf)
	//err = enc.Encode(c)
	//if err != nil {
	//	log.Fatal("encode error:", err)
	//}
	//conn.Write(buf.Bytes())
	//conn.Close()
}

func read(conn net.Conn) (*CmdWrap, error) {
	cmd := &CmdWrap{}
	dec := gob.NewDecoder(conn)
	// Will read from network.
	err := dec.Decode(cmd)
	fmt.Printf("bin:%v\n", cmd)
	return cmd, err
}
