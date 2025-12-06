<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="测试 websocket">
        尝试在下方输入框中输入任意消息内容，发送后 websocket 服务会原样返回
      </n-card>
    </div>
    <n-card :bordered="false" class="proCard">
      <n-space vertical>
        <n-input-group style="width: 520px">
          <n-input
            v-model:value="inputMessage"
            placeholder="请输入消息内容"
            :style="{ width: '78%' }"
            @keyup.enter="sendMessage"
            @focus="onFocus"
            @blur="onBlur"
          />
          <n-button type="primary" @click="sendMessage">发送消息</n-button>
        </n-input-group>

        <div class="mt-5"></div>

        <n-timeline :icon-size="20">
          <n-timeline-item v-if="isInput" color="grey" content="输入中...">
            <template #icon>
              <n-icon>
                <MessageOutlined />
              </n-icon>
            </template>
          </n-timeline-item>

          <n-timeline-item
            v-for="item in messages"
            :key="`${item.time}-${item.type}-${item.content}`"
            :type="item.type === Enum.SendType ? 'success' : 'info'"
            :title="item.type === Enum.SendType ? '发送消息' : '收到消息'"
            :content="item.content"
            :time="item.time"
          >
            <template #icon>
              <n-icon>
                <SendOutlined v-if="item.type === Enum.SendType" />
                <SoundOutlined v-else />
              </n-icon>
            </template>
          </n-timeline-item>
        </n-timeline>
      </n-space>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { MessageOutlined, SendOutlined, SoundOutlined } from '@vicons/antd';
import { format } from 'date-fns';
import { addOnMessage, removeOnMessage, sendMsg, WebSocketMessage } from '@/utils/websocket';
import { useMessage } from 'naive-ui';

const message = useMessage();
const messages = ref<Message[]>([]);
const inputMessage = ref('你好，HotGo');
const isInput = ref(false);
const testMessageEvent = 'admin/addons/hgexample/testMessage';

enum Enum {
  SendType = 1,
  ReceiveType = 2,
}

interface Message {
  type: Enum;
  content: string;
  time: string;
}

function onFocus() {
  isInput.value = true;
}

function onBlur() {
  isInput.value = false;
}

function insertMessage(msg: Message) {
  messages.value.unshift(msg);
  if (messages.value.length > 10) {
    messages.value = messages.value.slice(0, 10);
  }
}

function sendMessage() {
  if (!inputMessage.value.trim()) {
    message.error('消息内容不能为空');
    return;
  }

  sendMsg(testMessageEvent, { message: inputMessage.value });

  insertMessage({
    type: Enum.SendType,
    content: inputMessage.value,
    time: format(new Date(), 'yyyy-MM-dd HH:mm:ss'),
  });

  inputMessage.value = '';
}

const onMessage = (res: WebSocketMessage) => {
  const content = res?.data?.message ?? '';
  insertMessage({
    type: Enum.ReceiveType,
    content,
    time: format(new Date(), 'yyyy-MM-dd HH:mm:ss'),
  });
};

onMounted(() => {
  addOnMessage(testMessageEvent, onMessage);
});

onBeforeUnmount(() => {
  removeOnMessage(testMessageEvent);
});
</script>

<style scoped></style>
