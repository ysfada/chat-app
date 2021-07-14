import _debounce from 'lodash/debounce';
import { onBeforeUnmount, ref } from 'vue';

function makeOnline() {
  let usedCount = 0;

  const isOnline = ref(true);

  const onlineHandler = () => {
    isOnline.value = true;
    // location.reload()
  };
  const offlineHandler = () => {
    isOnline.value = false;
  };

  return () => {
    function init() {
      window.addEventListener('online', onlineHandler, false);
      window.addEventListener('offline', offlineHandler, false);
    }

    if (usedCount === 0) {
      init();
    }
    usedCount++;

    isOnline.value = window.navigator ? window.navigator.onLine : true;

    onBeforeUnmount(() => {
      usedCount--;
      if (usedCount === 0) {
        window.removeEventListener('online', onlineHandler);
        window.removeEventListener('offline', offlineHandler);
      }
    });

    return { isOnline };
  };
}

const useOnline = makeOnline();
export default useOnline;
