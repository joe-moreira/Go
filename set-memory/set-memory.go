package main

// sets instance 6379 max memory limit to half of server's physical memory

// #include <unistd.h>
import "C"
import (
	"fmt"
	"os/exec"
	"strconv"
)

func main() {
	totalMem := C.sysconf(C._SC_PHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE)
	_, err := exec.Command("/bin/bash", "-c", "redis-cli config set maxmemory "+strconv.Itoa(int(totalMem/2))).Output()
	exec.Command("/bin/bash", "-c", "redis-cli config rewrite")

	if err != nil {
		fmt.Println(err)
	}
}
