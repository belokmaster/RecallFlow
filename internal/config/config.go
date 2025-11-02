package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %v", filename, err)
	}
	defer file.Close()

	config := &Config{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// пустрая строка и комменты пропускаются
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// парсим на ключ и value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "host":
			config.Host = value
		case "port":
			config.Port = value
		case "user":
			config.User = value
		case "password":
			config.Password = value
		case "dbname":
			config.DBName = value
		case "sslmode":
			config.SSLMode = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("problem with reading file: %v", err)
	}

	return config, nil
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
