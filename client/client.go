package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Run: run,
		Use: "client host num1 num2",
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	host := args[0]
	num1, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("strconv.Atoi failed; %v", err.Error())
	}
	num2, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("strconv.Atoi failed; %v", err.Error())
	}
	if err := sendRequest(host, num1, num2); err != nil {
		log.Fatalf("sendRequest failed; %v", err.Error())
	}
}

type requestFormat struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

type responseFormat struct {
	Sum     int `json:"sum"`
	Product int `json:"product"`
}

func sendRequest(host string, num1, num2 int) error {
	url := fmt.Sprintf("http://%s", host)
	reqBody, err := json.Marshal(requestFormat{Num1: num1, Num2: num2})
	if err != nil {
		return fmt.Errorf("json.Marshal failed; %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("http.NewRequest failed; %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do failed; %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got non 200 HTTP status; %d; %s", resp.StatusCode, resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed; %w", err)
	}
	var respFormat responseFormat
	if err := json.Unmarshal(respBody, &respFormat); err != nil {
		return fmt.Errorf("json.Unmarshal failed; %s; %w", string(respBody), err)
	}
	fmt.Printf("%v\n", respFormat)
	return nil
}
