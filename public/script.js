const app = new Vue({
  el: "#app",
  data: {
    ws: null,
    message: "",
    messages: [],
    roomId: null,
    username: null,
    reconnectInterval: null,
  },
  watch: {
    messages: "scrollToBottom",
  },
  beforeMount() {
    this.roomId = this.randomId();
    this.username = `anonymous_${this.randomId()}`;
  },
  mounted() {
    // this.wsConnect()
    this.scrollToBottom();
  },
  methods: {
    randomId() {
      return crypto.getRandomValues(new Uint32Array(1))[0].toString(16);
    },
    wsConnect() {
      if (this.ws) this.ws.close(1000);

      this.ws = new WebSocket(`ws://localhost:8080/ws/${this.roomId}?v=1.0`);

      this.ws.onopen = (_ev) => {
        clearInterval(this.reconnectInterval); // clear interval after connection

        this.messages.push({
          username: this.username,
          message: "Welcome to the chat room!",
          type: "CONNECTED",
        });

        if (!this.ws) return;

        this.setUsername();
        this.getOldMessages();
      };

      this.ws.onclose = (ev) => {
        this.messages.push({
          username: this.username,
          message: "connection closed",
          type: "LEFT",
        });

        if (ev.code === 1000) return;

        if (this.reconnectInterval) return;
        this.messages.push({
          username: this.username,
          message: "trying to reconnect...",
          type: "CONNECTING",
        });
        this.wsConnect(); // immediately try to reconnect
        this.reconnectInterval = setInterval(() => {
          this.messages.push({
            username: this.username,
            message: "trying to reconnect...",
            type: "CONNECTING",
          });
          this.wsConnect();
        }, 10000); // try to reconnect after every 10sec
      };

      this.ws.onmessage = (ev) => {
        const res = JSON.parse(ev.data);
        switch (res.type) {
          case "JOIN":
          case "USERNAME_CHANGED":
          case "LEFT":
          case "MESSAGE":
            this.messages.push(res);
            break;

          case "MESSAGES":
            if (res.messages) this.messages.push(...res.messages);
            break;

          default:
            break;
        }
      };

      this.ws.onerror = (ev) => {
        // eslint-disable-next-line no-console
        console.error(ev);
        if (!this.ws) return;
        this.ws.close();
      };
    },
    sendMessage() {
      if (
        !this.ws ||
        this.ws.readyState !== WebSocket.OPEN ||
        this.message === "" ||
        this.username === ""
      )
        return;

      this.messages.push({
        username: this.username,
        message: this.message,
        type: "MESSAGE",
      });
      const msg = {
        message: this.message,
        type: "MESSAGE",
      };
      this.ws.send(JSON.stringify(msg));
      this.message = "";
    },
    setUsername() {
      if (
        !this.ws ||
        this.ws.readyState !== WebSocket.OPEN ||
        this.username === ""
      )
        return;

      const newUser = {
        message: this.username,
        type: "NEW_USER",
      };
      this.ws.send(JSON.stringify(newUser));
    },
    getOldMessages() {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;

      const oldMessages = {
        type: "MESSAGES",
      };
      this.ws.send(JSON.stringify(oldMessages));
    },
    scrollToBottom() {
      this.$nextTick(() => {
        if (this.$refs.messages == null) return;
        this.$refs.messages.scrollTop = this.$refs.messages.scrollHeight;
      });
    },
    clearMessages() {
      this.message = "";
      this.messages = [];
    },
  },
});
