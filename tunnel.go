package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/kardianos/osext"
	"github.com/mkideal/cli"
	"github.com/rs/cors"
	"github.com/speps/go-hashids"
)

// Blazon contains methods to publish final output
type Blazon struct {
	response http.ResponseWriter
	callback string
}

func (b Blazon) wrapper(content string) {
	if b.callback != "" {
		b.response.Header().Set("Content-Type", "text/javascript")
		jsonp := b.callback + "(" + content + ")"
		b.response.Write([]byte(jsonp))
	} else {
		b.response.Header().Set("Content-Type", "application/json")
		b.response.Write([]byte(content))
	}
}

func (b Blazon) publish(response string) string {
	dix, _ := json.Marshal(map[string]string{
		"status":   "success",
		"response": response})
	return string(dix)
}

func (b Blazon) trouble(response string) string {
	dix, _ := json.Marshal(map[string]string{
		"status":   "failure",
		"response": response})
	return string(dix)
}

// Console contains terminal input and output handlers.
type Console struct {
}

func (c Console) getCommand(cmd *exec.Cmd) string {
	return strings.Join(cmd.Args, " ")
}

func (c Console) getError(err error) string {
	if err != nil {
		return string(err.Error())
	}
	return ""
}

func (c Console) getOutput(outs []byte) string {
	if len(outs) > 0 {
		return string(outs)
	}
	return ""
}

func (c Console) process(input string) string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("bash", "-c", input)
	case "windows":
		cmd = exec.Command("cmd", "/C", input)
	default:
		cmd = exec.Command(input)
	}
	output, err := cmd.CombinedOutput()
	dix := map[string]string{
		"cmd": c.getCommand(cmd),
		"err": c.getError(err),
		"out": c.getOutput(output)}
	dixMap, _ := json.Marshal(dix)
	return string(dixMap)
}

// Utility functions
type Utility struct{}

func (u Utility) genAuthKey() string {
	hd := hashids.NewData()
	hd.Salt = time.Now().Format(time.RFC3339)
	h, _ := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{1, 2, 3})
	return id
}

func (u Utility) genApiKey(username string, password string) string {
	data := []byte(username + "|" + password)
	id := fmt.Sprintf("%x", md5.Sum(data))
	return id
}

// WebService contains browser specific commands.
type WebService struct {
	authKey string
}

// curl -X POST http://localhost:9999/session
func (ws WebService) sessionToken(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	callback := qp.Get("callback")
	blazon := Blazon{w, callback}
	nextStep := false
	clientApiKey := ""
	if callback != "" {
		clientApiKey = qp.Get("cak")
		nextStep = true
	} else {
		err := r.ParseForm()
		if err != nil {
			blazon.wrapper(blazon.trouble(string(err.Error())))
		} else {
			clientApiKey = r.PostFormValue("cak")
			nextStep = true
		}
	}
	if nextStep {
		if clientApiKey != "" {
			util := Utility{}
			serverApiKey := util.genApiKey(userId, userPw)
			if clientApiKey == serverApiKey {
				if authKey == "" {
					authKey = util.genAuthKey()
				}
				dix := map[string]string{
					"token": authKey}
				output, _ := json.Marshal(dix)
				blazon.wrapper(blazon.publish(string(output)))
			} else {
				blazon.wrapper(blazon.trouble("ClientAPIKey is invalid"))
			}
		} else {
			blazon.wrapper(blazon.trouble("ClientAPIKey is missing"))
		}
	}
}

// curl -d "token=value&cmd=ls" -X POST http://localhost:9999/terminal
func (ws WebService) terminalRun(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	callback := qp.Get("callback")
	blazon := Blazon{w, callback}
	if len(authKey) <= 0 {
		blazon.wrapper(blazon.trouble("Session is not established"))
	} else {
		konsole := Console{}
		nextStep := false
		command := ""
		token := ""
		if callback != "" {
			command = qp.Get("cmd")
			token = qp.Get("token")
			nextStep = true
		} else {
			err := r.ParseForm()
			if err != nil {
				blazon.wrapper(blazon.trouble(string(err.Error())))
			} else {
				command = r.PostFormValue("cmd")
				token = r.PostFormValue("token")
				nextStep = true
			}
		}
		if nextStep {
			if token != "" {
				if token == authKey {
					if command != "" {
						output := konsole.process(command)
						blazon.wrapper(blazon.publish(output))
					} else {
						blazon.wrapper(blazon.trouble("Command is missing"))
					}
				} else {
					blazon.wrapper(blazon.trouble("Token is invalid"))
				}
			} else {
				blazon.wrapper(blazon.trouble("Token is missing"))
			}
		}
	}
}

