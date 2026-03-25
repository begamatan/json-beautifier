<script setup lang="ts">
import { ref } from 'vue'
import * as api from '../api'
import type { ApiError } from '../api'

const input = ref('')
const output = ref('')
const indent = ref<2 | 4>(2)
const errorMsg = ref('')
const successMsg = ref('')
const isLoading = ref(false)

function clearMessages() {
  errorMsg.value = ''
  successMsg.value = ''
}

async function runAction(action: () => Promise<void>) {
  clearMessages()
  isLoading.value = true
  try {
    await action()
  } catch (err: unknown) {
    const apiErr = err as ApiError
    errorMsg.value = apiErr?.message ?? 'Unexpected error'
  } finally {
    isLoading.value = false
  }
}

async function onBeautify() {
  await runAction(async () => {
    const res = await api.beautify(input.value, indent.value)
    output.value = res.result
  })
}

async function onMinify() {
  await runAction(async () => {
    const res = await api.minify(input.value)
    output.value = res.result
  })
}

async function onValidate() {
  await runAction(async () => {
    const res = await api.validate(input.value)
    if (res.valid) {
      successMsg.value = res.message
    } else {
      errorMsg.value = res.message
    }
  })
}

async function onCopy() {
  if (!output.value) return
  await navigator.clipboard.writeText(output.value)
  successMsg.value = 'Copied to clipboard!'
  setTimeout(() => {
    if (successMsg.value === 'Copied to clipboard!') successMsg.value = ''
  }, 2000)
}

function onDownload() {
  if (!output.value) return
  const blob = new Blob([output.value], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'output.json'
  a.click()
  URL.revokeObjectURL(url)
}

function onClear() {
  input.value = ''
  output.value = ''
  clearMessages()
}
</script>

<template>
  <div class="tool-container">
    <h1 class="title">JSON Beautifier</h1>

    <div class="toolbar">
      <label class="indent-label">
        Indent:
        <select v-model="indent" data-testid="indent-select">
          <option :value="2">2 spaces</option>
          <option :value="4">4 spaces</option>
        </select>
      </label>

      <div class="actions">
        <button data-testid="btn-beautify" :disabled="isLoading" @click="onBeautify">Beautify</button>
        <button data-testid="btn-minify" :disabled="isLoading" @click="onMinify">Minify</button>
        <button data-testid="btn-validate" :disabled="isLoading" @click="onValidate">Validate</button>
        <button data-testid="btn-clear" @click="onClear">Clear</button>
      </div>
    </div>

    <div v-if="errorMsg" class="message error" data-testid="error-msg">{{ errorMsg }}</div>
    <div v-if="successMsg" class="message success" data-testid="success-msg">{{ successMsg }}</div>

    <div class="panels">
      <div class="panel">
        <label class="panel-label">Input</label>
        <textarea
          v-model="input"
          placeholder="Paste your JSON here…"
          data-testid="input-area"
          spellcheck="false"
        />
      </div>

      <div class="panel">
        <label class="panel-label">Output</label>
        <textarea
          :value="output"
          readonly
          placeholder="Output will appear here…"
          data-testid="output-area"
          spellcheck="false"
        />
        <div class="output-actions">
          <button data-testid="btn-copy" :disabled="!output" @click="onCopy">Copy</button>
          <button data-testid="btn-download" :disabled="!output" @click="onDownload">Download</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tool-container {
  max-width: 1100px;
  margin: 2rem auto;
  padding: 0 1rem;
  font-family: system-ui, sans-serif;
}

.title {
  text-align: center;
  margin-bottom: 1.5rem;
  font-size: 2rem;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
  margin-bottom: 1rem;
}

.indent-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.indent-label select {
  padding: 0.3rem 0.5rem;
  border-radius: 4px;
  border: 1px solid #ccc;
}

.actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

button {
  padding: 0.4rem 0.9rem;
  border-radius: 4px;
  border: 1px solid #aaa;
  cursor: pointer;
  background: #4f86f7;
  color: white;
  font-size: 0.9rem;
  transition: opacity 0.15s;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

button:not(:disabled):hover {
  opacity: 0.85;
}

.message {
  padding: 0.5rem 0.75rem;
  border-radius: 4px;
  margin-bottom: 0.75rem;
  font-size: 0.9rem;
}

.error {
  background: #fde8e8;
  color: #c0392b;
  border: 1px solid #f5a5a5;
}

.success {
  background: #e6f9ee;
  color: #1e7c3c;
  border: 1px solid #9dd9b2;
}

.panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

@media (max-width: 700px) {
  .panels {
    grid-template-columns: 1fr;
  }
}

.panel {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.panel-label {
  font-weight: 600;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #555;
}

textarea {
  width: 100%;
  min-height: 350px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 0.875rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 0.75rem;
  resize: vertical;
  box-sizing: border-box;
}

.output-actions {
  display: flex;
  gap: 0.5rem;
}
</style>
