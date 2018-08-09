package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVolumeCalculation(t *testing.T) {

	var calculated = AppendVolumeOptionsByVolumeName([]string{"mount"}, "simplevolume")
	var expected = []string{"mount", "--volfile-id=simplevolume"}
	if !reflect.DeepEqual(calculated, expected) {
		fmt.Errorf("%v didn't match expected", calculated)
		t.Fail()
	}
}

func TestVolumeCalculationOneLevel(t *testing.T) {

	var calculated = AppendVolumeOptionsByVolumeName([]string{"mount"}, "simplevolume/levelone")
	var expected = []string{"mount", "--volfile-id=simplevolume", "--subdir-mount=/levelone"}
	if !reflect.DeepEqual(calculated, expected) {
		fmt.Errorf("%v didn't match expected", calculated)
		t.Fail()
	}
}

func TestVolumeCalculationTwoLevels(t *testing.T) {

	var calculated = AppendVolumeOptionsByVolumeName([]string{"mount"}, "simplevolume/levelone/level2")
	var expected = []string{"mount", "--volfile-id=simplevolume", "--subdir-mount=/levelone/level2"}
	if !reflect.DeepEqual(calculated, expected) {
		fmt.Errorf("%v didn't match expected", calculated)
		t.Fail()
	}
}
