package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net"
	"os"
	"os/exec"
	"runtime"

	s "strings"
)

type Configuration struct {
	Versionmajor	string
	Versionminor	string
	Versionstate	string
	Port		string
}

const TEMPLATE_FILE = "dict.html";
const CONF_FILE = "dictdispatch.conf"

const X509CERTKEY = "server.key"
const X509CERT  = "server.crt"


func check(e error) {
	if e != nil {
		panic(e)
	}
}

var gConfig Configuration

var httpsmux map[string] func(http.ResponseWriter, *http.Request)

var dispatchDict = map[string] func(http.ResponseWriter, *http.Request) {
	"/" : prog_central,
	"/login" : prog_login,
}

func getIPAddr() (string) {
	var ipAddrStr string
	ifaces, err := net.Interfaces()
	check(err)

	for _, i := range ifaces {
		var ip net.IP
		addrs, err := i.Addrs()
		check(err)

		for _, addr := range addrs {
			switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			ipAddrStr = ip.String()
			if len(ipAddrStr) > 0 {
				ipAddrStr = s.Trim(ipAddrStr, "\n")
			}
		}
	}
	return ipAddrStr;
}


func prog_central(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := ioutil.ReadFile(TEMPLATE_FILE)
	check(err)
	htmlContentStr := string(htmlContent)

	io.WriteString(w, htmlContentStr)
}

func prog_login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "login")
}

func os_banner() {
	fmt.Printf("Startup on ")
	out, err := exec.Command("/bin/hostname").Output()
	check(err)
	fmt.Printf("%s\n", out)

	out, err = exec.Command("date").Output()
	check(err)
	fmt.Printf("%s", out)
	fmt.Printf("Welcome to the HTTP Dispatcher")
	fmt.Printf(" '%s.%s %s'\nRunning on %s:%s with %d CPUs\n", gConfig.Versionmajor, gConfig.Versionminor, gConfig.Versionstate, getIPAddr(), gConfig.Port, runtime.NumCPU())
}


func loadConfig() {
	file, _ := os.Open(CONF_FILE)
	defer file.Close()
	decoder := json.NewDecoder(file)
	gConfig = Configuration{}
	err := decoder.Decode(&gConfig)
	check(err)
}



func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	loadConfig()
	os_banner()

        cert, err := tls.LoadX509KeyPair(X509CERT, X509CERTKEY)
        if err != nil {
                log.Fatal("server: loadkeys: %s", err)
        } else {
                fmt.Println("certs loaded")
                cfg := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
                server := http.Server {
                        Addr:  ":" + gConfig.Port,
                        Handler: &myHandler{},
                        TLSConfig: &cfg,
                        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
                }


                httpsmux = dispatchDict

//              server.ListenAndServe()
                 server.ListenAndServeTLS(X509CERT, X509CERTKEY)
        }

/*
	httpsmux = dispatchDict
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Unable to start HTTP server: %s\n", err)
		os.Exit(1)
	}
*/
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	path := r.URL.Path[1]
fmt.Println(r.URL.RequestURI())

	if h, ok := httpsmux[r.URL.String()]; ok {	// allow list
		h(w,r)
		return
	}
	io.WriteString(w, "Page not found - " + r.URL.String())
}


