package main

import (
	"context"
	"extractor/internal/config"
	projectrepo "extractor/internal/repository/project"
	userrepo "extractor/internal/repository/user"
	projectservice "extractor/internal/service/project"
	userservice "extractor/internal/service/user"
	"flag"
	"log"
	"log/slog"
	"sync"
	"time"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "./config/config.yaml", "path to a config file")

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := userrepo.NewUserRepository(config.App.APIUrl, config.App.APIAuthToken)
	userSvc := userservice.NewUserService(userRepo, config.App.UserConfig)

	projectRepo := projectrepo.NewProjectRepository(config.App.APIUrl, config.App.APIAuthToken)
	projectSvc := projectservice.NewProjectService(projectRepo, config.App.ProjectConfig)

	userSvc.GetAll(context.Background(), config.App.WorkspaceGID)
	projectSvc.GetAll(context.Background(), config.App.WorkspaceGID)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Second * 30)
			if err := userSvc.GetAll(context.Background(), config.App.WorkspaceGID); err != nil {
				slog.Error("can't get all users", "err", err)
			}
			if err := projectSvc.GetAll(context.Background(), config.App.WorkspaceGID); err != nil {
				slog.Error("can't get all projects", "err", err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Minute * 5)
			if err := userSvc.GetAll(context.Background(), config.App.WorkspaceGID); err != nil {
				slog.Error("can't get all users", "err", err)
			}
			if err := projectSvc.GetAll(context.Background(), config.App.WorkspaceGID); err != nil {
				slog.Error("can't get all projects", "err", err)
			}
		}
	}()

	wg.Wait()
}
