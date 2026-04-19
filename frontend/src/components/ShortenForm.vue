<template>
  <section class="shorten-wrap">
    <!-- Label bar -->
    <div class="label-bar">
      <span class="label-tag">POST /shorten</span>
      <span class="cursor-blink">█</span>
    </div>

    <form class="shorten-form glass" @submit.prevent="submit">
      <div class="input-row">
        <span class="input-prefix">https://</span>
        <input
          ref="inputRef"
          v-model="rawUrl"
          class="url-input"
          type="url"
          placeholder="paste your long URL here..."
          autocomplete="off"
          spellcheck="false"
          :disabled="loading"
          @focus="focused = true"
          @blur="focused = false"
          :class="{ focused, error: !!errorMsg }"
        />
        <button class="submit-btn" type="submit" :disabled="loading || !rawUrl">
          <span v-if="!loading" class="btn-text">SHORTEN</span>
          <span v-else class="spinner" />
        </button>
      </div>

      <transition name="err">
        <p v-if="errorMsg" class="error-msg">
          <span class="err-icon">⚠</span> {{ errorMsg }}
        </p>
      </transition>
    </form>
  </section>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['result'])

const rawUrl    = ref('')
const loading   = ref(false)
const errorMsg  = ref('')
const focused   = ref(false)
const inputRef  = ref(null)

async function submit() {
  errorMsg.value = ''
  if (!rawUrl.value.trim()) return

  let url = rawUrl.value.trim()
  if (!/^https?:\/\//i.test(url)) url = 'https://' + url

  loading.value = true
  try {
    const res = await fetch('/shorten', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({ url }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Unknown error')
    emit('result', { short_url: data.short_url, original_url: url })
    rawUrl.value = ''
  } catch (e) {
    errorMsg.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.shorten-wrap {
  animation: slide-up 0.5s ease both;
  animation-delay: 0.3s;
  opacity: 0;
}

.label-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 11px;
  letter-spacing: 0.08em;
  color: var(--text-muted);
}

.label-tag {
  background: rgba(0,255,224,0.08);
  border: 1px solid var(--border);
  color: var(--cyan-dim);
  padding: 2px 8px;
  border-radius: 3px;
  font-size: 10px;
  letter-spacing: 0.1em;
}

.cursor-blink {
  color: var(--cyan);
  animation: blink 1.1s step-end infinite;
  font-size: 11px;
}

.shorten-form {
  padding: 6px;
  transition: border-color 0.25s, box-shadow 0.25s;
}

.shorten-form:has(.url-input.focused) {
  border-color: var(--border-hot);
  box-shadow: var(--glow-cyan);
}

.shorten-form:has(.url-input.error) {
  border-color: rgba(255, 62, 108, 0.5);
  box-shadow: 0 0 20px rgba(255, 62, 108, 0.15);
}

.input-row {
  display: flex;
  align-items: center;
  gap: 0;
}

.input-prefix {
  color: var(--text-muted);
  padding: 0 10px 0 14px;
  font-size: 13px;
  white-space: nowrap;
  user-select: none;
  display: none; /* hidden — full URL typed by user */
}

.url-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--text);
  font-family: var(--mono);
  font-size: 14px;
  padding: 14px 16px;
  caret-color: var(--cyan);
  min-width: 0;
}

.url-input::placeholder { color: var(--text-muted); }
.url-input:disabled     { opacity: 0.5; cursor: not-allowed; }

.submit-btn {
  flex-shrink: 0;
  background: var(--cyan);
  color: #000;
  border: none;
  padding: 12px 24px;
  font-family: var(--mono);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  cursor: pointer;
  border-radius: 3px;
  margin: 6px;
  transition: background 0.15s, box-shadow 0.15s, transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 96px;
  height: 42px;
}

.submit-btn:hover:not(:disabled) {
  background: #33ffe8;
  box-shadow: var(--glow-cyan);
}

.submit-btn:active:not(:disabled) { transform: scale(0.97); }
.submit-btn:disabled { opacity: 0.4; cursor: not-allowed; }

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(0,0,0,0.3);
  border-top-color: #000;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.error-msg {
  color: var(--red);
  font-size: 12px;
  padding: 6px 16px 10px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.err-icon { font-size: 13px; }

/* transitions */
.err-enter-active, .err-leave-active { transition: all 0.2s ease; }
.err-enter-from, .err-leave-to { opacity: 0; transform: translateY(-4px); }

@media (max-width: 640px) {
  .submit-btn { padding: 12px 16px; min-width: 80px; font-size: 11px; }
  .url-input  { font-size: 13px; padding: 12px 12px; }
}
</style>
