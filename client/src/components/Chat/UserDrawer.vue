<template>
  <aside :class="{
    open: isUserDrawerOpen,
  }">
    <div class="user-title">
      <h1>online users</h1>
      <div class="user-search-box">
        <input
          id="user-search-input"
          v-model="searchInput"
          type="search"
          name="userId"
          placeholder="search user"
        />
      </div>
    </div>

    <ul class="user-list">
      <li
        class="user-item"
        :style="{ backgroundColor: $route.query.roomId === user.id ? 'var(--clr-background-darker-1)' : '' }"
        v-for="user in orderedUsers"
        :key="user.id"
      >
        <div class="user-icon" :style="{ backgroundColor: stringToColor(user.id) }">
          <span class="user-letter">{{ user.username[0] }}</span>
        </div>
        <router-link class="user-link" :to="`/chat?roomId=${user.id}`">{{ user.username }}</router-link>
      </li>
    </ul>
  </aside>
</template>

<script lang='ts'>
import { computed, defineComponent, onBeforeMount, ref, watch } from 'vue'
import useBreakpoints from '../../composables/useBreakpoints';
import useChatState from '../../composables/useChatState';
import useDrawer from '../../composables/useDrawer';
import stringToColor from '../../utils/stringToColor';

export default defineComponent({
  name: 'UserDrawer',
  setup() {
    const { width } = useBreakpoints()
    const { isUserDrawerOpen, setUserDrawers } = useDrawer()
    const { users } = useChatState()

    const searchInput = ref('')

    const orderedUsers = computed(() => {
      return users.value
        .filter((user) =>
          user.username
            .toLocaleLowerCase()
            .includes(searchInput.value.toLocaleLowerCase())
        )
        .sort((a, b) => (a.username > b.username ? 1 : -1))
    })

    onBeforeMount(setUserDrawers)

    watch(width, setUserDrawers)

    return { isUserDrawerOpen, stringToColor, searchInput, orderedUsers }
  }
})
</script>

<style scoped>
.open {
  transform: translateX(0);
}

aside {
  position: fixed;
  top: var(--header-height);
  bottom: 0;
  right: 0;
  width: var(--room-drawer-width);
  height: calc(100vh - var(--header-height));
  overflow-y: overlay;
  z-index: var(--user-drawer--z-index);
  opacity: 0.95;
  will-change: transform, overflow-y;
  transform: translateX(var(--room-drawer-width));
  transition: transform var(--drawer-transition);
  background-color: var(--clr-background);
}

@media (hover: hover) {
  aside {
    /* visibility: hidden; */
    overflow-y: hidden;
  }
}

aside:hover,
aside:focus,
aside:focus-visible {
  /* visibility: visible; */
  overflow-y: overlay;
}

.user-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: sticky;
  top: 0;
  height: 5rem;
  padding: 0 var(--page-padding);
  box-shadow: 0px 0px 5px -2px var(--clr-background-darker-3);
  background-color: inherit;
}

.user-title > h1 {
  margin: 0;
}

#user-search-input {
  width: 100%;
  border: 0;
  border-radius: 5px;
  outline: none;
  font-size: 1.5rem;
  color: inherit;
  background-color: var(--clr-background-lighter-2);
}
#user-search-input:focus {
  box-shadow: 0 0 10px 1px hsl(210, 44%, 40%);
}

.user-list {
  min-height: calc(100% - 80px);
  margin: 0;
  /* visibility: visible; */
  list-style: none;
  padding: 0;
  background-color: inherit;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 0 var(--page-padding);
  transition: box-shadow 50ms ease-in;
}

.user-item:is(:hover, :focus-within) {
  box-shadow: inset 0 -5px 0 0 var(--clr-background-darker-2);
  background-color: var(--clr-background-darker-1);
}

.user-icon {
  padding: 0.5rem 1rem;
  border-radius: 5px;
  font-family: monospace;
  font-size: 1.5em;
  color: hsl(0, 0%, 100%);
}

.user-link {
  flex: 1;
  padding: 1rem 0 1rem 0.5rem;
  font-size: 1.2rem;
  text-decoration: none;
  color: inherit;
}
</style>
