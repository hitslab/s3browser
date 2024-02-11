<script setup>
import {computed} from "vue";

const props = defineProps({
  type: String,
})

const divClass = computed(() => {
  if (props.type === "danger") {
    return "alert danger"
  }

  if (props.type === "warning") {
    return "alert warning"
  }

  if (props.type === "success") {
    return "alert success"
  }

  return "alert"
})

</script>

<template>
  <div :class="divClass">
    <slot></slot>
  </div>
</template>

<style scoped>

</style>