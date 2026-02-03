package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Yoxi7/fckgoback/internal/archarchive"
	"github.com/Yoxi7/fckgoback/internal/mirrorlist"
	"github.com/Yoxi7/fckgoback/internal/utils"
)

func showWarning(lang string) bool {
	var message string
	if lang == "ru" {
		message = "ВНИМАНИЕ: Эта утилита изменит ваш mirrorlist и обновит пакеты.\n" +
			"Продолжить? (y/N): "
	} else {
		message = "WARNING: This utility will modify your mirrorlist and update packages.\n" +
			"Continue? (y/N): "
	}

	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func confirm(message string) bool {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func runPacmanUpdate() error {
	var cmd *exec.Cmd
	
	if os.Geteuid() == 0 {
		cmd = exec.Command("pacman", "-Syyuu")
	} else {
		cmd = exec.Command("sudo", "pacman", "-Syyuu")
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func getLocalizedMessage(lang, key string) string {
	messages := map[string]map[string]string{
		"ru": {
			"selecting_date":      "Выбор даты из архива...",
			"checking_availability": "Проверка доступности архива...",
			"backing_up":          "Создание резервной копии mirrorlist...",
			"writing_mirrorlist":  "Запись нового зеркала...",
			"updating":            "Обновление пакетов...",
			"restoring":           "Восстановление mirrorlist...",
			"success":             "Операция завершена успешно!",
			"error":               "Ошибка: %v\n",
			"restoring_on_exit":   "Восстановление mirrorlist при выходе...",
			"select_different_date": "Пожалуйста, выберите другую дату.",
		},
		"en": {
			"selecting_date":      "Selecting date from archive...",
			"checking_availability": "Checking archive availability...",
			"backing_up":          "Backing up mirrorlist...",
			"writing_mirrorlist":  "Writing new mirror...",
			"updating":            "Updating packages...",
			"restoring":           "Restoring mirrorlist...",
			"success":             "Operation completed successfully!",
			"error":               "Error: %v\n",
			"restoring_on_exit":   "Restoring mirrorlist on exit...",
			"select_different_date": "Please select a different date.",
		},
	}

	if msgs, ok := messages[lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	return messages["en"][key]
}

func main() {
	lang := utils.DetectLanguage()

	if err := mirrorlist.CheckRoot(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !showWarning(lang) {
		fmt.Println("Operation cancelled.")
		os.Exit(0)
	}

	// Handle Ctrl+C gracefully - restore mirrorlist on interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n" + getLocalizedMessage(lang, "restoring_on_exit"))
		if err := mirrorlist.RestoreMirrorlist(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to restore mirrorlist: %v\n", err)
		}
		os.Exit(1)
	}()

	archive := archarchive.NewArchArchive()

	fmt.Println(getLocalizedMessage(lang, "selecting_date"))
	url, err := archive.MenuRun()
	if err != nil {
		fmt.Fprintf(os.Stderr, getLocalizedMessage(lang, "error"), err)
		os.Exit(1)
	}

	fmt.Printf("Selected archive URL: %s\n", url)

	fmt.Println(getLocalizedMessage(lang, "checking_availability"))
	if err := archarchive.CheckArchiveAvailability(url); err != nil {
		fmt.Fprintf(os.Stderr, getLocalizedMessage(lang, "error"), err)
		fmt.Println("\n" + getLocalizedMessage(lang, "select_different_date"))
		os.Exit(1)
	}

	fmt.Println(getLocalizedMessage(lang, "backing_up"))
	if err := mirrorlist.BackupMirrorlist(); err != nil {
		fmt.Fprintf(os.Stderr, getLocalizedMessage(lang, "error"), err)
		os.Exit(1)
	}

	fmt.Println(getLocalizedMessage(lang, "writing_mirrorlist"))
	if err := mirrorlist.WriteMirrorlist(url); err != nil {
		fmt.Fprintf(os.Stderr, getLocalizedMessage(lang, "error"), err)
		mirrorlist.RestoreMirrorlist()
		os.Exit(1)
	}

	fmt.Println(getLocalizedMessage(lang, "updating"))
	if err := runPacmanUpdate(); err != nil {
		fmt.Fprintf(os.Stderr, getLocalizedMessage(lang, "error"), err)
		mirrorlist.RestoreMirrorlist()
		os.Exit(1)
	}

	fmt.Println(getLocalizedMessage(lang, "restoring"))
	if err := mirrorlist.RestoreMirrorlist(); err != nil {
		fmt.Fprintf(os.Stderr, getLocalizedMessage(lang, "error"), err)
		os.Exit(1)
	}

	fmt.Println(getLocalizedMessage(lang, "success"))
}

