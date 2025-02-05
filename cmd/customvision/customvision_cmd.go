package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/customvision"
	"github.com/ngothientuong/tngo-ai-svcs/internal/config"
	"github.com/ngothientuong/tngo-ai-svcs/pkg/errorcustom"
)

func main() {
	config.LoadEnv()

	training_key := os.Getenv("VISION_TRAINING_KEY")
	training_endpoint := os.Getenv("VISION_TRAINING_ENDPOINT")
	prediction_key := os.Getenv("VISION_PREDICTION_KEY")
	prediction_endpoint := os.Getenv("VISION_PREDICTION_ENDPOINT")
	prediction_resource_id := os.Getenv("VISION_PREDICTION_RESOURCE_ID")
	apiVersion := os.Getenv("VISION_TRAINING_API_VERSION")

	project_name := "Tuong Go Development Project"
	iteration_publish_name := "classifyModel"
	sampleDataDirectory := "./../../assets/images"

	if training_key == "" || training_endpoint == "" || prediction_key == "" || prediction_endpoint == "" || prediction_resource_id == "" {
		log.Println("One or more environment variables are not set.")
	}
	fmt.Printf("Tempoutput: %s\n%s\n%s", project_name, iteration_publish_name, sampleDataDirectory)

	fmt.Println("Checking if project exists...")

	projectURL := fmt.Sprintf("%s/customvision/%s/training/projects", training_endpoint, apiVersion)
	project, err := customvision.GetProjectByName(projectURL, training_key, project_name)
	if err != nil {
		fmt.Println("Project not found, creating new project...")

		projectParams := map[string]string{"name": project_name}
		project, err = customvision.CreateProject(projectURL, training_key, projectParams)
		if err != nil {
			log.Fatalf("failed to create project: %v", err)
		}

		fmt.Printf("Project created: %v\n", project)
	} else {
		fmt.Printf("Project found: %v\n", project)
	}

	// Check if the latest iteration exists
	iterationsURL := fmt.Sprintf("%s/customvision/%s/training/projects/%s/iterations", training_endpoint, apiVersion, project.ID)
	latestIteration, err := customvision.GetIterationLatest(iterationsURL, training_key)
	if err != nil {
		log.Fatalf("failed to get latest iteration: %v", err)
	}

	var iterationID *string
	if latestIteration != nil {
		fmt.Println("Latest iteration found:", latestIteration.ID)
		iterationIDStr := latestIteration.ID.String()
		iterationID = &iterationIDStr
	} else {
		fmt.Println("No latest iteration found.")
		iterationID = nil
	}

	// Check if the Hemlock tag exists
	tagURL := fmt.Sprintf("%s/customvision/%s/training/projects/%s/tags", training_endpoint, apiVersion, project.ID)
	hemlockTag, err := customvision.GetTagByName(tagURL, training_key, project.ID.String(), "Hemlock", iterationID)
	if err != nil {
		fmt.Printf("Hemlock tag not found, creating new tag... Error: %v\n", err)

		hemlockTagParams := map[string]string{"name": "Hemlock"}
		hemlockTag, err = customvision.CreateTag(tagURL, training_key, hemlockTagParams)
		if err != nil {
			log.Fatalf("failed to create tag Hemlock: %v", err)
		}
	}
	fmt.Printf("Hemlock tag: %v\n", hemlockTag)

	// Check if the Japanese Cherry tag exists
	cherryTag, err := customvision.GetTagByName(tagURL, training_key, project.ID.String(), "Japanese Cherry", iterationID)
	if err != nil {
		fmt.Printf("Japanese Cherry tag not found, creating new tag...Error: %v\n", err)

		cherryTagParams := map[string]string{"name": "Japanese Cherry"}
		cherryTag, err = customvision.CreateTag(tagURL, training_key, cherryTagParams)
		if err != nil {
			log.Fatalf("failed to create tag Japanese Cherry: %v", err)
		}
	}
	fmt.Printf("Japanese Cherry tag: %v\n", cherryTag)

	// Check if the Dandelion tag exists
	dandelionTag, err := customvision.GetTagByName(tagURL, training_key, project.ID.String(), "Dandelion", iterationID)
	if err != nil {
		fmt.Printf("Dandelion tag not found, creating new tag...Error: %v\n", err)

		dandelionTagParams := map[string]string{"name": "Dandelion"}
		dandelionTag, err = customvision.CreateTag(tagURL, training_key, dandelionTagParams)
		if err != nil {
			log.Fatalf("failed to create tag Dandelion: %v", err)
		}
	}
	fmt.Printf("Dandelion tag: %v\n", dandelionTag)

	fmt.Println("Adding images...")

	hemlockImages, err := os.ReadDir(path.Join(sampleDataDirectory, "Hemlock"))
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	var hemlockImageFiles [][]byte
	for _, file := range hemlockImages {
		imageFile, err := os.ReadFile(path.Join(sampleDataDirectory, "Hemlock", file.Name()))
		if err != nil {
			log.Fatalf("failed to read image file: %v", err)
		}
		hemlockImageFiles = append(hemlockImageFiles, imageFile)
	}
	uploadImagesURL := fmt.Sprintf("%s/customvision/%s/training/projects/%s/images", training_endpoint, apiVersion, project.ID)
	summary, err := customvision.CreateImagesFromData(uploadImagesURL, training_key, hemlockImageFiles, hemlockTag.ID)
	if err != nil {
		log.Fatalf("failed to upload images: %v", err)
	}
	fmt.Printf("Hemlock images upload summary: %+v\n", summary)

	japaneseCherryImages, err := os.ReadDir(path.Join(sampleDataDirectory, "Japanese Cherry"))
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	var japaneseCherryImageFiles [][]byte
	for _, file := range japaneseCherryImages {
		imageFile, err := os.ReadFile(path.Join(sampleDataDirectory, "Japanese Cherry", file.Name()))
		if err != nil {
			log.Fatalf("failed to read image file: %v", err)
		}
		japaneseCherryImageFiles = append(japaneseCherryImageFiles, imageFile)
	}
	uploadImagesURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/images", training_endpoint, apiVersion, project.ID)
	summary, err = customvision.CreateImagesFromData(uploadImagesURL, training_key, japaneseCherryImageFiles, cherryTag.ID)
	if err != nil {
		log.Fatalf("failed to upload images: %v", err)
	}
	fmt.Printf("Japanese Cherry images upload summary: %+v\n", summary)

	dandelionImages, err := os.ReadDir(path.Join(sampleDataDirectory, "Dandelion"))
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	var dandelionImageFiles [][]byte
	for _, file := range dandelionImages {
		imageFile, err := os.ReadFile(path.Join(sampleDataDirectory, "Dandelion", file.Name()))
		if err != nil {
			log.Fatalf("failed to read image file: %v", err)
		}
		dandelionImageFiles = append(dandelionImageFiles, imageFile)
	}
	uploadImagesURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/images", training_endpoint, apiVersion, project.ID)
	summary, err = customvision.CreateImagesFromData(uploadImagesURL, training_key, dandelionImageFiles, dandelionTag.ID)
	if err != nil {
		log.Fatalf("failed to upload images: %v", err)
	}
	fmt.Printf("Dandelion images upload summary: %+v\n", summary)

	fmt.Println("Training...")
	trainURL := fmt.Sprintf("%s/customvision/%s/training/projects/%s/train", training_endpoint, apiVersion, project.ID)
	trainParams := map[string]string{
		"notificationEmailAddress": "tuongdevops1@gmail.com",
	}
	iteration, err := customvision.TrainProject(trainURL, training_key, trainParams)
	if err != nil {
		if customVisionErr, ok := err.(*errorcustom.CustomVisionError); ok && strings.Contains(customVisionErr.Message, "BadRequestTrainingNotNeeded") {
			fmt.Println("Training not needed.")
		} else {
			log.Fatalf("failed to train project: %v", err)
		}
	}
	if iteration != nil {
		var iterationURL, performanceURL, iterPerformanceURL string
		for {
			if iteration.Status != "Training" {
				break
			}
			fmt.Println("Training status: " + iteration.Status)
			time.Sleep(1 * time.Second)
			iterationURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/iterations/%s", training_endpoint, apiVersion, project.ID, iteration.ID)
			iteration, err = customvision.GetIteration(iterationURL, training_key, nil)
			if err != nil {
				log.Fatalf("failed to get iteration: %v", err)
			}
		}
		fmt.Println("Training status: " + iteration.Status)
		// Retrieve the performance of the current iteration
		performanceURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/iterations/%s/performance", training_endpoint, apiVersion, project.ID, iteration.ID)
		performance, err := customvision.GetIterationPerformance(performanceURL, training_key, nil)
		if err != nil {
			log.Fatalf("failed to get iteration performance: %v", err)
		}

		// Retrieve all iterations for the project
		iterationURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/iterations", training_endpoint, apiVersion, project.ID)
		iterations, err := customvision.GetIterations(iterationURL, training_key)
		if err != nil {
			log.Fatalf("failed to retrieve iterations: %v", err)
		}

		fmt.Println("iterations: ", iterations)
		// Find the iteration with the highest precision
		highestPrecision := 0.0
		for _, iter := range iterations {
			iterPerformanceURL = fmt.Sprintf("%s/customvision/%s/training/projects/a5561243-2401-48f6-a695-eb014f17c1fb/iterations/%s/performance", training_endpoint, apiVersion, iter.ID)
			iterPerformanceParams := map[string]string{"threshold": "0.5"}
			iterPerformance, err := customvision.GetIterationPerformance(iterPerformanceURL, training_key, iterPerformanceParams)
			if err != nil {
				log.Fatalf("failed to get iteration performance: %v", err)
			}
			if iterPerformance.Precision > highestPrecision {
				highestPrecision = iterPerformance.Precision
			}
		}

		// Publish the current iteration if it has the highest precision
		if performance.Precision >= highestPrecision {
			fmt.Println("Publishing iteration...")
			publishURL := fmt.Sprintf("%s/customvision/%s/training/projects/%s/iterations/%s/publish", training_endpoint, apiVersion, project.ID, iteration.ID)
			publishParams := map[string]string{"publishName": iteration_publish_name + "-" + iteration.ID.String(), "predictionId": prediction_resource_id}
			err = customvision.PublishIteration(publishURL, training_key, publishParams)
			if err != nil {
				log.Fatalf("failed to publish iteration: %v", err)
			}
			fmt.Println("Iteration published.")
		} else {
			fmt.Println("Current iteration does not have the highest precision. Skipping publishing.")
		}
	}
	var prediction_iteration_id string
	if iteration != nil {
		prediction_iteration_id = iteration.ID.String()
	} else {
		prediction_iteration_id = latestIteration.ID.String()
	}

	fmt.Println("Predicting...")
	fmt.Println("Predicting with image HelmLock...")
	testImageURL := fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/image", training_endpoint, apiVersion, project.ID)
	testImageURLParams := map[string]string{"iterationId": prediction_iteration_id}
	results, err := customvision.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/hemlocktest1.jpg", testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	testImageURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/image", training_endpoint, apiVersion, project.ID)
	testImageURLParams = map[string]string{"iterationId": prediction_iteration_id}
	results, err = customvision.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/hemlocktest2.jpg", testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	fmt.Println("Predicting with image Dandelion...")
	testImageURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/image", training_endpoint, apiVersion, project.ID)
	testImageURLParams = map[string]string{"iterationId": prediction_iteration_id}
	results, err = customvision.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/dandeliontest1.jpg", testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	fmt.Println("Predicting with image Japanese Cherry...")
	testImageURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/image", training_endpoint, apiVersion, project.ID)
	testImageURLParams = map[string]string{"iterationId": prediction_iteration_id}
	results, err = customvision.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/japanesecherry2.jpg", testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	testImageURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/image", training_endpoint, apiVersion, project.ID)
	testImageURLParams = map[string]string{"iterationId": prediction_iteration_id}
	results, err = customvision.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/japanesecherry3.jpg", testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	testImageURL = fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/image", training_endpoint, apiVersion, project.ID)
	testImageURLParams = map[string]string{"iterationId": prediction_iteration_id}
	results, err = customvision.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/japanesecherry1.jpg", testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	fmt.Println("Predicting with URL...")
	fmt.Println("Predicting with image Cherry...")
	testImageUrl := "https://www.datocms-assets.com/101439/1700918559-cherry-blossom-at-hirosaki-park.jpg?auto=format&h=1000&w=2000"
	testImageUrlEndpoint := fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/url", training_endpoint, apiVersion, project.ID)
	results, err = customvision.QuickTestImageUrl(testImageUrlEndpoint, training_key, testImageUrl, testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image from URL: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	fmt.Println("Predicting with image Dandelion...")
	testImageUrl = "https://myfavouritepastime.com/wp-content/uploads/2018/10/img_2927.jpg"
	testImageUrlEndpoint = fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/url", training_endpoint, apiVersion, project.ID)
	results, err = customvision.QuickTestImageUrl(testImageUrlEndpoint, training_key, testImageUrl, testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image from URL: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	fmt.Println("Predicting with image Hemlock...")
	testImageUrl = "https://d18lev1ok5leia.cloudfront.net/chesapeakebay/critters/_700x600_fit_center-center_none/Eastern-Hemlock-20180420-IMG_5309.jpg"
	testImageUrlEndpoint = fmt.Sprintf("%s/customvision/%s/training/projects/%s/quicktest/url", training_endpoint, apiVersion, project.ID)
	results, err = customvision.QuickTestImageUrl(testImageUrlEndpoint, training_key, testImageUrl, testImageURLParams)
	if err != nil {
		log.Fatalf("failed to classify image from URL: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

}
