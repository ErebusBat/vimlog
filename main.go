package main

import (
	"fmt"
	"io"
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

////////////////////////////////////////////////////////////////////////////////
// Domain Specific Functionality
////////////////////////////////////////////////////////////////////////////////

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
	parentPath := viper.GetString("DateBasePath")

	for _, request := range days {
		if request == "--" {
			// We will get this because we tell cobra to not parse flags
			continue
		}

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

func runEditor(paths []string) {
	var editorPath string

	// Get editor path:
	//  1. Use VIMLOG_EDITOR
	//  2. Use $EDITOR
	//  3. Look for nvim in path
	if viper.IsSet("Editor") {
		editorPath = viper.GetString("Editor")
	} else if envPath := os.ExpandEnv("EDITOR"); len(envPath) > 0 {
		editorPath = envPath
	} else if neditorPath, err := exec.LookPath("nvim"); err != nil {
		editorPath = neditorPath
	} else {
		ensureOutput()
		log.Fatal("Could not find a valid editor!")
	}

	// options
	editor_options := viper.GetStringSlice("editor_options")

	// Build entire command
	args := []string{editorPath}
	args = append(args, editor_options...)
	args = append(args, paths...)
	log.Printf("Executing: %v", args)

	if viper.GetBool("NoEdit") {
		log.Fatal("VIMLOG_NOEDIT=1 exiting")
	}

	// This syscall should replace the process, assuming it succeeds
	err := syscall.Exec(editorPath, args, os.Environ())
	if err != nil {
		log.Printf("FATAL: spawning editor failed, make sure that %s exists and is executable", editorPath)
		log.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// Options + Logging
////////////////////////////////////////////////////////////////////////////////

func loadOptions() {
	//
	// Defaults
	//
	viper.SetDefault("DateBasePath", "")

	//
	// ENV overloading
	//
	viper.SetEnvPrefix("VIMLOG")
	viper.AutomaticEnv()

	//
	// Config file
	//
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

	//
	// Logging (default no logging)
	//
	logFlags := 0
	if viper.IsSet("LogFlags") {
		// We don't want LogFlags in the example config so set it this way
		logFlags = viper.GetInt("LogFlags")
	} else if viper.IsSet("Silent") && viper.GetBool("Silent") {
		// VIMLOG_LOGFLAGS will override VIMLOG_SILENT
		// But setting VIMLOG_LOGFLAGS=0 != VIMLOG_SILENT=1
		// as VIMLOG_SILENT=1 will not allow output to be enabled by certain
		// commands (like config print)
		logFlags = 0
	}

	log.SetFlags(logFlags)

	if logFlags == 0 {
		// LogFlags only control prefix (timestamp, etc)
		// if LogFlags == 0 then also make sure that we don't output message
		log.SetOutput(io.Discard)
	}

	//
	// Debug / output
	//
	if viper.GetBool("Debug") {
		log.Printf("Config: %v", viper.AllSettings())
	}
}

// Called by functions that output text by the nature of their design
// (i.e. config print).  This will ensure that logging output is
// actually enabled (default is disabled) for these commands
// Can be disabled via VIMLOG_SILENT=1 if you REALLY don't want output
func ensureOutput() {
	if log.Flags() > 0 || viper.GetBool("Silent") {
		return
	}

	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stderr)
}

////////////////////////////////////////////////////////////////////////////////
// Commands
////////////////////////////////////////////////////////////////////////////////

var configCmd = &cobra.Command{
	Use: "config",
	Run: func(_ *cobra.Command, args []string) {
		log.Printf("config %v", args)
	},
}

var configPrintCmd = &cobra.Command{
	Use: "print",
	Run: func(_ *cobra.Command, _ []string) {
		ensureOutput()
		log.Printf("Config file in use: %s", viper.ConfigFileUsed())
		log.Printf("Config: %v", viper.AllSettings())
	},
}

var configWriteCmd = &cobra.Command{
	Use: "write",
	Run: func(_ *cobra.Command, _ []string) {
		if viper.ConfigFileUsed() == "" {
			viper.SetConfigFile("./.vimlog.yaml")
			log.Printf("No config file found...")
		}
		viper.WriteConfig()
		log.Printf("Wrote config %s", viper.ConfigFileUsed())
	},
}

var rootCmd = &cobra.Command{
	Use: "vimlog",
	Run: func(_ *cobra.Command, args []string) {
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

		runEditor(outputPaths)
	},
	DisableFlagParsing: true,
}

////////////////////////////////////////////////////////////////////////////////
// Entry Point
////////////////////////////////////////////////////////////////////////////////

func main() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configPrintCmd)
	configCmd.AddCommand(configWriteCmd)

	// load options and THEN output log line (in case options suppress logging)
	loadOptions()
	log.Println("vimlog - starting")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
