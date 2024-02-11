<script setup>
import {onMounted, reactive, ref} from 'vue'
import {Connect, GetConfig} from '../../wailsjs/go/app/App'
import {config} from "../../wailsjs/go/models";
import Alert from "./Alert.vue";

const emit = defineEmits(['connected'])

const error = ref("")

const settings = reactive({
  base_url: "",
  region: "",
  bucket: "",
  access_key: "",
  secret_key: "",
  path_style: false,
  disable_ssl: false
})

function connect() {
  let s = new config.S3Settings();
  s.base_url = settings.base_url
  s.region = settings.region
  s.bucket = settings.bucket
  s.access_key = settings.access_key
  s.secret_key = settings.secret_key
  s.path_style = settings.path_style
  s.disable_ssl = settings.disable_ssl

  Connect(s).then(() => {
    emit('connected', s.bucket)
  }).catch((error) => {
    showError(error)
  });
}

onMounted(() => {
  GetConfig().then((result) => {
    settings.base_url = result.s3.base_url
    settings.region = result.s3.region
    settings.bucket = result.s3.bucket
    settings.access_key = result.s3.access_key
    settings.secret_key = result.s3.secret_key
    settings.path_style = result.s3.path_style
    settings.disable_ssl = result.s3.disable_ssl
  }).catch((error) => {
    showError(error)
  });
})

function showError(msg) {
  error.value = msg
  setTimeout(() => {error.value = ""}, 5000);
}
</script>

<template>
  <Alert type="danger" v-if="error.length">{{ error }}</Alert>

  <div class="settings-container">
    <div class="headers">S3 connection settings</div>
    <div class="fields">
      <div class="input-group">
        <label for="host" class="text-label">Base url</label>
        <input v-model="settings.base_url" id="host" type="text" class="text-input"
               placeholder="Example: http://localhost:9000">
      </div>
      <div class="input-group">
        <label for="region" class="text-label">Region</label>
        <input v-model="settings.region" id="region" type="text" class="text-input" placeholder="Example: us-east-1">
      </div>
      <div class="input-group">
        <label for="bucket" class="text-label">Bucket</label>
        <input v-model="settings.bucket" id="bucket" type="text" class="text-input">
      </div>
      <div class="input-group">
        <label for="access-key" class="text-label">Access key</label>
        <input v-model="settings.access_key" id="access-key" type="text" class="text-input">
      </div>
      <div class="input-group">
        <label for="secret-key" class="text-label">Secret key</label>
        <input v-model="settings.secret_key" id="secret-key" type="password" class="text-input">
      </div>
      <div class="input-group checkbox">
        <input v-model="settings.path_style" id="path-style" type="checkbox">
        <label for="path-style">Use path style endpoint</label>
      </div>
      <div class="input-group checkbox">
        <input v-model="settings.disable_ssl" id="disable-ssl" type="checkbox">
        <label for="disable-ssl">Disable ssl</label>
      </div>
      <div class="input-group">
        <button type="button" @click="connect">Connect</button>
      </div>
    </div>
  </div>
</template>
