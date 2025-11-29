// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package iplibrary

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"net"
	"sync"
)

//go:embed internal-ip-library.db
var ipLibraryData []byte

var defaultLibrary = NewIPLibrary()
var commonLibrary *IPLibrary

var libraryLocker = &sync.Mutex{} // 涓轰簡淇濇寔鍔犺浇椤哄簭鎬?

func DefaultIPLibraryData() []byte {
	return ipLibraryData
}

// InitDefault 鍔犺浇榛樿鐨処P搴?
func InitDefault() error {
	libraryLocker.Lock()
	defer libraryLocker.Unlock()

	if commonLibrary != nil {
		defaultLibrary = commonLibrary
		return nil
	}

	var library = NewIPLibrary()
	err := library.InitFromData(ipLibraryData, "", ReaderVersionV1)
	if err != nil {
		return err
	}

	commonLibrary = library
	defaultLibrary = commonLibrary
	return nil
}

// Lookup 鏌ヨIP淇℃伅
func Lookup(ip net.IP) *QueryResult {
	return defaultLibrary.Lookup(ip)
}

// LookupIP 鏌ヨIP淇℃伅
func LookupIP(ip string) *QueryResult {
	return defaultLibrary.LookupIP(ip)
}

// LookupIPSummaries 鏌ヨ涓€缁処P瀵瑰簲鐨勫尯鍩熸弿杩?
func LookupIPSummaries(ipList []string) map[string]string /** ip => summary **/ {
	var result = map[string]string{}
	for _, ip := range ipList {
		var region = LookupIP(ip)
		if region != nil && region.IsOk() {
			result[ip] = region.Summary()
		}
	}
	return result
}

type IPLibrary struct {
	reader ReaderInterface
}

func NewIPLibrary() *IPLibrary {
	return &IPLibrary{}
}

func NewIPLibraryWithReader(reader ReaderInterface) *IPLibrary {
	return &IPLibrary{reader: reader}
}

func (this *IPLibrary) InitFromData(data []byte, password string, version ReaderVersion) error {
	if len(data) == 0 || this.reader != nil {
		return nil
	}

	if len(password) > 0 {
		srcData, err := NewEncrypt().Decode(data, password)
		if err != nil {
			return err
		}
		data = srcData
	}

	var reader = bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer func() {
		_ = gzipReader.Close()
	}()

	var libReader ReaderInterface
	if version == ReaderVersionV2 {
		libReader, err = NewReaderV2(gzipReader)
	} else {
		libReader, err = NewReaderV1(gzipReader)
	}
	if err != nil {
		return err
	}
	this.reader = libReader

	return nil
}

func (this *IPLibrary) Lookup(ip net.IP) *QueryResult {
	if this.reader == nil {
		return &QueryResult{}
	}

	var result = this.reader.Lookup(ip)
	if result == nil {
		result = &QueryResult{}
	}

	return result
}

func (this *IPLibrary) LookupIP(ip string) *QueryResult {
	if this.reader == nil {
		return &QueryResult{}
	}
	return this.Lookup(net.ParseIP(ip))
}

func (this *IPLibrary) Destroy() {
	if this.reader != nil {
		this.reader.Destroy()
		this.reader = nil
	}
}
