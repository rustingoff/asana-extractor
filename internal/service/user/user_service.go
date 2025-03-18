package userservice

import (
	"context"
	"encoding/json"
	"errors"
	"extractor/internal/config"
	"fmt"
	"log/slog"
	"os"
)

var ErrUsersNotFound = errors.New("users not found")

type IUserRepository interface {
	GetAllUsers(ctx context.Context, path string) ([]byte, error)
	GetUserByGID(ctx context.Context, path string) ([]byte, error)
}

type UserService struct {
	userRepo IUserRepository
	cfg      config.UserConfig
}

func NewUserService(userRepo IUserRepository, cfg config.UserConfig) *UserService {
	return &UserService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (us *UserService) GetAll(ctx context.Context, workspaceID string) error {
	var path = us.cfg.Path
	var data map[string]any
	var users []map[string]any

	if workspaceID != "" {
		path = fmt.Sprintf("%s?workspace=%s&limit=%d", us.cfg.Path, workspaceID, us.cfg.Limit)
	}

	for {
		result, err := us.userRepo.GetAllUsers(ctx, path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(result, &data); err != nil {
			slog.Error("failed to unmarshal data to json", "err", err)
			return err
		}

		for _, user := range data["data"].([]interface{}) {
			users = append(users, user.(map[string]any))
		}

		if _, exists := data["next_page"]; !exists || data["next_page"] == nil {
			break
		}

		offset := data["next_page"].(map[string]any)["offset"]
		path = fmt.Sprintf("%s?workspace=%s&limit=%d&offset=%s", us.cfg.Path, workspaceID, us.cfg.Limit, offset.(string))
	}

	if err := us.storeUsersData(ctx, users); err != nil {
		return err
	}

	return nil
}

func (us *UserService) storeUsersData(ctx context.Context, users []map[string]any) error {
	if len(users) == 0 {
		return ErrUsersNotFound
	}

	for _, user := range users {
		var path = fmt.Sprintf("%s/%s", us.cfg.Path, user["gid"])
		userData, err := us.userRepo.GetUserByGID(ctx, path)
		if err != nil {
			return err
		}
		file, err := os.OpenFile("./store/users/"+user["gid"].(string)+".json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("failed to open file", "err", err)
			return err
		}

		_, err = file.WriteString(string(userData))
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
