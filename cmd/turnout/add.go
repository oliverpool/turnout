package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func add(cmd *cobra.Command, args []string) error {
	configDir, err := cmd.Flags().GetString("configdir")
	if err != nil {
		return err
	}

	host := args[0]
	content := args[1]
	path := filepath.Join(configDir, host)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(content), 0644)
}
