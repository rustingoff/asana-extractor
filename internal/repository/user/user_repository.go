package userrepo

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type UserRepository struct {
	httpClient *http.Client
	authToken  string
	apiUrl     string
}

func NewUserRepository(apiUrl, authToken string) *UserRepository {
	return &UserRepository{
		httpClient: &http.Client{},
		authToken:  authToken,
		apiUrl:     apiUrl,
	}
}

func (ur *UserRepository) GetAllUsers(ctx context.Context, path string) ([]byte, error) {
	var op = "GetAllUsers"

	url := fmt.Sprintf("%s/%s", ur.apiUrl, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get all users, err: %s, op: %s", err.Error(), op))
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ur.authToken))

	response, err := ur.httpClient.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintf("can't make a http request, err: %s, op: %s, status_code: %d", err.Error(), op, response.StatusCode))
		return nil, err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		retryAfterHeader, err := strconv.Atoi(response.Header.Get("Retry-After"))
		if err != nil {
			slog.Error("retry after header is not provided", "err", err)
			return nil, err
		}

		time.Sleep(time.Second * time.Duration(retryAfterHeader))

		response, err := ur.httpClient.Do(req)
		if err != nil {
			slog.Error(fmt.Sprintf("can't make a http request, err: %s, op: %s, status_code: %d", err.Error(), op, response.StatusCode))
			return nil, err
		}
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read response body, err: %s", err.Error()))
		return nil, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			slog.Error(fmt.Sprintf("request body can't be close, err: %s", err.Error()))
		}
	}()

	return data, nil
}

func (ur *UserRepository) GetUserByGID(ctx context.Context, path string) ([]byte, error) {
	var op = "GetUserByGID"

	url := fmt.Sprintf("%s/%s", ur.apiUrl, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get all users, err: %s, op: %s", err.Error(), op))
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ur.authToken))

	response, err := ur.httpClient.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintf("can't make a http request, err: %s, op: %s, status_code: %d", err.Error(), op, response.StatusCode))
		return nil, err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		retryAfterHeader, err := strconv.Atoi(response.Header.Get("Retry-After"))
		if err != nil {
			slog.Error("retry after header is not provided", "err", err)
			return nil, err
		}

		time.Sleep(time.Second * time.Duration(retryAfterHeader))

		response, err := ur.httpClient.Do(req)
		if err != nil {
			slog.Error(fmt.Sprintf("can't make a http request, err: %s, op: %s, status_code: %d", err.Error(), op, response.StatusCode))
			return nil, err
		}
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to read response body, err: %s", err.Error()))
		return nil, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			slog.Error(fmt.Sprintf("request body can't be close, err: %s", err.Error()))
		}
	}()

	return data, nil
}
