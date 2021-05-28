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

	if token == "" {
		log.Fatal("No github token found.")
	}

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oathClient := oauth2.NewClient(ctx, tokenSource)
	githubClient := github.NewClient(oathClient)

	r := gin.Default()
	r.GET("/:owner/:repo/workflows", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		workflows, _, err := githubClient.Actions.ListWorkflows(ctx, owner, repo, nil)
		if err != nil {
			c.String(200, err.Error())
		}
		c.JSON(200, workflows)
	})

	r.POST("/:owner/:repo/workflows/:name", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		fileName := c.Param("name")

		githubDispatchEvent := github.CreateWorkflowDispatchEventRequest{"main", nil}
		workflowResponse, err := githubClient.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, fileName, githubDispatchEvent)

		log.Print(err)
		if err != nil {
			c.String(200, err.Error())
		}
		c.JSON(200, workflowResponse)
	})

	r.GET("/:owner/:repo/workflows/:name", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		fileName := c.Param("name")
		workflowResponse, _, err := githubClient.Actions.GetWorkflowByFileName(ctx, owner, repo, fileName)
		if err != nil {
			c.String(200, err.Error())
		}
		c.JSON(200, workflowResponse)
	})

	r.Run("0.0.0.0:8080")
}
