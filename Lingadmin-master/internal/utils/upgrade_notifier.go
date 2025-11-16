// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/iwind/TeaGo/logs"
)

// NotifyLevel 通知级别
type NotifyLevel string

const (
	NotifyLevelInfo    NotifyLevel = "info"
	NotifyLevelWarning NotifyLevel = "warning"
	NotifyLevelError   NotifyLevel = "error"
	NotifyLevelSuccess NotifyLevel = "success"
)

// UpgradeNotification 升级通知
type UpgradeNotification struct {
	Level     NotifyLevel   `json:"level"`
	Component string        `json:"component"`
	Version   string        `json:"version"`
	Status    UpgradeStatus `json:"status"`
	Message   string        `json:"message"`
	Progress  float32       `json:"progress"`
	Error     string        `json:"error"`
	Timestamp time.Time     `json:"timestamp"`
}

// UpdateNotifier 更新通知器接口
type UpdateNotifier interface {
	NotifyStart(component, version string)
	NotifyProgress(component string, progress float32, message string)
	NotifySuccess(component, version string, duration time.Duration)
	NotifyFailed(component, version string, err error)
	NotifyCancelled(component, version string)
}

// MultiNotifier 多通道通知器
type MultiNotifier struct {
	notifiers []UpdateNotifier
	locker    sync.RWMutex
}

// NewMultiNotifier 创建多通道通知器
func NewMultiNotifier() *MultiNotifier {
	return &MultiNotifier{
		notifiers: []UpdateNotifier{},
	}
}

// AddNotifier 添加通知器
func (m *MultiNotifier) AddNotifier(notifier UpdateNotifier) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.notifiers = append(m.notifiers, notifier)
}

// NotifyStart 通知开始
func (m *MultiNotifier) NotifyStart(component, version string) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	for _, notifier := range m.notifiers {
		go func(n UpdateNotifier) {
			defer func() {
				if r := recover(); r != nil {
					logs.Println("[NOTIFIER]panic in NotifyStart:", r)
				}
			}()
			n.NotifyStart(component, version)
		}(notifier)
	}
}

// NotifyProgress 通知进度
func (m *MultiNotifier) NotifyProgress(component string, progress float32, message string) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	for _, notifier := range m.notifiers {
		go func(n UpdateNotifier) {
			defer func() {
				if r := recover(); r != nil {
					logs.Println("[NOTIFIER]panic in NotifyProgress:", r)
				}
			}()
			n.NotifyProgress(component, progress, message)
		}(notifier)
	}
}

// NotifySuccess 通知成功
func (m *MultiNotifier) NotifySuccess(component, version string, duration time.Duration) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	for _, notifier := range m.notifiers {
		go func(n UpdateNotifier) {
			defer func() {
				if r := recover(); r != nil {
					logs.Println("[NOTIFIER]panic in NotifySuccess:", r)
				}
			}()
			n.NotifySuccess(component, version, duration)
		}(notifier)
	}
}

// NotifyFailed 通知失败
func (m *MultiNotifier) NotifyFailed(component, version string, err error) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	for _, notifier := range m.notifiers {
		go func(n UpdateNotifier) {
			defer func() {
				if r := recover(); r != nil {
					logs.Println("[NOTIFIER]panic in NotifyFailed:", r)
				}
			}()
			n.NotifyFailed(component, version, err)
		}(notifier)
	}
}

// NotifyCancelled 通知取消
func (m *MultiNotifier) NotifyCancelled(component, version string) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	for _, notifier := range m.notifiers {
		go func(n UpdateNotifier) {
			defer func() {
				if r := recover(); r != nil {
					logs.Println("[NOTIFIER]panic in NotifyCancelled:", r)
				}
			}()
			n.NotifyCancelled(component, version)
		}(notifier)
	}
}

// LogNotifier 日志通知器
type LogNotifier struct{}

// NewLogNotifier 创建日志通知器
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

func (n *LogNotifier) NotifyStart(component, version string) {
	logs.Println(fmt.Sprintf("[UPGRADE]%s: starting upgrade to version %s", component, version))
}

func (n *LogNotifier) NotifyProgress(component string, progress float32, message string) {
	logs.Println(fmt.Sprintf("[UPGRADE]%s: progress %.1f%% - %s", component, progress*100, message))
}

