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
	apiVersion := os.Getenv("COMPUTER_VISION_API_VERSION")

	if subscriptionKey == "" || endpoint == "" {
		log.Println("One or more environment variables are not set.")
		return
	}

	indexName := "my-video-indexer"            // Acts like a namespace or container for videos
	ingestionName := "my-ingestion-or-video-1" // Represents the name of the ingestion operation
	videoUrl := "https://drive.google.com/file/d/173nr9-96GjeqrLbZce3v7ZWqPmnyVQuv/view?usp=sharing"

	// Step 1: Check if the index already exists
	indexURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s?api-version=%s", endpoint, indexName, apiVersion)

	// // Delete if exist
	// exists, err := computervision.CheckIndexExists(indexURL, subscriptionKey)
	// if err != nil {
	// 	if err.Error() == "status code 404: NotFound - {\"error\":{\"code\":\"NotFound\",\"message\":\"Index not found.\"}}" {
	// 		exists = false
	// 	} else {
	// 		log.Fatalf("failed to check if index exists: %v", err)
	// 	}
	// }
	// if exists {
	// 	// Delete the index
	// 	err := computervision.DeleteIndex(indexURL, subscriptionKey)
	// 	if err != nil {
	// 		log.Fatalf("failed to delete index: %v", err)
	// 	}
	// 	fmt.Println("Index deleted.")
	// }

	anotherexists, err := computervision.CheckIndexExists(indexURL, subscriptionKey)
	if err != nil {
		if err.Error() == "status code 404: NotFound - {\"error\":{\"code\":\"NotFound\",\"message\":\"Index not found.\"}}" {
			anotherexists = false
		} else {
			log.Fatalf("failed to check if index exists: %v", err)
		}
	}

	if !anotherexists {
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
					"domain": "surveillance",
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
				// Get the index to verify its status
				index, err := computervision.GetIndex(indexURL, subscriptionKey)
				if err != nil {
					log.Fatalf("failed to get index: %v", err)
				}
				fmt.Printf("Index status: %+v\n", index)
				break
			}
			fmt.Println("Waiting for index creation to complete...")
			time.Sleep(5 * time.Second)
		}
	} else {
		fmt.Println("Index already exists. Skipping creation.")
		// Get the index to verify its status
		index, err := computervision.GetIndex(indexURL, subscriptionKey)
		if err != nil {
			log.Fatalf("failed to get index: %v", err)
		}
		fmt.Printf("Index status: %+v\n", index)
	}

	// Step 4: Add video files to the index
	ingestionURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s/ingestions/%s?api-version=%s", endpoint, indexName, ingestionName, apiVersion)
	ingestionParams := computervision.IngestionRequest{
		Videos: []computervision.IngestionDocument{
			{
				Mode:        "add",
				DocumentId:  "thailand-night-video",
				DocumentUrl: videoUrl,
				Metadata: map[string]string{
					"cameraId":  "camera1",
					"timestamp": "2025-02-02T23:37:35Z",
				},
			},
		},
	}

	// Function to create ingestion
	createIngestion := func() (*computervision.Ingestion, error) {
		return computervision.CreateIngestion(ingestionURL, subscriptionKey, ingestionParams)
	}

	// Function to update ingestion
	updateIngestion := func() (*computervision.Ingestion, error) {
		return computervision.UpdateIngestion(ingestionURL, subscriptionKey, ingestionParams)
	}

	// Function to get ingestion status
	getIngestionStatus := func() (*computervision.Ingestion, error) {
		return computervision.GetIngestionStatus(ingestionURL, subscriptionKey)
	}

	// Step 5: Create ingestion and check status
	// Check if ingestion exists
	ingestionExists, err := computervision.CheckIngestionExists(ingestionURL, subscriptionKey)
	if err != nil {
		log.Fatalf("failed to check if ingestion exists: %v", err)
	}
	if !ingestionExists {
		ingestion, err := createIngestion()
		if err != nil {
			log.Fatalf("failed to create ingestion: %v", err)
		}
		fmt.Printf("Ingestion created: %v\n", ingestion)

		// Wait for ingestion to complete
		maxWaitTime := 5 * time.Minute
		startTime := time.Now()

		for {
			ingestionStatus, err := getIngestionStatus()
			if err != nil {
				log.Fatalf("failed to get ingestion status: %v", err)
			}
			if ingestionStatus.State == "Completed" {
				fmt.Println("Ingestion completed.")
				break
			}
			if ingestionStatus.State == "Failed" || ingestionStatus.State == "PartiallySucceeded" {
				fmt.Printf("Ingestion failed with state: %s, updating ingestion...\n", ingestionStatus.State)
				ingestion, err = updateIngestion()
				if err != nil {
					log.Fatalf("failed to update ingestion: %v", err)
				}
				fmt.Printf("Ingestion updated: %v\n", ingestion)
			}
			if time.Since(startTime) > maxWaitTime {
				log.Fatalf("Ingestion did not complete within 5 minutes.")
			}
			fmt.Println("Ingestion status:", ingestionStatus.State)
			time.Sleep(5 * time.Second)
		}
	} else {
		fmt.Println("Ingestion already exists. Skipping creation.")
		// Get the ingestion to verify its status
		ingestion, err := getIngestionStatus()
		if err != nil {
			log.Fatalf("failed to get ingestion status: %v", err)
		}
		fmt.Printf("Ingestion status: %+v\n", ingestion)
	}

	// List Ingestion
	listIngestionURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s/ingestions?api-version=%s", endpoint, indexName, apiVersion)
	ingestions, err := computervision.ListIngestions(listIngestionURL, subscriptionKey)
	if err != nil {
		log.Fatalf("failed to list ingestions: %v", err)
	}
	fmt.Printf("Ingestions: %+v\n", ingestions)

	// List Document
	listDocumentURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s/documents?api-version=%s", endpoint, indexName, apiVersion)
	documents, err := computervision.ListDocuments(listDocumentURL, subscriptionKey)
	if err != nil {
		log.Fatalf("failed to list documents: %v", err)
	}
	fmt.Printf("Documents: %+v\n", documents)

	// Step 6: Perform searches with metadata
	searchURL := fmt.Sprintf("%s/computervision/retrieval/indexes/%s:queryByText?api-version=%s", endpoint, indexName, apiVersion)
	searchParams := computervision.SearchQuery{
		QueryText: "woman",
		Filters: computervision.SearchFilters{
			StringFilters: []computervision.StringFilter{
				{
					FieldName: "cameraId",
					Values:    []string{"camera1"},
				},
			},
			FeatureFilters: []string{"vision", "speech"},
		},
	}
	searchResults, err := computervision.SearchByText(searchURL, subscriptionKey, searchParams)
	if err != nil {
		log.Fatalf("failed to search by text: %v", err)
	}

	if len(searchResults.Value) == 0 {
		fmt.Println("No search results found.")
	} else {
		fmt.Printf("Search results: %+v\n", searchResults)
	}
}
