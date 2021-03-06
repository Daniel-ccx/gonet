package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

import (
	"cfg"
	"helper"
	"misc/timer"
	. "types"
)

const (
	DEFAULT_MQ_SIZE   = 128
	DEFAULT_FLUSH_OPS = 10
)

func init() {
	log.SetPrefix("[GS]")
}

//----------------------------------------------- Start Agent when a client is connected
func StartAgent(in chan []byte, conn net.Conn) {
	defer helper.PrintPanicStack()

	config := cfg.Get()
	if config["profile"] == "true" {
		helper.SetMemProfileRate(1)
		defer func() {
			helper.GC()
			helper.DumpHeap()
			helper.PrintGCSummary()
		}()
	}

	var sess Session
	sess.IP = net.ParseIP(strings.Split(conn.RemoteAddr().String(), ":")[0])
	log.Println("connected from", sess.IP)
	sess.MQ = make(chan IPCObject, DEFAULT_MQ_SIZE)
	sess.ConnectTime = time.Now()
	sess.LastPacketTime = time.Now().Unix()
	sess.KickOut = false

	// standard 1-sec timer
	std_timer := make(chan int32, 1)
	timer.Add(1, time.Now().Unix()+1, std_timer)

	// write buffer
	bufctrl := make(chan bool)
	buf := NewBuffer(&sess, conn, bufctrl)
	go buf.Start()

	// max # of operartions allowed before flushing
	flush_ops, err := strconv.Atoi(config["flush_ops"])
	if err != nil {
		log.Println("cannot parse flush_ops from config", err)
		flush_ops = DEFAULT_FLUSH_OPS
	}

	// cleanup work
	defer func() {
		close_work(&sess)
		close(bufctrl)
	}()

	// the main message loop
	for {
		select {
		case msg, ok := <-in: // message from network
			if !ok {
				return
			}

			if result := UserRequestProxy(&sess, msg); result != nil {
				err := buf.Send(result)
				if err != nil {
					return
				}
			}
			sess.LastPacketTime = time.Now().Unix()

		case msg, ok := <-sess.MQ: // internal message -- IPCObject
			if !ok {
				return
			}

			if result := IPCRequestProxy(&sess, &msg); result != nil {
				err := buf.Send(result)
				if err != nil {
					return
				}
			}

		case <-std_timer: // per-session 1-sec timer
			timer_work(&sess)
			if session_timeout(&sess) {
				return
			}
			timer.Add(1, time.Now().Unix()+1, std_timer)
		}

		// Data Persistence #1: Too many operations , flush to db.
		if sess.OpCount > flush_ops {
			_flush(&sess)
		}

		// kicked out?
		if sess.KickOut {
			return
		}
	}
}
