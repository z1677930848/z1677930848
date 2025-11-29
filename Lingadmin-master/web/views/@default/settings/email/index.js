Tea.context(function () {
	this.success = NotifyReloadSuccess("保存成功")

	this.sendTest = function () {
		if (!this.testEmail) {
			teaweb.warn("请输入测试邮箱")
			return
		}

		this.$post("/settings/email/test")
			.params({
				smtpHost: this.smtpHost,
				smtpPort: this.smtpPort,
				smtpUsername: this.smtpUsername,
				smtpPassword: this.smtpPassword,
				fromEmail: this.fromEmail,
				useTLS: this.useTLS,
				testEmail: this.testEmail
			})
			.success(function () {
				teaweb.success("测试邮件发送成功")
			})
			.fail(function (resp) {
				teaweb.warn(resp.message)
			})
	}
})
