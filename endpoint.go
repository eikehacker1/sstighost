package main

import (
	"flag"
	"fmt"
	"sstighost/ssti"
)

func main() {
	
	endpoint := flag.String("e", "", "Endpoint to test for SSTI vulnerabilities")
	flag.Parse()

	
	if *endpoint == "" {
		fmt.Println("Error: Please provide an endpoint using the -e flag.")
		return
	}

	
	fmt.Printf("Testing endpoint: %s\n", *endpoint)
	for _, payload := range ssti.SSTIPayloads {
		result := ssti.SSTI(*endpoint, payload.Payload, payload.Expected, "", false)
		if result != "ERROR" {
			fmt.Println(result)
		}
	}
}