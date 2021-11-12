package main

import (
	"io"
	"log"
	"net"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	cmd := exec.Command("python", "-c", "import pty; pty.spawn(\"/bin/bash\")")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	// cmd := exec.Command("./uml-kernel", "name=t1", "title=t1", "umid=t1", "mem=132M", "ubd0=t1.disk,uml-fs", "root=98:0", "uml_dir=/run/user/1000/netkit/uml/GLOBAL", "con0=fd:0,fd:1", "con1=null", "SELINUX_INIT=0")
	l, err := net.Listen("unix", "/tmp/test.sock")
	if err != nil {
		log.Fatal(err)
	}
	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go io.Copy(ptmx, fd)
		go io.Copy(fd, ptmx)
	}
}
