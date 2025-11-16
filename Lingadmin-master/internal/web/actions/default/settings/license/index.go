// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package license

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	license, err := LoadLicense()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 如果存在授权，尝试从服务器刷新最新信息
	if license.Code != "" {
		refreshedLicense, refreshErr := RefreshLicenseFromServer(license.Code)
		if refreshErr == nil && refreshedLicense != nil {
			// 刷新成功，使用最新的授权信息
			license = refreshedLicense
			// 保存到本地
			_ = SaveLicense(license)
		} else {
			// 刷新失败，检查本地授权是否过期
			if license.ExpireAt != "" {
				expireTime, parseErr := time.Parse("2006-01-02 15:04:05", license.ExpireAt)
				if parseErr == nil {
					license.IsValid = time.Now().Before(expireTime)
				}
			}
		}
	}

	this.Data["license"] = license
	this.Data["isValid"] = license.IsValid

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	LicenseCode string
}) {
	// Generate unique system token
	systemToken := generateSystemToken()

	// Verify license with dl.lingcdn.cloud
	apiURL := "http://dl.lingcdn.cloud/api/license/verify?code=" + url.QueryEscape(params.LicenseCode) + "&token=" + systemToken

	resp, err := http.Get(apiURL)
	if err != nil {
		this.Fail("授权验证失败：无法连接到授权服务器")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		this.Fail("授权验证失败：读取响应失败")
		return
	}

	type VerifyResponse struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Data    *License `json:"data"`
	}

	var verifyResp VerifyResponse
	err = json.Unmarshal(body, &verifyResp)
	if err != nil {
		this.Fail("授权验证失败：解析响应失败")
		return
	}

	if verifyResp.Code != 200 {
		this.Fail("授权验证失败：" + verifyResp.Message)
		return
	}

	if verifyResp.Data == nil {
		this.Fail("授权验证失败：授权码无效")
		return
	}

	// Check if license is expired
	license := verifyResp.Data
	license.Code = params.LicenseCode

	// Set max nodes based on license type if not set by server
	if license.MaxNodes <= 0 {
		license.MaxNodes = GetMaxNodesByType(license.Type)
	}

	if license.ExpireAt != "" {
		expireTime, err := time.Parse("2006-01-02 15:04:05", license.ExpireAt)
		if err == nil {
			license.IsValid = time.Now().Before(expireTime)
		}
	}

	// Save license locally
	err = SaveLicense(license)
	if err != nil {
		this.Fail("授权保存失败：" + err.Error())
		return
	}

	// 返回授权数据
	this.Data["license"] = license
	this.Success()
}

// generateSystemToken generates a unique token for this system
func generateSystemToken() string {
	// Try to get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Generate token based on hostname and timestamp
	data := fmt.Sprintf("%s-%d", hostname, time.Now().Unix())
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// RefreshLicenseFromServer fetches the latest license info from server
func RefreshLicenseFromServer(licenseCode string) (*License, error) {
	systemToken := generateSystemToken()
	apiURL := "http://dl.lingcdn.cloud/api/license/verify?code=" + url.QueryEscape(licenseCode) + "&token=" + systemToken

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type VerifyResponse struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Data    *License `json:"data"`
	}

	var verifyResp VerifyResponse
	err = json.Unmarshal(body, &verifyResp)
	if err != nil {
		return nil, err
	}

	// 如果有数据返回（包括过期的情况410），都使用返回的数据
	if verifyResp.Data != nil {
		license := verifyResp.Data
		license.Code = licenseCode

		// Set max nodes based on license type if not set by server
		if license.MaxNodes <= 0 {
			license.MaxNodes = GetMaxNodesByType(license.Type)
		}

		// Check expiration
		if license.ExpireAt != "" {
			expireTime, parseErr := time.Parse("2006-01-02 15:04:05", license.ExpireAt)
			if parseErr == nil {
				license.IsValid = time.Now().Before(expireTime)
			}
		}

		return license, nil
	}

	// 其他错误情况（没有数据返回）
	return nil, fmt.Errorf("license verification failed: %s", verifyResp.Message)
}
