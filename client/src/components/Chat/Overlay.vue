<template>
  <div
    class="overlay"
    :class="{
      'overlay--opacity': showOverlay
    }"
    @click="hideDrawersOnOverlayClick"
  ></div>
</template>

<script lang="ts">import { computed, defineComponent, watch } from "@vue/runtime-core";
import useBreakpoints from "../../composables/useBreakpoints";
import useDrawer from "../../composables/useDrawer";


export default defineComponent({
  name: 'Overlay',
  setup() {
    const { widthType } = useBreakpoints()
    const { isRoomDrawerOpen, isUserDrawerOpen } = useDrawer()
    const { hideDrawersOnOverlayClick } = useDrawer()

    const showOverlay = computed(() => {
      if ((widthType.value === 'xs' || widthType.value === 'sm') && (isRoomDrawerOpen.value || isUserDrawerOpen.value)) {
        return true
      } else if ((widthType.value === 'md' || widthType.value === 'lg') && (isUserDrawerOpen.value)) {
        return true
      }
    })

    return { hideDrawersOnOverlayClick, showOverlay }
  },
})
</script>

<style scoped>
.overlay {
  position: fixed;
  visibility: hidden;
  top: var(--header-height);
  left: 0;
  right: 0;
  bottom: 0;
  height: calc(100vh - var(--header-height));
  width: 100vw;
  z-index: -1;
  opacity: 0;
  transition: opacity var(--drawer-transition);
  background-color: hsl(0, 0%, 13%);
}

.overlay--opacity {
  visibility: visible;
  opacity: 0.46;
  z-index: var(--overlay--z-index);
}
</style>
