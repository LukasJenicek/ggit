package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

type Config struct {
	User *User `config:"user"`
}

type User struct {
	Email string `config:"email"`
	Name  string `config:"name"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".config/git/config")

	cfgContent, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %s: %w", configPath, err)
	}

	c, err := ParseIniConfig(cfgContent)
	if err != nil {
		return nil, fmt.Errorf("parse git config content %s: %w", configPath, err)
	}

	cfg := &Config{}
	if err = setConfigValues(reflect.ValueOf(cfg), "", c); err != nil {
		return nil, fmt.Errorf("parse git config content %s: %w", configPath, err)
	}

	return nil, nil
}

func setConfigValues(v reflect.Value, section string, values map[string]map[string]any) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("%v is not a struct", v)
	}

	typ := v.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := v.Field(i)
		tag := field.Tag.Get("config")

		if tag == "" {
			continue
		}

		if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Struct {
			if err := setConfigValues(fieldValue, tag, values); err != nil {
				return fmt.Errorf("set config values %s: %w", field.Name, err)
			}
			continue
		}

		if section == "" {
			return fmt.Errorf("section empty")
		}

		val, exists := values[section][tag]
		if !exists {
			continue
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field %s", tag)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(val.(string))
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(val.(string))
			if err == nil {
				fieldValue.SetBool(boolVal)
			}
		default:
			return fmt.Errorf("field %s has unsupported type %s", field.Name, fieldValue.Kind())
		}
	}

	return nil
}
