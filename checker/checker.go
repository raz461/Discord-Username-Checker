package checker

import (
	"fmt"
	"users/globals"
	"users/logger"
)

func CheckerInit(username string, thread int) {
	if CheckBlacklist(username) {
		globals.InvalidUsernames++
		logger.Warn(fmt.Sprintf("Username [%s] is blacklisted, skipping check [%d]", username, thread+1))
		return
	}

	taken := CheckUsername(username)

	if taken {
		globals.InvalidUsernames++
		logger.Error(fmt.Sprintf("Username [%s] is taken [%d]", username, thread+1))
		if err := globals.SaveBlackList(username); err != nil {
			logger.Error(fmt.Sprintf("Error saving to blacklist: %v", err))
		}
		return
	}

	globals.ValidUsernames++
	logger.Success(fmt.Sprintf("Username [%s] is available [%d]", username, thread+1))

	if err := globals.SaveBlackList(username); err != nil {
		logger.Error(fmt.Sprintf("Error saving to blacklist: %v", err))
	}

	if err := globals.SaveValidUser(username); err != nil {
		logger.Error(fmt.Sprintf("Error saving valid user: %v", err))
	}

	if globals.Config.Webhook == "" {
		return
	}

	if err := globals.SendDiscordWebhook(username); err != nil {
		logger.Error(fmt.Sprintf("Error sending Discord webhook: %v", err))
		return
	}
}
