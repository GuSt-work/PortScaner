package infrastruction

import (
	"syscall"
	"unsafe"
)

var (
	ws2_32 = syscall.NewLazyDLL("ws2_32.dll")

	procWSAStartup = ws2_32.NewProc("WSAStartup")
	procWSACleanup = ws2_32.NewProc("WSACleanup")
)

type WSAData struct {
	Version      uint16
	HighVersion  uint16
	Description  [257]byte
	SystemStatus [129]byte
	VendorInfo   uintptr
}

func InitWinSock() error {

	var data WSAData

	ret, _, err := procWSAStartup.Call(
		uintptr(makeWord(2, 2)),
		uintptr(unsafe.Pointer(&data)),
	)

	if ret != 0 {
		return err
	}

	return nil
}

func makeWord(low byte, high byte) uint16 {
	return uint16(low) | uint16(high)<<8
}

func CleanupWinSock() {

	procWSACleanup.Call()

}
