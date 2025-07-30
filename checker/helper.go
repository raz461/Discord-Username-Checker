package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"users/globals"
	"users/logger"
	"users/types"
)

func GetRandomProxy() (string, error) {
	if len(globals.Proxies) == 0 {
		return "", fmt.Errorf("no proxies available")
	}

	randomIndex := rand.Intn(len(globals.Proxies))
	return globals.Proxies[randomIndex], nil
}

func CheckBlacklist(username string) bool {
	for _, blacklisted := range globals.BlackList {
		if username == blacklisted {
			return true
		}
	}
	return false
}

func CheckUsername(username string) bool {
	maxAttempts := globals.Config.Retry.MaxAttempts
	if !globals.Config.Retry.Enabled {
		maxAttempts = 1
	}

	requestBody := types.UsernameRequest{Username: username}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Error marshaling request body: %v", err))
		return true
	}

	proxy, _ := GetRandomProxy()
	useProxy := proxy != ""

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		var client *http.Client
		if useProxy {
			proxyURL, err := url.Parse(proxy)
			if err != nil {
				logger.Error(fmt.Sprintf("Invalid proxy URL: %v", err))
				if attempt < maxAttempts {
					logger.Info(fmt.Sprintf("Retrying request for username [%s], attempt %d", username, attempt+1))
					continue
				}
				return true
			}
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			client = &http.Client{Transport: transport, Timeout: time.Duration(globals.Config.Timeout) * time.Second}
		} else {
			client = &http.Client{Timeout: time.Duration(globals.Config.Timeout) * time.Second}
		}

		req, err := http.NewRequest(http.MethodPost, globals.DiscordUsernameCheckAPI, bytes.NewBuffer(jsonBody))
		if err != nil {
			logger.Error(fmt.Sprintf("Error creating request: %v", err))
			if attempt < maxAttempts {
				logger.Info(fmt.Sprintf("Retrying request for username [%s], attempt %d", username, attempt+1))
				continue
			}
			return true
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", globals.UserAgent)

		res, err := client.Do(req)
		if err != nil {
			logger.Error(fmt.Sprintf("Error making request: %v", err))
			if attempt < maxAttempts {
				logger.Info(fmt.Sprintf("Retrying request for username [%s], attempt %d", username, attempt+1))
				continue
			}
			return true
		}

		defer func() {
			if err := res.Body.Close(); err != nil {
				logger.Error(fmt.Sprintf("Error closing response body: %v", err))
			}
		}()

		if res.StatusCode != http.StatusOK {
			logger.Error(fmt.Sprintf("API returned status code %d for username [%s]", res.StatusCode, username))
			if attempt < maxAttempts {
				logger.Info(fmt.Sprintf("Retrying request for username [%s], attempt %d", username, attempt+1))
				continue
			}
			return true
		}

		var usernameResponse types.UsernameResponse
		if err := json.NewDecoder(res.Body).Decode(&usernameResponse); err != nil {
			logger.Error(fmt.Sprintf("Error decoding response: %v", err))
			if attempt < maxAttempts {
				logger.Info(fmt.Sprintf("Retrying request for username [%s], attempt %d", username, attempt+1))
				continue
			}
			return true
		}

		return usernameResponse.Taken
	}

	return true
}
