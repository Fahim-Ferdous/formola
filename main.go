/*
Copyright © 2023 Fahim Ferdous

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

var discordAuthToken, testChannelID, sqlite3DSN string

//go:embed assets/*.ico
var assets embed.FS

func init() {
	flag.StringVar(&discordAuthToken, "t", "", "Bot Token")
	flag.StringVar(&testChannelID, "c", "", "Test Channel")
	flag.StringVar(&sqlite3DSN, "d", "", "SQLite DSN")
	flag.Parse()
	if discordAuthToken == "" {
		discordAuthToken = os.Getenv("DISCORD_AUTH_TOKEN")
		if discordAuthToken == "" {
			panic("missing DISCORD_AUTH_TOKEN")
		}
	}

	if testChannelID == "" {
		panic("testChannelID is required")
	}

	if sqlite3DSN == "" {
		panic("sqlite3DSN is required")
	}
}

var msgqueue = make(chan *MessageQueue)

func main() {
	waits := []chan<- struct{}{
		runServer(),
		runBot(discordAuthToken),
	}
	sc := make(chan os.Signal, 1)
	fmt.Println("Press CTRL-C to exit.")
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	for _, w := range waits {
		w <- struct{}{}
	}
}
