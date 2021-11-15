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
)

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
	parentPath := "tguest/logs"

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
	vimPath, err := exec.LookPath("nvim")
	if err != nil {
		panic("Could not find NeoVim in path!")
	}

	args := []string{vimPath}
	args = append(args, paths...)
	log.Printf("%+v", args)

	syscall.Exec(vimPath, args, os.Environ())
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

		runVim(outputPaths)
	},
	DisableFlagParsing: true,
}

func main() {
	log.Println("vimlog - starting")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
