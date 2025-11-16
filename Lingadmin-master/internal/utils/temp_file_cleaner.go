// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/iwind/TeaGo/logs"
)

// TempFileInfo 临时文件信息
type TempFileInfo struct {
	Path      string
	IsDir     bool
	CreatedAt time.Time
	KeepUntil time.Time // 保留到什么时候（0表示立即删除）
}

// TempFileCleaner 临时文件清理器
type TempFileCleaner struct {
	files  []TempFileInfo
	locker sync.Mutex
}

// NewTempFileCleaner 创建临时文件清理器
func NewTempFileCleaner() *TempFileCleaner {
	return &TempFileCleaner{
		files: []TempFileInfo{},
	}
}

// AddFile 添加要清理的文件
func (c *TempFileCleaner) AddFile(path string) {
	c.AddFileWithDelay(path, 0)
}

// AddFileWithDelay 添加要清理的文件（延迟删除）
func (c *TempFileCleaner) AddFileWithDelay(path string, delay time.Duration) {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.files = append(c.files, TempFileInfo{
		Path:      path,
		IsDir:     false,
		CreatedAt: time.Now(),
		KeepUntil: time.Now().Add(delay),
	})

	logs.Println("[TEMP_CLEANER]registered file for cleanup:", path)
}

// AddDir 添加要清理的目录
func (c *TempFileCleaner) AddDir(path string) {
	c.AddDirWithDelay(path, 0)
}

// AddDirWithDelay 添加要清理的目录（延迟删除）
func (c *TempFileCleaner) AddDirWithDelay(path string, delay time.Duration) {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.files = append(c.files, TempFileInfo{
		Path:      path,
		IsDir:     true,
		CreatedAt: time.Now(),
		KeepUntil: time.Now().Add(delay),
	})

	logs.Println("[TEMP_CLEANER]registered directory for cleanup:", path)
}

// CleanupAll 清理所有临时文件
func (c *TempFileCleaner) CleanupAll() error {
	c.locker.Lock()
	defer c.locker.Unlock()

	var lastErr error
	now := time.Now()

	for _, info := range c.files {
		// 检查是否到清理时间
		if !info.KeepUntil.IsZero() && now.Before(info.KeepUntil) {
			logs.Println("[TEMP_CLEANER]skipping (not time yet):", info.Path)
			continue
		}

		// 检查文件/目录是否存在
		if _, err := os.Stat(info.Path); os.IsNotExist(err) {
			logs.Println("[TEMP_CLEANER]already removed:", info.Path)
			continue
		}

		// 删除文件或目录
		var err error
		if info.IsDir {
			err = os.RemoveAll(info.Path)
		} else {
			err = os.Remove(info.Path)
		}

		if err != nil {
			logs.Println("[TEMP_CLEANER]failed to remove:", info.Path, "error:", err)
			lastErr = err
		} else {
			logs.Println("[TEMP_CLEANER]successfully removed:", info.Path)
		}
	}

	// 清空列表
	c.files = []TempFileInfo{}

	return lastErr
}

// CleanupReady 清理已到期的临时文件
func (c *TempFileCleaner) CleanupReady() error {
	c.locker.Lock()
	defer c.locker.Unlock()

	var lastErr error
	now := time.Now()
	remaining := []TempFileInfo{}

	for _, info := range c.files {
		// 检查是否到清理时间
		if !info.KeepUntil.IsZero() && now.Before(info.KeepUntil) {
			remaining = append(remaining, info)
			continue
		}

		// 检查文件/目录是否存在
		if _, err := os.Stat(info.Path); os.IsNotExist(err) {
			logs.Println("[TEMP_CLEANER]already removed:", info.Path)
			continue
		}

		// 删除文件或目录
		var err error
		if info.IsDir {
			err = os.RemoveAll(info.Path)
		} else {
			err = os.Remove(info.Path)
		}

		if err != nil {
			logs.Println("[TEMP_CLEANER]failed to remove:", info.Path, "error:", err)
			lastErr = err
			// 删除失败的保留在列表中
			remaining = append(remaining, info)
		} else {
			logs.Println("[TEMP_CLEANER]successfully removed:", info.Path)
		}
	}

	c.files = remaining
	return lastErr
}

// Count 返回待清理的文件数量
func (c *TempFileCleaner) Count() int {
	c.locker.Lock()
	defer c.locker.Unlock()
	return len(c.files)
}

// Clear 清空清理列表（不执行删除）
func (c *TempFileCleaner) Clear() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.files = []TempFileInfo{}
}

// GlobalTempFileCleaner 全局临时文件清理器
var globalTempFileCleaner *TempFileCleaner
var globalCleanerOnce sync.Once

// SharedTempFileCleaner 获取全局临时文件清理器
func SharedTempFileCleaner() *TempFileCleaner {
	globalCleanerOnce.Do(func() {
		globalTempFileCleaner = NewTempFileCleaner()
	})
	return globalTempFileCleaner
}

// CleanupOldUpgradeFiles 清理旧的升级临时文件
func CleanupOldUpgradeFiles(olderThan time.Duration) error {
	tmpDir := os.TempDir()

	patterns := []string{
		filepath.Join(tmpDir, "edge-*-tmp"),
		filepath.Join(tmpDir, "ling-*.tar.gz"),
		filepath.Join(tmpDir, "ling-*.zip"),
		filepath.Join(tmpDir, "edge-*.tar.gz"),
		filepath.Join(tmpDir, "edge-*.zip"),
	}

	var lastErr error
	threshold := time.Now().Add(-olderThan)

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			logs.Println("[CLEANUP]glob pattern failed:", pattern, err)
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			// 只删除超过指定时间的文件
			if info.ModTime().Before(threshold) {
				var removeErr error
				if info.IsDir() {
					removeErr = os.RemoveAll(match)
				} else {
					removeErr = os.Remove(match)
				}

				if removeErr != nil {
					logs.Println("[CLEANUP]failed to remove old file:", match, removeErr)
					lastErr = removeErr
				} else {
					logs.Println("[CLEANUP]removed old file:", match)
				}
			}
		}
	}

	return lastErr
}

// ScheduleCleanupTask 启动定期清理任务
func ScheduleCleanupTask() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		// 启动时立即清理一次
		if err := CleanupOldUpgradeFiles(7 * 24 * time.Hour); err != nil {
			logs.Println("[CLEANUP_TASK]cleanup failed:", err)
		}

		for range ticker.C {
			// 清理7天前的临时文件
			if err := CleanupOldUpgradeFiles(7 * 24 * time.Hour); err != nil {
				logs.Println("[CLEANUP_TASK]cleanup failed:", err)
			}

			// 清理全局清理器中已到期的文件
			if err := SharedTempFileCleaner().CleanupReady(); err != nil {
				logs.Println("[CLEANUP_TASK]global cleaner failed:", err)
			}
		}
	}()

	logs.Println("[CLEANUP_TASK]scheduled cleanup task started")
}
