//go:build windows

package dbutils

import "github.com/iwind/TeaGo/types"

var HasFreeSpace = true
var IsLocalDatabase = false
var LocalDatabaseDataDir = ""

// StatDiskUsage is a Windows stub returning zero values.
func StatDiskUsage(path string) (totalBytes int64, usedBytes int64, freeBytes int64, err error) {
	return 0, 0, 0, nil
}

// StatDiskPercent is a Windows stub returning 0.
func StatDiskPercent(path string) (percent float64, err error) {
	return types.Float64(0), nil
}

// CheckHasFreeSpace stub for windows.
func CheckHasFreeSpace() bool { return true }
