<script setup>
import {computed, onMounted, ref, onBeforeMount} from 'vue'
import {List} from '../../wailsjs/go/app/App'

const props = defineProps({
  bucket: String,
})

const emit = defineEmits(['disconnect'])

const prefix = ref("")
const folders = ref(null)
const files = ref(null)

const showBack = computed(() => {
  return prefix.value !== "" && prefix.value !== "/"
})

const currentPath = computed(() => {
  return "s3://" + props.bucket + "/" + prefix.value
})

function list(p) {
  folders.value = null
  files.value = null
  prefix.value = p

  List(p).then(result => {
    folders.value = result.folders
    files.value = result.files
  })
}

function back() {
  let target = ""
  const splited = prefix.value.split('/');
  const parts = splited.filter((part) => part !== "");
  if (parts.length <= 1) {
    list(target)
    return
  }

  const needParts = parts.slice(0, parts.length - 1)
  target = needParts.join('/') + "/"
  list(target)
}

function withoutPrefix(v) {
  const splited = v.split('/');
  const parts = splited.filter((part) => part !== "");
  return parts[parts.length - 1];
}

function toSettings() {
  emit('disconnect')
}

onMounted(() => {
  list("")
})
</script>

<template>
  <div class="explorer-container">
    <div class="header">
      <div class="header-path">
        <div class="header-path-icon"></div>
        <div class="header-path-label">{{ currentPath }}</div>
      </div>
      <div class="settings-button" @click="toSettings">
        <div class="settings-icon"></div>
      </div>
    </div>
    <div class="files-list" id="files-list">
      <div class="files-item folder" v-if="showBack" @click="back()">
        <div class="files-item-icon folder"></div>
        <div class="files-item-label">..</div>
      </div>
      <div class="files-item folder" v-for="folder in folders" @click="list(folder)">
        <div class="files-item-icon folder"></div>
        <div class="files-item-label">{{ withoutPrefix(folder) }}</div>
      </div>
      <div class="files-item file" v-for="file in files">
        <div class="files-item-icon file"></div>
        <div class="files-item-label">{{ withoutPrefix(file) }}</div>
      </div>
    </div>
  </div>
</template>
