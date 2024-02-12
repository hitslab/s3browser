<script setup>
import {computed, ref} from "vue";

const show = ref(false)

const style = computed(() => {
  return {
    display: show.value === false ? 'none' : ''
  };
})

function open() {
  show.value = true
  setTimeout(() => {
    close()
  }, 10000);
}

function close() {
  show.value = false
}

defineExpose({
  open
})
</script>

<template>
  <div class="alert" :style="style">
    <div class="close" @click="close"></div>
    <slot></slot>
  </div>
</template>

<style scoped>

</style>