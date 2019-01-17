package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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

		var origin string
		if options.origin != "" {
			origin = options.origin
		} else {
			originURL := *dest
			if dest.Scheme == "wss" {
				originURL.Scheme = "https"
			} else {
				originURL.Scheme = "http"
			}
			origin = originURL.String()
		}

		var historyFile string
		user, err := user.Current()
		if err == nil {
			historyFile = filepath.Join(user.HomeDir, ".ws_history")
		}

		err = connect(username, dest.String(), origin, &readline.Config{
			Prompt:      ":: ",
			HistoryFile: historyFile,
		}, options.insecure)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if err != io.EOF && err != readline.ErrInterrupt {
				os.Exit(1)
			}
		}
	},
}

var username string

var options struct {
	origin       string
	printVersion bool
	insecure     bool
}

func init() {
	cmdClient.Flags().StringVarP(&username, "username", "u", "", "username in chartroom")
	cmdClient.MarkFlagRequired("username")
}

func connect(username, url, origin string, rlConf *readline.Config, allowInsecure bool) error {
	headers := make(http.Header)
	headers.Add("Origin", origin)

	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: allowInsecure,
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

	go sess.readConsole(username)
	go sess.readWebsocket(username)

	return <-sess.errChan
}

func (s *session) readConsole(username string) {
	for {
		line, err := s.rl.Readline()
		line = "[" + username + "]" + " : " + line
		if err != nil {
			s.errChan <- err
			return
		}

		err = s.ws.WriteMessage(+websocket.TextMessage, []byte(line))
		if err != nil {
			s.errChan <- err
			return
		}
	}
}

func (s *session) readWebsocket(username string) {
	rxSprintf := color.New(color.FgGreen).SprintfFunc()

	for {
		msgType, buf, err := s.ws.ReadMessage()
		if err != nil {
			s.errChan <- err
			return
		}

		var text string
		switch msgType {
		case websocket.TextMessage:
			text = string(buf)
		default:
			s.errChan <- fmt.Errorf("unknown websocket frame type: %d", msgType)
			return
		}
		if strings.Index(text, username) < 0 {
			fmt.Fprint(s.rl.Stdout(), rxSprintf(":: %s\n", text))
		}
	}
}
