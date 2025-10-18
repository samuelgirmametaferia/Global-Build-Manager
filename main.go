package main

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Config struct {
	BuildDir string
	BuildCmd string
	RunCmd   string
	Ignore   []string
}

func readConfig(path string) (Config, error) {
	file, err := os.Open(filepath.Join(path, "exec.gbf"))
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	cfg := Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "build="):
			cfg.BuildDir = strings.TrimSpace(strings.TrimPrefix(line, "build="))
		case strings.HasPrefix(line, "build_cmd="):
			cfg.BuildCmd = strings.TrimSpace(strings.TrimPrefix(line, "build_cmd="))
		case strings.HasPrefix(line, "run_cmd="):
			cfg.RunCmd = strings.TrimSpace(strings.TrimPrefix(line, "run_cmd="))
		case strings.HasPrefix(line, "ignore="):
			raw := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "ignore=")), ",")
			for _, s := range raw {
				if trimmed := strings.TrimSpace(s); trimmed != "" {
					cfg.Ignore = append(cfg.Ignore, trimmed)
				}
			}
		}
	}
	return cfg, scanner.Err()
}

func killProcess(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		color.Yellow("Stopping previous process PID %d...", cmd.Process.Pid)
		_ = cmd.Process.Kill()
		cmd.Wait()
	}
}

func containsIgnored(path string, ignores []string) bool {
	for _, ignore := range ignores {
		if ignore != "" && strings.Contains(path, ignore) {
			return true
		}
	}
	return false
}

func addDirRecursive(watcher *fsnotify.Watcher, path string) error {
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(p)
		}
		return nil
	})
}

func main() {
	watchPath := "."
	if len(os.Args) > 1 {
		watchPath = os.Args[1]
	}

	absPath, err := filepath.Abs(watchPath)
	if err != nil {
		color.Red("Failed to get absolute path: %v", err)
		return
	}
	color.Cyan("Watching folder: %s", absPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		color.Red("Failed to create watcher: %v", err)
		return
	}
	defer watcher.Close()

	if err := addDirRecursive(watcher, absPath); err != nil {
		color.Red("Failed to watch directories: %v", err)
		return
	}

	var mu sync.Mutex
	var currentProcess *exec.Cmd
	var lastChange time.Time
	var projectAverage = 1 * time.Second

	// Map to store idle times per file
	idleTimes := make(map[string]time.Duration)

	cfg, err := readConfig(absPath)
	if err != nil {
		color.Red("Error reading config: %v", err)
		return
	}

	triggerBuild := func() {
		mu.Lock()
		defer mu.Unlock()

		newCfg, err := readConfig(absPath)
		if err == nil {
			cfg = newCfg
		}

		killProcess(currentProcess)

		color.Green("Building...")
		buildParts := strings.Fields(cfg.BuildCmd)
		if len(buildParts) == 0 {
			color.Red("No build command specified")
			return
		}
		buildCmd := exec.Command(buildParts[0], buildParts[1:]...)
		buildCmd.Dir = cfg.BuildDir
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			color.Red("Build failed: %v", err)
			currentProcess = nil
			return
		}
		color.Green("Build successful!")

		color.Green("Running...")
		runParts := strings.Fields(cfg.RunCmd)
		if len(runParts) == 0 {
			color.Red("No run command specified")
			return
		}
		currentProcess = exec.Command(runParts[0], runParts[1:]...)
		currentProcess.Stdout = os.Stdout
		currentProcess.Stderr = os.Stderr
		currentProcess.Stdin = os.Stdin
		if err := currentProcess.Start(); err != nil {
			color.Red("Failed to run: %v", err)
			currentProcess = nil
		} else {
			color.Blue("Process PID: %d", currentProcess.Process.Pid)
		}
	}

	// Initial build/run
	triggerBuild()

	debounce := make(chan struct{}, 1)

	go func() {
		for range debounce {
			mu.Lock()
			since := time.Since(lastChange)
			delay := projectAverage - since
			if delay < 0 {
				delay = 0
			}
			mu.Unlock()

			time.Sleep(delay)

			triggerBuild()

			// Update adaptive projectAverage
			mu.Lock()
			total := time.Duration(0)
			count := 0
			for _, v := range idleTimes {
				total += v
				count++
			}
			if count > 0 {
				projectAverage = total / time.Duration(count)
			}
			mu.Unlock()
		}
	}()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Watch for new directories
			if event.Op&fsnotify.Create != 0 {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					_ = addDirRecursive(watcher, event.Name)
				}
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if strings.HasSuffix(event.Name, ".go") || filepath.Base(event.Name) == "exec.gbf" {
					if containsIgnored(event.Name, cfg.Ignore) {
						continue
					}

					now := time.Now()
					mu.Lock()
					if !lastChange.IsZero() {
						idleTimes[event.Name] = now.Sub(lastChange)
					} else {
						idleTimes[event.Name] = projectAverage
					}
					lastChange = now
					mu.Unlock()

					select {
					case debounce <- struct{}{}:
					default:
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			color.Red("Watcher error: %v", err)
		}
	}
}
