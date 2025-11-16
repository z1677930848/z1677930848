// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
)

// UpgradeStatus 升级状态
type UpgradeStatus string

const (
	StatusPending    UpgradeStatus = "pending"
	StatusDownloading UpgradeStatus = "downloading"
	StatusVerifying  UpgradeStatus = "verifying"
	StatusInstalling UpgradeStatus = "installing"
	StatusSuccess    UpgradeStatus = "success"
	StatusFailed     UpgradeStatus = "failed"
	StatusRollback   UpgradeStatus = "rollback"
	StatusCancelled  UpgradeStatus = "cancelled"
)

// UpgradeLog 升级日志
type UpgradeLog struct {
	ID            string        `json:"id"`            // 日志ID
	Component     string        `json:"component"`     // 组件名称: admin/api/node
	NodeID        int64         `json:"nodeId"`        // 节点ID（远程升级时使用）
	OldVersion    string        `json:"oldVersion"`    // 旧版本
	NewVersion    string        `json:"newVersion"`    // 新版本
	Status        UpgradeStatus `json:"status"`        // 状态
	StartTime     time.Time     `json:"startTime"`     // 开始时间
	EndTime       time.Time     `json:"endTime"`       // 结束时间
	Duration      int64         `json:"duration"`      // 持续时间（秒）
	DownloadSpeed float64       `json:"downloadSpeed"` // 下载速度（MB/s）
	DownloadSize  int64         `json:"downloadSize"`  // 下载大小（bytes）
	ErrorCode     int           `json:"errorCode"`     // 错误码
	ErrorMessage  string        `json:"errorMessage"`  // 错误信息
	ErrorStage    string        `json:"errorStage"`    // 错误阶段
	RetryCount    int           `json:"retryCount"`    // 重试次数
	BackupPath    string        `json:"backupPath"`    // 备份路径
	DownloadURL   string        `json:"downloadUrl"`   // 下载地址
	Metadata      string        `json:"metadata"`      // 元数据（JSON格式）
}

// UpgradeLogManager 升级日志管理器
type UpgradeLogManager struct {
	logFile string
	locker  sync.RWMutex
}

var sharedUpgradeLogManager *UpgradeLogManager
var logManagerOnce sync.Once

// SharedUpgradeLogManager 获取单例
func SharedUpgradeLogManager() *UpgradeLogManager {
	logManagerOnce.Do(func() {
		sharedUpgradeLogManager = &UpgradeLogManager{
			logFile: Tea.ConfigFile("upgrade_logs.json"),
		}
	})
	return sharedUpgradeLogManager
}

// CreateLog 创建升级日志
func (m *UpgradeLogManager) CreateLog(log *UpgradeLog) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	if log.ID == "" {
		log.ID = m.generateID()
	}
	if log.StartTime.IsZero() {
		log.StartTime = time.Now()
	}

	logs.Println("[UPGRADE_LOG]create log:", log.ID, "component:", log.Component, "version:", log.OldVersion, "->", log.NewVersion)

	return m.saveLog(log)
}

// UpdateLog 更新升级日志
func (m *UpgradeLogManager) UpdateLog(log *UpgradeLog) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	if !log.EndTime.IsZero() && !log.StartTime.IsZero() {
		log.Duration = int64(log.EndTime.Sub(log.StartTime).Seconds())
	}

	logs.Println("[UPGRADE_LOG]update log:", log.ID, "status:", log.Status)

	return m.saveLog(log)
}

// GetLog 获取升级日志
func (m *UpgradeLogManager) GetLog(id string) (*UpgradeLog, error) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	allLogs, err := m.loadAllLogs()
	if err != nil {
		return nil, err
	}

	for _, log := range allLogs {
		if log.ID == id {
			return &log, nil
		}
	}

	return nil, nil
}

// GetLatestLog 获取最新的升级日志
func (m *UpgradeLogManager) GetLatestLog(component string) (*UpgradeLog, error) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	allLogs, err := m.loadAllLogs()
	if err != nil {
		return nil, err
	}

	var latestLog *UpgradeLog
	for i := len(allLogs) - 1; i >= 0; i-- {
		if allLogs[i].Component == component {
			latestLog = &allLogs[i]
			break
		}
	}

	return latestLog, nil
}

// GetLogs 获取升级日志列表
func (m *UpgradeLogManager) GetLogs(component string, limit int) ([]UpgradeLog, error) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	allLogs, err := m.loadAllLogs()
	if err != nil {
		return nil, err
	}

	var result []UpgradeLog
	for i := len(allLogs) - 1; i >= 0 && len(result) < limit; i-- {
		if component == "" || allLogs[i].Component == component {
			result = append(result, allLogs[i])
		}
	}

	return result, nil
}

// CleanOldLogs 清理旧日志（保留最近N条）
func (m *UpgradeLogManager) CleanOldLogs(keepCount int) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	allLogs, err := m.loadAllLogs()
	if err != nil {
		return err
	}

	if len(allLogs) <= keepCount {
		return nil
	}

	// 保留最新的N条
	allLogs = allLogs[len(allLogs)-keepCount:]

	return m.saveAllLogs(allLogs)
}

// saveLog 保存单条日志
func (m *UpgradeLogManager) saveLog(log *UpgradeLog) error {
	allLogs, err := m.loadAllLogs()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 查找是否已存在
	found := false
	for i, existingLog := range allLogs {
		if existingLog.ID == log.ID {
			allLogs[i] = *log
			found = true
			break
		}
	}

	// 如果不存在，添加到末尾
	if !found {
		allLogs = append(allLogs, *log)
	}

	return m.saveAllLogs(allLogs)
}

// loadAllLogs 加载所有日志
func (m *UpgradeLogManager) loadAllLogs() ([]UpgradeLog, error) {
	data, err := os.ReadFile(m.logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []UpgradeLog{}, nil
		}
		return nil, err
	}

	var logs []UpgradeLog
	err = json.Unmarshal(data, &logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// saveAllLogs 保存所有日志
func (m *UpgradeLogManager) saveAllLogs(logs []UpgradeLog) error {
	// 确保目录存在
	dir := filepath.Dir(m.logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.logFile, data, 0644)
}

// generateID 生成日志ID
func (m *UpgradeLogManager) generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		time.Sleep(time.Nanosecond)
	}
	return string(result)
}
