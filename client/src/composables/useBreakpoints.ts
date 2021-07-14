import _debounce from 'lodash/debounce';
import { computed, onBeforeUnmount, ref } from 'vue';

function makeBreakpoints() {
  let usedCount = 0;

  const windowWidth = ref(window.innerWidth);

  const widthType = computed(() => {
    if (windowWidth.value < 358) return 'xs';
    else if (windowWidth.value > 357 && windowWidth.value < 608) return 'sm';
    else if (windowWidth.value > 607 && windowWidth.value < 968) return 'md';
    else if (windowWidth.value > 967 && windowWidth.value < 1272) return 'lg';
    else if (windowWidth.value > 1271 && windowWidth.value < 1912) return 'xl';
    else return 'xxl'; // xxl > 1911
  });

  const width = computed(() => windowWidth.value);

  const onWidthChange = _debounce(() => {
    windowWidth.value = window.innerWidth;
  }, 100);

  return () => {
    function init() {
      window.addEventListener('resize', onWidthChange);
    }

    if (usedCount === 0) {
      init();
    }
    usedCount++;

    onBeforeUnmount(() => {
      usedCount--;
      if (usedCount === 0) {
        window.removeEventListener('resize', onWidthChange);
      }
    });

    return { width, widthType };
  };
}

const useBreakpoints = makeBreakpoints();
export default useBreakpoints;
