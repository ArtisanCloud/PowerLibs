package tests

import (
	. "github.com/ArtisanCloud/go-libs/carbon"
	"github.com/golang-module/carbon"
	"testing"
)

var CurrentRange *CarbonPeriod

const duration = 3

func init() {
	startDate := carbon.Now()
	endDate := startDate.AddDays(duration)
	CurrentRange = CreateCarbonPeriod().
		SetStartDate(&startDate, nil).
		SetEndDate(&endDate, nil)

}

func Test_Overlaps_RangeOverlapsTestRange(t *testing.T) {
	now := carbon.Now()
	testStartDate := now.AddDay()
	testEndDate := testStartDate.AddDay()
	testRange := CreateCarbonPeriod().
		SetStartDate(&testStartDate, nil).
		SetEndDate(&testEndDate, nil)

	//helper.Dump(CurrentRange)
	//helper.Dump(testRange)

	if !CurrentRange.Overlaps(testRange) {
		t.Error("current range contains test range")
		t.Error("current range does not overlaps test range")
	}

}

func Test_Overlaps_RangeLeftOverlapsTestRange(t *testing.T) {
	now := carbon.Now()
	testStartDate := now.AddDays(-duration)
	testEndDate := testStartDate.AddDay()
	testRange := CreateCarbonPeriod().
		SetStartDate(&testStartDate, nil).
		SetEndDate(&testEndDate, nil)

	//helper.Dump(CurrentRange)
	//helper.Dump(testRange)

	if !CurrentRange.Overlaps(testRange) {
		t.Error("current range left intersect test range")
		t.Error("current range does not overlaps test range")
	}

}

func Test_Overlaps_RangeRightOverlapsTestRange(t *testing.T) {
	now := carbon.Now()
	testStartDate := now.AddDays(1)
	testEndDate := testStartDate.AddDays(duration + 1)
	testRange := CreateCarbonPeriod().
		SetStartDate(&testStartDate, nil).
		SetEndDate(&testEndDate, nil)

	//helper.Dump(CurrentRange)
	//helper.Dump(testRange)

	if !CurrentRange.Overlaps(testRange) {
		t.Error("current range right intersect test range")
		t.Error("current range does not overlaps test range")
	}

}

func Test_Overlaps_RangeIsLeftSideOfTestRange(t *testing.T) {
	now := carbon.Now()
	testStartDate := now.AddDays(-1)
	testEndDate := testStartDate.AddDays(1)
	testRange := CreateCarbonPeriod().
		SetStartDate(&testStartDate, nil).
		SetEndDate(&testEndDate, nil)

	//helper.Dump(CurrentRange)
	//helper.Dump(testRange)

	if !CurrentRange.Overlaps(testRange) {
		t.Error("current range left side test range")
		t.Error("current range does not overlaps test range")
	}
}

func Test_Overlaps_RangeIsRightSideOfTestRange(t *testing.T) {
	now := carbon.Now()
	testStartDate := now.AddDays(duration)
	testEndDate := testStartDate.AddDays(1)
	testRange := CreateCarbonPeriod().
		SetStartDate(&testStartDate, nil).
		SetEndDate(&testEndDate, nil)

	//helper.Dump(CurrentRange)
	//helper.Dump(testRange)

	if !CurrentRange.Overlaps(testRange) {
		t.Error("current range right side test range")
		t.Error("current range does not overlaps test range")
	}
}
