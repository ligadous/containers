package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func run() {
	fmt.Printf("running %v\n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	//Only in linux
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	check(cmd.Run())
}

func child() {
	fmt.Printf("running %v as pid %d \n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	check(syscall.Chroot("/root/rootfs"))
	fmt.Println("Chroot")

	check(os.Chdir("/"))
	fmt.Println("Chdir")

	check(syscall.Mount("proc", "proc", "proc", 0, ""))
	fmt.Println("Mount")

	check(cmd.Run())
	fmt.Println("Run")
}

func check(err error) {
	if err != nil {
		panic(err)
	}

}

//docker run <container> command args
//go run container.go run command args
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what?")
	}
}
