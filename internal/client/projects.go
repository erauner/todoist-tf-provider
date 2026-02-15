package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	Color          string `json:"color"`
	IsShared       bool   `json:"is_shared"`
	Order          int    `json:"order"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"`
	URL            string `json:"url"`
	ParentID       string `json:"parent_id"`
}

// projectV1 matches the Unified API "project" response shape (subset used by this provider).
// Some fields differ from REST v2, so we map them into the legacy Project struct above to
// preserve the Terraform schema.
type projectV1 struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Color        any    `json:"color"` // string name or integer id
	IsShared     bool   `json:"is_shared"`
	ChildOrder   int    `json:"child_order"`
	IsFavorite   bool   `json:"is_favorite"`
	InboxProject bool   `json:"inbox_project"`
	IsTeamInbox  bool   `json:"is_team_inbox"`
	ViewStyle    string `json:"view_style"`
	URL          string `json:"url"`
	ParentID     string `json:"parent_id"`
}

func colorToString(v any) string {
	switch c := v.(type) {
	case nil:
		return ""
	case string:
		return c
	case float64:
		// JSON numbers decode as float64.
		if name, ok := todoistColorIDToName[int(c)]; ok {
			return name
		}
		return strconv.Itoa(int(c))
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Todoist color IDs table (subset).
// This is used when the Unified API returns an integer color id instead of a name.
var todoistColorIDToName = map[int]string{
	30: "berry_red",
	31: "red",
	32: "orange",
	33: "yellow",
	34: "olive_green",
	35: "lime_green",
	36: "green",
	37: "mint_green",
	38: "teal",
	39: "sky_blue",
	40: "light_blue",
	41: "blue",
	42: "grape",
	43: "violet",
	44: "lavender",
	45: "magenta",
	46: "salmon",
	47: "charcoal",
	48: "grey",
	49: "taupe",
}

func (c *Client) GetProject(ctx context.Context, projectId string) (*Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s", c.BaseURL, projectId), nil)
	log.WithFields(log.Fields{
		"projectId": projectId,
	}).Info("Reading project")
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	var v1 projectV1
	_, _, err = c.sendRequest(req, &v1)
	if err != nil {
		return nil, err
	}

	res := Project{
		ID:             v1.ID,
		Name:           v1.Name,
		Color:          colorToString(v1.Color),
		IsShared:       v1.IsShared,
		Order:          v1.ChildOrder,
		IsFavorite:     v1.IsFavorite,
		IsInboxProject: v1.InboxProject,
		IsTeamInbox:    v1.IsTeamInbox,
		ViewStyle:      v1.ViewStyle,
		URL:            v1.URL,
		ParentID:       v1.ParentID,
	}

	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project read")
	return &res, nil
}

type CreateProject struct {
	// Required fields
	Name *string `json:"name"`
	// Optional fields
	ParentID   *string `json:"parent_id,omitempty"`
	Color      *string `json:"color,omitempty"`
	IsFavorite *bool   `json:"is_favorite,omitempty"`
	ViewStyle  *string `json:"view_style,omitempty"`
}

func (c *Client) CreateProject(ctx context.Context, createProject CreateProject) (*Project, error) {
	payload, err := json.Marshal(createProject)
	log.WithFields(log.Fields{
		"payload": string(payload),
	}).Info("Creating project")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects", c.BaseURL), bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res := Project{}
	_, _, err = c.sendRequest(req, &res)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project created")
	return &res, nil
}

type UpdateProject struct {
	ID         *string
	Name       *string `json:"name,omitempty"`
	Color      *string `json:"color,omitempty"`
	IsFavorite *bool   `json:"is_favorite,omitempty"`
	ViewStyle  *string `json:"view_style,omitempty"`
}

func (c *Client) UpdateProject(ctx context.Context, updateProject UpdateProject) (*Project, error) {
	payload, err := json.Marshal(updateProject)
	if updateProject.ID == nil {
		return nil, fmt.Errorf("missing project id")
	}
	log.WithFields(log.Fields{
		"payload":   string(payload),
		"projectId": updateProject.ID,
	}).Info("Updating project")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s", c.BaseURL, *updateProject.ID), bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res := Project{}
	_, _, err = c.sendRequest(req, &res)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project updated")
	return &res, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectId string) (statusCode int, body string, err error) {
	log.WithFields(log.Fields{
		"projectId": projectId,
	}).Info("Deleting project")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", c.BaseURL, projectId), nil)
	if err != nil {
		return 0, "", err
	}
	req = req.WithContext(ctx)
	statusCode, body, err = c.sendRequest(req, nil)
	if err != nil {
		return -1, "", err
	}
	return statusCode, body, err
}
