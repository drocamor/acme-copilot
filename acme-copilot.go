package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

var currentWindow int

func listWindows() string {
	var buf bytes.Buffer
	windows, err := acme.Windows()
	if err != nil {
		return fmt.Sprintf("ls error: %s", err)
	}
	for _, w := range windows {
		entry := fmt.Sprintf("[%d] %s\n", w.ID, w.Name)
		buf.WriteString(entry)
	}

	return buf.String()
}

func changeWindow(id string) string {
	i, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Sprintf("cd error: %s", err)
	}

	currentWindow = i

	return fmt.Sprintf("cd to window %d", i)
}

func sendCtl(cmd string) string {
	w, err := acme.Open(currentWindow, nil)
	if err != nil {
		return fmt.Sprintf("error getting window: %s", err)
	}
	err = w.Ctl("%s", cmd)
	if err != nil {
		return fmt.Sprintf("error sending ctl: %s", err)
	}
	return fmt.Sprintf("sent: %s", cmd)
}

func setAddr(addr string) string {
	w, err := acme.Open(currentWindow, nil)
	if err != nil {
		return fmt.Sprintf("error getting window: %s", err)
	}
	err = w.Addr("%s", addr)
	if err != nil {
		return fmt.Sprintf("error setting addr: %s", err)
	}
	return fmt.Sprintf("set addr: %s", addr)
}

func runCommand(s string) string {
	trimmedCommand := strings.TrimSpace(s)
	entry := strings.SplitN(trimmedCommand, " ", 2)

	switch entry[0] {
	case "ls":
		return listWindows()
	case "cd":
		if len(entry) > 1 {
			return changeWindow(entry[1])
		} else {
			return "usage: cd <window number>"
		}
	case "ctl":
		if len(entry) > 1 {
			return sendCtl(entry[1])
		} else {
			return "usage: ctl <ctl command>"
		}
	case "addr":
		if len(entry) > 1 {
			return setAddr(entry[1])
		} else {
			return "usage: addr <address>"
		}
	case "log":
		go eventLogger()
		return "logging on"
	default:
		return "need command: ls, cd, ctl, addr, log"
	}
}

func eventLogger() {
	w, err := acme.Open(currentWindow, nil)
	if err != nil {
		log.Println("Could not open window:", err)
	}

	events := w.EventChan()

	for {
		e := <-events
		log.Printf("%+v", e)
		var buf bytes.Buffer
		for _, r := range e.Text {
			buf.WriteRune(rune(r))
		}
		log.Println("Text:", buf.String())
	}
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%d > ", currentWindow)

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Got error: ", err)
		}

		fmt.Println(runCommand(text))

	}
}
