package cmd

import (
	"os"
	"path/filepath"
)

var homeDir, _ = os.UserHomeDir()
var defaultConfigFile = filepath.Join(homeDir, ".koffer", "config.yml")
