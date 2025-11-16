# 从 dl.lingcdn.cloud 获取更新内容 - 完整说明

## 🎉 好消息

**是的！可以完整地从 dl.lingcdn.cloud 获取到更新内容！**

## ✅ 已验证的功能

### API 端点

```
http://dl.lingcdn.cloud/api/boot/versions
```

### 请求参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|-----|------|------|------|------|
| os | string | 否 | 操作系统 | linux, windows, darwin |
| arch | string | 否 | 架构 | amd64, arm64, 386 |
| component | string | 否 | 组件代码 | admin, api, node |
| current_version | string | 否 | 当前版本 | 1.0.7 |

### 返回数据示例

#### 1. 获取所有组件版本列表

```bash
curl "http://dl.lingcdn.cloud/api/boot/versions?os=linux&arch=amd64"
```

**返回数据**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "host": "http://dl.lingcdn.cloud",
        "versions": [
            {
                "code": "admin",
                "name": "LingCDN管理系统",
                "version": "1.0.10",
                "url": "/updates/admin/linux/amd64/ling-admin-v1.0.10-linux-amd64.zip",
                "size": 14418272,
                "md5": "b198790470446eaca84b10243929b69b",
                "sha256": "09035e50b9d8e57480ba02945778839bab9271026f01c84903a2774ac7eecebd",
                "releaseTime": "2025-10-25 16:28:17",
                "description": "LingCDN管理系统 v1.0.10 - Plus专业版",
                "changelog": "新增功能:\n- 自动更新检测功能(每6小时检查一次)\n- 完整的下载和安装更新流程\n- MD5文件完整性验证\n- 更新前自动备份当前版本\n- 更新后自动重启服务\n- Web目录自动更新支持\n\n改进:\n- 默认启用Plus专业版功能\n- 修改服务名称为ling-admin\n- 优化日志输出\n- 增强错误处理\n\nPlus专业版功能:\n- 完整的多租户管理\n- 高级WAF防护\n- 更多管理功能",
                "isRequired": false
            }
        ]
    }
}
```

#### 2. 检查特定组件更新

```bash
curl "http://dl.lingcdn.cloud/api/boot/versions?component=admin&os=linux&arch=amd64&current_version=1.0.7"
```

**返回数据**:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "host": "http://dl.lingcdn.cloud",
        "currentVersion": "1.0.7",
        "latestVersion": "1.0.10",
        "needUpdate": true,
        "version": {
            "code": "admin",
            "name": "LingCDN管理系统",
            "version": "1.0.10",
            "url": "/updates/admin/linux/amd64/ling-admin-v1.0.10-linux-amd64.zip",
            "size": 14418272,
            "md5": "b198790470446eaca84b10243929b69b",
            "sha256": "09035e50b9d8e57480ba02945778839bab9271026f01c84903a2774ac7eecebd",
            "releaseTime": "2025-10-25 16:28:17",
            "description": "LingCDN管理系统 v1.0.10 - Plus专业版",
            "changelog": "新增功能:\n- 自动更新检测功能(每6小时检查一次)\n- 完整的下载和安装更新流程\n...",
            "isRequired": false,
            "needUpdate": true,
            "downloadCount": 0
        }
    }
}
```

## 📋 返回字段说明

### 版本信息字段

| 字段 | 类型 | 说明 |
|-----|------|------|
| code | string | 组件代码 (admin/api/node) |
| name | string | 组件名称 |
| version | string | 版本号 |
| url | string | 下载相对路径 |
| size | int | 文件大小（字节） |
| **md5** | string | **MD5 校验值** ✓ |
| **sha256** | string | **SHA-256 校验值** ✓ |
| releaseTime | string | 发布时间 |
| **description** | string | **版本描述** ✓ |
| **changelog** | string | **更新日志（变更内容）** ✓ |
| isRequired | boolean | 是否强制更新 |
| needUpdate | boolean | 是否需要更新 |

## 🎯 完整的更新内容

### ✅ 可以获取的信息

1. **版本号** - `version: "1.0.10"`
2. **变更日志** - `changelog: "新增功能:\n- xxx\n改进:\n- yyy"`
3. **版本描述** - `description: "LingCDN管理系统 v1.0.10 - Plus专业版"`
4. **下载地址** - `url: "/updates/admin/linux/amd64/ling-admin-v1.0.10-linux-amd64.zip"`
5. **文件大小** - `size: 14418272`
6. **MD5 校验** - `md5: "b198790470446eaca84b10243929b69b"`
7. **SHA-256 校验** - `sha256: "09035e50b9d8e57480ba02945778839bab9271026f01c84903a2774ac7eecebd"`
8. **发布时间** - `releaseTime: "2025-10-25 16:28:17"`
9. **是否需要更新** - `needUpdate: true`

## 🔧 已完成的配置

### 1. API 端点已存在

```
✅ /www/wwwroot/dl.lingcdn.cloud/api/boot/versions.php
```

### 2. 数据库已配置

```sql
✅ versions 表结构包含所有必要字段
✅ file_sha256 字段已添加
✅ changelog 字段已添加
✅ description 字段已添加
```

### 3. SHA-256 已计算

```bash
✅ admin v1.0.10: 09035e50b9d8e57480ba02945778839bab9271026f01c84903a2774ac7eecebd
✅ api v1.0.4:   b0ed720dc8ba3df3c25e234ebeb2624e9043cb1561ff35447c13ed50f3a32a4a
✅ node v1.0.0:  6b37cf8a368008b88c54f7ec9beb9b1573fbdefb7f955160db44c6203a2521c9
```

### 4. 更新脚本已创建

```bash
✅ /www/wwwroot/dl.lingcdn.cloud/update-sha256.php
```

## 📝 使用示例

### 在代码中调用

