package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/ndk123-web/trak/internal/models"
	"github.com/ndk123-web/trak/internal/ui"
	"github.com/spf13/cobra"
)

var doneCmd = cobra.Command{
	Use:   "done",
	Short: "To Mark the Module as true",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]

		if moduleName == "" {
			ui.Error(fmt.Sprintf("Error: %v", errors.New("Module Name Can't Be Empty")))
			return
		}

		// if all right then fetch all trak.json from current directory
		// 1. first check if there is trak.json
		// 2. if there trak.json then fetch module_breakdown
		// 3. if valid then update the module_breakdown , if its invalid then error

		cwd, err := os.Getwd()
		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		trakPath := path.Join(cwd, "trak.json")
		_, err = os.Stat(trakPath)

		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		dataBytes, err := os.ReadFile(trakPath)
		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		var trakStruct models.WorkspaceMetadata
		err = json.NewDecoder(bytes.NewReader(dataBytes)).Decode(&trakStruct)
		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		moduleBreakdown := trakStruct.ModuleBreakdown

		_, exist := moduleBreakdown[moduleName]
		if !exist {
			ui.Error(fmt.Sprintf("%v", "Module name doesn't Exist"))
			return
		}

		// now it exist then
		moduleBreakdown[moduleName] = true

		trakStruct.ModuleBreakdown = moduleBreakdown

		newDataBytes, err := json.Marshal(trakStruct)
		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		err = os.WriteFile(trakPath, newDataBytes, 0755)
		if err != nil {
			ui.Error(fmt.Sprintf("Error: %v", err.Error()))
			return
		}

		ui.Success("Mark it Sir..")
	},
}

func init() {
	rootCmd.AddCommand(&doneCmd)
}
