<template>
  <div id="log" ref="messagesDivRef">
    <ul class="messages">
      <li v-for="(msg, index) in messages" :key="index">
        <p
          class="message"
          :class="{
            me: msg.username === username,
            connecting: msg.type === 'CONNECTING',
            connected: msg.type === 'CONNECTED',
            join: msg.type === 'JOIN',
            left: msg.type === 'LEFT',
          }"
        >
          <b
            >{{
              msg.username === username ? `Me (${msg.username})` : msg.username
            }}:</b
          >
          {{
            msg.type === "MESSAGE"
              ? msg.message
              : `&lt;&lt; ${msg.message}
              &gt;&gt;`
          }}
        </p>
      </li>
    </ul>
  </div>

  <form id="form" @submit.prevent="sendMessage">
    <input type="button" value="Username" @click="setUsername" />
    <input
      id="username"
      v-model="username"
      type="text"
      :size="16"
      autocomplete="off"
      @keydown.enter="setUsername"
    />
    <input type="submit" value="Send" />
    <input
      id="msg"
      ref="messageInputRef"
      v-model="message"
      type="text"
      :size="64"
      placeholder="Type here..."
      autocomplete="off"
    />
    <input type="button" value="Connect (Room ID)" @click="wsConnect" />
    <input
      id="roomId"
      v-model="roomId"
      type="text"
      :size="64"
      autocomplete="off"
      @keydown.enter="wsConnect"
    />
    <button @click.prevent="clearMessages">Clear</button>
  </form>
</template>

<script lang="ts">
import {
  defineComponent,
  onMounted,
  onBeforeMount,
  watch,
  ref,
  nextTick,
} from "vue";

interface IMessage {
  username: string;
  message: string;
  type: string;
}

export default defineComponent({
  name: "Chat",
  setup() {
    const ws = ref<WebSocket | null>(null);
    const message = ref<string>("");
    const messages = ref<IMessage[]>([]);
    const roomId = ref<string>("");
    const username = ref<string>("");
    const reconnectInterval = ref<number | undefined>(undefined);
    const messagesDivRef = ref<HTMLDivElement>();
    const messageInputRef = ref<HTMLInputElement>();

    const randomId = () => {
      return crypto.getRandomValues(new Uint32Array(1))[0].toString(16);
    };

    const wsConnect = () => {
      if (ws.value) ws.value.close(1000);

      ws.value = new WebSocket(`ws://localhost:8080/ws/${roomId.value}?v=1.0`);

      ws.value.onopen = (_ev) => {
        clearInterval(reconnectInterval.value); // clear interval after connection

        messages.value.push({
          username: username.value,
          message: "Welcome to the chat room!",
          type: "CONNECTED",
        });

        if (!ws.value) return;

        setUsername();
        getOldMessages();
      };

      ws.value.onclose = (ev) => {
        messages.value.push({
          username: username.value,
          message: "connection closed",
          type: "LEFT",
        });

        if (ev.code === 1000) return;

        if (reconnectInterval.value) return;
        messages.value.push({
          username: username.value,
          message: "trying to reconnect...",
          type: "CONNECTING",
        });
        wsConnect(); // immediately try to reconnect
        reconnectInterval.value = setInterval(() => {
          messages.value.push({
            username: username.value,
            message: "trying to reconnect...",
            type: "CONNECTING",
          });
          wsConnect();
        }, 10000); // try to reconnect after every 10sec
      };

      ws.value.onmessage = (ev) => {
        const res = JSON.parse(ev.data);
        switch (res.type) {
          case "JOIN":
          case "USERNAME_CHANGED":
          case "LEFT":
          case "MESSAGE":
            messages.value.push(res);
            break;

          case "MESSAGES":
            if (res.messages) messages.value.push(...res.messages);
            break;

          default:
            break;
        }
      };

      ws.value.onerror = (ev) => {
        // eslint-disable-next-line no-console
        console.error(ev);
        if (!ws.value) return;
        ws.value.close();
      };
    };

    const sendMessage = () => {
      if (
        !ws.value ||
        ws.value.readyState !== WebSocket.OPEN ||
        message.value === "" ||
        username.value === ""
      )
        return;

      messages.value.push({
        username: username.value,
        message: message.value,
        type: "MESSAGE",
      });
      const msg = {
        message: message.value,
        type: "MESSAGE",
      };
      ws.value.send(JSON.stringify(msg));
      message.value = "";
    };

    const setUsername = () => {
      if (
        !ws.value ||
        ws.value.readyState !== WebSocket.OPEN ||
        username.value === ""
      )
        return;

      const newUser = {
        message: username.value,
        type: "NEW_USER",
      };
      ws.value.send(JSON.stringify(newUser));
    };

    const getOldMessages = () => {
      if (!ws.value || ws.value.readyState !== WebSocket.OPEN) return;

      const oldMessages = {
        type: "MESSAGES",
      };
      ws.value.send(JSON.stringify(oldMessages));
    };

    const scrollToBottom = () => {
      nextTick(() => {
        if (messagesDivRef.value == null) return;

        messagesDivRef.value.scrollTop = messagesDivRef.value.scrollHeight;
      });
    };

    const clearMessages = () => {
      message.value = "";
      messages.value = [];
    };

    onBeforeMount(() => {
      roomId.value = randomId();
      username.value = `anonymous_${randomId()}`;
    });

    onMounted(scrollToBottom);
    onMounted(() => {
      // wsConnect()
      messageInputRef.value?.focus();
    });

    watch(messages.value, () => scrollToBottom());

    return {
      messages,
      message,
      username,
      roomId,
      messagesDivRef,
      messageInputRef,
      wsConnect,
      sendMessage,
      setUsername,
      clearMessages,
    };
  },
});
</script>

<style>
html {
  overflow: hidden;
}

body {
  overflow: hidden;
  padding: 0;
  margin: 0;
  width: 100%;
  height: 100%;
  background: gray;
}

#log {
  background: white;
  margin-bottom: 4rem;
  padding: 0 0.5em 0 0.5em;
  position: absolute;
  top: 1.5em;
  left: 0.5em;
  right: 0.5em;
  bottom: 3em;
  overflow: auto;
}

#form {
  box-sizing: border-box;
  padding: 0 0.5em 0 0.5em;
  margin: 0;
  position: absolute;
  bottom: 1em;
  left: 0px;
  width: 100%;
  overflow: hidden;
}

#msg,
#roomId,
#username {
  width: 80%;
}

input[type="submit"],
input[type="button"],
label {
  width: 15%;
}

.messages {
  list-style: none;
  padding: 0;
  margin: 0;
}

.message {
  overflow-wrap: break-word;
  padding: 0.5rem;
  background-color: #b9b6b6;
}

.me {
  color: white;
  background-color: #05c535;
}

.connecting {
  background-color: rgb(231, 114, 36);
}

.connected {
  background-color: rgb(52, 118, 204);
}

.join {
  background-color: rgb(204, 223, 35);
}

.left {
  background-color: rgb(209, 50, 50);
}
</style>
