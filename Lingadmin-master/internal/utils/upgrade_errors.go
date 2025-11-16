// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"fmt"
)

// UpgradeStage 升级阶段
type UpgradeStage string

const (
	StageCheckVersion UpgradeStage = "check_version"
	StageDownload     UpgradeStage = "download"
	StageVerify       UpgradeStage = "verify"
	StageUnzip        UpgradeStage = "unzip"
	StageBackup       UpgradeStage = "backup"
	StageInstall      UpgradeStage = "install"
	StageRestart      UpgradeStage = "restart"
	StageCleanup      UpgradeStage = "cleanup"
)

// UpgradeErrorCode 错误码
type UpgradeErrorCode int

const (
	ErrCodeUnknown           UpgradeErrorCode = 1000
	ErrCodeNetworkFailed     UpgradeErrorCode = 1001
	ErrCodeInvalidResponse   UpgradeErrorCode = 1002
	ErrCodeNoNewVersion      UpgradeErrorCode = 1003
	ErrCodeDownloadFailed    UpgradeErrorCode = 1004
	ErrCodeVerifyFailed      UpgradeErrorCode = 1005
	ErrCodeUnzipFailed       UpgradeErrorCode = 1006
	ErrCodeBackupFailed      UpgradeErrorCode = 1007
	ErrCodeInstallFailed     UpgradeErrorCode = 1008
	ErrCodeRestartFailed     UpgradeErrorCode = 1009
	ErrCodeCleanupFailed     UpgradeErrorCode = 1010
	ErrCodeInsufficientSpace UpgradeErrorCode = 1011
	ErrCodePermissionDenied  UpgradeErrorCode = 1012
	ErrCodeCancelled         UpgradeErrorCode = 1013
	ErrCodeTimeout           UpgradeErrorCode = 1014
	ErrCodeAlreadyRunning    UpgradeErrorCode = 1015
)

// UpgradeError 升级错误类型
type UpgradeError struct {
	Stage   UpgradeStage
	Code    UpgradeErrorCode
	Message string
	Err     error
}

func NewUpgradeError(stage UpgradeStage, code UpgradeErrorCode, message string, err error) *UpgradeError {
	return &UpgradeError{
		Stage:   stage,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *UpgradeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s][%d] %s: %v", e.Stage, e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s][%d] %s", e.Stage, e.Code, e.Message)
}

func (e *UpgradeError) Unwrap() error {
	return e.Err
}

// IsUpgradeError 判断是否为升级错误
func IsUpgradeError(err error) bool {
	_, ok := err.(*UpgradeError)
	return ok
}

// GetUpgradeError 获取升级错误
func GetUpgradeError(err error) *UpgradeError {
	if ue, ok := err.(*UpgradeError); ok {
		return ue
	}
	return nil
}

// IsRetryable 判断错误是否可重试
func (e *UpgradeError) IsRetryable() bool {
	switch e.Code {
	case ErrCodeNetworkFailed, ErrCodeDownloadFailed, ErrCodeTimeout:
		return true
	case ErrCodeCancelled, ErrCodeAlreadyRunning, ErrCodePermissionDenied, ErrCodeInsufficientSpace:
		return false
	default:
		return false
	}
}

// IsFatal 判断是否为致命错误
func (e *UpgradeError) IsFatal() bool {
	switch e.Code {
	case ErrCodePermissionDenied, ErrCodeInsufficientSpace:
		return true
	default:
		return false
	}
}
