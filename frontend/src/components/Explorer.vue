<script setup>
import {computed, onMounted, ref} from 'vue'
import {List} from '../../wailsjs/go/app/App'
import ContextMenu from "./ContextMenu.vue";
import {Download} from "../../wailsjs/go/app/App.js";
import Alert from "./Alert.vue";

const props = defineProps({
  bucket: String,
})

const emit = defineEmits(['disconnect'])

const prefix = ref("")
const folders = ref(null)
const files = ref(null)
const menu = ref(null)
const contextItemText = ref("Download")
const contextItemPrefix = ref(null)
const error = ref("")
const alert = ref(null)
const isDownloading = ref(false)
const downloadPercent = ref(0)

const showBack = computed(() => {
  return prefix.value !== "" && prefix.value !== "/"
})

const currentPath = computed(() => {
  return "s3://" + props.bucket + "/" + prefix.value
})

const processStyle = computed(() => {
  return {
    width: downloadPercent.value + '%',
  };
})

function list(p) {
  folders.value = null
  files.value = null
  prefix.value = p

  List(p).then(result => {
    folders.value = result.folders
    files.value = result.files
  }).catch(error => {
    showError(error)
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

function download() {
  isDownloading.value = true
  Download(contextItemPrefix.value).then(result => {
    isDownloading.value = false
    console.log(result)
  }).catch(error => {
    isDownloading.value = false
    showError(error)
  })
}

function openContextMenu(e, type, prefix) {
  if (type === "file") {
    contextItemText.value = "Download file"
  } else if (type === "folder") {
    contextItemText.value = "Download folder"
  }

  contextItemPrefix.value = prefix
  e.stopPropagation()
  e.preventDefault()
  menu.value.open(e);
}

function showError(msg) {
  error.value = msg
  alert.value.open();
}

onMounted(() => {
  list("")

  window.runtime.EventsOn("download-process", (percent) => {
    downloadPercent.value = percent
  });
})
</script>

<template>
  <Alert ref="alert">
    {{ error }}
  </Alert>

  <ContextMenu ref="menu">
    <div class="context-item" @click="download">
      <div class="context-item-icon download"></div>
      <div class="context-item-text">{{ contextItemText }}</div>
    </div>
  </ContextMenu>

  <div class="explorer">
    <div class="explorer-header">
      <div class="explorer-header-path">{{ currentPath }}</div>
      <div class="explorer-header-controls">
        <div class="icon-button settings" @click="toSettings"></div>
      </div>
    </div>
    <div class="files-list">
      <div class="files-item folder" v-if="showBack" @click="back()">
        <div class="files-item-icon folder"></div>
        <div class="files-item-name">..</div>
        <div class="icon-button more"></div>
      </div>
      <div class="files-item folder" v-for="folder in folders" @click="list(folder)" @contextmenu="openContextMenu($event, 'folder', folder)">
        <div class="files-item-icon folder"></div>
        <div class="files-item-name">{{ withoutPrefix(folder) }}</div>
        <div class="icon-button more" @click="openContextMenu($event, 'folder', folder)"></div>
      </div>
      <div class="files-item file" v-for="file in files" @contextmenu="openContextMenu($event, 'file', file)">
        <div class="files-item-icon file"></div>
        <div class="files-item-name">{{ withoutPrefix(file) }}</div>
        <div class="icon-button more" @click="openContextMenu($event, 'file', file)"></div>
      </div>
    </div>
  </div>

  <div class="progress" v-if="isDownloading">
    <div class="progress-label">
      Downloading
    </div>
    <div class="progress-bar">
      <div class="progress-fill" :style="processStyle"></div>
    </div>
    <div class="progress-percent">
      {{ downloadPercent }}%
    </div>
  </div>
</template>
