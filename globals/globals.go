package globals

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"users/logger"
	"users/types"

	discordwebhook "github.com/bensch777/discord-webhook-golang"
)

const (
	// API endpoints
	DiscordUsernameCheckAPI = "https://discord.com/api/v9/unique-username/username-attempt-unauthed"

	// File paths
	ConfigFile    = "data/config.json"
	ProxiesFile   = "data/proxies.txt"
	UsernamesFile = "data/usernames.txt"
	BlacklistFile = "data/blacklist.txt"
	ValidsFile    = "data/valids.txt"
)

var (
	Config    types.Config
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"

	Proxies   = []string{}
	Usernames = []string{}
	BlackList = []string{}

	ValidUsernames   int64
	InvalidUsernames int64

	blacklistMutex     sync.Mutex
	validUsernameMutex sync.Mutex
)

func GenerateRandomUsername(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	username := make([]byte, length)
	for i := range username {
		username[i] = charset[rand.Intn(len(charset))]
	}
	return string(username), nil
}

func LoadConfig() error {
	jsonFile, err := os.Open(ConfigFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error closing config file: %v", err))
		}
	}()

	if err := json.NewDecoder(jsonFile).Decode(&Config); err != nil {
		return err
	}

	return nil
}

func LoadProxies() error {
	file, err := os.ReadFile(ProxiesFile)
	if err != nil {
		// If file doesn't exist, that's okay - we can run without proxies
		if os.IsNotExist(err) {
			logger.Warn("No proxies.txt file found, running without proxies")
			return nil
		}
		return err
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "Proxies goes here" {
			continue
		}
		// Validate proxy format
		if !strings.Contains(trimmed, ":") {
			logger.Warn(fmt.Sprintf("Invalid proxy format: %s", trimmed))
			continue
		}
		if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
			trimmed = "http://" + trimmed
		}
		Proxies = append(Proxies, trimmed)
	}

	return nil
}

func LoadUsernames() error {
	if !Config.Usernames.Custom {

		for i := 0; i < Config.Usernames.Amount; i++ {
			username, err := GenerateRandomUsername(Config.Usernames.Length)
			if err != nil {
				return err
			}
			Usernames = append(Usernames, username)
		}

	} else {
		file, err := os.ReadFile(UsernamesFile)
		if err != nil {
			return err
		}

		lines := strings.Split(string(file), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			Usernames = append(Usernames, line)
		}
	}

	return nil
}

func LoadBlackList() error {
	file, err := os.ReadFile(BlacklistFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		BlackList = append(BlackList, line)
	}

	return nil
}

func SaveBlackList(username string) error {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	for _, blacklisted := range BlackList {
		if username == blacklisted {
			return nil
		}
	}

	BlackList = append(BlackList, username)

	file, err := os.Create(BlacklistFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error closing blacklist file: %v", err))
		}
	}()

	for _, blacklisted := range BlackList {
		if _, err := file.WriteString(blacklisted + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func SaveValidUser(username string) error {
	validUsernameMutex.Lock()
	defer validUsernameMutex.Unlock()

	file, err := os.OpenFile(ValidsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error closing valids file: %v", err))
		}
	}()

	if _, err := file.WriteString(username + "\n"); err != nil {
		return err
	}

	return nil
}

func SendDiscordWebhook(username string) error {
	embed := discordwebhook.Embed{
		Title:     "Valid Username Found",
		Color:     rand.Intn(0xFFFFFF),
		Url:       "https://undesync.com",
		Timestamp: time.Now(),
		Thumbnail: discordwebhook.Thumbnail{
			Url: "https://i.imgur.com/Z6K0B0s.png",
		},
		Fields: []discordwebhook.Field{
			{
				Name:   "Username",
				Value:  username,
				Inline: true,
			},
		},
		Footer: discordwebhook.Footer{
			Text:     "Undesync Name Checker",
			Icon_url: "https://i.imgur.com/Z6K0B0s.png",
		},
	}

	if err := SendEmbed(Config.Webhook, embed); err != nil {
		logger.Error(fmt.Sprintf("Error sending embed: %v", err))
		return err
	}
	return nil
}

func SendEmbed(link string, embeds discordwebhook.Embed) error {

	hook := discordwebhook.Hook{
		Username:   "Username Checker",
		Avatar_url: "https://i.imgur.com/Z6K0B0s.png",
		Embeds:     []discordwebhook.Embed{embeds},
	}

	payload, err := json.Marshal(hook)
	if err != nil {
		logger.Error(err.Error())
	}
	err = discordwebhook.ExecuteWebhook(link, payload)
	return err

}
