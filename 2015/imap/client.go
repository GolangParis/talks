package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xarg/imap"
)

func main() {
	// Setup logging
	imap.DefaultLogger = log.New(os.Stdout, "", 0)
	imap.DefaultLogMask = imap.LogConn | imap.LogRaw

	// Allocates imap.Client struct - opens TCP connection to the server
	c := Dial(Addr)
	defer func() { ReportOK(c.Logout(30 * time.Second)) }()

	// Check if server supports TLS and use it
	if c.Caps["STARTTLS"] {
		ReportOK(c.StartTLS(nil))
	}

	// My name is Bond..
	if c.Caps["ID"] {
		ReportOK(c.ID("name", "goimap"))
	}

	ReportOK(Login(c, "username", "password"))

	// Do your stuff ...
}

func Dial(addr string) (c *imap.Client) {
	var err error
	if strings.HasSuffix(addr, ":993") {
		c, err = imap.DialTLS(addr, nil)
	} else {
		c, err = imap.Dial(addr)
	}
	if err != nil {
		panic(err)
	}
	return c
}

func Login(c *imap.Client, user, pass string) (cmd *imap.Command, err error) {
	defer c.SetLogMask(Sensitive(c, "LOGIN"))
	return c.Login(user, pass)
}

func Sensitive(c *imap.Client, action string) imap.LogMask {
	mask := c.SetLogMask(imap.LogConn)
	hide := imap.LogCmd | imap.LogRaw
	if mask&hide != 0 {
		c.Logln(imap.LogConn, "Raw logging disabled during", action)
	}
	c.SetLogMask(mask &^ hide)
	return mask
}

func ReportOK(cmd *imap.Command, err error) *imap.Command {
	var rsp *imap.Response
	if cmd == nil {
		fmt.Printf("--- ??? ---\n%v\n\n", err)
		panic(err)
	} else if err == nil {
		rsp, err = cmd.Result(imap.OK)
	}
	if err != nil {
		fmt.Printf("--- %s ---\n%v\n\n", cmd.Name(true), err)
		panic(err)
	}
	c := cmd.Client()
	fmt.Printf("--- %s ---\n"+
		"%d command response(s), %d unilateral response(s)\n"+
		"%s %s\n\n",
		cmd.Name(true), len(cmd.Data), len(c.Data), rsp.Status, rsp.Info)
	c.Data = nil
	return cmd
}
