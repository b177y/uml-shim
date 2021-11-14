package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
	"github.com/docker/docker/pkg/reexec"
)

func handleConnection(c net.Conn, cOut, cIn chan []byte) {
	go func() {
		for {
			buf := make([]byte, 2048)
			c.Read(buf)
			cIn <- buf
		}
	}()
	go func() {
		for {
			buf := <-cOut
			br := bytes.NewReader(buf)
			io.Copy(c, br)
		}
	}()
}

type Broadcaster struct {
	clients []chan []byte
}

func (b *Broadcaster) SendAll(msg []byte) {
	for _, c := range b.clients {
		c <- msg
	}
}

func (b *Broadcaster) AddClient(newChan chan []byte) {
	b.clients = append(b.clients, newChan)
}

func runShim() {
	// cmd := exec.Command("python", "-c", "import pty; pty.spawn(\"/bin/bash\")")
	// cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd := exec.Command("/home/billy/netkit-jh/kernel/netkit-kernel", "name=testmachine3", "title=testmachine3", "umid=testmachine3", "mem=132M", "ubd0=/home/billy/.local/share/netkit/uml/overlay/GLOBAL/testmachine3.disk,/home/billy/netkit-jh/fs/netkit-fs", "root=98:0", "uml_dir=/run/user/1000/netkit/uml/GLOBAL", "ssl0=fd:3,fd:1", "con1=null", "SELINUX_INIT=0")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	// TODO set DIR/status to booting
	l, err := net.Listen("unix", "/tmp/test.sock")
	if err != nil {
		log.Fatal(err)
	}
	cIn := make(chan []byte)
	bc := new(Broadcaster)
	// TODO add logger chan, to send to MACHINE.log
	// if "Welcome To Netkit" seen, change machine status to running
	go func() {
		for {
			buf := make([]byte, 1)
			io.ReadAtLeast(ptmx, buf, 1)
			bc.SendAll(buf)
		}
	}()
	go func() {
		for {
			buf := <-cIn
			br := bytes.NewReader(buf)
			io.Copy(ptmx, br)
		}
	}()
	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("[INFO] New connection from %s.\n", fd.LocalAddr())
			newChan := make(chan []byte)
			bc.AddClient(newChan)
			go handleConnection(fd, newChan, cIn)
		}
	}()
	cmd.Wait()
	l.Close()
	// ec := cmd.ProcessState.ExitCode()
	// TODO write exit code to a file
}

func main() {
	// cmdline := []string{"/home/billy/netkit-jh/kernel/netkit-kernel", "name=testmachine3", "title=testmachine3", "umid=testmachine3", "mem=132M", "ubd0=/home/billy/.local/share/netkit/uml/overlay/GLOBAL/testmachine3.disk,/home/billy/netkit-jh/fs/netkit-fs", "root=98:0", "uml_dir=/run/user/1000/netkit/uml/GLOBAL", "ssl0=fd:3,fd:1", "con1=null", "SELINUX_INIT=0"}
	// cmdline := []string{"python", "-c", "import pty; pty.spawn(\"/bin/bash\")"}
	// runShim(cmdline)
	cmd := reexec.Command("childProcess")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to run command: %s", err)
	}
}

func init() {
	log.Printf("init start, os.Args = %+v\n", os.Args)
	reexec.Register("childProcess", runShim)
	if reexec.Init() {
		os.Exit(0)
	}
}
