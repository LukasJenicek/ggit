package config

import (
	"fmt"
	"strings"
)

func ParseIniConfig(content []byte) (map[string]map[string]any, error) {
	cfgStruct := map[string]map[string]any{}

	section := ""

	var err error

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// load section
		if line[0] == '[' {
			section, err = loadSection(line)
			if err != nil {
				return nil, fmt.Errorf("parse section: %w", err)
			}

			_, ok := cfgStruct[section]
			if !ok {
				cfgStruct[section] = make(map[string]any)
			}

			continue
		}

		key, value, err := loadSectionKeyAndValue(line)
		if err != nil {
			return nil, fmt.Errorf("parse section key and value: %w", err)
		}

		cfgStruct[section][key] = value
	}

	return cfgStruct, nil
}

func loadSectionKeyAndValue(line string) (string, any, error) {
	// load section values
	eq := strings.Index(line, "=")
	if eq == -1 {
		return "", nil, fmt.Errorf("parse section values: invalid line %s", line)
	}

	key := strings.TrimSpace(line[:eq])
	value := strings.TrimSpace(line[eq+1:])

	return key, value, nil
}

func loadSection(line string) (string, error) {
	end := strings.Index(line, "]")
	if end == -1 {
		return "", fmt.Errorf(`section end "%s" not found`, line)
	}

	section := strings.TrimSpace(line[1:end])

	if section == "" {
		return "", fmt.Errorf(`section %s is empty`, line)
	}

	return section, nil
}
