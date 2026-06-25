package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"zerodhaapi-go/kiteconnect"
)

func main() {
	client := kiteconnect.NewClient("x9hd5jluppbcqmeq", "ir07d4c9gedzhwcgj81tvrjczcfk6f2d", kiteconnect.WithBaseURL(kiteconnect.DefaultBaseURL))

	fmt.Println("Open this URL and login:")
	fmt.Println(client.LoginURL())

	fmt.Print("Paste request_token or redirected URL: ")
	reader := bufio.NewReader(os.Stdin)
	requestToken, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	requestToken = strings.TrimSpace(requestToken)

	session, err := client.GenerateSession(context.Background(), requestToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Access token:")
	fmt.Println(session.AccessToken)
}
