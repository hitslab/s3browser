<script setup>
import Explorer from './components/Explorer.vue'
import Settings from './components/Settings.vue'
import {onMounted, ref} from "vue";

const bucket = ref("")
const showExplorer = ref(false)

function connected(b) {
  bucket.value = b
  showExplorer.value = true
}

function disconnect(b) {
  bucket.value = ""
  showExplorer.value = false
}

onMounted(() => {
  window.runtime.EventsOn("test-event", (data) => {
    console.log("!!!!test-event", data)
  });
})
</script>

<template>
  <Settings @connected="connected" v-if="!showExplorer"/>
  <Explorer :bucket="bucket" v-if="showExplorer" @disconnect="disconnect"/>
</template>

<style>
</style>
