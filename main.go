package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var addr string = ""
var port string = ""
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
		_, err := writer.WriteString(line)
		if err != nil {
			fmt.Print(err)
		}
	}
	writer.Flush()

	apt_write_cmd:=exec.Command("sudo","mv","/tmp/temp.txt","/etc/apt/apt.conf")
	apt_write_cmd.Run()

	proxies=proxies[:0]

	for _, val := range kind {
		proxies = append(proxies, fmt.Sprintf("%v_proxy=\"http://%v/\"\n", val, proxy))
		proxies = append(proxies, fmt.Sprintf("%v_PROXY=\"http://%v/\"\n", strings.ToUpper(val), proxy))
	}
	proxies=append(proxies,"no_proxy=localhost,127.0.0.0/8")
	proxies=append(proxies,"NO_PROXY=localhost,127.0.0.0/8")
	file_profile, err := os.Open("/etc/environment")
	defer file_profile.Close()
	if err != nil {
		fmt.Print(err)
	}
	scanner_profile := bufio.NewScanner(file_profile)
	cache =cache[:0]
	for scanner.Scan() {
		cache = append(cache, scanner_profile.Text())
	}
	for _, val := range proxies {
		cache = append(cache, val)
	}
	tmp_profile, err := os.OpenFile("/tmp/temp.txt", os.O_RDWR|os.O_CREATE, 0644)
	writer_profile := bufio.NewWriter(tmp_profile)
	for _, line := range cache {
		_, err := writer_profile.WriteString(line)
		if err != nil {
			fmt.Print(err)
			return
		}
	}
	writer_profile.Flush()

	profile_write_cmd:=exec.Command("sudo","mv","/tmp/temp.txt","/etc/environment")

	profile_write_cmd.Run()

//	GTK 3 applications
	gtk_cmd:=exec.Command("gsettings","set","org.gnome.system.proxy","mode","'manual'")
	gtk_cmd.Run()
	gtk_cmd=exec.Command("gsettings","set","org.gnome.system.proxy.http","host",addr)
	gtk_cmd.Run()
	gtk_cmd=exec.Command("gsettings","set","org.gnome.system.proxy.http","port",port)
	gtk_cmd.Run()
	gtk_cmd=exec.Command("gsettings","set","org.gnome.system.proxy.https","host",addr)
	gtk_cmd.Run()
	gtk_cmd=exec.Command("gsettings","set","org.gnome.system.proxy.https","port",port)
	gtk_cmd.Run()
	gtk_cmd=exec.Command("gsettings","set","org.gnome.system.proxy","ignore-hosts","localhost","127.0.0.0/8")
	gtk_cmd.Run()
}
func proxy_off() {
	apt_conf,_:=os.Open("/etc/apt/apt.conf")
	apt_conf_reader:=bufio.NewScanner(apt_conf)
	cache:=make([]string,10)
	for apt_conf_reader.Scan(){
		cache=append(cache,apt_conf_reader.Text())
	}
	outp:=make([]string,10)
	for i := range cache{
		if !strings.Contains(strings.ToLower(cache[i]),"proxy") {
			outp=append(outp,cache[i])
		}
	}
	apt_conf_temp,_:=os.OpenFile("/tmp/temp.txt",os.O_RDWR|os.O_CREATE,0644)
	apt_conf_writer:=bufio.NewWriter(apt_conf_temp)
	for _,val := range outp{
		_,err:=apt_conf_writer.WriteString(val)
		if err!=nil{
			fmt.Print(err)
		}
	}
	apt_conf_writer.Flush()

	apt_cmd:=exec.Command("sudo","mv","/tmp/temp.txt","/etc/apt/apt.conf")
	apt_cmd.Run()

	fmt.Println("Apt settings resetted")

	env_conf,_:=os.Open("/etc/environment")
	env_conf_reader:=bufio.NewScanner(env_conf)
	cache=cache[:0]
	for env_conf_reader.Scan(){
		cache=append(cache,env_conf_reader.Text())
	}
	outp=outp[:0]
	for i := range cache{
		if !strings.Contains(strings.ToLower(cache[i]),"proxy") {
			outp=append(outp,cache[i])
		}
	}
	env_conf_temp,_:=os.OpenFile("/tmp/temp.txt",os.O_RDWR|os.O_CREATE,0644)
	env_conf_writer:=bufio.NewWriter(env_conf_temp)
	for _,val := range outp {
		_, err := env_conf_writer.WriteString(val)
		if err != nil {
			fmt.Print(err)
		}
	}
	env_conf_writer.Flush()

	env_cmd:=exec.Command("sudo","mv","/tmp/temp.txt","/etc/environment")
	env_cmd.Run()
	fmt.Println("Environment settings resetted")

	gtk_cmd:=exec.Command("gsettings set org.gnome.system.proxy","mode","'none'")
	gtk_cmd.Run()
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
