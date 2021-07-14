<template>
  <aside :class="{
    open: isRoomDrawerOpen,
  }">
    <div class="room-title">
      <h1>chat rooms</h1>
      <div class="room-search-box">
        <input
          id="room-search-input"
          v-model="searchInput"
          type="search"
          name="roomId"
          placeholder="search room"
        />
      </div>
    </div>

    <ul class="room-list">
      <li
        class="room-item"
        :style="{ backgroundColor: $route.query.roomId === room.id ? 'var(--clr-background-darker-1)' : '' }"
        v-for="room in orderedRooms"
        :key="room.id"
      >
        <div class="room-icon" :style="{ backgroundColor: stringToColor(room.id) }">
          <span class="room-letter">{{ room.name[0] }}</span>
        </div>
        <router-link class="room-link" :to="`/chat?roomId=${room.id}`">{{ room.name }}</router-link>
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
  name: 'RoomDrawer',
  setup() {
    const { width } = useBreakpoints()
    const { isRoomDrawerOpen, setRoomDrawers } = useDrawer()
    const { rooms } = useChatState()

    const searchInput = ref('')

    const orderedRooms = computed(() => {
      return rooms.value
        .filter((room) =>
          room.name
            .toLocaleLowerCase()
            .includes(searchInput.value.toLocaleLowerCase())
        )
        .sort((a, b) => (a.name > b.name ? 1 : -1))
    })

    onBeforeMount(setRoomDrawers)

    watch(width, setRoomDrawers)

    return { isRoomDrawerOpen, stringToColor, searchInput, orderedRooms }
  }
})
</script>

<style scoped>
.open {
  transform: translateX(0);
}

aside {
  position: absolute;
  top: var(--header-height);
  bottom: 0;
  left: 0;
  width: var(--room-drawer-width);
  height: calc(100vh - var(--header-height));
  overflow-y: overlay;
  z-index: var(--room-drawer--z-index);
  opacity: 0.95;
  will-change: transform, overflow-y;
  transform: translateX(calc(var(--room-drawer-width) * -1));
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

.room-title {
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

.room-title > h1 {
  margin: 0;
}

#room-search-input {
  width: 100%;
  border: 0;
  border-radius: 5px;
  outline: none;
  font-size: 1.5rem;
  color: inherit;
  background-color: var(--clr-background-lighter-2);
}
#room-search-input:focus {
  box-shadow: 0 0 10px 1px hsl(210, 44%, 40%);
}

.room-list {
  min-height: calc(100% - 80px);
  margin: 0;
  /* visibility: visible; */
  list-style: none;
  padding: 0;
  background-color: inherit;
}

.room-item {
  display: flex;
  align-items: center;
  padding: 0 var(--page-padding);
  transition: box-shadow 50ms ease-in;
}

.room-item:is(:hover, :focus-within) {
  box-shadow: inset 0 -5px 0 0 var(--clr-background-darker-2);
  background-color: var(--clr-background-darker-1);
}

.room-icon {
  padding: 0.5rem 1rem;
  border-radius: 5px;
  font-family: monospace;
  font-size: 1.5em;
  color: hsl(0, 0%, 100%);
}

.room-link {
  flex: 1;
  padding: 1rem 0 1rem 0.5rem;
  font-size: 1.2rem;
  text-decoration: none;
  color: inherit;
}
</style>
