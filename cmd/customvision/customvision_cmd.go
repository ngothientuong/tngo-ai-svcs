package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/ngothientuong/tngo-ai-svcs/internal/ai"
	"github.com/ngothientuong/tngo-ai-svcs/internal/config"
)

func main() {
	config.LoadEnv()

	training_key := os.Getenv("VISION_TRAINING_KEY")
	training_endpoint := os.Getenv("VISION_TRAINING_ENDPOINT")
	prediction_key := os.Getenv("VISION_PREDICTION_KEY")
	prediction_endpoint := os.Getenv("VISION_PREDICTION_ENDPOINT")
	prediction_resource_id := os.Getenv("VISION_PREDICTION_RESOURCE_ID")

	project_name := "Tuong Go Development Project"
	iteration_publish_name := "classifyModel"
	sampleDataDirectory := "./../../assets/images"

	if training_key == "" || training_endpoint == "" || prediction_key == "" || prediction_endpoint == "" || prediction_resource_id == "" {
		log.Println("One or more environment variables are not set.")
	}

	fmt.Println("Creating project...")

	projectURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects", training_endpoint)
	projectParams := map[string]string{"name": project_name}
	project, err := ai.CreateProject(projectURL, training_key, projectParams)
	if err != nil {
		log.Fatalf("failed to create project: %v", err)
	}

	fmt.Printf("Project created: %v\n", project)

	tagURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/tags", training_endpoint, project.ID)
	hemlockTagParams := map[string]string{"name": "Hemlock"}
	hemlockTag, err := ai.CreateTag(tagURL, training_key, hemlockTagParams)
	if err != nil {
		log.Fatalf("failed to create tag: %v", err)
	}

	cherryTagParams := map[string]string{"name": "Japanese Cherry"}
	cherryTag, err := ai.CreateTag(tagURL, training_key, cherryTagParams)
	if err != nil {
		log.Fatalf("failed to create tag: %v", err)
	}

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
	uploadImagesURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/images", training_endpoint, project.ID)
	err = ai.UploadImages(uploadImagesURL, training_key, hemlockImageFiles, hemlockTag.ID)
	if err != nil {
		log.Fatalf("failed to upload images: %v", err)
	}

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
	uploadImagesURL = fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/images", training_endpoint, project.ID)
	err = ai.UploadImages(uploadImagesURL, training_key, japaneseCherryImageFiles, cherryTag.ID)
	if err != nil {
		log.Fatalf("failed to upload images: %v", err)
	}

	fmt.Println("Training...")
	trainURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/train", training_endpoint, project.ID)
	trainParams := map[string]string{
		"notificationEmailAddress": "tuongdevops1@gmail.com",
	}
	iteration, err := ai.TrainProject(trainURL, training_key, trainParams)
	if err != nil {
		log.Fatalf("failed to train project: %v", err)
	}

	for {
		if iteration.Status != "Training" {
			break
		}
		fmt.Println("Training status: " + iteration.Status)
		time.Sleep(1 * time.Second)
		iterationURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/iterations/%s", training_endpoint, project.ID, iteration.ID)
		iteration, err = ai.GetIteration(iterationURL, training_key, nil)
		if err != nil {
			log.Fatalf("failed to get iteration: %v", err)
		}
	}
	fmt.Println("Training status: " + iteration.Status)

	publishURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/iterations/%s/publish", training_endpoint, project.ID, iteration.ID)
	publishParams := map[string]string{"publishName": iteration_publish_name, "predictionId": prediction_resource_id}
	err = ai.PublishIteration(publishURL, training_key, publishParams)
	if err != nil {
		log.Fatalf("failed to publish iteration: %v", err)
	}

	fmt.Println("Predicting...")
	testImageURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/quicktest/image", training_endpoint, project.ID)
	results, err := ai.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/test_image.jpg", nil)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}

	fmt.Println("Predicting with URL...")
	testImageUrl := "https://mortonarb.org/app/uploads/2020/12/Japanese-Flowering-Cherry_5144776054_ac5340eb34_o-1920x1440-c-default.jpg"
	testImageUrlEndpoint := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/quicktest/url", training_endpoint, project.ID)
	results, err = ai.QuickTestImageUrl(testImageUrlEndpoint, training_key, testImageUrl, nil)
	if err != nil {
		log.Fatalf("failed to classify image from URL: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}
}
