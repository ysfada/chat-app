<template>
  <main
    :class="{
      roomDrawerOpened: isRoomDrawerOpen,
      userDrawerOpened: isUserDrawerOpen && currentRoom,
    }"
    ref="mainRef"
    @scroll="infiniteScrollTop"
  >
    <template v-if="currentRoom">
      <div class="chat-title">
        <h1>{{ isRoom(currentRoom) ? currentRoom.name : currentRoom.username }}</h1>
      </div>
      <ul class="message-list">
        <li
          v-for="message in messages"
          :key="message.id"
          class="message-item"
          :class="{
            'message-item--me': message.type === MessageType.NEW_MESSAGE && message.user.id === me.id,
            'message-item--other': message.type === MessageType.NEW_MESSAGE && message.user.id !== me.id,
            'message-item--notification': message.type !== MessageType.NEW_MESSAGE,
            'message-item--notification--me': message.type !== MessageType.NEW_MESSAGE && message.user.id === me.id,
            'message-item--notification--other': message.type !== MessageType.NEW_MESSAGE && message.user.id !== me.id,
          }"
        >
          <ChatAvatar
            v-if="message.type === MessageType.NEW_MESSAGE"
            :src="message.user.avatar || 'https://picsum.photos/56/56'"
            :username="message.user.username"
          />
          <div
            v-if="message.type === MessageType.NEW_MESSAGE"
            class="username"
          >{{ message.user.username }}</div>
          <div
            class="message"
            :class="{
              'message--me': message.type === MessageType.NEW_MESSAGE && message.user.id === me.id,
              'message--other': message.type === MessageType.NEW_MESSAGE && message.user.id !== me.id,
              'message--notification': message.type !== MessageType.NEW_MESSAGE,
            }"
          >
            <p>{{ message.message }}</p>
          </div>
          <div v-if="message.type === MessageType.NEW_MESSAGE" class="timestamp">
            <small>
              <i>{{ formatDate(message.timestamp) }}</i>
            </small>
          </div>
        </li>
      </ul>

      <!-- TODO: fix position when drawers open/closed-->
      <button
        class="scroll-bottom"
        :class="{
          'scroll-bottom--hide': hideScrollButton,
          'scroll-bottom--roomDrawerClosed': !isRoomDrawerOpen,
          'scroll-bottom--userDrawerOpened': isUserDrawerOpen && currentRoom,
          'scroll-bottom--bothClosed': !isRoomDrawerOpen && (!isUserDrawerOpen || !currentRoom),
        }"
        type="button"
        aria-label="Scroll to bottom"
        title="Scroll to bottom"
        @click="scrollToBottom"
      >
        <svg style="width:24px;height:24px" viewBox="0 0 24 24">
          <path fill="currentColor" d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z" />
        </svg>
      </button>

      <form
        :class="{
          'form--RoomDrawerOpened': isRoomDrawerOpen,
          'form--UserDrawerOpened': isUserDrawerOpen,
        }"
        @submit.prevent="onSubmit"
      >
        <textarea
          ref="textareaRef"
          v-model="messageInput"
          cols="1"
          rows="1"
          placeholder="Type message..."
          required
          @input="onType"
          @keydown.esc.exact="clearText"
          @keydown.enter.prevent.exact="onSubmit"
          @keydown.enter.shift.exact="addNewLine"
        ></textarea>
        <button type="submit">Send</button>
      </form>
    </template>

    <slot v-else></slot>
  </main>
</template>

<script lang='ts'>
import _debounce from 'lodash/debounce'
import _throttle from 'lodash/throttle'
import { defineComponent, nextTick, onMounted, ref, watch } from 'vue'
import useChat from '../../composables/useChat';
import useChatState, { isRoom, MessageType } from '../../composables/useChatState';
import useDrawer from '../../composables/useDrawer';
import { formatDate } from '../../utils/formatDate';
import ChatAvatar from './ChatAvatar.vue';

