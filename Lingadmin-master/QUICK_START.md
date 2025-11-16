# TeaUI 用户端快速开始指南

## 🚀 5 分钟快速上手

本指南帮助你快速了解和使用 TeaUI 用户端组件系统。

---

## 📦 文件引入

### 1. 在页面模板中引入 CSS

```html
<!-- TeaUI 用户端主题 -->
<link rel="stylesheet" href="/css/teaui-user-theme.css"/>
```

### 2. 引入用户端组件

```html
<!-- 用户端专属组件 -->
<script src="/js/components/user/user-dashboard-card.js"></script>
<script src="/js/components/user/user-domain-card.js"></script>
<script src="/js/components/user/user-stats-chart.js"></script>
<script src="/js/components/user/user-quick-action.js"></script>
```

---

## 🎨 基础样式使用

### 布局系统

```html
<!-- 4 列网格 -->
<div class="tea-grid tea-grid-cols-4">
    <div>列 1</div>
    <div>列 2</div>
    <div>列 3</div>
    <div>列 4</div>
</div>

<!-- 2 列网格 -->
<div class="tea-grid tea-grid-cols-2">
    <div>列 1</div>
    <div>列 2</div>
</div>
```

### 卡片组件

```html
<div class="tea-card">
    <div class="tea-card-header">
        <h2 class="tea-card-title">卡片标题</h2>
    </div>
    <div class="tea-card-body">
        卡片内容
    </div>
    <div class="tea-card-footer">
        卡片底部
    </div>
</div>
```

### 按钮组件

```html
<!-- 主要按钮 -->
<button class="tea-button tea-button-primary">主要按钮</button>

<!-- 次要按钮 -->
<button class="tea-button tea-button-secondary">次要按钮</button>

<!-- 成功按钮 -->
<button class="tea-button tea-button-success">成功按钮</button>

<!-- 危险按钮 -->
<button class="tea-button tea-button-danger">危险按钮</button>

<!-- 小按钮 -->
<button class="tea-button tea-button-sm">小按钮</button>

<!-- 大按钮 -->
<button class="tea-button tea-button-lg">大按钮</button>
```

### 表单组件

```html
<div class="tea-form-group">
    <label class="tea-form-label">用户名</label>
    <input type="text" class="tea-form-input" placeholder="请输入用户名">
</div>

<div class="tea-form-group">
    <label class="tea-form-label">描述</label>
    <textarea class="tea-form-textarea" placeholder="请输入描述"></textarea>
</div>

<div class="tea-form-group">
    <label class="tea-form-label">选择</label>
    <select class="tea-form-select">
        <option>选项 1</option>
        <option>选项 2</option>
    </select>
</div>
```

### 徽章组件

```html
<span class="tea-badge tea-badge-success">成功</span>
<span class="tea-badge tea-badge-warning">警告</span>
<span class="tea-badge tea-badge-error">错误</span>
<span class="tea-badge tea-badge-info">信息</span>
```

### 消息提示

```html
<div class="tea-message tea-message-success">
    操作成功！
</div>

<div class="tea-message tea-message-warning">
    请注意！
</div>

<div class="tea-message tea-message-error">
    操作失败！
</div>

<div class="tea-message tea-message-info">
    提示信息
</div>
```

### 加载状态

```html
<div class="tea-loading">
    <div class="tea-spinner"></div>
</div>
```

### 空状态

```html
<div class="tea-empty">
    <div class="tea-empty-icon">📭</div>
    <div class="tea-empty-text">暂无数据</div>
    <button class="tea-button tea-button-primary">添加数据</button>
</div>
```

---

## 🧩 Vue 组件使用

### 1. 仪表盘卡片

```html
<tea-user-dashboard-card
    title="总域名数"
    :value="12"
    trend="运行中 10"
    trend-type="up"
    icon="🌐"
    icon-color="blue"
    link="/user/domains">
</tea-user-dashboard-card>
```

**Props 说明**:
- `title`: 卡片标题 (String)
- `value`: 数值 (String | Number)
- `trend`: 趋势文本 (String)
- `trend-type`: 趋势类型 (up/down/neutral)
- `icon`: 图标 (String, 支持 emoji 或 HTML)
- `icon-color`: 图标颜色 (blue/green/purple/orange/red)
- `link`: 点击跳转链接 (String, 可选)

### 2. 域名卡片

```html
<tea-user-domain-card
    :domain="{
        id: 1,
        name: 'example.com',
        status: 'running',
        requests: 125680,
        traffic: '45.2 GB'
    }"
    :show-actions="true"
    @delete="handleDelete">
</tea-user-domain-card>
```

**Props 说明**:
- `domain`: 域名对象 (Object)
  - `id`: 域名 ID
  - `name`: 域名
  - `status`: 状态 (running/stopped/error)
  - `requests`: 请求数
  - `traffic`: 流量
  - `quota`: 配额 (可选)
  - `usedQuota`: 已用配额 (可选)
- `show-actions`: 是否显示操作按钮 (Boolean)

**Events**:
- `@delete`: 删除事件，参数为域名 ID

### 3. 统计图表

```html
<tea-user-stats-chart
    :data="{
        labels: ['1/1', '1/2', '1/3', '1/4', '1/5'],
        values: [1000, 2000, 1500, 3000, 2500]
    }"
    :loading="false"
    type="line"
    height="300px">
</tea-user-stats-chart>
```

**Props 说明**:
- `data`: 图表数据 (Object)
  - `labels`: X 轴标签数组
  - `values`: Y 轴数值数组
