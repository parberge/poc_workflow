package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	githubOwner := "parberge"
	githubRepo := "poc_workflow"
	githubWorkflowFileName := "poc_workflow.yml"

	githubDispatchEvent := github.CreateWorkflowDispatchEventRequest{"main", nil}

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oathClient := oauth2.NewClient(ctx, tokenSource)

	githubClient := github.NewClient(oathClient)

	// Trigger workflow
	workflowResponse, err := githubClient.Actions.CreateWorkflowDispatchEventByFileName(ctx, githubOwner, githubRepo, githubWorkflowFileName, githubDispatchEvent)

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/trigger", func(c *gin.Context) {
		c.JSON(200, workflowResponse)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
