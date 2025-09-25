package filestorage

import (
	"os"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr"
	"gopkg.in/yaml.v3"
)

const (
	roleFileType = "yaml"
)

type (
	roleConfig struct {
		Privileges  []string `yaml:"privileges"`
		Permissions []string `yaml:"permissions"`
	}
)

// ErrParsingRoleFileFailed - parsing role file failed.
var ErrParsingRoleFileFailed = mrerr.NewKindInternal("parsing role file failed: '{Path}'")

func loadRoleConfig(filePath string) (*roleConfig, error) {
	cfg := roleConfig{}

	if err := parseFile(filePath, &cfg); err != nil {
		return nil, ErrParsingRoleFileFailed.Wrap(err, filePath)
	}

	return &cfg, nil
}

func parseFile(path string, data any) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	defer f.Close()

	return yaml.NewDecoder(f).Decode(data)
}

func getFilePath(dirPath, name string) string {
	// dir/role.ext: ./roles/administrators.yaml
	return strings.TrimRight(dirPath, "/") + "/" + strings.Trim(name, "/") + "." + roleFileType
}
