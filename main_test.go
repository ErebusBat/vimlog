package main

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

func dateStr(date time.Time) string {
	return date.Format("2006-01-02")
}

func assertPreviousWorkingDay(t *testing.T, today, expected time.Time, annotation string) {
	x := NewGomegaWithT(t)
	actual := getPreviousWorkingDay(today)
	x.Expect(actual).To(Equal(expected), annotation)
}

func TestGetPreviousWorkingDay(t *testing.T) {
	p_fri := time.Date(2021, 10, 29, 0, 0, 0, 0, time.Local)
	p_sat := time.Date(2021, 10, 30, 0, 0, 0, 0, time.Local)
	p_sun := time.Date(2021, 10, 31, 0, 0, 0, 0, time.Local)
	t_mon := time.Date(2021, 11, 01, 0, 0, 0, 0, time.Local)
	t_tue := time.Date(2021, 11, 02, 0, 0, 0, 0, time.Local)
	t_wed := time.Date(2021, 11, 03, 0, 0, 0, 0, time.Local)
	t_thu := time.Date(2021, 11, 04, 0, 0, 0, 0, time.Local)
	t_fri := time.Date(2021, 11, 05, 0, 0, 0, 0, time.Local)
	t_sat := time.Date(2021, 11, 06, 0, 0, 0, 0, time.Local)

	// Sunday - Should return preceeding Friday
	{
		today := p_sun
		expected := p_fri
		assertPreviousWorkingDay(t, today, expected, "Sunday - Should return preceeding Friday")
	}

	// Saturday - Should return preceeding Friday
	{
		today := p_sat
		expected := p_fri
		assertPreviousWorkingDay(t, today, expected, "Saturday - Should return preceeding Friday")
	}

	// Monday - Should return preceeding Friday
	{
		today := t_mon
		expected := p_fri
		assertPreviousWorkingDay(t, today, expected, "Monday - Should return preceeding Friday")
	}

	// Tuesday - Should return Monday
	{
		today := t_tue
		expected := t_mon
		assertPreviousWorkingDay(t, today, expected, "Tuesday - Should return Monday")
	}

	// Wednesday - Should return Tuesday
	{
		today := t_wed
		expected := t_tue
		assertPreviousWorkingDay(t, today, expected, "Wednesday - Should return Tuesday")
	}

	// Thursday - Should return Wednesday
	{
		today := t_thu
		expected := t_wed
		assertPreviousWorkingDay(t, today, expected, "Thursday - Should return Wednesday")
	}

	// Friday - Should return Thursday
	{
		today := t_fri
		expected := t_thu
		assertPreviousWorkingDay(t, today, expected, "Friday - Should return Thursday")
	}

	// Saturday - Should return Friday
	{
		today := t_sat
		expected := t_fri
		assertPreviousWorkingDay(t, today, expected, "Saturday - Should return Friday")
	}

}

// func dateOffsetsToPaths(today time.Time, days []string) (outputPaths []string) {
func TestDateOffsetsToPaths(t *testing.T) {
	x := NewGomegaWithT(t)
	today := time.Date(2021, 11, 15, 0, 0, 0, 0, time.Local)
	expected := []string{
		"test/journal/2021-11-12.md",
		"test/journal/2021-11-12.md",
		"pass_through.md",
		"test/journal/2021-11-13.md",
		"test/journal/2021-11-14.md",
		"test/journal/2021-11-15.md",
	}
	days := []string{
		"-3",
		"y",
		"pass_through.md",
		"--", // Should not propigate
		"-2",
		"-1",
		"0",
	}

	// Set a specific testing path, but also end with a trailing slash to
	// ensure that a proper path.Join is being performed
	viper.Set("DateBasePath", "test/journal/")

	actual := dateOffsetsToPaths(today, days)
	x.Expect(actual).To(Equal(expected))
}

func TestFileNameFormat(t *testing.T) {
	x := NewGomegaWithT(t)
	today := time.Date(2021, 11, 15, 0, 0, 0, 0, time.Local)
	viper.Set("DateBasePath", "test/path/")

	// Start without a specific viper setting for FileNameFormat
	x.Expect(dateOffsetsToPaths(
		today,
		[]string{"0"},
	)).To(Equal([]string{"test/path/2021-11-15.md"}))

	// Explicit default format
	viper.Set("FileNameFormat", "YYYY-MM-DD")
	x.Expect(dateOffsetsToPaths(
		today,
		[]string{"0"},
	)).To(Equal([]string{"test/path/2021-11-15.md"}))

	// Using All Known Formats, in a sub path
	viper.Set("FileNameFormat", "YYYY/MM-MMM/YYYY-MM-DD-ddd")
	x.Expect(dateOffsetsToPaths(
		today,
		[]string{"0"},
	)).To(Equal([]string{"test/path/2021/11-Nov/2021-11-15-Mon.md"}))
}

func TestGetFormattedFileName(t *testing.T) {
	x := NewGomegaWithT(t)
	today := time.Date(2021, 11, 14, 0, 0, 0, 0, time.Local)

	// Default, no viper setting
	viper.Set("FileNameFormat", nil)
	x.Expect(getFormattedFileName(today)).To(Equal("2021-11-14.md"))

	// Explicit default format
	viper.Set("FileNameFormat", "YYYY-MM-DD")
	x.Expect(getFormattedFileName(today)).To(Equal("2021-11-14.md"))

	// All Supported Formats
	viper.Set("FileNameFormat", "YYYY/MM-MMM/YYYY-MM-DD-ddd")
	x.Expect(getFormattedFileName(today)).To(Equal("2021/11-Nov/2021-11-14-Sun.md"))

	// Explicitly stating .md extension in config
	viper.Set("FileNameFormat", "YYYY-MM-DD.md")
	x.Expect(getFormattedFileName(today)).To(Equal("2021-11-14.md"))

	// Embedded .md in config works
	viper.Set("FileNameFormat", "f.md/YYYY-MM-DD")
	x.Expect(getFormattedFileName(today)).To(Equal("f.md/2021-11-14.md"))

	// Explicitly stating a different extension in config
	viper.Set("FileNameFormat", "YYYY-MM-DD.txt")
	x.Expect(getFormattedFileName(today)).To(Equal("2021-11-14.txt.md"))

	// Test with single digit month/day
	today = time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local)
	viper.Set("FileNameFormat", "YYYY-MM-DD")
	x.Expect(getFormattedFileName(today)).To(Equal("2021-01-04.md"))
}
