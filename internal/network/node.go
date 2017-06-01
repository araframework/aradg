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
	"io"
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/consts/code"
	"encoding/binary"
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
func (node *Node) Start() {
	go node.listen()
	time.Sleep(time.Second)
	go node.join()
}

// stop this cluster
func (node *Node) Stop() {
	if node.listener != nil {
		node.listener.Close()
	}

	if node.leaderConn != nil {
		node.leaderConn.Close()
	}
}

func (node *Node) listen() {
	port := strconv.FormatUint(uint64(node.option.Network.Listen.Port), 10)
	laddr := node.option.Network.Listen.Ip + ":" + port
	fmt.Println("Starting ", laddr)
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	node.listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go node.handleConnection(conn)
	}
}

func (node *Node) join() {
	timeout := node.option.Network.Members.Timeout
	for _, addr := range node.option.Network.Members.Init {
		ip := net.ParseIP(addr.Ip)
		if ip == nil || bytes.Equal(ip, node.me.Ip) {
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
		cmdJoin := NewCmdJoin(node.me)
		fmt.Printf("cmdJoin:%x\n", cmdJoin)
		conn.Write(cmdJoin)

		// TODO
		ch := make(chan []byte)
		eCh := make(chan error)
		// Start a goroutine to read from our net connection
		go func(ch chan []byte, eCh chan error) {
			for {

				// try to read the data
				buff := make([]byte, 8)
				n, err := conn.Read(buff)
				fmt.Println("frome server: ", n)
				if err != nil {
					// send an error if it's encountered
					fmt.Println(err)
					eCh <- err
					return
				}

				if n <= 0 {
					eCh <- io.EOF
					return
				}
				// send data if we read some.
				ch <- buff
			}
		}(ch, eCh)

		ticker := time.Tick(time.Second * 10)
		// continuously read from the connection

		bodyBuf := bytes.NewBuffer(make([]byte, 0))
		for {
			select {
			// This case means we recieved data on the connection
			case data := <-ch:
				fmt.Printf("from ch:%x\n", data)
				bodyBuf.Write(data)
				fmt.Printf("from buf:%x\n", bodyBuf)
				// Do something with the data
				// This case means we got an error and the goroutine has finished
			case err := <-eCh:
				fmt.Printf("eCh:%x\n", err)
				// handle our error then exit for loop
				break
				// This will timeout on the read.
			case <-ticker:
				// do nothing? this is just so we can time out if we need to.
				// you probably don't even need to have this here unless you want
				// do something specifically on the timeout.
			}
		}

		fmt.Printf("Done: %x\n", bodyBuf)

		conn.Close()
	}
}

func (node *Node) handleConnection(conn net.Conn) {
	buff := bytes.NewBuffer(nil)
	var buf [8]byte
	first := true
	var len, bodyLen, totalLen uint32 = 0, 0, 0
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err != io.EOF {
				log.Println("connection read err:", err)
				// break loop if not an EOF err raised
				break
			}

			log.Println("connection closed by remote")
			if n == 0 {
				// n == 0 with an EOF err means conn closed and no more data available
				break
			}
		}

		len += uint32(n)
		buff.Write(buf[0:n])

		if first {
			first = false
			if buff.Len() < 7 {
				log.Println("Invalid len:", n)
				break
			}
			magic := binary.LittleEndian.Uint16(buff.Bytes()[:2])
			if magic != consts.Magic {
				log.Println("Invalid message")
				break
			}
			bodyLen = binary.LittleEndian.Uint32(buff.Bytes()[3:7])
			totalLen = bodyLen + 7
		}

		// reached the msg end
		if len >= totalLen {
			buff.Truncate(int(totalLen)) // TODO uint32 to int?
			go node.handleMessage(conn, buff.Bytes())

			// reset for next msg
			len, bodyLen = 0, 0
			first = true
			buff.Reset()

			if len > totalLen {
				buff.Write(buf[(uint32(n) - (len - totalLen)):])
			}
		}
	}

	// old node will be leader
	//if node.me.StartTime < cmd {
	//	node.me.Status = status.Leader
	//} else {
	//	node.me.Status = status.Member
	//}
	//
	//var buf bytes.Buffer
	//enc := gob.NewEncoder(&buf)
	//err = enc.Encode(node)
	//if err != nil {
	//	log.Fatal("encode error:", err)
	//}
	//conn.Write(buf.Bytes())
	//conn.Close()
}
func (node *Node) handleMessage(conn net.Conn, msg []byte) {
	fmt.Printf("handleMessage %x\n", msg)
	conn.Write([]byte{0xaa, 0xBB, 0xCC})

	cmdCode := msg[2]
	switch cmdCode {
	case code.Join:
		fmt.Println("Join")
		break
	case code.JoinAck:
		fmt.Println("Join Ack")
		break
	default:
		fmt.Println("Unkown: ", cmdCode)
	}
}
