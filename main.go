package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// type Options struct {
//   DateParentPath
// }

// func parseConfig(path string) (config *photorg.Options, err error) {
// 	config = new(photorg.Options)
// 	log.Printf("Loading config %s\n", path)
// 	if _, err = toml.DecodeFile(path, &config); err != nil {
// 		return nil, err
// 	}
// 	log.Printf("  SourcePath = %s\n", config.SourcePath)
// 	log.Printf("    DestRoot = %s\n", config.DestRoot)
// 	return config, nil
// }

func getPreviousWorkingDay(today time.Time) time.Time {
	baseDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	for i := -1; i > -14; i-- {
		adjustedDate := baseDate.AddDate(0, 0, i)
		dow := adjustedDate.Weekday()
		if dow >= time.Monday && dow <= time.Friday {
			return adjustedDate
		}
	}
	return today
}

func dateOffsetsToPaths(today time.Time, days []string) (outputPaths []string) {
	// TODO: This shouldn't be hardcoded
	parentPath := viper.GetString("DateBasePath")

	for _, request := range days {
		offset, err := strconv.Atoi(request)
		var offDate time.Time

		if err == nil {
			offDate = today.AddDate(0, 0, offset)
		} else {
			if request == "y" {
				offDate = getPreviousWorkingDay(today)
			} else {
				outputPaths = append(outputPaths, request)
				continue
			}
		}

		fileName := offDate.Format("2006-01-02")
		fileName += ".md"
		path := filepath.Join(parentPath, fileName)

		outputPaths = append(outputPaths, path)
	}
	return
}

func runVim(paths []string) {
	vimPath := viper.GetString("Vim")

	args := []string{vimPath}
	args = append(args, paths...)
	log.Printf("Executing: %v", args)

	// This syscall should replace the process, assuming it succeeds
	err := syscall.Exec(vimPath, args, os.Environ())
	if err != nil {
		log.Printf("FATAL: spawning editor failed, make sure that %s exists and is executable", vimPath)
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use: "vimlog",
	Run: func(cmd *cobra.Command, args []string) {
		today := time.Now()
		outputPaths := dateOffsetsToPaths(today, args)

		if len(outputPaths) <= 0 {
			log.Println("No args, opening todays log")
			outputPaths = dateOffsetsToPaths(today, []string{"0"})
		}

		if viper.GetBool("Debug") {
			log.Printf("Debug mode enabled, not running vim")
			log.Printf("%v", outputPaths)
			os.Exit(0)
		}

		runVim(outputPaths)
	},
	DisableFlagParsing: true,
}

func loadOptions() {
	// Defaults
	viper.SetDefault("DateBasePath", "")
	viper.SetDefault("Debug", false)
	if vimPath, err := exec.LookPath("nvim"); err == nil {
		viper.SetDefault("Vim", vimPath)
	}

	// Setup ENV overloading
	viper.SetEnvPrefix("VIMLOG")
	viper.AutomaticEnv()

	// Setup config file
	viper.SetConfigName(".vimlog")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// No config file... that is okay, we will use defaults
		} else {
			log.Fatalf("Config file error %v", err)
		}
	}

	// Debug / output
	if viper.GetBool("Debug") {
		log.Printf("Config: %v", viper.AllSettings())

		// TODO: Maybe add a cobra flag to write this?
		// if viper.ConfigFileUsed() == "" {
		// 	viper.SetConfigFile("./.vimlog.yaml")
		// }
		// viper.WriteConfig()
		// log.Printf("Wrote config %s", viper.ConfigFileUsed())
	}
}

func main() {
	log.Println("vimlog - starting")

	loadOptions()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
