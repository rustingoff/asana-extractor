package projectservice

import (
	"context"
	"encoding/json"
	"errors"
	"extractor/internal/config"
	"fmt"
	"log/slog"
	"os"
)

var ErrProjectsNotFound = errors.New("projects not found")

type IProjectRepository interface {
	GetAllProjects(ctx context.Context, path string) ([]byte, error)
	GetProjectByGID(ctx context.Context, path string) ([]byte, error)
}

type ProjectService struct {
	projectRepo IProjectRepository
	cfg         config.ProjectConfig
}

func NewProjectService(projectRepo IProjectRepository, cfg config.ProjectConfig) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		cfg:         cfg,
	}
}

func (us *ProjectService) GetAll(ctx context.Context, workspaceID string) error {
	var path = us.cfg.Path
	var data map[string]any
	var projects []map[string]any

	if workspaceID != "" {
		path = fmt.Sprintf("%s?workspace=%s&limit=%d", us.cfg.Path, workspaceID, us.cfg.Limit)
	}

	for {
		result, err := us.projectRepo.GetAllProjects(ctx, path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(result, &data); err != nil {
			slog.Error("failed to unmarshal data to json", "err", err)
			return err
		}

		for _, project := range data["data"].([]interface{}) {
			projects = append(projects, project.(map[string]any))
		}

		if _, exists := data["next_page"]; !exists || data["next_page"] == nil {
			break
		}

		offset := data["next_page"].(map[string]any)["offset"]
		path = fmt.Sprintf("%s?workspace=%s&limit=%d&offset=%s", us.cfg.Path, workspaceID, us.cfg.Limit, offset.(string))
	}

	if err := us.storeProjectsData(ctx, projects); err != nil {
		return err
	}

	return nil
}

func (us *ProjectService) storeProjectsData(ctx context.Context, projects []map[string]any) error {
	if len(projects) == 0 {
		return ErrProjectsNotFound
	}

	for _, project := range projects {
		var path = fmt.Sprintf("%s/%s", us.cfg.Path, project["gid"])
		projectData, err := us.projectRepo.GetProjectByGID(ctx, path)
		if err != nil {
			return err
		}
		file, err := os.OpenFile("./store/projects/"+project["gid"].(string)+".json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("failed to open file", "err", err)
			return err
		}

		_, err = file.WriteString(string(projectData))
		if err != nil {
			slog.Error("failed to write data to file", "err", err)
			return err
		}
		if err := file.Close(); err != nil {
			slog.Error("failed to close file", "err", err)
			return err
		}
	}

	return nil
}
