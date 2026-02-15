package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetLevel(logrus.DebugLevel)
	token := os.Getenv("TODOIST_TOKEN")
	if token == "" {
		t.Skip("TODOIST_TOKEN not set")
	}
	c, err := NewClient(token)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	name := fmt.Sprint(time.Now().Unix())
	projectId := ""
	t.Run("create project", func(t *testing.T) {
		project := &CreateProject{
			Name: &name,
		}
		res, err := c.CreateProject(ctx, *project)
		projectId = res.ID
		assert.Nil(t, err, "expecting nil error")
		assert.NotNil(t, res, "expecting non-nil result")
	})

	t.Run("read project", func(t *testing.T) {
		res, err := c.GetProject(ctx, projectId)
		assert.Nil(t, err, "expecting nil error")
		assert.NotNil(t, res, "expecting non-nil result")
		assert.Equal(t, res.Name, name, fmt.Sprintf("expecting %s, got %s", name, res.Name))
	})

	t.Run("update project", func(t *testing.T) {
		name = "new name"
		project := &UpdateProject{
			ID:   &projectId,
			Name: &name,
		}
		res, err := c.UpdateProject(ctx, *project)
		assert.Nil(t, err, "expecting nil error")
		assert.Equal(t, res.Name, "new name", "expecting project name to be %s, got %s", &name, res.Name)

	})

	t.Run("delete project", func(t *testing.T) {
		statuscode, _, err := c.DeleteProject(ctx, projectId)
		assert.Nil(t, err, "expecting nil error")
		assert.True(t, statuscode == 200 || statuscode == 204, "expecting statuscode to be 200 or 204, got %d", statuscode)
	})

	t.Run("check that deleted projects is not active", func(t *testing.T) {
		res, err := c.GetProject(ctx, projectId)
		assert.NotNil(t, err, "expecting nil error")
		assert.Nil(t, res, "expecting nil result")
	})

}
