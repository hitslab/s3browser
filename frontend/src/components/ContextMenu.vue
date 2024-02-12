<script setup>
import {computed, ref, nextTick} from "vue";

const show = ref(false)
const left = ref(0)
const top = ref(0)

const context = ref(null);

const style = computed(() => {
  return {
    top: top.value + 'px',
    left: left.value + 'px',
    visibility: show.value === false ? 'hidden' : ''
  };
})

function close() {
  show.value = false;
  left.value = 0;
  top.value = 0;
}

function open(evt) {
  let maxLeft = window.innerWidth - context.value.clientWidth - 10
  let leftValue = evt.pageX || evt.clientX
  if (leftValue > maxLeft) {
    leftValue = maxLeft
  }
  left.value = leftValue;

  let maxTop = window.innerHeight - context.value.clientHeight - 10
  let topValue = evt.pageY || evt.clientY;
  if (topValue > maxTop) {
    topValue = maxTop
  }

  top.value = topValue;

  show.value = true

  nextTick(() => context.value.focus());
}

defineExpose({
  open
})
</script>

<template>
  <div class="context-menu" ref="context" :style="style" tabindex="0" @blur="close">
    <slot></slot>
  </div>
</template>

<style scoped>

</style>