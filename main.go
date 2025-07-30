package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"
	"users/checker"
	"users/globals"
	"users/logger"
)

var mu sync.Mutex

func init() {
	if err := globals.LoadConfig(); err != nil {
		logger.Error(fmt.Sprintf("Error loading config: %v", err))
		logger.Error("Please ensure data/config.json exists and is valid")
		os.Exit(1)
	}
	logger.Success("Config loaded successfully")

	if err := globals.LoadProxies(); err != nil {
		logger.Error(fmt.Sprintf("Error loading proxies: %v", err))
		// Don't exit here, we can run without proxies
	} else {
		logger.Info(fmt.Sprintf("Proxies loaded successfully: %d", len(globals.Proxies)))
	}

	if err := globals.LoadUsernames(); err != nil {
		logger.Error(fmt.Sprintf("Error loading usernames: %v", err))
		logger.Error("Please check your username configuration")
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Usernames loaded successfully: %d", len(globals.Usernames)))

	if len(globals.Usernames) == 0 {
		logger.Error("No usernames to check. Please configure usernames in config.json or data/usernames.txt")
		os.Exit(1)
	}

	if err := globals.LoadBlackList(); err != nil {
		logger.Warn(fmt.Sprintf("Error loading blacklist: %v (will create new one)", err))
	} else {
		logger.Info(fmt.Sprintf("Blacklist loaded successfully: %d", len(globals.BlackList)))
	}
}

func main() {
	logger.Title("Checker Started", "cyan")

	userCount := len(globals.Usernames)

	go func() {
		setTitle := func(title string) {
			if runtime.GOOS == "windows" {
				kernel32, err := syscall.LoadLibrary("kernel32.dll")
				if err != nil {
					logger.Error(fmt.Sprintf("Error loading kernel32.dll: %v", err))
					return
				}
				defer func() {
					if err := syscall.FreeLibrary(kernel32); err != nil {
						logger.Error(fmt.Sprintf("Error freeing kernel32.dll: %v", err))
					}
				}()

				setConsoleTitle, err := syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
				if err != nil {
					logger.Error(fmt.Sprintf("Error getting SetConsoleTitleW proc: %v", err))
					return
				}

				ptr, err := syscall.UTF16PtrFromString(title)
				if err != nil {
					logger.Error(fmt.Sprintf("Error converting title to UTF16: %v", err))
					return
				}

				if _, _, err := syscall.SyscallN(setConsoleTitle, uintptr(unsafe.Pointer(ptr)), 0, 0); err != 0 {
					logger.Error(fmt.Sprintf("Error setting console title: %v", err))
				}
			} else {
				fmt.Printf("\033]0;%s\007", title)
			}
		}

		for {
			mu.Lock()
			remainingChecks := int64(userCount) - (globals.ValidUsernames + globals.InvalidUsernames)
			title := fmt.Sprintf(".gg/undesync | Valid: %d | Failed: %d | Remaining: %d",
				globals.ValidUsernames, globals.InvalidUsernames, remainingChecks)
			mu.Unlock()

			setTitle(title)

			time.Sleep(1 * time.Second)
		}
	}()

	usernameChannel := make(chan string, globals.Config.Threads)
	var wg sync.WaitGroup

	for i := 1; i <= globals.Config.Threads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for username := range usernameChannel {
				checker.CheckerInit(username, workerID)
			}
		}(i)
	}

	go func() {
		for _, username := range globals.Usernames {
			usernameChannel <- username
		}
		close(usernameChannel)
	}()

	wg.Wait()
	logger.Info("All checks completed.")
}
