package mirrorlist

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	MirrorlistPath  = "/etc/pacman.d/mirrorlist"
	MirrorlistBackup = "/etc/pacman.d/mirrorlist.bak"
)

func CheckRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}
	return fmt.Errorf("this program must be run as root. Please run with: sudo %s", os.Args[0])
}

func BackupMirrorlist() error {
	if _, err := os.Stat(MirrorlistPath); os.IsNotExist(err) {
		return fmt.Errorf("mirrorlist not found: %s", MirrorlistPath)
	}

	srcFile, err := os.Open(MirrorlistPath)
	if err != nil {
		return fmt.Errorf("failed to open mirrorlist: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(MirrorlistBackup)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// WriteMirrorlist writes archive URL to mirrorlist with $repo/os/$arch format
func WriteMirrorlist(url string) error {
	file, err := os.Create(MirrorlistPath)
	if err != nil {
		return fmt.Errorf("failed to create mirrorlist: %w", err)
	}
	defer file.Close()

	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	
	mirrorURL := url + "$repo/os/$arch"
	
	_, err = fmt.Fprintf(file, "Server = %s\n", mirrorURL)
	if err != nil {
		return fmt.Errorf("failed to write mirrorlist: %w", err)
	}

	return nil
}

func RestoreMirrorlist() error {
	if _, err := os.Stat(MirrorlistBackup); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", MirrorlistBackup)
	}

	if err := os.Remove(MirrorlistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove mirrorlist: %w", err)
	}

	srcFile, err := os.Open(MirrorlistBackup)
	if err != nil {
		return fmt.Errorf("failed to open backup: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(MirrorlistPath)
	if err != nil {
		return fmt.Errorf("failed to create mirrorlist: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to restore mirrorlist: %w", err)
	}

	return nil
}



