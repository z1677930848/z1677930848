<template>
  <div>
    <n-card title="Example Portal" class="proCard">
      <n-tabs type="line" animated>
        <n-tab-pane name="admin" tab="Admin API">
          <Form path="/admin/hgexample/index/test?name=lingcdn" />
        </n-tab-pane>
        <n-tab-pane name="api" tab="Public API">
          <Form path="/api/hgexample/index/test?name=lingcdn" />
        </n-tab-pane>
        <n-tab-pane name="home" tab="Home">
          <Form path="/home/hgexample/index/test?name=lingcdn" />
        </n-tab-pane>
        <n-tab-pane name="static" tab="Static Assets">
          <Form path="/addons/hgexample/default" />
        </n-tab-pane>
        <n-tab-pane name="websocket" tab="Websocket">
          <Form path="/socket/hgexample/index/test?name=lingcdn" :showQrcode="false" />
          <WebsocketTest />
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import Form from './form.vue';
  import WebsocketTest from './websocketTest.vue';
</script>

<style scoped>
  ::v-deep(.card-tabs .n-tabs-nav--bar-type) {
    padding-left: 4px;
  }
</style>
