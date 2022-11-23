package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(3),
		Run:  run,
		Use:  "cronclient customer-service-address repeats interval-seconds",
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	customerAddress := args[0]
	repeats, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("strconv.Atoi failed; %v", err.Error())
	}
	intervalSeconds, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("strconv.Atoi failed; %v", err.Error())
	}
	for i := 0; i < repeats; i++ {
		if err := accessCustomerService(customerAddress); err != nil {
			log.Printf("%v\n", err.Error())
		}
		time.Sleep(time.Second * time.Duration(intervalSeconds))
	}
}

func accessCustomerService(customerAddress string) error {
	customerID := rand.Intn(100)
	url := fmt.Sprintf("http://%s/%d", customerAddress, customerID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http.Get failed; %s; %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non 200 HTTP status; %d; %s", resp.StatusCode, resp.Status)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed; %w", err)
	}
	log.Printf("%s\n", string(respBytes))
	return nil
}
