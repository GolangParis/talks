package main

import (
	"log"
	"os"

	"github.com/xarg/imap"
)

func main() {
	c := imap.NewClient(
		imap.Logger(
			log.New(os.Stdout, "", 0),
			imap.LogConn|imap.LogRaw,
		))

	// Figure out if it's SSL/TLS (:993) or plaintext/STARTTLS (:143)
	err := c.Dial(Addr)
	if err != nil {
		log.Fatal("Failed connection", err)
	}
}
