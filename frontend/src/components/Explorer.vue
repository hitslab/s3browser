<script setup>
import {computed, onMounted, ref} from 'vue'
import {List, UploadFile, UploadFolder} from '../../wailsjs/go/app/App'
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
const contextItems = ref([])
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

function uploadFile() {
  UploadFile(prefix.value).then(result => {
    list(prefix.value)
  }).catch(error => {
    showError(error)
  })
}

function uploadFolder() {
  UploadFolder(prefix.value).then(result => {
    list(prefix.value)
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

function download(target) {
  isDownloading.value = true
  Download(target).then(result => {
    isDownloading.value = false
    console.log(result)
  }).catch(error => {
    isDownloading.value = false
    showError(error)
  })
}

function openDownloadContext(e, type, prefix) {
  let itemText = "";
  if (type === "file") {
    itemText = "Download file"
  } else if (type === "folder") {
    itemText = "Download folder"
  }

  contextItems.value = [
    {
      text: itemText,
      iconClass: "context-item-icon download",
      target: () => {
        download(prefix)
      },
    }
  ]

  e.stopPropagation()
  e.preventDefault()
  menu.value.open(e);
}

function openUploadContext(e, type, prefix) {
  contextItems.value = [
    {
      text: "Upload file",
      iconClass: "context-item-icon",
      target: () => {
        uploadFile(prefix)
      },
    },
    {
      text: "Upload folder",
      iconClass: "context-item-icon",
      target: () => {
        uploadFolder(prefix)
      },
    }
  ]

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
    <div v-for="contextItem in contextItems" class="context-item" @click="contextItem.target">
      <div :class="contextItem.iconClass"></div>
      <div class="context-item-text">{{ contextItem.text }}</div>
    </div>
  </ContextMenu>

  <div class="explorer">
    <div class="explorer-header">
      <div class="explorer-header-path">{{ currentPath }}</div>
      <div class="explorer-header-controls">
        <div class="icon-button upload" @click="openUploadContext"></div>
        <div class="icon-button settings" @click="toSettings"></div>
      </div>
    </div>
    <div class="files-list">
      <div class="files-item folder" v-if="showBack" @click="back()">
        <div class="files-item-icon folder"></div>
        <div class="files-item-name">..</div>
        <div class="icon-button more"></div>
      </div>
      <div class="files-item folder" v-for="folder in folders" @click="list(folder)" @contextmenu="openDownloadContext($event, 'folder', folder)">
        <div class="files-item-icon folder"></div>
        <div class="files-item-name">{{ withoutPrefix(folder) }}</div>
        <div class="icon-button more" @click="openDownloadContext($event, 'folder', folder)"></div>
      </div>
      <div class="files-item file" v-for="file in files" @contextmenu="openDownloadContext($event, 'file', file)">
        <div class="files-item-icon file"></div>
        <div class="files-item-name">{{ withoutPrefix(file) }}</div>
        <div class="icon-button more" @click="openDownloadContext($event, 'file', file)"></div>
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
