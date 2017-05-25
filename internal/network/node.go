package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
	me *Member
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

	c.me = &Member{status.New, time.Now().UnixNano(), c.option.Network.Interface}
	return c
}

// start this Cluster
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

func (c *Node) join() {
	for _, value := range c.option.Network.Join["tcp-ip"] {
		// skip self
		if value == c.option.Network.Interface {
			//continue
		}
		conn, err := net.Dial("tcp", value)
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(newJoin(c.me))
		if err != nil {
			log.Fatal("encode error:", err)
		}
		conn.Write(buf.Bytes())
	}
}

func (c *Node) handleConnection(conn net.Conn) {
	d, err := read(conn)
	if err != nil {
		log.Panicln("decode error:", err)
		// TODO handle read err
	}

	// old node will be leader
	if c.me.StartTime < d.Me.StartTime {
		c.me.Status = status.Leader
	} else {
		c.me.Status = status.Member
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(c)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	conn.Write(buf.Bytes())
	conn.Close()
}

func read(conn net.Conn) (CmdJoin, error) {
	cmd := CmdJoin{}
	dec := gob.NewDecoder(conn)
	// Will read from network.
	err := dec.Decode(&cmd)
	fmt.Printf("bin:%d\n", cmd)
	return cmd, err
}
