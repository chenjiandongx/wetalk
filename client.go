package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

type session struct {
	ws      *websocket.Conn
	rl      *readline.Instance
	errChan chan error
}

var cmdClient = &cobra.Command{
	Use:   "client <addr>",
	Short: "start websockets client",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			cmd.Help()
			os.Exit(1)
		}

		dest, err := url.Parse(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		originURL := *dest
		if dest.Scheme == "wss" {
			originURL.Scheme = "https"
		} else {
			originURL.Scheme = "http"
		}
		origin := originURL.String()

		err = connect(name, dest.String(), origin, &readline.Config{Prompt: ":: "})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if err != io.EOF && err != readline.ErrInterrupt {
				os.Exit(1)
			}
		}
	},
}

var name string

func init() {
	cmdClient.Flags().StringVarP(&name, "name", "u", "", "nickname int the chatroom")
	cmdClient.MarkFlagRequired("name")
}

func connect(name, url, origin string, rlConf *readline.Config) error {
	headers := make(http.Header)
	headers.Add("Origin", origin)

	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	ws, _, err := dialer.Dial(url, headers)
	if err != nil {
		return err
	}

	rl, err := readline.NewEx(rlConf)
	if err != nil {
		return err
	}
	defer rl.Close()

	sess := &session{
		ws:      ws,
		rl:      rl,
		errChan: make(chan error),
	}

	go sess.readConsole(name)
	go sess.readWebsocket(name)

	return <-sess.errChan
}

func (s *session) readConsole(name string) {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			s.errChan <- err
			return
		}

		msg, err := json.Marshal(Message{Name: name, Msg: line})
		if err != nil {
			s.errChan <- err
			return
		}
		err = s.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			s.errChan <- err
			return
		}
	}
}

func (s *session) readWebsocket(name string) {
	rxSprintf := color.New(color.FgGreen).SprintfFunc()

	for {
		msgType, buf, err := s.ws.ReadMessage()
		if err != nil {
			s.errChan <- err
			return
		}
		var msg Message
		switch msgType {
		case websocket.TextMessage:
			err := json.Unmarshal(buf, &msg)
			if err != nil {
				s.errChan <- err
				return
			}
		default:
			s.errChan <- fmt.Errorf("unknown websocket frame type: %d", msgType)
			return
		}
		// 不需要给自己广播
		if msg.Name != name {
			fmt.Fprint(s.rl.Stdout(), rxSprintf(":: [%s]: %s\n", msg.Name, msg.Msg))
		}
	}
}
