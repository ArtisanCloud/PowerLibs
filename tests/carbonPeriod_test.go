package tests

import (
	"fmt"
	. "github.com/ArtisanCloud/PowerLibs/v2/datetime/carbon"
	"github.com/golang-module/carbon"
	"testing"
)

var CurrentRange *CarbonPeriod

const duration = 3

func init() {

	// current range (today , today+3)
	startDate := carbon.Now()
	endDate := carbon.Now().AddDays(duration)
	CurrentRange = CreateCarbonPeriod().
		SetStartDate(&startDate, nil).
		SetEndDate(&endDate, nil)

}

func Test_Overlaps_RangeOverlapsTestRange(t *testing.T) {
	// test range (today+1 , today+2)
	fmt.Println("test range (today+1 , today+2) RangeOverlapsTestRange")
	testStartDate := carbon.Now().AddDay()
	testEndDate := carbon.Now().AddDays(2)
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
	// test range (today-1 , today_2)
	fmt.Println("test range (today-1 , today+2) RangeLeftOverlapsTestRange")
	testStartDate := carbon.Now().AddDays(-1)
	testEndDate := carbon.Now().AddDays(-1 + duration)
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
	// test range (today+1 , today+1+3)
	fmt.Println("test range (today+1 , today+1+3) RangeRightOverlapsTestRange")
	testStartDate := carbon.Now().AddDays(1)
	testEndDate := carbon.Now().AddDays(duration + 1)
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
	// test range (today-1 , today)
	fmt.Println("test range (today-1 , today) RangeIsLeftSideOfTestRange")
	testStartDate := carbon.Now().AddDays(-1)
	testEndDate := carbon.Now()
	testRange := CreateCarbonPeriod().
		SetStartDate(&testStartDate, nil).
		SetEndDate(&testEndDate, nil)

	//helper.Dump(CurrentRange)
	//helper.Dump(testRange)

	if CurrentRange.Overlaps(testRange) {
		t.Error("current range left side test range")
		t.Error("current range does not overlaps test range")
	}
}

func Test_Overlaps_RangeIsRightSideOfTestRange(t *testing.T) {
	// test range (today+3 , today+3+1)
	fmt.Println("test range (today+3 , today+3+1) RangeIsRightSideOfTestRange")
	testStartDate := carbon.Now().AddDays(duration)
	testEndDate := carbon.Now().AddDays(1 + duration)
	testRange := CreateCarbonPeriod().
		SetStartDate(&testStartDate, nil).
		SetEndDate(&testEndDate, nil)

	//helper.Dump(CurrentRange)
	//helper.Dump(testRange)

	if CurrentRange.Overlaps(testRange) {
		t.Error("current range right side test range")
		t.Error("current range does not overlaps test range")
	}
}
