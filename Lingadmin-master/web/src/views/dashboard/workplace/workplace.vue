<template>
  <div>
    <div class="mt-2">
      <n-card :bordered="false" title="工作台">
        <n-grid cols="2 s:1 m:1 l:2 xl:2 2xl:2" responsive="screen">
          <n-gi>
            <div class="flex items-center">
              <div>
                <n-avatar circle :size="64" :src="schoolboy" />
              </div>
              <div>
                <p class="px-4 text-xl">早安，{{ username }}，开始你一天的工作吧！</p>
                <p class="px-4 text-gray-400">集群运行正常 · 最新配置已同步到所有节点</p>
              </div>
            </div>
          </n-gi>
          <n-gi>
            <div class="flex justify-end w-full flex-wrap gap-6">
              <div class="flex flex-col text-right">
                <span class="text-secondary">实时带宽</span>
                <span class="text-2xl font-semibold">36 Gbps</span>
                <span class="text-xs text-emerald-500">+4.2% 今日峰值</span>
              </div>
              <div class="flex flex-col text-right">
                <span class="text-secondary">今日请求</span>
                <span class="text-2xl font-semibold">1.82 亿</span>
                <span class="text-xs text-emerald-500">命中率 92.6%</span>
              </div>
              <div class="flex flex-col text-right">
                <span class="text-secondary">在线节点</span>
                <span class="text-2xl font-semibold">128</span>
                <span class="text-xs text-gray-400">可用率 99.95%</span>
              </div>
            </div>
          </n-gi>
        </n-grid>
      </n-card>
    </div>

    <n-grid class="mt-4" cols="2 s:1 m:1 l:2 xl:2 2xl:2" responsive="screen" :x-gap="12" :y-gap="9">
      <n-gi>
        <n-card
          :segmented="{ content: true }"
          content-style="padding: 0;"
          :bordered="false"
          size="small"
          title="运行概览"
        >
          <n-grid cols="2 s:1 m:2 l:2 xl:4 2xl:4" responsive="screen" :x-gap="10" :y-gap="10" class="p-3">
            <n-gi v-for="item in overviewCards" :key="item.label">
              <div class="metric-card">
                <div>
                  <p class="text-sm text-gray-400">{{ item.label }}</p>
                  <p class="text-2xl font-semibold">{{ item.value }}</p>
                  <p class="text-xs text-gray-400">{{ item.desc }}</p>
                </div>
                <n-icon size="30" :color="item.color" :component="item.icon" />
              </div>
            </n-gi>
          </n-grid>
        </n-card>

        <n-card
          :segmented="{ content: true }"
          content-style="padding-top: 0;padding-bottom: 0;"
          :bordered="false"
          size="small"
          title="运维动态"
          class="mt-4"
        >
          <template #header-extra><a href="javascript:;">更多</a></template>
          <n-list>
            <n-list-item v-for="item in activities" :key="item.title + item.time">
              <template #prefix>
                <n-avatar circle :size="40" :src="schoolboy" />
              </template>
              <n-thing :title="item.title">
                <template #description>
                  <p class="text-xs text-gray-500">{{ item.time }} · {{ item.desc }}</p>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card
          :segmented="{ content: true }"
          content-style="padding: 0;"
          :bordered="false"
          size="small"
          title="快捷操作"
        >
          <div class="flex flex-wrap project-card">
            <n-card
              v-for="action in quickActions"
              :key="action.label"
              size="small"
              class="cursor-pointer project-card-item"
              hoverable
            >
              <div class="flex flex-col justify-center text-gray-500">
                <span class="text-center">
                  <n-icon size="30" :color="action.color" :component="action.icon" />
                </span>
                <span class="text-lx text-center">{{ action.label }}</span>
              </div>
            </n-card>
          </div>
        </n-card>
        <n-card :segmented="{ content: true }" :bordered="false" size="small" class="mt-4">
          <img src="~@/assets/images/Business.svg" class="w-full" />
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script lang="ts" setup>
  import schoolboy from '@/assets/images/schoolboy.png';
  import { useUserStore } from '@/store/modules/user';
  import {
    DashboardOutlined,
    ProfileOutlined,
    FileProtectOutlined,
    SettingOutlined,
    ApartmentOutlined,
    Html5Outlined,
  } from '@vicons/antd';

  const userStore = useUserStore();
  const username = userStore.info.username;

  const overviewCards = [
    {
      label: '今日请求',
      value: '1.82 亿',
      desc: '同比 +12%，命中率 92.6%',
      icon: ProfileOutlined,
      color: '#18a058',
    },
    {
      label: '带宽峰值',
      value: '42 Gbps',
      desc: '过去 5 分钟窗口',
      icon: DashboardOutlined,
      color: '#2d8cf0',
    },
    {
      label: '在线节点',
      value: '128 台',
      desc: '可用率 99.95%',
      icon: ApartmentOutlined,
      color: '#ff9800',
    },
    {
      label: '告警/工单',
      value: '3 条',
      desc: 'WAF 告警 2 · 证书到期 1',
      icon: FileProtectOutlined,
      color: '#ff5c93',
    },
  ];

  const quickActions = [
    { label: '接入域名', icon: DashboardOutlined, color: '#18a058' },
    { label: '证书管理', icon: FileProtectOutlined, color: '#1890ff' },
    { label: '访问日志', icon: ProfileOutlined, color: '#f06b96' },
    { label: '节点监控', icon: ApartmentOutlined, color: '#7238d1' },
    { label: 'WAF 策略', icon: SettingOutlined, color: '#ff9800' },
    { label: '配置发布', icon: Html5Outlined, color: '#009688' },
  ];

  const activities = [
    {
      title: `${username} 提交了新的工作台改版`,
      time: '2021-07-04 22:37:16',
      desc: '已合并到主干，等待发布。',
    },
    {
      title: '广州 BGP 节点恢复，带宽回切至主链路',
      time: '2021-07-04 20:12:10',
      desc: '链路监控正常。',
    },
    {
      title: '开启 WAF CC 规则模板「防护-严」',
      time: '2021-07-04 18:05:44',
      desc: '来源：运维自动化任务。',
    },
    {
      title: '推送证书到华南集群，5 个域名更新完成',
      time: '2021-07-04 16:20:00',
      desc: '签发渠道：ACME 自动续期。',
    },
  ];
</script>

<style lang="less" scoped>
  .project-card {
    margin-right: -6px;

    &-item {
      margin: -1px;
      width: 33.333333%;
    }
  }

  .metric-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px;
    border-radius: 12px;
    background-color: #f7f8fa;
  }
</style>
