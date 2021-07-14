<template>
  <header>
    <nav>
      <button class="toggle-room-drawer" aria-label="toggle room drawer" @click="toggleRoomDrawer">
        <svg style="width:24px;height:24px" viewBox="0 0 24 24">
          <path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z" />
        </svg>
      </button>
      <router-link class="brand" :to="{ name: 'Chat' }">Chat App</router-link>
      <router-link class="user" to="/me">
        <svg style="width:12px;height:12px;margin-right:0.25em" viewBox="0 0 24 24">
          <path
            :fill="isOnline ? 'limegreen' : 'indianred'"
            d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"
          />
        </svg>
        {{ me.username }}
      </router-link>
      <button
        v-show="currentRoom"
        class="toggle-user-drawer"
        aria-label="toggle user drawer"
        @click="toggleUserDrawer"
      >
        <svg style="width:24px;height:24px" viewBox="0 0 24 24">
          <path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z" />
        </svg>
      </button>
    </nav>
  </header>
</template>

<script lang='ts'>
import { defineComponent } from 'vue'
import useChatState from '../../composables/useChatState';
import useDrawer from '../../composables/useDrawer';
import useOnline from '../../composables/useOnline';

export default defineComponent({
  name: 'ChatHeader',
  setup() {
    const { toggleRoomDrawer, toggleUserDrawer } = useDrawer()
    const { isOnline } = useOnline()
    const { currentRoom, me } = useChatState()

    return { toggleRoomDrawer, toggleUserDrawer, isOnline, currentRoom, me }
  }
})
</script>

<style scoped>
header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: var(--header-height);
  z-index: var(--header--z-index);
  /* opacity: 0.85; */
  background-color: var(--clr-background);
}

header > nav {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 var(--page-padding);
  box-shadow: 0px 0px 5px -2px var(--clr-background-darker-3);
}

.brand {
  margin-left: 0.5rem;
}

.user {
  display: flex;
  align-items: center;
  margin-right: 0.5rem;
  margin-left: auto;
}

.brand,
.user {
  position: relative;
  color: inherit;
  text-decoration: none;
  white-space: nowrap;
  overflow: hidden;
}

.user:is(:hover, :focus-visible) {
  text-decoration: underline;
}

.toggle-room-drawer,
.toggle-user-drawer {
  position: relative;
  display: flex;
  align-items: center;
  justify-items: center;
  aspect-ratio: 1;
  border: 0;
  border-radius: 50%;
  cursor: pointer;
  color: inherit;
  background-color: inherit;
}

:is(.toggle-room-drawer, .toggle-user-drawer)::before {
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

:is(.toggle-room-drawer, .toggle-user-drawer):is(:hover, :focus-visible)::before {
  opacity: 0.14;
}
</style>
