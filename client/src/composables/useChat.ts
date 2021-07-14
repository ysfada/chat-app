import { ref } from 'vue';
import randomId from '../utils/randomId';
import useWebSocket from './useWebSocket';
import useChatState, { IMessage, MessageType } from './useChatState';

enum ResponseEvents {
  ERROR,
  CONNECTED,
  TOPIC_ROOMS,
  ME_CHANGED_USERNAME,
  OTHER_CHANGED_USERNAME,
  ME_JOINED_CHAT,
  OTHER_JOINED_CHAT,
  ME_LEFT_CHAT,
  OTHER_LEFT_CHAT,
  ME_MESSAGE_SEND,
  OTHER_MESSAGE_SEND,
  OLD_MESSAGES
}

enum RequestEvents {
  GET_ROOMS,
  CHANGE_USERNAME,
  JOIN_CHAT,
  LEFT_CHAT,
  SEND_MESSAGE,
  GET_OLD_MESSAGES
}

const url = ref('ws://localhost:8080/ws/chat');
const ws = ref<WebSocket>();
const isConnecting = ref(false);

export default function useChat() {
  const {
    rooms,
    currentRoom,
    me,
    messageInput,
    messages,
    users
  } = useChatState();

  async function connectChat(): Promise<void> {
    if (isConnecting.value) return;

    const { connect, connections } = useWebSocket();

    isConnecting.value = true;
    await connect(url.value)
      .catch(err => console.error(err))
      .catch(() => connectChat())
      .finally(() => (isConnecting.value = false));

    ws.value = connections.value.get(url.value);

    addListener(ws.value);
  }

  function addListener(ws: WebSocket | undefined): void {
    if (!ws) return;

    // ws.addEventListener('open', (e) => {
    //   console.log('wsChat(open):', e);
    // });
    ws.addEventListener('close', e => {
      console.log('wsChat(close):', e);
      if (e.code !== 1000) connectChat();
    });
    ws.addEventListener('error', e => {
      console.error('wsChat(error):', e);
    });
    ws.addEventListener('message', e => {
      const res = JSON.parse(e.data);
      switch (res.type) {
        case ResponseEvents.ERROR:
          console.log('wsChat(message): ResponseEvents.ERROR:');
          console.log(JSON.stringify(res, null, 2));
          break;

        case ResponseEvents.CONNECTED:
          console.log('wsChat(message): ResponseEvents.CONNECTED:');
          getRooms();
          me.value = res.body.data;
          const username = window.localStorage.getItem('username');
          if (username) {
            changeUsername(username);
          }
          break;

        case ResponseEvents.TOPIC_ROOMS:
          console.log('wsChat(message): ResponseEvents.TOPIC_ROOMS:');
          if (res.body && res.body.data) {
            rooms.value = res.body.data;
          }
          break;

        case ResponseEvents.ME_CHANGED_USERNAME:
          console.log('wsChat(message): ResponseEvents.ME_CHANGED_USERNAME:');
          me.value = res.body.data.user;
          window.localStorage.setItem('username', me.value.username);
          // messages.value.push({
          //   id: randomId(),
          //   message: 'you changed username',
          //   timestamp: Date.now(),
          //   type: MessageType.ME_CHANGED_USERNAME,
          //   user: res.body.data,
          // });
          messages.value = [
            ...messages.value,
            {
              id: randomId(),
              message: 'you changed username',
              timestamp: Date.now(),
              type: MessageType.ME_CHANGED_USERNAME,
              user: res.body.data
            }
          ];
          break;

        case ResponseEvents.OTHER_CHANGED_USERNAME:
          console.log(
            'wsChat(message): ResponseEvents.OTHER_CHANGED_USERNAME:'
          );
          // messages.value.push({
          //   id: randomId(),
          //   message: `${res.body.data.username} changed username`,
          //   timestamp: Date.now(),
          //   type: MessageType.OTHER_CHANGED_USERNAME,
          //   user: res.body.data,
          // });
          messages.value = [
            ...messages.value,
            {
              id: randomId(),
              message: `${res.body.data.username} changed username`,
              timestamp: Date.now(),
              type: MessageType.OTHER_CHANGED_USERNAME,
              user: res.body.data
            }
          ];
          console.log(JSON.stringify(res, null, 2));
          break;

        case ResponseEvents.ME_JOINED_CHAT:
          console.log('wsChat(message): ResponseEvents.ME_JOINED_CHAT:');
          if (res.body.data.messages) {
            res.body.data.messages.map(
              (msg: IMessage) => (msg.type = MessageType.NEW_MESSAGE)
            );
            messages.value = res.body.data.messages;
          }
          if (res.body.data.users) {
            users.value = res.body.data.users;
          }
          if (res.body.data.room) {
            currentRoom.value = res.body.data.room;
          }
          // messages.value.push({
          //   id: randomId(),
          //   message: res.body.message,
          //   timestamp: Date.now(),
          //   type: MessageType.ME_JOINED,
          //   user: me.value,
          // });
          messages.value = [
            ...messages.value,
            {
              id: randomId(),
              message: res.body.message,
              timestamp: Date.now(),
              type: MessageType.ME_JOINED,
              user: me.value
            }
          ];
          break;

        case ResponseEvents.OTHER_JOINED_CHAT:
          console.log('wsChat(message): ResponseEvents.OTHER_JOINED_CHAT:');
          users.value.push(res.body.data);
          // messages.value.push({
          //   id: randomId(),
          //   message: `${res.body.data.username} joined chat`,
          //   timestamp: Date.now(),
          //   type: MessageType.OTHER_JOINED,
          //   user: res.body.data,
          // });
          messages.value = [
            ...messages.value,
            {
              id: randomId(),
              message: `${res.body.data.username} joined chat`,
              timestamp: Date.now(),
              type: MessageType.OTHER_JOINED,
              user: res.body.data
            }
          ];
          break;

        case ResponseEvents.ME_LEFT_CHAT:
          console.log('wsChat(message): ResponseEvents.ME_LEFT_CHAT:');
          console.log(JSON.stringify(res, null, 2));
          currentRoom.value = undefined;
          break;

        case ResponseEvents.OTHER_LEFT_CHAT:
          console.log('wsChat(message): ResponseEvents.OTHER_LEFT_CHAT:');
          users.value = users.value.filter(
            user => user.id !== res.body.data.id
          );
          // messages.value.push({
          //   id: randomId(),
          //   message: `${res.body.data.username} left chat`,
          //   timestamp: Date.now(),
          //   type: MessageType.OTHER_LEFT,
          //   user: res.body.data,
          // });
          messages.value = [
            ...messages.value,
            {
              id: randomId(),
              message: `${res.body.data.username} left chat`,
              timestamp: Date.now(),
              type: MessageType.OTHER_LEFT,
              user: res.body.data
            }
          ];
          break;

        case ResponseEvents.ME_MESSAGE_SEND:
          console.log('wsChat(message): ResponseEvents.ME_MESSAGE_SEND:');
          res.body.data.type = MessageType.NEW_MESSAGE;
          res.body.data.user = me.value;
          // messages.value.push(res.body.data);
          messages.value = [...messages.value, res.body.data];
          break;

        case ResponseEvents.OTHER_MESSAGE_SEND:
          console.log('wsChat(message): ResponseEvents.OTHER_MESSAGE_SEND:');
          res.body.data.type = MessageType.NEW_MESSAGE;
          // messages.value.push(res.body.data);
          messages.value = [...messages.value, res.body.data];
          break;

        case ResponseEvents.OLD_MESSAGES:
          console.log('wsChat(message): ResponseEvents.OLD_MESSAGES:');
          if (res.body.data.messages) {
            res.body.data.messages.map(
              (msg: IMessage) => (msg.type = MessageType.NEW_MESSAGE)
            );
            // messages.value.unshift(...res.body.data.messages);
            messages.value = [...res.body.data.messages, ...messages.value];
          } else {
            if (currentRoom.value) {
              currentRoom.value.doneLoading = true;
            }
          }
          break;

        default:
          console.log('wsChat(message): Unknown:');
          console.log(JSON.stringify(res, null, 2));
          break;
      }
    });
  }

  function changeUsername(username: string) {
    if (!username) return;

    try {
      ws.value?.send(
        JSON.stringify({
          type: RequestEvents.CHANGE_USERNAME,
          body: {
            username
          }
        })
      );
    } catch (error) {
      console.error(error);
    }
  }

  function getRooms() {
    try {
      ws.value?.send(
        JSON.stringify({
          type: RequestEvents.GET_ROOMS
        })
      );
    } catch (error) {
      console.error(error);
    }
  }

  function joinChat(roomId: string | (string | null)[]) {
    if (!roomId) return;

    try {
      ws.value?.send(
        JSON.stringify({
          type: RequestEvents.JOIN_CHAT,
          body: {
            roomId
          }
        })
      );
    } catch (error) {
      console.error(error);
    }
  }

  function leftChat(roomId: string | (string | null)[]) {
    if (!roomId) return;

    messages.value = [];
    users.value = [];

    try {
      ws.value?.send(
        JSON.stringify({
          type: RequestEvents.LEFT_CHAT,
          body: {
            roomId
          }
        })
      );
    } catch (error) {
      console.error(error);
    }
  }

  function sendMessage(msg: string, roomId: string | (string | null)[]) {
    if (!roomId || !msg || ws.value?.readyState !== ws.value?.OPEN) return;

    try {
      ws.value?.send(
        JSON.stringify({
          type: RequestEvents.SEND_MESSAGE,
          body: {
            message: msg,
            roomId
          }
        })
      );

      messageInput.value = '';
    } catch (error) {
      console.error(error);
    }
  }

  function getOldMessages(
    roomId: string | (string | null)[],
    oldestMsgId: string
  ) {
    if (!roomId || !oldestMsgId) return;

    try {
      ws.value?.send(
        JSON.stringify({
          type: RequestEvents.GET_OLD_MESSAGES,
          body: {
            roomId,
            oldestMsgId
          }
        })
      );
    } catch (error) {
      console.error(error);
    }
  }

  return {
    connectChat,
    ws,
    url,
    changeUsername,
    getRooms,
    joinChat,
    leftChat,
    sendMessage,
    getOldMessages
  };
}
