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
	// precision, err := time.ParseDuration("24h")
	// if err != nil {
	// 	panic("Get your code right!")
	// }
	// today = today.Truncate(precision)
	baseDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	// baseDateString := baseDate.Format("2006-01-02")

	for i := -1; i > -14; i-- {
		adjustedDate := baseDate.AddDate(0, 0, i)
		dow := adjustedDate.Weekday()
		if dow >= time.Monday && dow <= time.Friday {
			return adjustedDate
		}
		// adjustedString := adjustedDate.Format("2006-01-02")
		// log.Printf("getPreviousWorkingDay[%s %d]: %s %v", baseDateString, i, adjustedString, adjustedDate.Weekday())
	}
	return today
}

func dateOffsetsToPaths(today time.Time, days []string) (outputPaths []string) {
	// days := []string{"0", "-1", "y", "-2"}
	parentPath := "tguest/logs"

	// today := time.Now()
	// today = today.AddDate(0, 0, -4)

	// outputPaths := make([]string, 0)

	// var days []string
	// if len(os.Args) > 1 {
	// 	// Get rid of the exe path from $0
	// 	days = os.Args[1:]
	// } else {
	// 	days = []string{"0", "y"}
	// }

	for _, request := range days {
		offset, err := strconv.Atoi(request)
		var offDate time.Time

		if err == nil {
			offDate = today.AddDate(0, 0, offset)
		} else {
			if request == "y" {
				offDate = getPreviousWorkingDay(today)
			} else {
				// panic(fmt.Sprintf("Unknown offset %s", request))
				outputPaths = append(outputPaths, request)
				continue
			}
		}

		fileName := offDate.Format("2006-01-02")
		fileName += ".md"
		path := filepath.Join(parentPath, fileName)

		// log.Printf("%d = %s %d %v", reqnum, request, offset, offDate.Format("2006-01-02"))
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

		// log.Printf("Got %d paths:", len(outputPaths))
		// for x, path := range outputPaths {
		// 	log.Printf("[%d] %s", x, path)
		// }
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
