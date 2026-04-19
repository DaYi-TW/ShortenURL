<template>
  <div class="stats-row">
    <div class="stat-card glass" :class="{ loading: loadingStats }">
      <div class="stat-header">
        <span class="stat-icon cyan">◈</span>
        <span class="stat-title">TOTAL LINKS</span>
      </div>
      <div class="stat-value" :key="total">
        <AnimatedNumber :to="total" />
      </div>
      <div class="stat-bar">
        <div class="stat-bar-fill cyan-fill" :style="{ width: total > 0 ? '100%' : '0%' }" />
      </div>
    </div>

    <div class="stat-card glass" :class="{ loading: loadingStats }">
      <div class="stat-header">
        <span class="stat-icon acid">◈</span>
        <span class="stat-title">TODAY</span>
      </div>
      <div class="stat-value acid" :key="today">
        <AnimatedNumber :to="today" />
      </div>
      <div class="stat-bar">
        <div
          class="stat-bar-fill acid-fill"
          :style="{ width: total > 0 ? Math.min((today / total) * 100, 100) + '%' : '0%' }"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, defineComponent, h, watch } from 'vue'

// ── tiny animated number counter ─────────────────────────────
const AnimatedNumber = defineComponent({
  props: { to: { type: Number, default: 0 } },
  setup(props) {
    const displayed = ref(0)
    let raf = null
    watch(() => props.to, (target) => {
      const start = displayed.value
      const delta = target - start
      if (delta === 0) return
      const duration = 600
      const startTime = performance.now()
      const step = (now) => {
        const t = Math.min((now - startTime) / duration, 1)
        const ease = 1 - Math.pow(1 - t, 3)
        displayed.value = Math.round(start + delta * ease)
        if (t < 1) raf = requestAnimationFrame(step)
      }
      cancelAnimationFrame(raf)
      raf = requestAnimationFrame(step)
    }, { immediate: true })
    return () => h('span', displayed.value.toLocaleString())
  },
})

const total        = ref(0)
const today        = ref(0)
const loadingStats = ref(true)
let   pollTimer    = null

async function fetchStats() {
  try {
    const [r1, r2] = await Promise.all([
      fetch('/stats'),
      fetch('/stats/today'),
    ])
    const d1 = await r1.json()
    const d2 = await r2.json()
    total.value = d1.shortened_url_count   ?? 0
    today.value = d2.shortened_url_count_today ?? 0
  } catch { /* silent */ } finally {
    loadingStats.value = false
  }
}

onMounted(() => {
  fetchStats()
  pollTimer = setInterval(fetchStats, 15_000)
})

onUnmounted(() => clearInterval(pollTimer))

// expose for parent to trigger a manual refresh
defineExpose({ fetchStats })
</script>

<style scoped>
.stats-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  animation: slide-up 0.5s ease both;
  animation-delay: 0.45s;
  opacity: 0;
}

.stat-card {
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: border-color 0.3s, box-shadow 0.3s;
  position: relative;
  overflow: hidden;
}

.stat-card.loading::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(0,255,224,0.04) 50%, transparent 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  from { background-position: -200% 0; }
  to   { background-position: 200% 0; }
}

.stat-card:hover {
  border-color: var(--border-hot);
  box-shadow: var(--glow-cyan);
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 7px;
}

.stat-icon {
  font-size: 11px;
}
.stat-icon.cyan  { color: var(--cyan); text-shadow: 0 0 8px var(--cyan); }
.stat-icon.acid  { color: var(--acid); text-shadow: 0 0 8px var(--acid); }

.stat-title {
  font-size: 9px;
  letter-spacing: 0.2em;
  color: var(--text-muted);
}

.stat-value {
  font-family: var(--sans);
  font-size: 32px;
  font-weight: 800;
  color: var(--cyan);
  line-height: 1;
  animation: count-up 0.4s ease;
}

.stat-value.acid { color: var(--acid); }

.stat-bar {
  height: 2px;
  background: rgba(255,255,255,0.05);
  border-radius: 1px;
  overflow: hidden;
}

.stat-bar-fill {
  height: 100%;
  border-radius: 1px;
  transition: width 0.8s cubic-bezier(0.16,1,0.3,1);
}

.cyan-fill { background: var(--cyan); box-shadow: 0 0 6px var(--cyan); }
.acid-fill { background: var(--acid); box-shadow: 0 0 6px var(--acid); }

@media (max-width: 480px) {
  .stat-value { font-size: 26px; }
  .stat-card  { padding: 14px 16px; }
}
</style>
