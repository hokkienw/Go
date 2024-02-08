package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/json"
)

func main() {
	candyType := flag.String("k", "", "Candy type (two-letter abbreviation)")
	candyCount := flag.Int("c", 0, "Count of candy to buy")
	money := flag.Int("m", 0, "Amount of money given to the machine")
	flag.Parse()

	clientCert, err := tls.LoadX509KeyPair("client_cert.pem", "client_key.pem")
	if err != nil {
		fmt.Println("Error loading client certificate and key:", err)
		return
	}

	caCert, err := ioutil.ReadFile("minica.pem")
	if err != nil {
		fmt.Println("Error loading CA certificate:", err)
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	url := "http://127.0.0.1:3333/buy_candy"
	requestBody := fmt.Sprintf(`{"money": %d, "candyType": "%s", "candyCount": %d}`, *money, *candyType, *candyCount)
	resp, err := client.Post(url, "application/json", strings.NewReader(requestBody))
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
if err != nil {
    fmt.Println("Error reading response:", err)
    return
}

var response struct {
    Change  int    `json:"change"`
    Thanks  string `json:"thanks"`
}

if err := json.Unmarshal(responseBody, &response); err != nil {
    fmt.Println("Error parsing response:", err)
    return
}

if resp.StatusCode == http.StatusCreated {
    fmt.Printf("Thank you! Your change is %d\n", response.Change)
} else {
    fmt.Printf("Server returned an error: %s\n", resp.Status)
}
}
