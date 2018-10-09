// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"bytes"

	"log"
	"net/http"
	"time"
	"tsEngine/tsJson"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

/*
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
*/

// Client is a middleman between the websocket connection and the hub.
type Message struct {
	Type    int
	User    string
	Channel string
	Content string
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	//用户昵称
	channel string

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func init() {

}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(user, channel string) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		temp := Message{}
		temp.Type = 2
		temp.User = user
		temp.Channel = channel
		temp.Content = string(message)

		msg, err := tsJson.ToString(temp)

		if err != nil {
			beego.Error(err)
			continue
		}
		c.hub.broadcast <- []byte(msg)
		beego.Trace("消息：", msg)
		/*
			if temp.Type == -3 {
				c.hub.broadcast <- []byte(msg)
			}
		*/
		//c.hub.broadcast <- []byte(msg)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			//根据消息判断是哪个频道
			msg := Message{}
			err := tsJson.FromString(string(message), &msg)
			if err != nil {
				continue
			}
			if msg.Channel != c.channel {
				continue
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(user, channel string, hub *Hub, w http.ResponseWriter, r *http.Request) {
	//conn, err := upgrader.Upgrade(w, r, nil)
	conn, err := websocket.Upgrade(w, r, nil, 65536, 65536)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{channel: channel, hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump(user, channel)
}
