package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func createAIProject(projectName, hubName string) {
	cmd := exec.Command("python3", "/home/tngo/ngo/projects/tngo-ai-svcs/python/createaifoundryproject.py", projectName, hubName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error creating AI Foundry project:", err)
		return
	}

	fmt.Println("✅ AI Foundry project created successfully!")
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("/home/tngo/ngo/projects/tngo-ai-svcs/.env")
	if err != nil {
		fmt.Println("❌ Error loading .env file:", err)
		return
	}

	projectName := "tngodemo1-ai-foundry-project"
	hubName := "tngodemo1aifoundryuseast2"

	createAIProject(projectName, hubName)
}
