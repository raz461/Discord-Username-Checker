package logger

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

var (
	infoColor     = color.New(color.FgCyan).SprintFunc()
	errorColor    = color.New(color.FgRed).SprintFunc()
	warnColor     = color.New(color.FgYellow).SprintFunc()
	successColor  = color.New(color.FgGreen).SprintFunc()
	debugColor    = color.New(color.FgBlue).SprintFunc()
	threadIDColor = color.New(color.FgMagenta, color.Bold).SprintFunc()
)

var threadIDPattern = regexp.MustCompile(`\[(\d+)\]`)

func getTime() string {
	return time.Now().Format("15:04")
}

func colorizeThreadID(msg string) string {
	return threadIDPattern.ReplaceAllStringFunc(msg, func(match string) string {
		id := threadIDPattern.FindStringSubmatch(match)[1]
		return "[" + threadIDColor(id) + "]"
	})
}

func Info(msg string) {
	coloredMsg := colorizeThreadID(msg)
	if _, err := fmt.Fprintf(os.Stdout, "[%s] [%s] %s\n", getTime(), infoColor("INFO"), coloredMsg); err != nil {
		fmt.Fprintf(os.Stderr, "[LOGGER ERROR] %v\n", err)
	}
}

func Error(msg string) {
	coloredMsg := colorizeThreadID(msg)
	if _, err := fmt.Fprintf(os.Stderr, "[%s] [%s] %s\n", getTime(), errorColor("ERROR"), coloredMsg); err != nil {
		fmt.Fprintf(os.Stderr, "[LOGGER ERROR] %v\n", err)
	}
}

func Warn(msg string) {
	coloredMsg := colorizeThreadID(msg)
	if _, err := fmt.Fprintf(os.Stdout, "[%s] [%s] %s\n", getTime(), warnColor("WARN"), coloredMsg); err != nil {
		fmt.Fprintf(os.Stderr, "[LOGGER ERROR] %v\n", err)
	}
}

func Success(msg string) {
	coloredMsg := colorizeThreadID(msg)
	if _, err := fmt.Fprintf(os.Stdout, "[%s] [%s] %s\n", getTime(), successColor("SUCCESS"), coloredMsg); err != nil {
		fmt.Fprintf(os.Stderr, "[LOGGER ERROR] %v\n", err)
	}
}

func Debug(msg string) {
	coloredMsg := colorizeThreadID(msg)
	if _, err := fmt.Fprintf(os.Stdout, "[%s] [%s] %s\n", getTime(), debugColor("DEBUG"), coloredMsg); err != nil {
		fmt.Fprintf(os.Stderr, "[LOGGER ERROR] %v\n", err)
	}
}

func Title(title, colorType string) {
	fig := figure.NewFigure(title, "", true)
	coloredText := fig.String()

	switch colorType {
	case "info":
		coloredText = infoColor(coloredText)
	case "error":
		coloredText = errorColor(coloredText)
	case "warn":
		coloredText = warnColor(coloredText)
	case "success":
		coloredText = successColor(coloredText)
	case "debug":
		coloredText = debugColor(coloredText)
	}

	fmt.Println(coloredText)
	fmt.Println()
}
