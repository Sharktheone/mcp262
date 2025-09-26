package runner

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

const DEFAULT_WORKERS = 256
const DEFAULT_TEST_ROOT = "./test262/test"
const REPO_PATH = "./"

type Config struct {
	RepoPath    string `toml:"repo_path"`
	Workers     int    `toml:"workers"`
	TestRootDir string `toml:"test_root_dir"`
}

func NewConfig() *Config {
	return &Config{
		RepoPath:    REPO_PATH,
		Workers:     DEFAULT_WORKERS,
		TestRootDir: DEFAULT_TEST_ROOT,
	}
}

func NewFromEnv() *Config {
	config := NewConfig()

	if repoPath, exists := os.LookupEnv("REPO_PATH"); exists {
		config.RepoPath = repoPath
	}

	if workers, exists := os.LookupEnv("WORKERS"); exists {
		var w int
		_, err := fmt.Sscanf(workers, "%d", &w)
		if err == nil && w > 0 {
			config.Workers = w
		}
	}

	if testRoot, exists := os.LookupEnv("TEST_ROOT_DIR"); exists {
		config.TestRootDir = testRoot
	}

	return config
}

func LoadConfig() *Config {
	config := NewFromEnv()

	configFile := flag.String("config", "config.toml", "Path to TOML config file")
	repoPath := flag.String("repo", config.RepoPath, "Path to external repository for CI results")
	workers := flag.Int("workers", config.Workers, "Number of workers")
	testRootDir := flag.String("test_root", config.TestRootDir, "Path to test root directory")

	flag.Parse()

	if *configFile != "" {
		if err := loadConfigFile(*configFile, config); err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("Failed to load config file %s: %v", *configFile, err)
			}
		}
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "repo":
			config.RepoPath = *repoPath
		case "workers":
			config.Workers = *workers
		case "test_root":
			config.TestRootDir = *testRootDir
		}
	})

	return config
}

func loadConfigFile(filename string, config *Config) error {
	_, err := toml.DecodeFile(filename, config)
	return err
}
