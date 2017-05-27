package network

import (
	"bytes"
	"fmt"
	"github.com/araframework/aradg/internal/consts/status"
	"github.com/araframework/aradg/internal/utils/conf"
	"log"
	"net"
	"strconv"
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
	c.me = Member{status.New, uint64(time.Now().UnixNano()), ip, port}
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
	port := strconv.FormatUint(uint64(c.option.Network.Listen.Port), 10)
	laddr := c.option.Network.Listen.Ip + ":" + port
	fmt.Println("Starting ", laddr)
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
		port := strconv.FormatUint(uint64(addr.Port), 10)
		raddr := addr.Ip + ":" + port
		// dial timeout in milliseconds
		duration := time.Millisecond * time.Duration(timeout)
		fmt.Println("Dialing ", raddr)
		conn, err := net.DialTimeout("tcp", raddr, duration)
		if err != nil {
			// TODO start a task to dial with interval
			log.Println(err)
			continue
		}

		// send join request
		cmdJoin := newCmdJoin(c.me).encode()
		fmt.Printf("cmdJoin:%x\n", cmdJoin)
		conn.Write(cmdJoin)

		// TODO
		ch := make(chan []byte)
		eCh := make(chan error)
		// Start a goroutine to read from our net connection
		go func(ch chan []byte, eCh chan error) {
			for {
				// try to read the data
				buff := make([]byte, 30)
				_, err := conn.Read(buff)
				fmt.Println(":::", buff)
				if err != nil {
					// send an error if it's encountered
					eCh <- err
					return
				}
				// send data if we read some.
				ch <- buff
			}
		}(ch, eCh)

		ticker := time.Tick(time.Second)
		// continuously read from the connection
		for {
			select {
			// This case means we recieved data on the connection
			case data := <-ch:
				fmt.Printf("received 1:%x\n", data)
				// Do something with the data
				// This case means we got an error and the goroutine has finished
			case err := <-eCh:
				fmt.Printf("err 1:%x\n", err)
				// handle our error then exit for loop
				break
				// This will timeout on the read.
			case <-ticker:
				// do nothing? this is just so we can time out if we need to.
				// you probably don't even need to have this here unless you want
				// do something specifically on the timeout.
			}
		}

		conn.Close()
	}
}

func (c *Node) handleConnection(conn net.Conn) {

	fmt.Println("got 2 ", conn.RemoteAddr().String())
	buff := make([]byte, 30)
	_, err := conn.Read(buff)
	fmt.Printf("fff:%v\n", buff)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("received:%x\n", buff)
	conn.Write([]byte{1, 2, 3})

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
