<template>
  <div class="wrapper">
    <img alt="Vue logo" src="../assets/logo.png" />
    <HelloWorld msg="Hello Vue 3 + TypeScript + Vite" />
  </div>
</template>

<script lang='ts'>
import { defineComponent } from 'vue'
import HelloWorld from '../components/HelloWorld.vue'

export default defineComponent({
  name: 'Home',
  components: {
    HelloWorld
  },
})
</script>

<style scoped>
:root {
  --font-family: Avenir, Helvetica, Arial, sans-serif;
  --clr-primary: hsl(210, 29%, 24%);
  --clr-secondary: hsl(110, 29%, 24%);
}
.wrapper {
  font-family: var(--font-family);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: var(--clr-primary);
}
</style>
