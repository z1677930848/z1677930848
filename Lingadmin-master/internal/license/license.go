// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package license

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type License struct {
	Code      string `json:"code"`
	Type      string `json:"type"`       // trial(10), standard(30), professional(60), enterprise(100), offline(unlimited)
	Company   string `json:"company"`
	Email     string `json:"email"`
	MaxNodes  int    `json:"maxNodes"`   // 0 means unlimited
	ExpireAt  string `json:"expireAt"`
	CreatedAt string `json:"createdAt"`
	IsValid   bool   `json:"isValid"`
}

var licensePath = "/opt/lingcdn/configs/license.json"

// GetMaxNodesByType returns the max nodes limit for a license 
type
func GetMaxNodesByType(licenseType string) int {
	switch licenseType {
	case "trial":
		return 10
	case "standard":
		return 30
	case "professional":
		return 60
	case "enterprise":
		return 100
	case "offline":
		return 0 // 0 means unlimited
	default:
		return 0
	}
}

// LoadLicense loads the license from local file
func LoadLicense() (*License, error) {
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return &License{}, nil
	}

	data, err := ioutil.ReadFile(licensePath)
	if err != nil {
		return nil, err
	}

	var license License
	err = json.Unmarshal(data, &license)
	if err != nil {
		return nil, err
	}

	// Check if license is expired
	if license.ExpireAt != "" {
		expireTime, err := time.Parse("2006-01-02 15:04:05", license.ExpireAt)
		if err == nil {
			license.IsValid = time.Now().Before(expireTime)
		}
	}

	return &license, nil
}

// SaveLicense saves the license to local file
func SaveLicense(license *License) error {
	// Create directory if not exists
	dir := filepath.Dir(licensePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(licensePath, data, 0644)
}
