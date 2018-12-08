//package rawsocket
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func main() {
	// Socket is defined as:
	// func Socket(domain, typ, proto int) (fd int, err error)
	// Domain specifies the protocol family to be used - this should be AF_PACKET
	// to indicate we want the low level packet interface
	// Type specifies the semantics of the socket
	// Protocol specifies the protocol to use - kept here as ETH_P_ALL to
	// indicate all protocols over Ethernet
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW,
		syscall.ETH_P_ALL)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	fmt.Println("Obtained fd ", fd)
	defer syscall.Close(fd)
	packet := C.GoBytes(unsafe.Pointer(C.FillRequestPacketFields(iface_cstr, ip_cstr)),
		C.int(size))

	var addr syscall.SockaddrLinklayer
	addr.Protocol = syscall.ETH_P_ARP
	addr.Ifindex = interf.Index
	addr.Hatype = syscall.ARPHRD_ETHER
	// Send the packet
	err = syscall.Sendto(fd, packet, 0, &addr)
	// Do something with fd here

}
