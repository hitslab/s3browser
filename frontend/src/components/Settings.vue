<script setup>
import {onMounted, reactive, ref} from 'vue'
import {Connect, GetConfig} from '../../wailsjs/go/app/App'
import {config} from "../../wailsjs/go/models";
import Alert from "./Alert.vue";

const emit = defineEmits(['connected'])

const error = ref("")
const alert = ref(null)

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
  alert.value.open();
}
</script>

<template>
  <Alert ref="alert">
    {{ error }}
  </Alert>

  <div class="settings">
    <div class="settings-header">
      S3 connection settings
    </div>
    <div class="settings-fields">
      <div class="input-group field">
        <label for="host">Base url</label>
        <input id="host" type="text" placeholder="Example: http://localhost:9000" v-model="settings.base_url">
      </div>
      <div class="input-group field">
        <label for="region">Region</label>
        <input id="region" type="text" placeholder="Example: us-east-1" v-model="settings.region">
      </div>
      <div class="input-group field">
        <label for="bucket">Bucket</label>
        <input id="bucket" type="text" v-model="settings.bucket">
      </div>
      <div class="input-group field">
        <label for="access-key">Access key</label>
        <input id="access-key" type="text" v-model="settings.access_key">
      </div>
      <div class="input-group field">
        <label for="secret-key">Secret key</label>
        <input id="secret-key" type="password" v-model="settings.secret_key">
      </div>
      <div class="input-group checkbox">
        <input id="path-style" type="checkbox" v-model="settings.path_style">
        <label for="path-style">Use path style endpoint</label>
      </div>
      <div class="input-group checkbox">
        <input id="disable-ssl" type="checkbox" v-model="settings.disable_ssl">
        <label for="disable-ssl">Disable ssl</label>
      </div>
      <div class="input-group">
        <button type="button" @click="connect">Connect</button>
      </div>
    </div>
  </div>
</template>