export default defineComponent({
  name: 'ChatMain',
  components: {
    ChatAvatar,
  },
  setup() {
    const { isRoomDrawerOpen, isUserDrawerOpen } = useDrawer()
    const { sendMessage, getOldMessages } = useChat()
    const { me, currentRoom, messages, messageInput } = useChatState()

    const mainRef = ref()
    const textareaRef = ref<HTMLTextAreaElement>()
    const hideScrollButton = ref(true)

    const addNewLine = (_e: KeyboardEvent) => {
      messageInput.value += '\n'
    }

    const onType = _debounce(() => {
      // TODO: inform users in chat about typing
      // console.log('typing... ', e.data)
      console.log('typed... ', messageInput.value)
    }, 200)

    const clearText = () => {
      messageInput.value = ''
    }

    const onSubmit = (_e: Event) => {
      if (messageInput.value && currentRoom.value) {
        sendMessage(messageInput.value, currentRoom.value.id)
      }
    }

    const scrollToBottom = async () => {
      await nextTick()
      if (mainRef.value == null) return
      mainRef.value.scrollTop = mainRef.value.scrollHeight
    }

    onMounted(() => textareaRef.value?.focus())

    onMounted(scrollToBottom)

    const setScrollButton = _throttle((e) => {
      if (e.target.scrollHeight - e.target.scrollTop < 750) {
        hideScrollButton.value = true
      } else {
        hideScrollButton.value = false
      }
    }, 300)

    onMounted(() => {
      mainRef.value?.addEventListener('scroll', setScrollButton)
    })

    watch(messages, (value, oldValue) => {
      if (
        oldValue.length > 0 &&
        oldValue[0] &&
        value[0] &&
        oldValue[0].id !== value[0].id
      ) return // in this case new messages pushed to array not unshifted

      scrollToBottom()
    })

    const infiniteScrollTop = _throttle(() => {
      // await nextTick()
      if (
        currentRoom.value?.doneLoading ||
        mainRef.value == null ||
        mainRef.value.scrollTop / mainRef.value.clientHeight > 1
      )
        return

      console.log("Get Messages")
      if (currentRoom.value && messages.value[0]) {
        getOldMessages(currentRoom.value.id, messages.value[0].id)
      }
    }, 300)

    return {
      mainRef,
      textareaRef,
      MessageType,
      isRoom,
      isRoomDrawerOpen,
      isUserDrawerOpen,
      me,
      currentRoom,
      messages,
      messageInput,
      addNewLine,
      onType,
      clearText,
      onSubmit,
      scrollToBottom,
      infiniteScrollTop,
      formatDate,
      hideScrollButton,
    }
  }
})
</script>

<style scoped>
main {
  height: calc(100vh - (var(--header-height) + var(--textarea--height)));
  margin-top: var(--header-height);
  padding: 0 var(--page-padding);
  overflow-y: auto;
  will-change: padding, margin;
  transition: padding var(--drawer-transition), margin var(--drawer-transition);
}

.chat-title {
  display: flex;
  justify-content: center;
  align-items: center;
  position: sticky;
  top: 0;
  opacity: 0.7;
  background-color: var(--clr-background);
}

.chat-title > h1 {
  margin: 0;
}

.message-list {
  max-width: 35rem;
  margin: 0 auto;
  padding: 0;
  list-style: none;
}

.message-item {
  /* width: fit-content; */
  /* max-width: 26rem; */
  margin-bottom: 0.5rem;
}

.message-item--me {
  margin-left: auto;
}

.message-item--other {
  margin-right: auto;
}

.message-item--notification {
  display: flex !important;
}

.message-item--notification--me {
  justify-content: flex-end;
}

.message-item--notification--other {
  justify-content: flex-start;
}

.username {
  overflow: hidden;
  white-space: nowrap;
}

.message {
  display: flex;
  align-items: center;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  white-space: pre-line;
  /* line-break: anywhere; */
  hyphens: auto;
}

.message--me {
  border-bottom-left-radius: 1rem;
  color: var(--clr-background);
  background-color: limegreen;
}

.message--other {
  border-bottom-right-radius: 1rem;
  color: var(--clr-background);
  background-color: gainsboro;
}

