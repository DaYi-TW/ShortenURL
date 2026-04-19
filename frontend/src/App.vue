<template>
  <ParticleBackground />

  <div class="layout">
    <!-- ── HEADER ── -->
    <header class="site-header">
      <div class="logo-block">
        <span class="logo-bracket">[</span>
        <span class="logo-text" aria-label="SHRT">
          <span class="glitch-wrap" data-text="SHRT">SHRT</span>
        </span>
        <span class="logo-bracket">]</span>
        <span class="logo-sub">URL_SHORTENER.v2</span>
      </div>
      <div class="header-meta">
        <span class="health-dot" :class="healthClass" :title="healthStatus" />
        <span class="health-label">{{ healthStatus }}</span>
      </div>
    </header>

    <!-- ── MAIN ── -->
    <main class="main-content">
      <ShortenForm @result="onResult" />

      <ResultCard
        :result="latestResult"
        @clear="latestResult = null"
      />

      <StatsPanel ref="statsPanel" />

      <HistoryList ref="historyList" />
    </main>

    <!-- ── FOOTER ── -->
    <footer class="site-footer">
      <span>GO + GIN + POSTGRESQL + REDIS</span>
      <span class="footer-dot">◆</span>
      <span>BASE62 · 8-CHAR CODE</span>
      <span class="footer-dot">◆</span>
      <span>7-DAY CACHE</span>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import ParticleBackground from './components/ParticleBackground.vue'
import ShortenForm        from './components/ShortenForm.vue'
import ResultCard         from './components/ResultCard.vue'
import StatsPanel         from './components/StatsPanel.vue'
import HistoryList        from './components/HistoryList.vue'

const latestResult = ref(null)
const statsPanel   = ref(null)
const historyList  = ref(null)
const healthStatus = ref('CONNECTING...')
const healthClass  = ref('dot-warn')

onMounted(async () => {
  try {
    const r = await fetch('/health')
    if (r.ok) { healthStatus.value = 'ONLINE'; healthClass.value = 'dot-ok' }
    else       { healthStatus.value = 'DEGRADED'; healthClass.value = 'dot-warn' }
  } catch {
    healthStatus.value = 'OFFLINE'; healthClass.value = 'dot-err'
  }
})

function onResult(result) {
  const time = new Date().toLocaleTimeString('zh-TW', {
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  })
  latestResult.value = result
  historyList.value?.addItem({ ...result, time })
  statsPanel.value?.fetchStats()
}
</script>

<style scoped>
/* ── Layout ────────────────────────────────────────────────────── */
.layout {
  position: relative;
  z-index: 1;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  max-width: 760px;
  margin: 0 auto;
  padding: 0 20px;
}

/* ── Header ────────────────────────────────────────────────────── */
.site-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32px 0 28px;
  animation: fade-in 0.6s ease both;
}

.logo-block {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.logo-bracket {
  font-family: var(--mono);
  font-size: 28px;
  color: var(--text-muted);
  font-weight: 300;
  line-height: 1;
}

.logo-text {
  font-family: var(--sans);
  font-size: 30px;
  font-weight: 800;
  letter-spacing: -0.02em;
  line-height: 1;
  color: var(--cyan);
  text-shadow: var(--glow-cyan);
  animation: flicker 8s ease infinite;
}

/* glitch layers */
.glitch-wrap {
  position: relative;
  display: inline-block;
}
.glitch-wrap::before,
.glitch-wrap::after {
  content: attr(data-text);
  position: absolute;
  inset: 0;
  opacity: 0;
}
.glitch-wrap:hover::before {
  opacity: 1;
  color: var(--red);
  animation: glitch-1 0.35s steps(1) 1;
}
.glitch-wrap:hover::after {
  opacity: 1;
  color: var(--acid);
  animation: glitch-2 0.35s steps(1) 1;
}

.logo-sub {
  font-family: var(--mono);
  font-size: 9px;
  color: var(--text-muted);
  letter-spacing: 0.14em;
  margin-left: 8px;
  align-self: flex-end;
  padding-bottom: 4px;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 7px;
}

.health-dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot-ok   { background: var(--acid);  box-shadow: 0 0 8px var(--acid);  animation: pulse-border 2s ease infinite; }
.dot-warn { background: #ffba00;       box-shadow: 0 0 8px #ffba00; }
.dot-err  { background: var(--red);   box-shadow: 0 0 8px var(--red); }

.health-label {
  font-size: 10px;
  letter-spacing: 0.14em;
  color: var(--text-muted);
}

/* ── Main ──────────────────────────────────────────────────────── */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-bottom: 40px;
}

/* ── Footer ────────────────────────────────────────────────────── */
.site-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 20px 0 28px;
  font-size: 9px;
  letter-spacing: 0.14em;
  color: var(--text-muted);
  border-top: 1px solid var(--border);
  animation: fade-in 0.8s ease both;
  animation-delay: 0.7s;
  opacity: 0;
  flex-wrap: wrap;
}

.footer-dot { color: var(--cyan); opacity: 0.4; font-size: 6px; }

/* ── Responsive ────────────────────────────────────────────────── */
@media (max-width: 640px) {
  .layout       { padding: 0 14px; }
  .site-header  { padding: 22px 0 20px; }
  .logo-text    { font-size: 24px; }
  .logo-bracket { font-size: 22px; }
  .logo-sub     { display: none; }
  .site-footer  { gap: 8px; }
}

@media (max-width: 380px) {
  .logo-text { font-size: 20px; }
}
</style>
