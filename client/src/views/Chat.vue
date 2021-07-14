<template>
  <ChatHeader />
  <RoomDrawer />
  <ChatMain>
    <div id="welcome">
      <h1 id="welcome-title">Welcome to chat app</h1>
      <p
        id="welcome-text"
      >You can join a chat room from left panel and once you join can see the online users from right panel</p>
    </div>
  </ChatMain>
  <UserDrawer v-show="currentRoom" />
  <Overlay />
</template>

<script lang='ts'>
import { computed, defineComponent, onBeforeMount, onBeforeUnmount, watch } from 'vue'
import ChatHeader from '../components/Chat/Header.vue'
import RoomDrawer from '../components/Chat/RoomDrawer.vue'
import UserDrawer from '../components/Chat/UserDrawer.vue'
import ChatMain from '../components/Chat/Main.vue'
import Overlay from '../components/Chat/Overlay.vue'
import useChat from '../composables/useChat'
import useChatState from '../composables/useChatState'
import { useRoute } from 'vue-router'

export default defineComponent({
  name: 'Chat',
  components: {
    ChatHeader,
    RoomDrawer,
    UserDrawer,
    ChatMain,
    Overlay,
  },
  setup() {
    const route = useRoute()
    const { connectChat, joinChat, leftChat, ws } = useChat()
    const { currentRoom } = useChatState()

    const roomId = computed(() => route.query.roomId)

    watch(roomId, (value, oldValue) => {
      if (oldValue) {
        leftChat(oldValue)
      }
      if (value) {
        joinChat(value)
      }
    })

    onBeforeMount(() => {
      window.onbeforeunload = (): void => {
        if (currentRoom.value) {
          leftChat(currentRoom.value.id)
        }
      }
    })

    onBeforeMount(async () => {
      await connectChat().catch((err) => console.error(err))
      if (roomId.value) {
        joinChat(roomId.value)
      }
    })

    onBeforeUnmount(() => {
      ws.value?.close(1000)
    })

    return { currentRoom }
  }
})
</script>

<style>
:root {
  --scrollbar--with: 8px;
  --scrollbar-thumb: hsl(210, 2%, 48%);
  --scrollbar-thumb-hover: hsl(210, 2%, 43%);
  --scrollbar-track: hsl(213, 12%, 14%);
  --font-family: Avenir, Helvetica, Arial, sans-serif;
  --header-height: 3.5rem;
  --room-drawer-width: 16.875rem;
  --user-drawer-width: 16.875rem;
  --textarea--height: 2rem;
  --drawer-transition: 300ms cubic-bezier(0.25, 0.8, 0.5, 1);
  --btn-transition: 200ms cubic-bezier(0.4, 0, 0.6, 1);
  --page-padding: 0.75rem;
  --header--z-index: 9998;
  --room-drawer--z-index: 9997;
  --user-drawer--z-index: 9997;
  --overlay--z-index: 9996;
}

@media (prefers-color-scheme: dark) {
  :root {
    --clr-background-lighter-3: hsl(210, 14%, 23%);
    --clr-background-lighter-2: hsl(210, 14%, 21%);
    --clr-background-lighter-1: hsl(210, 14%, 19%);
    --clr-background: hsl(210, 14%, 17%);
    --clr-background-darker-1: hsl(210, 14%, 15%);
    --clr-background-darker-2: hsl(210, 14%, 13%);
    --clr-background-darker-3: hsl(210, 14%, 11%);

    --clr-foreground-lighter-3: hsl(0, 0%, 86%);
    --clr-foreground-lighter-2: hsl(0, 0%, 84%);
    --clr-foreground-lighter-1: hsl(0, 0%, 82%);
    --clr-foreground: hsl(0, 0%, 80%);
    --clr-foreground-darker-1: hsl(0, 0%, 78%);
    --clr-foreground-darker-2: hsl(0, 0%, 76%);
    --clr-foreground-darker-3: hsl(0, 0%, 74%);

    --gradient-dot: hsl(256, 33%, 70%);
    --gradient-line: hsl(257, 20%, 85%);
    --gradient-bg: hsl(210, 14%, 17%);
  }
}

@media (prefers-color-scheme: light) {
  :root {
    --clr-background-lighter-3: hsl(0, 0%, 86%);
    --clr-background-lighter-2: hsl(0, 0%, 84%);
    --clr-background-lighter-1: hsl(0, 0%, 82%);
    --clr-background: hsl(0, 0%, 80%);
    --clr-background-darker-1: hsl(0, 0%, 78%);
    --clr-background-darker-2: hsl(0, 0%, 76%);
    --clr-background-darker-3: hsl(0, 0%, 74%);

    --clr-foreground-lighter-3: hsl(210, 14%, 23%);
    --clr-foreground-lighter-2: hsl(210, 14%, 21%);
    --clr-foreground-lighter-1: hsl(210, 14%, 19%);
    --clr-foreground: hsl(210, 14%, 17%);
    --clr-foreground-darker-1: hsl(210, 14%, 15%);
    --clr-foreground-darker-2: hsl(210, 14%, 13%);
    --clr-foreground-darker-3: hsl(210, 14%, 11%);

    --gradient-dot: hsl(256, 33%, 70%);
    --gradient-line: hsl(210, 14%, 17%);
    --gradient-bg: hsl(257, 20%, 85%);
  }
}

* {
  scrollbar-width: thin;
  scrollbar-color: var(--scrollbar-thumb) var(--scrollbar-track);
}

::-webkit-scrollbar {
  width: var(--scrollbar--with);
}

::-webkit-scrollbar-track {
  background-color: var(--scrollbar-track);
}

::-webkit-scrollbar-thumb {
  background-color: var(--scrollbar-thumb);
}

::-webkit-scrollbar-thumb:hover {
  background-color: var(--scrollbar-thumb-hover);
}

body {
  font-family: var(--font-family);
  color: var(--clr-foreground);
  background-color: var(--clr-background);
}

body::before {
  content: " ";
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  will-change: transform;
  z-index: -1;

  background: radial-gradient(var(--gradient-dot) 3px, transparent 4px),
    radial-gradient(var(--gradient-dot) 3px, transparent 4px),
    linear-gradient(var(--gradient-bg) 4px, transparent 0),
    linear-gradient(
      45deg,
      transparent 74px,
      transparent 75px,
      var(--gradient-line) 75px,
      var(--gradient-line) 76px,
      transparent 77px,
      transparent 109px
    ),
    linear-gradient(
      -45deg,
      transparent 75px,
      transparent 76px,
      var(--gradient-line) 76px,
      var(--gradient-line) 77px,
      transparent 78px,
      transparent 109px
    ),
    var(--gradient-bg);
  background-size: 109px 109px, 109px 109px, 100% 6px, 109px 109px, 109px 109px;
  background-position: 54px 55px, 0px 0px, 0px 0px, 0px 0px, 0px 0px;
}

#welcome {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  /* opacity: 0.75;
  background-color: var(--clr-background-darker-3); */
}

#welcome-title {
  padding: 0.25em 0.5em;
  opacity: 0.75;
  background-color: var(--clr-background-darker-3);
}

#welcome-text {
  padding: 0.25em 0.5em;
  opacity: 0.75;
  background-color: var(--clr-background-darker-3);
}
</style>
