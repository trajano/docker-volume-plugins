package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVolumeCalculation(t *testing.T) {

	var calculated = AppendBucketOptionsByVolumeName([]string{"mount"}, "mybucket")
	var expected = []string{"mount", "bucket=mybucket"}
	if !reflect.DeepEqual(calculated, expected) {
		fmt.Errorf("%v didn't match expected", calculated)
		t.Fail()
	}
}

func TestVolumeCalculationOneLevel(t *testing.T) {

	var calculated = AppendBucketOptionsByVolumeName([]string{"mount"}, "mybucket/levelone")
	var expected = []string{"mount", "bucket=mybucket", "servicepath=/levelone"}
	if !reflect.DeepEqual(calculated, expected) {
		fmt.Errorf("%v didn't match expected", calculated)
		t.Fail()
	}
}

func TestVolumeCalculationTwoLevels(t *testing.T) {

	var calculated = AppendBucketOptionsByVolumeName([]string{"mount"}, "mybucket/levelone/level2")
	var expected = []string{"mount", "bucket=mybucket", "servicepath=/levelone/level2"}
	if !reflect.DeepEqual(calculated, expected) {
		fmt.Errorf("%v didn't match expected", calculated)
		t.Fail()
	}
}
