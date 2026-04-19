<template>
  <transition name="result-slide">
    <div v-if="result" class="result-wrap glass" :class="{ 'copy-flash': flashing }">
      <!-- header strip -->
      <div class="result-header">
        <span class="status-dot" />
        <span class="result-label">200 OK — Short URL Generated</span>
        <span class="created-at">{{ timeAgo }}</span>
      </div>

      <div class="result-body">
        <!-- URL display -->
        <div class="url-display-block">
          <p class="original-url" :title="result.original_url">
            <span class="field-label">ORIGINAL</span>
            <span class="url-text muted">{{ truncate(result.original_url, 60) }}</span>
          </p>
          <div class="short-url-row">
            <span class="field-label">SHORT</span>
            <a :href="result.short_url" target="_blank" rel="noopener" class="short-url-link">
              {{ result.short_url }}
            </a>
          </div>
        </div>

        <!-- QR Code -->
        <div class="qr-block">
          <canvas ref="qrCanvas" class="qr-canvas" width="96" height="96" />
          <p class="qr-label">SCAN</p>
        </div>
      </div>

      <!-- actions -->
      <div class="result-actions">
        <button class="action-btn copy-btn" @click="copy">
          <span v-if="!copied">⎘ COPY LINK</span>
          <span v-else class="copied-text">✓ COPIED</span>
        </button>
        <button class="action-btn qr-dl-btn" @click="downloadQR">
          ↓ QR CODE
        </button>
        <button class="action-btn clear-btn" @click="$emit('clear')">
          ✕ CLEAR
        </button>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import QRCode from 'qrcode'

const props = defineProps({
  result: { type: Object, default: null },
})
defineEmits(['clear'])

const qrCanvas = ref(null)
const copied   = ref(false)
const flashing = ref(false)
const createdTime = ref(null)

watch(() => props.result, async (val) => {
  if (!val) return
  createdTime.value = new Date()
  await new Promise(r => setTimeout(r, 50)) // wait for DOM
  if (qrCanvas.value) {
    await QRCode.toCanvas(qrCanvas.value, val.short_url, {
      width: 96,
      margin: 1,
      color: { dark: '#00ffe0', light: '#06060d' },
    })
  }
})

const timeAgo = computed(() => {
  if (!createdTime.value) return ''
  return createdTime.value.toLocaleTimeString('zh-TW', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
})

function truncate(str, max) {
  return str.length > max ? str.slice(0, max) + '…' : str
}

async function copy() {
  try {
    await navigator.clipboard.writeText(props.result.short_url)
    copied.value   = true
    flashing.value = true
    setTimeout(() => { copied.value   = false }, 2200)
    setTimeout(() => { flashing.value = false }, 600)
  } catch {
    // fallback
    const ta = document.createElement('textarea')
    ta.value = props.result.short_url
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    ta.remove()
    copied.value = true
    setTimeout(() => { copied.value = false }, 2200)
  }
}

function downloadQR() {
  if (!qrCanvas.value) return
  const link = document.createElement('a')
  link.download = 'qr-code.png'
  link.href = qrCanvas.value.toDataURL()
  link.click()
}
</script>

<style scoped>
.result-wrap {
  overflow: hidden;
  animation: slide-up 0.4s ease both;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.result-wrap.copy-flash {
  animation: copy-flash 0.5s ease forwards;
}

/* header */
.result-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  border-bottom: 1px solid var(--border);
  background: rgba(0,255,224,0.04);
}

.status-dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: var(--acid);
  box-shadow: 0 0 8px var(--acid);
  animation: pulse-border 2s ease infinite;
  flex-shrink: 0;
}

.result-label {
  font-size: 11px;
  color: var(--acid-dim);
  letter-spacing: 0.08em;
  flex: 1;
}

.created-at {
  font-size: 10px;
  color: var(--text-muted);
}

/* body */
.result-body {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 20px 18px;
}

.url-display-block { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 10px; }

.original-url, .short-url-row {
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.field-label {
  font-size: 9px;
  letter-spacing: 0.14em;
  color: var(--text-muted);
  flex-shrink: 0;
  min-width: 54px;
}

.url-text {
  font-size: 12px;
  color: var(--text-dim);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.short-url-link {
  color: var(--cyan);
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  letter-spacing: 0.02em;
  word-break: break-all;
  transition: color 0.15s, text-shadow 0.15s;
}
.short-url-link:hover {
  color: #fff;
  text-shadow: var(--glow-cyan);
}

/* QR */
.qr-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.qr-canvas {
  border: 1px solid var(--border);
  border-radius: 4px;
  display: block;
}

.qr-label {
  font-size: 9px;
  letter-spacing: 0.2em;
  color: var(--text-muted);
}

/* actions */
.result-actions {
  display: flex;
  gap: 8px;
  padding: 12px 18px;
  border-top: 1px solid var(--border);
  flex-wrap: wrap;
}

.action-btn {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-dim);
  font-family: var(--mono);
  font-size: 11px;
  letter-spacing: 0.1em;
  padding: 7px 14px;
  border-radius: 3px;
  cursor: pointer;
  transition: all 0.15s;
}

.copy-btn:hover  { border-color: var(--acid); color: var(--acid); box-shadow: var(--glow-acid); }
.qr-dl-btn:hover { border-color: var(--cyan-dim); color: var(--cyan); }
.clear-btn:hover { border-color: rgba(255,62,108,0.5); color: var(--red); }

.copied-text { color: var(--acid); }

/* transition */
.result-slide-enter-active { transition: all 0.35s cubic-bezier(0.16, 1, 0.3, 1); }
.result-slide-leave-active { transition: all 0.2s ease; }
.result-slide-enter-from  { opacity: 0; transform: translateY(16px) scale(0.98); }
.result-slide-leave-to    { opacity: 0; transform: translateY(-8px); }

@media (max-width: 640px) {
  .qr-block { display: none; }
  .result-body { padding: 16px 14px; }
  .short-url-link { font-size: 13px; }
  .result-actions { gap: 6px; }
  .action-btn { font-size: 10px; padding: 6px 10px; }
}
</style>
