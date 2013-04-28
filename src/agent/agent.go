package agent

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

import (
	"cfg"
	"hub/names"
	. "types"
)

func send(conn net.Conn, p []byte) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(len(p)))
	_, err := conn.Write(header)
	if err != nil {
		log.Println("Error send reply header:", err.Error())
		return err
	}

	_, err = conn.Write(p)
	if err != nil {
		log.Println("Error send reply msg:", err.Error())
		return err
	}

	return nil
}


func _timer(interval int, ch chan string) {
	defer func() {
		recover()
	}()

	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch <- "ding!dong!"
	}
}

func StartAgent(in chan []byte, conn net.Conn) {
	var sess Session
	sess.MQ = make(chan interface{}, 128)
	sess.ServerMQ = make(chan []byte, 128)

	config := cfg.Get()

	// session timeout
	timer_ch_session := make(chan string)
	session_timeout := 30 // sec
	if config["session_timeout"] != "" {
		session_timeout, _ = strconv.Atoi(config["session_timeout"])
	}

	go _timer(session_timeout, timer_ch_session)

L:
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				break L
			}

			if result := UserRequestProxy(&sess, msg); result != nil {
				fmt.Println(result)
				err := send(conn, result)
				if err != nil {
					break L
				}
			}

		case msg, ok := <-sess.MQ:
			if !ok {
				break L
			}

			if result := IPCRequestProxy(&sess, msg); result != nil{
				fmt.Println(result)
				err := send(conn, []byte(result))
				if err != nil {
					break L
				}
			}

		case msg, ok := <-sess.ServerMQ:
			if !ok {
				break L
			}

			err := send(conn, msg)
			if err != nil {
				break L
			}

		case _ = <-timer_ch_session:
			if session_work(&sess,session_timeout) {
				conn.Close()
			}
		}
	}

	// cleanup
	names.Unregister(sess.User.Id)
	close(timer_ch_session)
}