// Server is an application server
type Server struct {
	docRoot string
	url     string
}

func (s Server) waitServer() bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(s.url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}

func (s Server) startBrowser() bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], s.url)...)
	return cmd.Start() == nil
}

func (s Server) probeDocRoot() string {
	serverRoot, err := osext.ExecutableFolder()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if appRoot == true {
		s.docRoot = serverRoot
		if docPath != "" {
			s.docRoot += docPath
		}
	} else {
		if docPath != "" {
			s.docRoot = docPath
		} else {
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			s.docRoot = pwd
		}
	}
	return s.docRoot
}

func (s Server) initialize() {
	httpAddr := hostIP + ":" + strconv.Itoa(portNum)

	s.url = "http://" + httpAddr
	s.docRoot = s.probeDocRoot()

	timestamp := time.Now()
	fmt.Println(appName, "configuration \n  Root \t", s.docRoot, "\n  URL \t", s.url, "\n  Time \t", timestamp.Format(time.RFC1123), "\n")

	go func() {
		fmt.Println(appName, "status: STARTED")
		if s.waitServer() && openBrowser && s.startBrowser() {
			fmt.Println("A browser window should open. If not, visit the link.")
		} else {
			fmt.Println("Please open your web browser and visit the link.")
		}
		fmt.Println("Please hit 'ctrl + C' to STOP the server.")
	}()

	ws := WebService{}
	if corsEnabled {
		fmt.Println("CORS: enabled")
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowedHeaders:   []string{"Content-Type"},
			AllowCredentials: true,
		})
		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.Dir(s.docRoot)))
		mux.HandleFunc("/session", ws.sessionToken)
		mux.HandleFunc("/terminal", ws.terminalRun)
		handler := c.Handler(mux)
		http.ListenAndServe(httpAddr, handler)
	} else {
		http.Handle("/", http.FileServer(http.Dir(s.docRoot)))
		http.HandleFunc("/session", ws.sessionToken)
		http.HandleFunc("/terminal", ws.terminalRun)
		http.ListenAndServe(httpAddr, nil)
	}
}

var (
	appName     = "Tunnel"
	version     = "1.0.1"
	docPath     = ""
	hostIP      = "127.0.0.1"
	portNum     = 9999
	appRoot     = false
	openBrowser = false
	corsEnabled = true
	authKey     = ""
	userId      = "admin"
	userPw      = "123456"
)

type argT struct {
	cli.Helper
	Port    int    `cli:"p,port" usage:"set custom port number" dft:"9999"`
	Host    string `cli:"u,host" usage:"set host IP or server address" dft:"127.0.0.1"`
	DocPath string `cli:"d,docpath" usage:"set document directory's path" dft:""`
	Browser bool   `cli:"b,browser" usage:"open browser on server start" dft:"false"`
	AppRoot bool   `cli:"a,approot" usage:"serve from application's root" dft:"false"`
	Cors    bool   `cli:"x,cors" usage:"allows cross domain requests" dft:"false"`
	User    string `cli:"i,user" usage:"username of account" dft:"admin"`
	Pass    string `cli:"c,pass" usage:"password of account" dft:"123456"`
}

func main() {
	today := time.Now()
	fmt.Printf("\n%s (Version %s) \nCopyright (c) 2017-%s Abhishek Kumar. \nLicensed under MIT License. \n\n", appName, version, strconv.Itoa(today.Year()))

	mode := false
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		docPath = argv.DocPath
		hostIP = argv.Host
		portNum = argv.Port
		openBrowser = argv.Browser
		appRoot = argv.AppRoot
		corsEnabled = argv.Cors
		userId = argv.User
		userPw = argv.Pass
		mode = true
		return nil
	})

	if mode {
		server := Server{}
		server.initialize()
	}

	fmt.Println("\nDone!\n")
}
