/*
 * Copyright © Siemens 2024 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package utils

import (
	"errors"
	"systemservice/internal/common/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initializeDiskOutput() (*mocks.MUtil, StorageDiskTypeOutput) {
	tUtil := new(mocks.MUtil)

	lsblkOutput := StorageDiskTypeOutput{}
	return tUtil, lsblkOutput
}

func Test_DiskUsage(t *testing.T) {
	path := "/"

	diskStatus, _ := DiskUsage(path)
	if diskStatus.All == 0 {
		t.Errorf("Expected disk total %v, but got %d", diskStatus, diskStatus.All)
	}
}

func Test_getDiskType(t *testing.T) {
	notExpectedValue := ""

	testDiskType := StorageDiskTypeOutput{}
	testDiskType.BlockDevices = append(testDiskType.BlockDevices, BlockDevice{Name: "sda", Type: "disk", Rotational: false, Removable: false})

	diskType, diskName, _ := getDiskInfo(testDiskType)
	if diskType == notExpectedValue {
		t.Errorf("Expected disk type %s, but got %s", notExpectedValue, diskType)
	}
	if diskName == notExpectedValue {
		t.Errorf("Expected disk name %s, but got %s", notExpectedValue, diskName)
	}
}

func Test_getDiskSpeed(t *testing.T) {
	_, diskName, _ := diskCommander()
	readSectors, writeSectors, _ := getDiskSpeed(diskName)
	if readSectors == 0 {
		t.Errorf("Expected read sectors to be %d, but got 0", readSectors)
	}

	if writeSectors == 0 {
		t.Errorf("Expected write sectors to be %d, but got 0", writeSectors)
	}
}

func Test_DiskUsage_Error(t *testing.T) {
	path := "/abc"

	diskStatus, _ := DiskUsage(path)
	if diskStatus.All != 0 {
		t.Errorf("Expected disk total %v, but got %d", diskStatus, diskStatus.All)
	}
}

func Test_identifyNonRotationalDiskType_mmcblk(t *testing.T) {
	deviceName := "mmcblk0"
	expected := emmcDisk

	result := identifyNonRotationalDiskType(deviceName)
	if result != expected {
		t.Errorf("For device name %s, expected %s, but got %s", deviceName, expected, result)
	}
}

func Test_identifyNonRotationalDiskType_sd(t *testing.T) {
	deviceName := "sda"
	expected := ssdDisk

	result := identifyNonRotationalDiskType(deviceName)
	if result != expected {
		t.Errorf("For device name %s, expected %s, but got %s", deviceName, expected, result)
	}
}

func Test_identifyNonRotationalDiskType_default(t *testing.T) {
	deviceName := " "
	expected := defaultDisk
	result := identifyNonRotationalDiskType(deviceName)
	if result != expected {
		t.Errorf("For device name %s, expected %s, but got %s", deviceName, expected, result)
	}
}

func Test_getDiskInfo_SSD(t *testing.T) {
	t.Parallel()

	tUtil, _ := initializeDiskOutput()
	tUtil.CommandList = make([]mocks.CmdContainer, 0)
	s1 := mocks.CmdContainer{CommandVal: []byte(`{"blockdevices":[{"name":"sda","type":"disk","rota":false,"rm":false}]}`), CommandErr: nil}
	tUtil.CommandList = append(tUtil.CommandList, s1)

	mockOutput := StorageDiskTypeOutput{}
	mockOutput.BlockDevices = append(mockOutput.BlockDevices, BlockDevice{Name: "sda", Type: "disk", Rotational: false, Removable: false})

	diskType, diskName, err := getDiskInfo(mockOutput)
	assert.Nil(t, err, "Did not get expected result. Wanted: %q, got: %q", nil, err)
	assert.Equal(t, "SSD", diskType, "Did not get expected result. Wanted: %q, got: %q", "SSD", diskType)
	assert.Equal(t, "sda", diskName, "Did not get expected result. Wanted: %q, got: %q", "sda", diskName)
}

func Test_getDiskInfo_HDD(t *testing.T) {
	t.Parallel()

	tUtil, _ := initializeDiskOutput()
	tUtil.CommandList = make([]mocks.CmdContainer, 0)
	s1 := mocks.CmdContainer{CommandVal: []byte(`{"blockdevices":[{"name":"sdb","type":"disk","rota":true,"rm":false}]}`), CommandErr: nil}
	tUtil.CommandList = append(tUtil.CommandList, s1)

	mockOutput := StorageDiskTypeOutput{}
	mockOutput.BlockDevices = append(mockOutput.BlockDevices, BlockDevice{Name: "sdb", Type: "disk", Rotational: true, Removable: false})

	diskType, diskName, err := getDiskInfo(mockOutput)
	assert.Nil(t, err, "Did not get expected result. Wanted: %q, got: %q", nil, err)
	assert.Equal(t, "HDD", diskType, "Did not get expected result. Wanted: %q, got: %q", "HDD", diskType)
	assert.Equal(t, "sdb", diskName, "Did not get expected result. Wanted: %q, got: %q", "sdb", diskName)
}

func Test_getDiskInfo_EMMC(t *testing.T) {
	t.Parallel()

	tUtil, _ := initializeDiskOutput()
	tUtil.CommandList = make([]mocks.CmdContainer, 0)
	s1 := mocks.CmdContainer{CommandVal: []byte(`{"blockdevices":[{"name":"mmcblk0","type":"disk","rota":false,"rm":false}]}`), CommandErr: nil}
	tUtil.CommandList = append(tUtil.CommandList, s1)

	mockOutput := StorageDiskTypeOutput{}
	mockOutput.BlockDevices = append(mockOutput.BlockDevices, BlockDevice{Name: "mmcblk0", Type: "disk", Rotational: false, Removable: false})

	diskType, diskName, err := getDiskInfo(mockOutput)
	assert.Nil(t, err, "Did not get expected result. Wanted: %q, got: %q", nil, err)
	assert.Equal(t, "eMMC", diskType, "Did not get expected result. Wanted: %q, got: %q", "eMMC", diskType)
	assert.Equal(t, "mmcblk0", diskName, "Did not get expected result. Wanted: %q, got: %q", "mmcblk0", diskName)
}

// func Test_getDiskInfo_Error(t *testing.T) {
// 	t.Parallel()
// 	expectedError := errors.New("Utils:getDiskInfo(), No block devices found")

// 	tUtil, _ := initializeDiskOutput()
// 	tUtil.CommandList = make([]mocks.CmdContainer, 0)
// 	s1 := mocks.CmdContainer{CommandVal: nil, CommandErr: expectedError}
// 	tUtil.CommandList = append(tUtil.CommandList, s1)

// 	// []byte(`{"blockdevices":[{"name":"mmcblk0","type":"","rota":false,"rm":false}]}`)

// 	mockOutput := StorageDiskTypeOutput{}
// 	// mockOutput.BlockDevices = append(mockOutput.BlockDevices, BlockDevice{Name: "mmcblk0", Type: "disk", Rotational: false, Removable: false
// 	_, _, err := getDiskInfo(mockOutput)
// 	log.Println("Error: ", err)
// 	assert.Equal(t, err.Error(), expectedError.Error())
// 	// assert.Equal(t, "eMMC", diskType, "Did not get expected result. Wanted: %q, got: %q", "eMMC", diskType)
// 	// assert.Equal(t, "mmcblk0", diskName, "Did not get expected result. Wanted: %q, got: %q", "mmcblk0", diskName)
// }

func Test_getDiskInfo_NoBlockDevices(t *testing.T) {
	t.Parallel()
	expectedError := errors.New("Utils:getDiskInfo(), No block devices found")

	mockOutput := StorageDiskTypeOutput{}
	diskType, diskName, err := getDiskInfo(mockOutput)
	assert.Equal(t, err.Error(), expectedError.Error())
	assert.Equal(t, "", diskType, "Did not get expected result. Wanted: %q, got: %q", "", diskType)
	assert.Equal(t, "", diskName, "Did not get expected result. Wanted: %q, got: %q", "", diskName)
}

func Test_getDiskInfo_NonDiskType(t *testing.T) {
	t.Parallel()

	mockOutput := StorageDiskTypeOutput{}
	mockOutput.BlockDevices = append(mockOutput.BlockDevices, BlockDevice{Name: "sda", Type: "partition", Rotational: false, Removable: false})

	diskType, diskName, err := getDiskInfo(mockOutput)
	assert.Nil(t, err, "Did not get expected result. Wanted: %q, got: %q", nil, err)
	assert.Equal(t, "", diskType, "Did not get expected result. Wanted: %q, got: %q", "", diskType)
	assert.Equal(t, "", diskName, "Did not get expected result. Wanted: %q, got: %q", "", diskName)
}

func Test_getDiskInfo_RemovableDisk(t *testing.T) {
	t.Parallel()

	mockOutput := StorageDiskTypeOutput{}
	mockOutput.BlockDevices = append(mockOutput.BlockDevices, BlockDevice{Name: "sda", Type: "disk", Rotational: false, Removable: true})

	diskType, diskName, err := getDiskInfo(mockOutput)
	assert.Nil(t, err, "Did not get expected result. Wanted: %q, got: %q", nil, err)
	assert.Equal(t, "", diskType, "Did not get expected result. Wanted: %q, got: %q", "", diskType)
	assert.Equal(t, "", diskName, "Did not get expected result. Wanted: %q, got: %q", "", diskName)
}
