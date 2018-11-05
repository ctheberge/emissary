package dns

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func Flush() {
	pids, err := getPIDs()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Flushing DNS: %v\n", pids)
	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil { log.Println(err) }
		err = proc.Signal(syscall.SIGHUP)
		if err != nil { log.Println(err) }
	}
}

func getPIDs() (pids []int, err error) {
	cmd := exec.Command("ps", "-axo", "pid=,command=")
	out, err := cmd.CombinedOutput()
	if err != nil { return }
	if !cmd.ProcessState.Success() {
		err = fmt.Errorf("%s", out)
		return
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(strings.ToLower(line), "mdnsresponder") {
			parts := strings.Fields(line)
			var pid int
			pid, err = strconv.Atoi(parts[0])
			if err != nil { return }
			pids = append(pids, pid)
		}
	}

	return
}
