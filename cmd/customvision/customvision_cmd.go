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

	projectURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects?name=%s", training_endpoint, project_name)
	project, err := ai.CreateProject(projectURL, training_key, project_name)
	if err != nil {
		log.Fatalf("failed to create project: %v", err)
	}

	fmt.Printf("Project created: %v\n", project)

	tagURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/tags", training_endpoint, project.ID)
	hemlockTag, err := ai.CreateTag(tagURL, training_key, project.ID, "Hemlock")
	if err != nil {
		log.Fatalf("failed to create tag: %v", err)
	}

	cherryTag, err := ai.CreateTag(tagURL, training_key, project.ID, "Japanese Cherry")
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
		imageFile, _ := os.ReadFile(path.Join(sampleDataDirectory, "Hemlock", file.Name()))
		hemlockImageFiles = append(hemlockImageFiles, imageFile)
	}
	uploadImagesURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/images?tagIds=%s", training_endpoint, project.ID, hemlockTag.ID)
	err = ai.UploadImages(uploadImagesURL, training_key, project.ID, hemlockImageFiles, hemlockTag.ID)
	if err != nil {
		log.Fatalf("failed to upload images: %v", err)
	}

	japaneseCherryImages, err := os.ReadDir(path.Join(sampleDataDirectory, "Japanese Cherry"))
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	var japaneseCherryImageFiles [][]byte
	for _, file := range japaneseCherryImages {
		imageFile, _ := os.ReadFile(path.Join(sampleDataDirectory, "Japanese Cherry", file.Name()))
		japaneseCherryImageFiles = append(japaneseCherryImageFiles, imageFile)
	}
	uploadImagesURL = fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/images?tagIds=%s", training_endpoint, project.ID, cherryTag.ID)
	err = ai.UploadImages(uploadImagesURL, training_key, project.ID, japaneseCherryImageFiles, cherryTag.ID)
	if err != nil {
		log.Fatalf("failed to upload images: %v", err)
	}

	fmt.Println("Training...")
	trainURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/train", training_endpoint, project.ID)
	iteration, err := ai.TrainProject(trainURL, training_key, project.ID)
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
		iteration, err = ai.GetIteration(iterationURL, training_key, project.ID, iteration.ID)
		if err != nil {
			log.Fatalf("failed to get iteration: %v", err)
		}
	}
	fmt.Println("Training status: " + iteration.Status)

	publishURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/iterations/%s/publish?publishName=%s&predictionId=%s", training_endpoint, project.ID, iteration.ID, iteration_publish_name, prediction_resource_id)
	err = ai.PublishIteration(publishURL, training_key, iteration_publish_name, prediction_resource_id, project.ID, iteration.ID)
	if err != nil {
		log.Fatalf("failed to publish iteration: %v", err)
	}

	fmt.Println("Predicting...")
	testImageURL := fmt.Sprintf("%s/customvision/v3.4-preview/training/projects/%s/quicktest/image", training_endpoint, project.ID)
	results, err := ai.QuickTestImage(testImageURL, training_key, sampleDataDirectory, "Test/test_image.jpg", project.ID)
	if err != nil {
		log.Fatalf("failed to classify image: %v", err)
	}

	for _, prediction := range results.Predictions {
		fmt.Printf("%s: %.2f%%\n", prediction.TagName, prediction.Probability*100)
	}
}
