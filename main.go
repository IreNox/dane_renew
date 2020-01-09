package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func configPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, err := filepath.Abs(filepath.Dir(exe))
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), nil
}

func main() {
	configPath, err := configPath()
	if err != nil {
		log.Fatal(err)
		return
	}

	var cfg config
	if err := readConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
		return
	}

	rest := newhostingDeRestAPI(cfg.URL, cfg.AuthToken)

	if err := createAuthRecord(rest, "sys.ioniel.net", "..."); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("hello, world: %s\n", configPath)
	fmt.Printf("config: %s\n", cfg)

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.\n", os)
	}
}
