package network

import (
	"bytes"
	"encoding/gob"
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

// request for join and cluster members
type Data struct {
	Magic     uint16
	Leader    bool
	StartTime int64
	Members   []Member
}

type Member struct {
	Leader    bool
	StartTime int64
	Interface string
}

// new Cluster instance
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

// start this Cluster
func (c *Cluster) Start() {
	go c.listen()
	time.Sleep(time.Second)
	go c.join()
}

// stop this cluster
func (c *Cluster) Stop() {
	if c.listener != nil {
		c.listener.Close()
	}
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
		// skip self
		if value == c.option.Network.Interface {
			//continue
		}
		conn, err := net.Dial("tcp", value)
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		memberSelf := Member{false, c.startTime, c.option.Network.Interface}
		data := Data{consts.CmdMagic, false, c.startTime, []Member{memberSelf}}
		enc := gob.NewEncoder(&buf) // Will write to network.

		// Encode (send) some values.
		err = enc.Encode(data)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		conn.Write(buf.Bytes())
	}
}

func (c *Cluster) handleConnection(conn net.Conn) {
	d := Data{}
	dec := gob.NewDecoder(conn) // Will read from network.
	err := dec.Decode(&d)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}
	fmt.Printf("%q: {%d, %d}\n", d.StartTime, d.Leader, d.Magic)

	fmt.Printf("bin:%v\n", d)
	conn.Close()
}
