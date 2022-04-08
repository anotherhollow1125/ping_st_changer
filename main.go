package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	var device_ip string
	var connect_url string
	var disconnect_url string
	var timeout time.Duration
	var ipv6 bool

	flag.StringVar(&device_ip, "dev", "", "device ip")
	flag.StringVar(&connect_url, "con", "", "API URL access when device connected to LAN")
	flag.StringVar(&disconnect_url, "dis", "", "API URL access when device disconnected from LAN")
	flag.DurationVar(&timeout, "timeout", time.Microsecond, "timeout")
	flag.Parse()
	if device_ip == "" {
		log.Fatal("device ip is required")
	}

	// log.Printf("device_ip: %s\n", device_ip)

	cache, err := readCache()
	if err != nil {
		log.Fatal(err)
	}

	is_connected := cache == "con"

	proto := "ip4"
	if ipv6 {
		proto = "ip6"
	}

	ip, err := net.ResolveIPAddr(proto, device_ip)
	if err != nil {
		log.Fatalf("ResolveIPAddr: %v", err)
	}

	c, err := icmp.ListenPacket(proto+":icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer c.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(""),
		},
	}
	wb, err := wm.Marshal(nil)

	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	if _, err := c.WriteTo(wb, &net.IPAddr{IP: ip.IP}); err != nil {
		log.Fatalf("WriteTo: %v", err)
	}

	c.SetReadDeadline(time.Now().Add(timeout))
	_, _, err = c.ReadFrom(wb)
	if err != nil {
		if is_connected {
			_, err = http.Get(disconnect_url)
			if err == nil {
				writeCache("dis")
				log.Printf("Change State: Disconnect")
			}
		}
	} else {
		if !is_connected {
			_, err = http.Get(connect_url)
			if err == nil {
				writeCache("con")
				log.Printf("Change State: Connect")
			}
		}
	}
}

func readCache() (string, error) {
	if _, err := os.Stat("/var/cache/ping_st_changer"); os.IsNotExist(err) {
		// create file & initialize with "dis" string
		file, err := os.Create("/var/cache/ping_st_changer")
		if err != nil {
			return "", err
		}
		defer file.Close()
		file.WriteString("dis")
		return "dis", nil
	}

	content, err := os.ReadFile("/var/cache/ping_st_changer")
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func writeCache(content string) error {
	file, err := os.Create("/var/cache/ping_st_changer")
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(content)
	return nil
}
