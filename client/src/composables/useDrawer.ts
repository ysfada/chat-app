import { ref } from 'vue';
import useBreakpoints from './useBreakpoints';
import useChatState from './useChatState';

const isRoomDrawerOpen = ref(false);
const isUserDrawerOpen = ref(false);

export default function useDrawer() {
  const { widthType } = useBreakpoints();
  const { currentRoom } = useChatState();

  const toggleRoomDrawer = () => {
    if (
      isUserDrawerOpen.value &&
      (widthType.value === 'xs' || widthType.value === 'sm')
    ) {
      isUserDrawerOpen.value = false;
    }
    isRoomDrawerOpen.value = !isRoomDrawerOpen.value;
  };

  const toggleUserDrawer = () => {
    if (
      isRoomDrawerOpen.value &&
      (widthType.value === 'xs' || widthType.value === 'sm')
    ) {
      isRoomDrawerOpen.value = false;
    } else {
      // if (!currentRoom.value) {
      //   isUserDrawerOpen.value = false;
      //   return;
      // }
      isUserDrawerOpen.value = !isUserDrawerOpen.value;
    }
  };

  const setRoomDrawers = () => {
    if (widthType.value === 'xs' || widthType.value === 'sm') {
      isRoomDrawerOpen.value = false;
    } else if (widthType.value === 'md' || widthType.value === 'lg') {
      isRoomDrawerOpen.value = true;
    } else {
      isRoomDrawerOpen.value = true;
    }
  };

  const setUserDrawers = () => {
    if (widthType.value === 'xs' || widthType.value === 'sm') {
      isUserDrawerOpen.value = false;
    } else if (widthType.value === 'md' || widthType.value === 'lg') {
      isUserDrawerOpen.value = false;
    } else {
      // if (!currentRoom.value) {
      //   isUserDrawerOpen.value = false;
      //   return;
      // }
      isUserDrawerOpen.value = true;
    }
  };

  const hideDrawersOnOverlayClick = () => {
    if (widthType.value === 'xs' || widthType.value === 'sm') {
      isRoomDrawerOpen.value = false;
      isUserDrawerOpen.value = false;
    } else if (widthType.value === 'md' || widthType.value === 'lg') {
      isUserDrawerOpen.value = false;
    }
  };

  return {
    isRoomDrawerOpen,
    isUserDrawerOpen,
    toggleRoomDrawer,
    toggleUserDrawer,
    setRoomDrawers,
    setUserDrawers,
    hideDrawersOnOverlayClick,
  };
}
