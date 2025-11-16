package userui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/actions"
	"io"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	config, err := configloaders.LoadUserUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["config"] = config

	// 时区
	this.Data["timeZoneGroups"] = nodeconfigs.FindAllTimeZoneGroups()
	this.Data["timeZoneLocations"] = nodeconfigs.FindAllTimeZoneLocations()

	if len(config.TimeZone) == 0 {
		config.TimeZone = nodeconfigs.DefaultTimeZoneLocation
	}
	this.Data["timeZoneLocation"] = nodeconfigs.FindTimeZoneLocation(config.TimeZone)

	// 带宽单位
	this.Data["bandwidthUnits"] = []map[string]interface{}{
		{"name": "Mbps", "code": systemconfigs.BandwidthUnitBit},
		{"name": "MB/s", "code": systemconfigs.BandwidthUnitByte},
	}

	// 带宽日期范围
	this.Data["bandwidthDateRanges"] = systemconfigs.FindAllBandwidthDateRanges()

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ProductName        string
	UserSystemName     string
	ShowPageFooter     bool
	PageFooterHTML     string
	ShowVersion        bool
	Version            string
	ShowFinance        bool
	FaviconFile        *actions.File
	LogoFile           *actions.File
	TimeZone           string
	ClientIPHeaderNames string
	CheckCNAME         bool
	BandwidthUnit      systemconfigs.BandwidthUnit
	ShowTrafficCharts  bool
	ShowCacheInfoInTrafficCharts bool
	ShowBandwidthCharts          bool
	BandwidthPercentile          int32
	DefaultBandwidthDateRange    string
	BandwidthAlgo                systemconfigs.BandwidthAlgo

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改用户界面设置")

	params.Must.
		Field("productName", params.ProductName).
		Require("请输入产品名称").
		Field("userSystemName", params.UserSystemName).
		Require("请输入用户系统名称")

	config, err := configloaders.LoadUserUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	config.ProductName = params.ProductName
	config.UserSystemName = params.UserSystemName
	config.ShowPageFooter = params.ShowPageFooter
	config.PageFooterHTML = params.PageFooterHTML
	config.ShowVersion = params.ShowVersion
	config.Version = params.Version
	config.ShowFinance = params.ShowFinance
	config.TimeZone = params.TimeZone
	config.ClientIPHeaderNames = params.ClientIPHeaderNames
	config.Server.CheckCNAME = params.CheckCNAME
	config.BandwidthUnit = params.BandwidthUnit
	config.ShowTrafficCharts = params.ShowTrafficCharts
	config.ShowCacheInfoInTrafficCharts = params.ShowCacheInfoInTrafficCharts
	config.ShowBandwidthCharts = params.ShowBandwidthCharts
	config.TrafficStats.BandwidthPercentile = params.BandwidthPercentile
	config.TrafficStats.DefaultBandwidthDateRange = params.DefaultBandwidthDateRange
	config.TrafficStats.BandwidthAlgo = params.BandwidthAlgo

	// 上传Favicon文件
	if params.FaviconFile != nil {
		createResp, err := this.RPC().FileRPC().CreateFile(this.AdminContext(), &pb.CreateFileRequest{
			Filename: params.FaviconFile.Filename,
			Size:     params.FaviconFile.Size,
			IsPublic: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		fileId := createResp.FileId

		// 上传内容
		buf := make([]byte, 512*1024)
		reader, err := params.FaviconFile.OriginFile.Open()
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				_, err = this.RPC().FileChunkRPC().CreateFileChunk(this.AdminContext(), &pb.CreateFileChunkRequest{
					FileId: fileId,
					Data:   buf[:n],
				})
				if err != nil {
					this.Fail("上传失败：" + err.Error())
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				this.Fail("上传失败：" + err.Error())
			}
		}

		// 置为已完成
		_, err = this.RPC().FileRPC().UpdateFileFinished(this.AdminContext(), &pb.UpdateFileFinishedRequest{FileId: fileId})
		if err != nil {
			this.ErrorPage(err)
		}
		config.FaviconFileId = fileId
	}

	// 上传Logo文件
	if params.LogoFile != nil {
		createResp, err := this.RPC().FileRPC().CreateFile(this.AdminContext(), &pb.CreateFileRequest{
			Filename: params.LogoFile.Filename,
			Size:     params.LogoFile.Size,
			IsPublic: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		fileId := createResp.FileId

		// 上传内容
		buf := make([]byte, 512*1024)
		reader, err := params.LogoFile.OriginFile.Open()
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				_, err = this.RPC().FileChunkRPC().CreateFileChunk(this.AdminContext(), &pb.CreateFileChunkRequest{
					FileId: fileId,
					Data:   buf[:n],
				})
				if err != nil {
					this.Fail("上传失败：" + err.Error())
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				this.Fail("上传失败：" + err.Error())
			}
		}

		// 置为已完成
		_, err = this.RPC().FileRPC().UpdateFileFinished(this.AdminContext(), &pb.UpdateFileFinishedRequest{FileId: fileId})
		if err != nil {
			this.ErrorPage(err)
		}
		config.LogoFileId = fileId
	}

	err = configloaders.UpdateUserUIConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