func (n *LogNotifier) NotifySuccess(component, version string, duration time.Duration) {
	logs.Println(fmt.Sprintf("[UPGRADE]%s: successfully upgraded to version %s (took %s)", component, version, duration))
}

func (n *LogNotifier) NotifyFailed(component, version string, err error) {
	logs.Println(fmt.Sprintf("[UPGRADE]%s: failed to upgrade to version %s: %v", component, version, err))
}

func (n *LogNotifier) NotifyCancelled(component, version string) {
	logs.Println(fmt.Sprintf("[UPGRADE]%s: upgrade to version %s was cancelled", component, version))
}

// WebhookNotifier Webhook通知器
type WebhookNotifier struct {
	URL     string
	Timeout time.Duration
	client  *http.Client
}

// NewWebhookNotifier 创建Webhook通知器
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL:     url,
		Timeout: 10 * time.Second,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (n *WebhookNotifier) sendNotification(notification *UpgradeNotification) {
	if n.URL == "" {
		return
	}

	notification.Timestamp = time.Now()

	data, err := json.Marshal(notification)
	if err != nil {
		logs.Println("[WEBHOOK]failed to marshal notification:", err)
		return
	}

	resp, err := n.client.Post(n.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		logs.Println("[WEBHOOK]failed to send notification:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logs.Println("[WEBHOOK]webhook returned error status:", resp.StatusCode)
	}
}

func (n *WebhookNotifier) NotifyStart(component, version string) {
	n.sendNotification(&UpgradeNotification{
		Level:     NotifyLevelInfo,
		Component: component,
		Version:   version,
		Status:    StatusDownloading,
		Message:   fmt.Sprintf("Starting upgrade to version %s", version),
	})
}

func (n *WebhookNotifier) NotifyProgress(component string, progress float32, message string) {
	n.sendNotification(&UpgradeNotification{
		Level:     NotifyLevelInfo,
		Component: component,
		Status:    StatusDownloading,
		Progress:  progress,
		Message:   message,
	})
}

func (n *WebhookNotifier) NotifySuccess(component, version string, duration time.Duration) {
	n.sendNotification(&UpgradeNotification{
		Level:     NotifyLevelSuccess,
		Component: component,
		Version:   version,
		Status:    StatusSuccess,
		Message:   fmt.Sprintf("Successfully upgraded to version %s (took %s)", version, duration),
	})
}

func (n *WebhookNotifier) NotifyFailed(component, version string, err error) {
	n.sendNotification(&UpgradeNotification{
		Level:     NotifyLevelError,
		Component: component,
		Version:   version,
		Status:    StatusFailed,
		Message:   "Upgrade failed",
		Error:     err.Error(),
	})
}

func (n *WebhookNotifier) NotifyCancelled(component, version string) {
	n.sendNotification(&UpgradeNotification{
		Level:     NotifyLevelWarning,
		Component: component,
		Version:   version,
		Status:    StatusCancelled,
		Message:   "Upgrade was cancelled",
	})
}

// ConsoleNotifier 控制台通知器（带彩色输出）
type ConsoleNotifier struct{}

// NewConsoleNotifier 创建控制台通知器
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

func (n *ConsoleNotifier) NotifyStart(component, version string) {
	fmt.Printf("\033[34m[INFO]\033[0m %s: Starting upgrade to version %s\n", component, version)
}

func (n *ConsoleNotifier) NotifyProgress(component string, progress float32, message string) {
	fmt.Printf("\033[36m[PROGRESS]\033[0m %s: %.1f%% - %s\n", component, progress*100, message)
}

func (n *ConsoleNotifier) NotifySuccess(component, version string, duration time.Duration) {
	fmt.Printf("\033[32m[SUCCESS]\033[0m %s: Successfully upgraded to version %s (took %s)\n", component, version, duration)
}

func (n *ConsoleNotifier) NotifyFailed(component, version string, err error) {
	fmt.Printf("\033[31m[ERROR]\033[0m %s: Failed to upgrade to version %s: %v\n", component, version, err)
}

func (n *ConsoleNotifier) NotifyCancelled(component, version string) {
	fmt.Printf("\033[33m[WARNING]\033[0m %s: Upgrade to version %s was cancelled\n", component, version)
}
