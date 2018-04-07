package main

import (
	"bufio"
	"fmt"
	"os"
)

var addr string = "172.16.19.10"
var port string = "80"
var no_proxy string
var kind []string = []string{"https", "http", "socks", "ftp"}

func print_help() {
	fmt.Println("Help:")
	fmt.Println("	proxy_manager addr <address> : To set the address")
	fmt.Println("	proxy_manager port <address> : To set the port")
	fmt.Println("	proxy_manager on : To activate proxies")
	fmt.Println("	proxy_manager off : To deactivate the proxies")
}
func set_addr() {
	addr = os.Args[2]
	return
}
func set_port() {
	port = os.Args[2]
	return
}

func proxy_on() {
	if addr == "" || port == "" {
		print_help()
		return
	}
	proxy := fmt.Sprintf("%v:%v", addr, port)
	var proxies []string
	for _, val := range kind {
		proxies = append(proxies, fmt.Sprintf("Acquire::%v::Proxy \"%v://%v\";\n", val, val, proxy))
	}
	file, err := os.Open("/etc/apt/apt.conf")
	defer file.Close()
	if err != nil {
		fmt.Print(err)
	}
	scanner := bufio.NewScanner(file)
	cache := make([]string, 30)
	for scanner.Scan() {
		cache = append(cache, scanner.Text())
	}
	for _, val := range proxies {
		cache = append(cache, val)
	}
	tmp, err := os.OpenFile("/tmp/temp.txt", os.O_RDWR|os.O_CREATE, 0644)
	writer := bufio.NewWriter(tmp)
	for _, line := range cache {
		fmt.Println(line)
		_, err := writer.WriteString(line)
		if err != nil {
			fmt.Print(err)
		}
	}
	writer.Flush()

}
func proxy_off() {
	if addr == "" || port == "" {
		print_help()
		return
	}
	proxy := fmt.Sprintf("%v:%v", addr, port)
	fmt.Print(proxy)
}
func main() {
	if len(os.Args) <= 1 {
		print_help()
		return
	}
	switch {
	case os.Args[1] == "on":
		proxy_on()
	case os.Args[1] == "off":
		proxy_off()
	case os.Args[1] == "addr":
		set_addr()
	case os.Args[1] == "port":
		set_port()
	default:
		print_help()
	}

}