- `loading`: 加载状态 (Boolean)
- `type`: 图表类型 (line/bar)
- `height`: 图表高度 (String)

### 4. 快速操作

```html
<tea-user-quick-action
    icon="➕"
    label="添加域名"
    href="/user/domains/create"
    color="blue">
</tea-user-quick-action>
```

**Props 说明**:
- `icon`: 图标 (String, 支持 emoji 或 HTML)
- `label`: 标签文本 (String)
- `href`: 跳转链接 (String)
- `color`: 颜色主题 (blue/green/purple/orange/red)

---

## 💡 完整示例

### 仪表盘页面示例

```html
{$layout}
{$template "header"}
<script src="/js/echarts/echarts.min.js"></script>
{$end}

<div v-cloak>
    <!-- 页面标题 -->
    <h1 style="font-size: 1.875rem; font-weight: 700; margin-bottom: 2rem;">
        仪表盘 📊
    </h1>

    <!-- 统计卡片 -->
    <div class="tea-grid tea-grid-cols-4" style="margin-bottom: 2rem;">
        <tea-user-dashboard-card
            title="总域名数"
            :value="stats.totalDomains"
            trend="运行中 10"
            trend-type="up"
            icon="🌐"
            icon-color="blue">
        </tea-user-dashboard-card>

        <tea-user-dashboard-card
            title="今日请求"
            :value="stats.todayRequests"
            trend="↑ 12.5%"
            trend-type="up"
            icon="📈"
            icon-color="green">
        </tea-user-dashboard-card>

        <tea-user-dashboard-card
            title="今日流量"
            :value="stats.todayTraffic"
            trend="总计 1.2 TB"
            trend-type="neutral"
            icon="📦"
            icon-color="purple">
        </tea-user-dashboard-card>

        <tea-user-dashboard-card
            title="SSL证书"
            :value="stats.sslCerts"
            trend="2 个即将过期"
            trend-type="down"
            icon="🔒"
            icon-color="orange">
        </tea-user-dashboard-card>
    </div>

    <!-- 图表 -->
    <div class="tea-card">
        <div class="tea-card-header">
            <h2 class="tea-card-title">流量趋势</h2>
        </div>
        <div class="tea-card-body">
            <tea-user-stats-chart
                :data="trafficData"
                :loading="false"
                type="line"
                height="300px">
            </tea-user-stats-chart>
        </div>
    </div>
</div>

{$template "footer"}
<script>
Tea.context(function () {
    // 统计数据
    this.stats = {
        totalDomains: 12,
        todayRequests: 125680,
        todayTraffic: '45.2 GB',
        sslCerts: 8
    };

    // 流量数据
    this.trafficData = {
        labels: ['1/1', '1/2', '1/3', '1/4', '1/5', '1/6', '1/7'],
        values: [5000, 8000, 6500, 9000, 7500, 10000, 8500]
    };
});
</script>
{$end}
```

---

## 🎨 CSS 变量使用

### 在自定义样式中使用设计令牌

```css
.my-custom-card {
    background: white;
    border-radius: var(--tea-radius-lg);
    padding: var(--tea-space-lg);
    box-shadow: var(--tea-shadow-md);
    transition: all var(--tea-transition-base);
}

.my-custom-card:hover {
    box-shadow: var(--tea-shadow-lg);
    transform: translateY(-4px);
}

.my-custom-button {
    background: var(--tea-brand-gradient);
    color: white;
    border-radius: var(--tea-radius-md);
    padding: var(--tea-space-md) var(--tea-space-xl);
}
```

---

## 📱 响应式设计

### 网格系统自动响应

```html
<!-- 桌面端 4 列，平板 2 列，移动端 1 列 -->
<div class="tea-grid tea-grid-cols-4">
    <div>项目 1</div>
    <div>项目 2</div>
    <div>项目 3</div>
    <div>项目 4</div>
</div>
```

**断点说明**:
- 桌面端 (>1024px): 保持原有列数
- 平板端 (768px-1024px): 4 列变 2 列
- 移动端 (<768px): 全部变 1 列

---

## 🔧 常见问题

### Q1: 组件不显示？
**A**: 检查是否正确引入了组件 JS 文件和 TeaUI 主题 CSS。

### Q2: 样式不生效？
**A**: 确保 `teaui-user-theme.css` 在 Semantic UI 之后引入。

### Q3: 图表不显示？
**A**: 确保引入了 ECharts 库：
```html
<script src="/js/echarts/echarts.min.js"></script>
```

### Q4: 移动端样式错乱？
**A**: 检查是否添加了 viewport meta 标签：
```html
<meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=0">
```

---

## 📚 更多资源

- **完整文档**: `TEAUI_OPTIMIZATION_SUMMARY.md`
- **优化方案**: `OPTIMIZATION_PLAN.md`
- **组件源码**: `web/public/js/components/user/`
- **主题样式**: `web/public/css/teaui-user-theme.css`

---

## 🎉 开始使用

现在你已经掌握了 TeaUI 用户端的基础使用方法，可以开始构建你的页面了！

**推荐学习路径**:
1. ✅ 阅读本快速开始指南
2. 📖 查看完整示例页面 (`web/views/@user/dashboard/index.html`)
3. 🎨 了解设计令牌系统 (`teaui-user-theme.css`)
4. 🧩 学习组件开发 (`web/public/js/components/user/`)
5. 📚 阅读完整文档 (`TEAUI_OPTIMIZATION_SUMMARY.md`)

**祝你使用愉快！** 🚀
