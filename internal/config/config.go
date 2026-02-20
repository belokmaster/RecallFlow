package config

import "os"

func GetDatabasePath() string {
	if path := os.Getenv("RECALL_FLOW_DB"); path != "" {
		return path
	}

	return "./recall_flow.db"
}
