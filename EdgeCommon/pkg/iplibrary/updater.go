// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package iplibrary

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdaterSource interface {
	// DataDir 鏂囦欢鐩綍
	DataDir() string

	// FindLatestFile 妫€鏌ユ渶鏂扮殑IP搴撴枃浠?
	FindLatestFile() (code string, fileId int64, err error)

	// DownloadFile 涓嬭浇鏂囦欢
	DownloadFile(fileId int64, writer io.Writer) error

	// LogInfo 鏅€氭棩蹇?
	LogInfo(message string)

	// LogError 閿欒鏃ュ織
	LogError(err error)
}

type Updater struct {
	source UpdaterSource

	currentCode string
	ticker      *time.Ticker

	isUpdating bool
}

func NewUpdater(source UpdaterSource, interval time.Duration) *Updater {
	return &Updater{
		source: source,
		ticker: time.NewTicker(interval),
	}
}

func (this *Updater) Start() {
	// 鍒濆鍖?
	err := this.Init()
	if err != nil {
		this.source.LogError(err)
	}

	// 鍏堣繍琛屼竴娆?
	err = this.Loop()
	if err != nil {
		this.source.LogError(err)
	}

	// 寮€濮嬪畾鏃惰繍琛?
	for range this.ticker.C {
		err = this.Loop()
		if err != nil {
			this.source.LogError(err)
		}
	}
}

func (this *Updater) Init() error {
	// 妫€鏌ュ綋鍓嶆鍦ㄤ娇鐢ㄧ殑IP搴?
	var path = this.source.DataDir() + "/ip-library.db"
	fp, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("read ip library file failed '%w'", err)
	}
	defer func() {
		_ = fp.Close()
	}()

	return this.loadFile(fp)
}

func (this *Updater) Loop() error {
	if this.isUpdating {
		return nil
	}

	this.isUpdating = true

	defer func() {
		this.isUpdating = false
	}()

	code, fileId, err := this.source.FindLatestFile()
	if err != nil {
		// 涓嶆彁绀鸿繛鎺ラ敊璇?
		if this.isConnError(err) {
			return nil
		}
		return err
	}
	if len(code) == 0 || fileId <= 0 {
		// 杩樺師鍒板唴缃甀P搴?
		if len(this.currentCode) > 0 {
			this.currentCode = ""
			this.source.LogInfo("resetting to default ip library ...")

			var defaultPath = this.source.DataDir() + "/ip-library.db"
			_, err = os.Stat(defaultPath)
			if err == nil {
				err = os.Remove(defaultPath)
				if err != nil {
					this.source.LogError(errors.New("can not remove default 'ip-library.db'"))
				}
			}

			err = InitDefault()
			if err != nil {
				this.source.LogError(errors.New("initialize default ip library failed: " + err.Error()))
			}
		}

		return nil
	}

	// 涓嬭浇
	if this.currentCode == code {
		// 涓嶅啀閲嶅涓嬭浇
		return nil
	}

	// 妫€鏌ユ槸鍚﹀瓨鍦?
	var dir = this.source.DataDir()
	var path = dir + "/ip-" + code + ".db"
	stat, err := os.Stat(path)
	if err == nil && !stat.IsDir() && stat.Size() > 0 {
		fp, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func() {
			_ = fp.Close()
		}()

		err = this.loadFile(fp)
		if err != nil {
			// 灏濊瘯鍒犻櫎
			_ = os.Remove(path)
		} else {
			this.currentCode = code

			// 鎷疯礉鍒?ip-library.db
			err = this.createDefaultFile(path, dir)
			if err != nil {
				this.source.LogError(err)
			}
		}
		return err
	}

	// write to file
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("create ip library file failed: %w", err)
	}

	var isOk = false
	defer func() {
		if !isOk {
			_ = os.Remove(path)
		}
	}()

	err = this.source.DownloadFile(fileId, fp)
	if err != nil {
		_ = fp.Close()
		return err
	}
	err = fp.Close()
	if err != nil {
		return nil
	}

	// load library from file
	fp, err = os.Open(path)
	if err != nil {
		return nil
	}
	err = this.loadFile(fp)
	_ = fp.Close()
	if err != nil {
		return fmt.Errorf("load file failed: %w", err)
	}

	isOk = true
	this.currentCode = code

	// 鎷疯礉鍒?ip-library.db
	err = this.createDefaultFile(path, dir)
	if err != nil {
		this.source.LogError(err)
	}

	return nil
}

func (this *Updater) loadFile(fp *os.File) error {
	this.source.LogInfo("load ip library from '" + fp.Name() + "' ...")

	var version = ReaderVersionV1
	if strings.HasSuffix(fp.Name(), ".v2.db") {
		version = ReaderVersionV2
	}

	fileReader, err := NewFileDataReader(fp, "", version)
	if err != nil {
		return fmt.Errorf("load ip library from reader failed: %w", err)
	}

	var reader = fileReader.RawReader()
	defaultLibrary = NewIPLibraryWithReader(reader)
	this.currentCode = reader.Meta().Code
	return nil
}

func (this *Updater) createDefaultFile(sourcePath string, dir string) error {
	sourceFp, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("prepare to copy file to 'ip-library.db' failed: %w", err)
	}
	defer func() {
		_ = sourceFp.Close()
	}()

	dstFp, err := os.Create(dir + "/ip-library.db")
	if err != nil {
		return fmt.Errorf("prepare to copy file to 'ip-library.db' failed: %w", err)
	}
	defer func() {
		_ = dstFp.Close()
	}()
	_, err = io.Copy(dstFp, sourceFp)
	if err != nil {
		return fmt.Errorf("copy file to 'ip-library.db' failed: %w", err)
	}
	return nil
}

// isConnError 鏄惁涓鸿繛鎺ラ敊璇?
func (this *Updater) isConnError(err error) bool {
	if err == nil {
		return false
	}

	// 妫€鏌ユ槸鍚︿负杩炴帴閿欒
	statusErr, ok := status.FromError(err)
	if ok {
		var errorCode = statusErr.Code()
		return errorCode == codes.Unavailable || errorCode == codes.Canceled
	}

	if strings.Contains(err.Error(), "code = Canceled") {
		return true
	}

	return false
}
