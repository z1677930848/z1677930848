Tea.context(function () {
	this.$delay(function () {
		this.checkUpdate()
	})
})

Tea.context(function () {
	this.isChecking = false
	this.hasUpdate = false
	this.currentVersion = ""
	this.checkTime = ""
	this.updateInfo = {
		version: "",
		currentVersion: "",
		changelog: "",
		description: "",
		checkTime: "",
		downloadURL: ""
	}

	this.isUpgrading = false
	this.upgradeStatus = "idle" // idle, downloading, verifying, installing, success, failed
	this.upgradeProgress = 0
	this.upgradeStatusText = ""
	this.progressInfo = {}
	this.progressTimer = null

	this.history = []
	this.isLoadingHistory = false

	/**
	 * 检查更新
	 */
	this.checkUpdate = function () {
		this.isChecking = true

		this.$post("/settings/upgrade/check")
			.params({})
			.success(function (resp) {
				this.isChecking = false

				if (resp.data.hasUpdate) {
					this.hasUpdate = true
					this.updateInfo = resp.data.updateInfo
				} else {
					this.hasUpdate = false
					// 获取当前版本信息
					this.getCurrentVersion()
				}
			})
			.fail(function () {
				this.isChecking = false
			})
	}

	/**
	 * 获取当前版本
	 */
	this.getCurrentVersion = function () {
		// 从页面元素或其他方式获取当前版本
		var versionElement = document.querySelector('[data-version]')
		if (versionElement) {
			this.currentVersion = versionElement.getAttribute('data-version')
		} else {
			this.currentVersion = "1.0.7" // 默认版本
		}
		this.checkTime = new Date().toLocaleString('zh-CN')
	}

	/**
	 * 开始升级
	 */
	this.startUpgrade = function () {
		var that = this

		// 确认对话框
		if (!confirm("确定要升级到版本 " + this.updateInfo.version + " 吗？\n\n升级过程中系统将暂时不可用，升级完成后会自动重启。")) {
			return
		}

		this.isUpgrading = true
		this.upgradeStatus = "downloading"
		this.upgradeProgress = 0
		this.upgradeStatusText = "正在准备下载..."

		// 发送升级请求
		this.$post("/settings/upgrade/install")
			.params({})
			.success(function (resp) {
				teaweb.success("升级任务已启动", function () {
					// 开始轮询进度
					that.startProgressPolling()
				})
			})
			.fail(function (resp) {
				that.isUpgrading = false
				that.upgradeStatus = "failed"
				teaweb.warn("启动升级任务失败：" + resp.message)
			})
	}

	/**
	 * 开始轮询进度
	 */
	this.startProgressPolling = function () {
		var that = this

		// 清除之前的定时器
		if (this.progressTimer) {
			clearInterval(this.progressTimer)
		}

		// 每秒轮询一次
		this.progressTimer = setInterval(function () {
			that.getProgress()
		}, 1000)
	}

	/**
	 * 停止轮询进度
	 */
	this.stopProgressPolling = function () {
		if (this.progressTimer) {
			clearInterval(this.progressTimer)
			this.progressTimer = null
		}
	}

	/**
	 * 获取升级进度
	 */
	this.getProgress = function () {
		var that = this

		this.$post("/settings/upgrade/progress")
			.params({})
			.success(function (resp) {
				var status = resp.data.status
				that.upgradeStatus = status
				that.progressInfo = resp.data.progress || {}

				// 根据状态更新进度
				switch (status) {
					case "pending":
						that.upgradeProgress = 5
						that.upgradeStatusText = "准备中..."
						break
					case "downloading":
						that.upgradeProgress = 30
						that.upgradeStatusText = "下载中..."
						break
					case "verifying":
						that.upgradeProgress = 70
						that.upgradeStatusText = "验证文件..."
						break
					case "installing":
						that.upgradeProgress = 85
						that.upgradeStatusText = "安装中..."
						break
					case "success":
						that.upgradeProgress = 100
						that.upgradeStatusText = "升级完成！"
						that.stopProgressPolling()
						that.isUpgrading = false

						// 5秒后刷新页面
						setTimeout(function () {
							window.location.reload()
						}, 5000)
						break
					case "failed":
						that.upgradeProgress = 0
						that.upgradeStatusText = "升级失败：" + (that.progressInfo.errorMessage || "未知错误")
						that.stopProgressPolling()
						that.isUpgrading = false
						break
					case "cancelled":
						that.upgradeProgress = 0
						that.upgradeStatusText = "升级已取消"
						that.stopProgressPolling()
						that.isUpgrading = false
						break
					default:
						// idle 状态
						break
				}
			})
	}

	/**
	 * 显示升级历史
	 */
	this.showHistory = function () {
		this.isLoadingHistory = true
		this.history = []

		var that = this

		// 打开模态框
		$('#historyModal').modal({
			closable: true,
			onShow: function () {
				// 加载历史数据
				that.$post("/settings/upgrade/history")
					.params({
						limit: 20
					})
					.success(function (resp) {
						that.isLoadingHistory = false
						that.history = resp.data.history || []
					})
					.fail(function () {
						that.isLoadingHistory = false
					})
			}
		}).modal('show')
	}

	/**
	 * 关闭历史对话框
	 */
	this.closeHistory = function () {
		$('#historyModal').modal('hide')
	}

	/**
	 * 格式化速度
	 */
	this.formatSpeed = function (speed) {
		if (!speed || speed <= 0) {
			return "-"
		}
		return speed.toFixed(2) + " MB/s"
	}

	/**
	 * 格式化时长
	 */
	this.formatDuration = function (seconds) {
		if (!seconds || seconds <= 0) {
			return "-"
		}
		if (seconds < 60) {
			return seconds + " 秒"
		}
		var minutes = Math.floor(seconds / 60)
		var secs = seconds % 60
		return minutes + " 分 " + secs + " 秒"
	}

	/**
	 * 获取状态样式类
	 */
	this.getStatusClass = function (status) {
		switch (status) {
			case "success":
				return "green"
			case "failed":
				return "red"
			case "downloading":
			case "verifying":
			case "installing":
				return "blue"
			case "cancelled":
				return "grey"
			default:
				return ""
		}
	}

	/**
	 * 获取状态文本
	 */
	this.getStatusText = function (status) {
		var statusMap = {
			"pending": "准备中",
			"downloading": "下载中",
			"verifying": "验证中",
			"installing": "安装中",
			"success": "成功",
			"failed": "失败",
			"rollback": "已回滚",
			"cancelled": "已取消"
		}
		return statusMap[status] || status
	}

	/**
	 * 组件销毁时清理定时器
	 */
	this.$on("$destroy", function () {
		this.stopProgressPolling()
	})
})
