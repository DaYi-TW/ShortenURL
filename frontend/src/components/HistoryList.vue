<template>
  <section v-if="history.length > 0" class="history-wrap">
    <div class="history-header">
      <span class="section-label">// HISTORY</span>
      <span class="history-count">{{ history.length }} record{{ history.length !== 1 ? 's' : '' }}</span>
      <button class="clear-all-btn" @click="clearAll">CLEAR ALL</button>
    </div>

    <div class="history-list glass">
      <TransitionGroup name="hist-item" tag="div">
        <div
          v-for="item in history"
          :key="item.short_url"
          class="hist-row"
        >
          <div class="hist-info">
            <div class="hist-short-row">
              <span class="hist-code">{{ extractCode(item.short_url) }}</span>
              <a :href="item.short_url" target="_blank" class="hist-short-link">{{ item.short_url }}</a>
            </div>
            <p class="hist-original" :title="item.original_url">
              {{ truncate(item.original_url, 55) }}
            </p>
          </div>
          <div class="hist-time">{{ item.time }}</div>
          <div class="hist-actions">
            <button class="hist-btn" @click="copyItem(item)" :class="{ copied: item._copied }">
              {{ item._copied ? '✓' : '⎘' }}
            </button>
            <button class="hist-btn del" @click="removeItem(item.short_url)">✕</button>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </section>
</template>

<script setup>
import { ref, watch } from 'vue'

const STORAGE_KEY = 'shrt_history'
const MAX_ITEMS   = 20

const history = ref(load())

function load() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]') }
  catch { return [] }
}

function save() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(history.value))
}

function addItem(item) {
  // dedupe by short_url
  history.value = history.value.filter(h => h.short_url !== item.short_url)
  history.value.unshift({ ...item, _copied: false })
  if (history.value.length > MAX_ITEMS) history.value.pop()
  save()
}

function removeItem(shortUrl) {
  history.value = history.value.filter(h => h.short_url !== shortUrl)
  save()
}

function clearAll() {
  history.value = []
  save()
}

async function copyItem(item) {
  try { await navigator.clipboard.writeText(item.short_url) } catch {}
  item._copied = true
  setTimeout(() => { item._copied = false }, 2000)
}

function truncate(str, max) {
  return str && str.length > max ? str.slice(0, max) + '…' : str
}

function extractCode(url) {
  return url.split('/').pop()
}

defineExpose({ addItem })
</script>

<style scoped>
.history-wrap {
  animation: slide-up 0.5s ease both;
  animation-delay: 0.55s;
  opacity: 0;
}

.history-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.section-label {
  font-size: 11px;
  color: var(--text-muted);
  letter-spacing: 0.1em;
  flex: 1;
}

.history-count {
  font-size: 10px;
  color: var(--text-muted);
  background: rgba(255,255,255,0.04);
  padding: 2px 8px;
  border-radius: 20px;
  border: 1px solid var(--border);
}

.clear-all-btn {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-muted);
  font-family: var(--mono);
  font-size: 9px;
  letter-spacing: 0.12em;
  padding: 3px 10px;
  border-radius: 3px;
  cursor: pointer;
  transition: all 0.15s;
}
.clear-all-btn:hover { border-color: rgba(255,62,108,0.5); color: var(--red); }

.history-list {
  overflow: hidden;
}

.hist-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
  transition: background 0.15s;
}
.hist-row:last-child  { border-bottom: none; }
.hist-row:hover       { background: rgba(0,255,224,0.03); }

.hist-info { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px; }

.hist-short-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hist-code {
  font-size: 10px;
  background: rgba(0,255,224,0.08);
  color: var(--cyan-dim);
  padding: 1px 6px;
  border-radius: 2px;
  letter-spacing: 0.05em;
  flex-shrink: 0;
}

.hist-short-link {
  font-size: 12px;
  color: var(--cyan);
  text-decoration: none;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.15s;
}
.hist-short-link:hover { color: #fff; }

.hist-original {
  font-size: 11px;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.hist-time {
  font-size: 10px;
  color: var(--text-muted);
  white-space: nowrap;
  flex-shrink: 0;
}

.hist-actions { display: flex; gap: 6px; flex-shrink: 0; }

.hist-btn {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-muted);
  font-family: var(--mono);
  font-size: 11px;
  padding: 4px 8px;
  border-radius: 3px;
  cursor: pointer;
  transition: all 0.15s;
  min-width: 28px;
}
.hist-btn:hover       { border-color: var(--cyan-dim); color: var(--cyan); }
.hist-btn.del:hover   { border-color: rgba(255,62,108,0.5); color: var(--red); }
.hist-btn.copied      { border-color: var(--acid); color: var(--acid); }

/* list transition */
.hist-item-enter-active { transition: all 0.3s cubic-bezier(0.16,1,0.3,1); }
.hist-item-leave-active { transition: all 0.2s ease; position: absolute; width: 100%; }
.hist-item-enter-from   { opacity: 0; transform: translateX(-12px); }
.hist-item-leave-to     { opacity: 0; transform: translateX(12px); }
.hist-item-move         { transition: transform 0.3s ease; }

@media (max-width: 640px) {
  .hist-time    { display: none; }
  .hist-original{ display: none; }
  .hist-row     { padding: 10px 12px; }
}
</style>
