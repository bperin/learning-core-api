package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

func main() {
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		log.Fatalf("Failed to read swagger.json: %v", err)
	}

	var doc2 openapi2.T
	if err := json.Unmarshal(data, &doc2); err != nil {
		log.Fatalf("Failed to unmarshal swagger.json: %v", err)
	}

	doc3, err := openapi2conv.ToV3(&doc2)
	if err != nil {
		log.Fatalf("Failed to convert to V3: %v", err)
	}

	data3, err := json.MarshalIndent(doc3, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal openapi3.json: %v", err)
	}

	if err := os.WriteFile("docs/openapi3.json", data3, 0644); err != nil {
		log.Fatalf("Failed to write openapi3.json: %v", err)
	}

	fmt.Println("Converted docs/swagger.json to docs/openapi3.json")
}