.message--notification {
  width: fit-content;
  padding: 0.25rem 0.5rem;
  background-color: hsla(0, 0%, 15%, 0.8);
}

.message > p {
  margin: 0;
}

.timestamp {
  text-align: right;
}

.scroll-bottom {
  display: flex;
  align-items: center;
  justify-content: center;
  position: fixed;
  right: var(--page-padding);
  bottom: calc(var(--textarea--height) * 2);
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  border: 0;
  cursor: pointer;
  transform: scale(1);
  transition: var(--drawer-transition);
}

.scroll-bottom--hide {
  transform: scale(0);
}

.scroll-bottom::before {
  content: "";
  position: absolute;
  inset: 0;
  border-radius: inherit;
  pointer-events: none;
  transition: opacity var(--btn-transition);
  opacity: 0;
  color: inherit;
  background-color: currentColor;
}

.scroll-bottom:is(:hover, :focus-visible)::before {
  opacity: 0.14;
}

form {
  display: flex;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: var(--textarea--height);
  max-width: 30rem;
  margin-left: auto;
  margin-right: auto;
  box-shadow: 0px 0px 5px 0px #787a7d;
  transition: box-shadow 50ms ease-in, left var(--drawer-transition),
    right var(--drawer-transition);
  background-color: hsl(210, 14%, 17%);
}

/* form:before,
form:after {
  content: '';
  position: absolute;
  height: 100%;
  width: 100%;
  opacity: 0.46;
  background-color: hsl(0, 0%, 13%);
}
form:after {
  transform: translateX(-100%);
}
form:before {
  transform: translateX(100%);
} */

form:focus-within {
  box-shadow: 0 0 10px 3px hsl(210, 44%, 40%) !important;
}

textarea {
  flex: 1;
  padding: 3px 3px 0 3px;
  overflow: hidden;
  resize: none;
  font-size: 1.5rem;
  cursor: auto;
  border: 0;
  outline: none;
  border-radius: 3px 0 0 0;
  background-color: hsl(100, 10%, 80%);
}

button[type="submit"] {
  border: 0;
  border-radius: 0 3px 0 0;
  border-left: 2px solid hsl(210, 14%, 17%);
  outline: none;
  background-color: hsl(100, 10%, 80%);
}
button[type="submit"]:is(:hover, :focus-within) {
  background-color: hsl(100, 10%, 75%);
}

@media only screen and (min-width: 358px) {
  .message-item {
    display: grid;
    grid-template-areas:
      ". username"
      "avatar message"
      ". timestamp";
    grid-template-rows: 1rem auto 1rem;
    grid-template-columns: 3.5rem fit-content(26rem);
    gap: 0 0.5rem;
  }

  .message-item--me {
    justify-content: end;
  }

  .message-item--other {
    justify-content: start;
  }

  .username {
    grid-area: username;
  }

  .message {
    grid-area: message;
  }

  .message--notification {
    grid-row: 2;
    grid-column: 1/3;
  }

  .timestamp {
    grid-area: timestamp;
  }
}

@media only screen and (min-width: 608px) {
  .roomDrawerOpened {
    padding-left: calc(var(--room-drawer-width) + var(--page-padding));
  }

  .scroll-bottom--roomDrawerClosed {
    /* margin-right: var(--room-drawer-width); */
    margin-right: 0;
  }

  form {
    left: 0;
  }

  .form--RoomDrawerOpened {
    left: var(--user-drawer-width);
  }
}

@media only screen and (min-width: 1272px) {
  .userDrawerOpened {
    margin-right: var(--user-drawer-width);
  }

  .scroll-bottom--userDrawerOpened {
    margin-right: var(--user-drawer-width);
  }

  .scroll-bottom--roomDrawerClosed {
    margin-right: calc(
      var(--room-drawer-width) + var(--user-drawer-width)
    ) !important;
  }

  .scroll-bottom--bothClosed {
    margin-right: var(--user-drawer-width) !important;
  }

  form {
    left: 0;
  }

  .form--UserDrawerOpened {
    right: var(--user-drawer-width);
  }
}
</style>
