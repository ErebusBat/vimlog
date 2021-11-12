package main

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVimLog(t *testing.T) {
	// gomega.RegisterFailHandler(ginkgo.Fail)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vimlog Suite")
}

func dateStr(date time.Time) string {
	return date.Format("2006-01-02")
}

// func assertPreviousWorkingDay(t *testing.T, today, expected time.Time) {
// 	actual := getPreviousWorkingDay(today)
// 	if actual != expected {
// 		t.Fatalf(
// 			"%v[%s] - Should return %v [%s]!  Got %v[%s]",
// 			today.Weekday(), dateStr(today),
// 			expected.Weekday(), dateStr(expected),
// 			actual.Weekday(), dateStr(actual))
// 	} else {
// 		t.Logf(
// 			"%v[%s] Returns %v[%s]",
// 			today.Weekday(), dateStr(today),
// 			expected.Weekday(), dateStr(expected),
// 		)
// 	}
// }

// func TestGetPreviousWorkingDay(t *testing.T) {
// 	p_fri := time.Date(2021, 10, 29, 0, 0, 0, 0, time.Local)
// 	p_sat := time.Date(2021, 10, 30, 0, 0, 0, 0, time.Local)
// 	p_sun := time.Date(2021, 10, 31, 0, 0, 0, 0, time.Local)
// 	t_mon := time.Date(2021, 11, 01, 0, 0, 0, 0, time.Local)
// 	t_tue := time.Date(2021, 11, 02, 0, 0, 0, 0, time.Local)
// 	t_wed := time.Date(2021, 11, 03, 0, 0, 0, 0, time.Local)
// 	t_thu := time.Date(2021, 11, 04, 0, 0, 0, 0, time.Local)
// 	t_fri := time.Date(2021, 11, 05, 0, 0, 0, 0, time.Local)
// 	t_sat := time.Date(2021, 11, 06, 0, 0, 0, 0, time.Local)

// 	// Sunday - Should return preceeding Friday
// 	{
// 		today := p_sun
// 		expected := p_fri
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// 	// Saturday - Should return preceeding Friday
// 	{
// 		today := p_sat
// 		expected := p_fri
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// 	// Monday - Should return preceeding Friday
// 	{
// 		today := t_mon
// 		expected := p_fri
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// 	// Tuesday - Should return Monday
// 	{
// 		today := t_tue
// 		expected := t_mon
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// 	// Wednesday - Should return Tuesday
// 	{
// 		today := t_wed
// 		expected := t_tue
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// 	// Thursday - Should return Wednesday
// 	{
// 		today := t_thu
// 		expected := t_wed
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// 	// Friday - Should return Thursday
// 	{
// 		today := t_fri
// 		expected := t_thu
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// 	// Saturday - Should return Friday
// 	{
// 		today := t_sat
// 		expected := t_fri
// 		assertPreviousWorkingDay(t, today, expected)
// 	}

// }

// func dateOffsetsToPaths(today time.Time, days []string) (outputPaths []string) {
// func TestDateOffsetsToPaths(t *testing.T) {
// 	// x := gomega.NewGomegaWithT(t)
// 	x := NewGomegaWithT(t)
// 	today := time.Date(2021, 11, 12, 0, 0, 0, 0, time.Local)
// 	expected := []string{
// 		"tguest/logs/2021-11-09.md",
// 		"tguest/logs/2021-11-10.md",
// 		"tguest/logs/2021-11-11.md",
// 		"tguest/logs/2021-11-11.md",
// 	}
// 	days := []string{
// 		"-3",
// 		"-2",
// 		"-1",
// 		"y",
// 	}

// 	actual := dateOffsetsToPaths(today, days)
// 	x.Expect(actual).To(Equal(expected))
// 	// if actual != expected {
// 	// 	t.Fatalf("Didn't get expected return; got: %+v", actual)
// 	// }
// }

var _ = Describe("dateOffsetsToPaths", func() {
	// today := time.Date(2021, 11, 12, 0, 0, 0, 0, time.Local)
	today := time.Date(2021, 8, 12, 0, 0, 0, 0, time.Local)
	expected := []string{
		"tguest/logs/2021-11-09.md",
		"tguest/logs/2021-11-10.md",
		"tguest/logs/2021-11-11.md",
		"tguest/logs/2021-11-11.md",
	}
	days := []string{
		"-3",
		"-2",
		"-1",
		"y",
	}

	actual := dateOffsetsToPaths(today, days)
	Expect(actual).To(Equal(expected))
})
