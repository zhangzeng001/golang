# 获取服务器ip地址

* 拿到net.interface的结构体对象列表

  `ifaces, err := net.Interfaces()`

  *返回：*

  `func Interfaces() ([]Interface, error)`

  *源码：*

  ```go
  type Interface struct {
      Index        int          // positive integer that starts at one, zero is never used
      MTU          int          // maximum transmission unit
      Name         string       // e.g., "en0", "lo0", "eth0.100"
      HardwareAddr HardwareAddr // IEEE MAC-48, EUI-48 and EUI-64 form
      Flags        Flags        // e.g., FlagUp, FlagLoopback, FlagMulticast
  }
  ```

* 获取em2网卡ip地址

  ```go
  package main
  
  import (
  	"errors"
  	"fmt"
  	"net"
  )
  
  func externalIP() (net.IP, error) {
  	// 拿到net.interface的结构体对象
  	ifaces, err := net.Interfaces()
  	if err != nil {
  		return nil, err
  	}
  	for _, iface := range ifaces {
  		if iface.Flags&net.FlagUp == 0 {
  			continue // interface down
  		}
  		if iface.Flags&net.FlagLoopback != 0 {
  			continue // loopback interface
  		}
  		// 返回地址列表
  		/*
  		[172.21.0.14/20 fe80::5054:ff:fe7a:c3c/64]
  		[192.168.0.14/24 fe80::1224:ff:fe7a:dsb/64]
  		*/
  		addrs, err := iface.Addrs()
  		if err != nil {
  			return nil, err
  		}
  		fmt.Println("输出addrs列表",addrs)
  		// 输出指定网卡IP地址
  		if iface.Name == "em2"{
  			for _, addr := range addrs {
  				ip := getIpFromAddr(addr)
  				if ip == nil {
  					continue
  				}
  				return ip, nil
  			}
  		}
  	}
  	return nil, errors.New("connected to the network?")
  }
  // 从获取的列表中取出指定网卡ip
  func getIpFromAddr(addr net.Addr) net.IP {
  	var ip net.IP
  	switch v := addr.(type) {
  	case *net.IPNet:
  		ip = v.IP
  	case *net.IPAddr:
  		ip = v.IP
  	}
  	if ip == nil || ip.IsLoopback() {
  		return nil
  	}
  	ip = ip.To4()
  	if ip == nil {
  		return nil // not an ipv4 address
  	}
  	return ip
  }
  
  func main() {
  	ip, err := externalIP()
  	if err != nil {
  		fmt.Println(err)
  	}
  
  	fmt.Println(ip.String())
  }
  ```

  

