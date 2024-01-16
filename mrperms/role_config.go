package mrperms

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	rolesPathPattern = "%s/%s.%s" // dir/role.ext, ./roles/administrators.yaml
)

type (
	roleConfig struct {
		Privileges  []string `yaml:"privileges"`
		Permissions []string `yaml:"permissions"`
	}
)

func loadRoleConfig(roleName, dirPath, fileType string) (*roleConfig, error) {
	cfg := roleConfig{}
	filePath := fmt.Sprintf(rolesPathPattern, dirPath, roleName, fileType)

	if err := cleanenv.ReadConfig(filePath, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing role file '%s': %w", filePath, err)
	}

	return &cfg, nil
}