```go
// 1. 检查更新
apiURL := "http://dl.lingcdn.cloud/api/boot/versions?os=linux&arch=amd64"
resp, err := http.Get(apiURL)

// 2. 解析返回数据
var result struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        Host     string `json:"host"`
        Versions []struct {
            Code        string `json:"code"`
            Version     string `json:"version"`
            Changelog   string `json:"changelog"`    // ← 变更日志
            Description string `json:"description"`  // ← 版本描述
            SHA256      string `json:"sha256"`       // ← SHA-256
            URL         string `json:"url"`
        } `json:"versions"`
    } `json:"data"`
}
json.NewDecoder(resp.Body).Decode(&result)

// 3. 显示更新内容
for _, ver := range result.Data.Versions {
    if ver.Code == "admin" {
        fmt.Println("新版本:", ver.Version)
        fmt.Println("变更日志:", ver.Changelog)
        fmt.Println("描述:", ver.Description)
    }
}
```

### 在 Web 界面显示

```javascript
// 1. 检查更新
fetch('http://dl.lingcdn.cloud/api/boot/versions?os=linux&arch=amd64')
    .then(res => res.json())
    .then(data => {
        const adminVersion = data.data.versions.find(v => v.code === 'admin');

        // 2. 显示版本信息
        console.log('新版本:', adminVersion.version);
        console.log('变更日志:', adminVersion.changelog);
        console.log('描述:', adminVersion.description);

        // 3. 在界面上展示
        document.getElementById('version').textContent = adminVersion.version;
        document.getElementById('changelog').textContent = adminVersion.changelog;
        document.getElementById('description').textContent = adminVersion.description;
    });
```

## 🔄 完整的更新流程

### 当前系统的更新流程

```
1. 系统启动 / 每6小时
   ↓
2. 调用 API: http://dl.lingcdn.cloud/api/boot/versions
   ↓
3. 获取返回数据：
   - version: "1.0.10"
   - changelog: "新增功能:\n- xxx"
   - description: "LingCDN管理系统 v1.0.10"
   - sha256: "09035e50b9d8e574..."
   - url: "/updates/admin/linux/amd64/..."
   ↓
4. 保存到 update_info.json
   ↓
5. Web 界面读取 update_info.json
   ↓
6. 显示在升级页面：
   - 显示新版本号
   - 显示变更日志 ← 这里就是从 API 获取的！
   - 显示版本描述 ← 这里也是从 API 获取的！
   - 显示一键升级按钮
   ↓
7. 用户点击升级
   ↓
8. 下载文件: http://dl.lingcdn.cloud + url
   ↓
9. 验证 SHA-256 ← 使用 API 返回的 sha256
   ↓
10. 安装并重启
```

## 🧪 测试验证

### 测试 1：获取所有版本

```bash
curl -s "http://dl.lingcdn.cloud/api/boot/versions?os=linux&arch=amd64" \
  | python3 -m json.tool \
  | grep -A 20 '"admin"'
```

**结果**: ✅ 成功返回 admin 的完整信息，包括 changelog

### 测试 2：检查是否需要更新

```bash
curl -s "http://dl.lingcdn.cloud/api/boot/versions?component=admin&os=linux&arch=amd64&current_version=1.0.7" \
  | python3 -m json.tool \
  | grep -E '(needUpdate|version|changelog)'
```

**结果**: ✅ 正确返回 needUpdate: true 和完整的变更日志

### 测试 3：验证 SHA-256

```bash
# 下载文件
wget http://dl.lingcdn.cloud/updates/admin/linux/amd64/ling-admin-v1.0.10-linux-amd64.zip

# 计算 SHA-256
sha256sum ling-admin-v1.0.10-linux-amd64.zip

# 对比 API 返回的值
curl -s "http://dl.lingcdn.cloud/api/boot/versions?component=admin" | grep sha256
```

**结果**: ✅ 校验值完全匹配

## 📚 相关文件

### API 文件
```
/www/wwwroot/dl.lingcdn.cloud/api/boot/versions.php    - 版本查询 API
```

### 工具脚本
```
/www/wwwroot/dl.lingcdn.cloud/update-sha256.php        - SHA-256 计算工具
```

### 数据库
```
Database: lingcdn
Table: versions
Fields: component_code, version, changelog, description, file_sha256
```

## 🎯 总结

### ✅ 完全可以从 dl.lingcdn.cloud 获取更新内容！

**包括:**
- ✅ 版本号
- ✅ **变更日志（changelog）**
- ✅ **版本描述（description）**
- ✅ 下载地址
- ✅ SHA-256 校验值
- ✅ MD5 校验值
- ✅ 文件大小
- ✅ 发布时间
- ✅ 是否需要更新判断

### 🎉 现在的一键升级功能

用户访问升级页面时：

1. **自动检查更新** ← 从 dl.lingcdn.cloud 获取
2. **显示新版本** ← 从 API 获取
3. **显示更新内容** ← **从 API 的 changelog 字段获取** ✓
4. **显示版本描述** ← **从 API 的 description 字段获取** ✓
5. **一键升级按钮** ← 下载并安装
6. **SHA-256 验证** ← 使用 API 返回的 sha256 值

**所有数据都来自 dl.lingcdn.cloud，无需手动配置！**

## 🚀 下次发布新版本时

只需要：

1. 上传新版本文件到 `/www/wwwroot/dl.lingcdn.cloud/updates/admin/linux/amd64/`
2. 在数据库中添加版本记录（包括 changelog 和 description）
3. 运行 `php update-sha256.php` 计算 SHA-256
4. 完成！所有客户端自动检测到新版本和更新内容

---

**创建时间**: 2025-10-31
**测试状态**: ✅ 全部通过
**API 状态**: ✅ 正常运行
