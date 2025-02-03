package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/computervision"
	"github.com/ngothientuong/tngo-ai-svcs/internal/config"
)

func main() {
	config.LoadEnv()

	subscriptionKey := os.Getenv("COMPUTER_VISION_KEY")
	endpoint := os.Getenv("COMPUTER_VISION_ENDPOINT")

	if subscriptionKey == "" || endpoint == "" {
		log.Println("One or more environment variables are not set.")
		return
	}

	indexName := "my-video-indexer"          // Acts like a namespace or container for videos
	ingestionName := "my-ingestion-or-video" // Represents the name of the ingestion operation
	videoUrl := "https://www.youtube.com/watch?v=wuy3kQPsEWs"

	// Step 1: Check if the index already exists
	indexURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s?api-version=2023-05-01-preview", endpoint, indexName)
	// // Delete the index if it already exists
	// err := computervision.DeleteIndex(indexURL, subscriptionKey)
	// if err != nil {
	// 	log.Fatalf("failed to delete index: %v", err)
	// }
	exists, err := computervision.CheckIndexExists(indexURL, subscriptionKey)
	if err != nil {
		if err.Error() == "status code 404: NotFound - {\"error\":{\"code\":\"NotFound\",\"message\":\"Index not found.\"}}" {
			exists = false
		} else {
			log.Fatalf("failed to check if index exists: %v", err)
		}
	}

	if !exists {
		// Step 2: Create an Index
		indexParams := map[string]interface{}{
			"metadataSchema": map[string]interface{}{
				"fields": []map[string]interface{}{
					{
						"name":       "cameraId",
						"searchable": false,
						"filterable": true,
						"type":       "string",
					},
					{
						"name":       "timestamp",
						"searchable": false,
						"filterable": true,
						"type":       "datetime",
					},
				},
			},
			"features": []map[string]interface{}{
				{
					"name":   "vision",
					"domain": "generic",
				},
				{
					"name":   "speech",
					"domain": "generic",
				},
			},
		}
		index, err := computervision.CreateIndex(indexURL, subscriptionKey, indexParams)
		if err != nil {
			log.Fatalf("failed to create index: %v", err)
		}
		fmt.Printf("Index created: %v\n", index)

		// Step 3: Wait for the index to be completely created
		for {
			exists, err := computervision.CheckIndexExists(indexURL, subscriptionKey)
			if err != nil {
				log.Fatalf("failed to check if index exists: %v", err)
			}
			if exists {
				fmt.Println("Index creation completed.")
				break
			}
			fmt.Println("Waiting for index creation to complete...")
			time.Sleep(5 * time.Second)
		}
	} else {
		fmt.Println("Index already exists. Skipping creation.")
	}

	// Step 4: Add video files to the index
	ingestionURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s/ingestions/%s?api-version=2023-05-01-preview", endpoint, indexName, ingestionName)
	ingestionParams := computervision.IngestionRequest{
		Videos: []computervision.IngestionDocument{
			{
				Mode:        "add",
				DocumentId:  "youtube-video",
				DocumentUrl: videoUrl,
				Metadata: map[string]string{
					"cameraId":  "camera1",
					"timestamp": "2025-02-02T23:37:35Z",
				},
			},
		},
	}
	// Check if the ingestion already exists
	ingestionStatusURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s/ingestions/%s?api-version=2023-05-01-preview", endpoint, indexName, ingestionName)
	ingestionExists, err := computervision.CheckIngestionExists(ingestionStatusURL, subscriptionKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		log.Fatalf("failed to check if ingestion exists: %v", err)
	}

	var ingestion *computervision.Ingestion
	if !ingestionExists {
		ingestion, err = computervision.CreateIngestion(ingestionURL, subscriptionKey, ingestionParams)
		if err != nil {
			log.Fatalf("failed to create ingestion: %v", err)
		}
		fmt.Printf("Ingestion created: %v\n", ingestion)
	} else {
		fmt.Println("Ingestion already exists. Skipping creation.")
		ingestion, err = computervision.GetIngestionStatus(ingestionStatusURL, subscriptionKey)
		if err != nil {
			log.Fatalf("failed to get ingestion status: %v", err)
		} else {
			fmt.Printf("Ingestion status: %+v\n", ingestion)
		}

	}

	// Step 5: Wait for ingestion to complete
	maxWaitTime := 5 * time.Minute
	startTime := time.Now()

	for {
		ingestionStatus, err := computervision.GetIngestionStatus(ingestionStatusURL, subscriptionKey)
		if err != nil {
			log.Fatalf("failed to get ingestion status: %v", err)
		}
		if ingestionStatus.State == "Completed" {
			fmt.Println("Ingestion completed.")
			break
		}
		if ingestionStatus.State != "Running" {
			log.Fatalf("Ingestion failed with state: %s", ingestionStatus.State)
		}
		if time.Since(startTime) > maxWaitTime {
			log.Fatalf("Ingestion did not complete within 5 minutes.")
		}
		fmt.Println("Ingestion status:", ingestionStatus.State)
		time.Sleep(5 * time.Second)
	}

	// Step 6: Perform searches with metadata
	searchURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s:queryByText?api-version=2023-05-01-preview", endpoint, indexName)
	searchParams := computervision.SearchQuery{
		QueryText: "screen",
		Filters: computervision.SearchFilters{
			StringFilters: []computervision.StringFilter{
				{
					FieldName: "cameraId",
					Values:    []string{"camera1"},
				},
			},
			FeatureFilters: []string{"speech", "vision"},
		},
	}
	searchResults, err := computervision.SearchByText(searchURL, subscriptionKey, searchParams)
	if err != nil {
		log.Fatalf("failed to search by text: %v", err)
	}
	fmt.Printf("Search results: %+v\n", searchResults)
}
