<template>
  <canvas ref="canvas" class="particle-canvas" aria-hidden="true" />
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const canvas = ref(null)
let animId = null

onMounted(() => {
  const c = canvas.value
  const ctx = c.getContext('2d')
  let W, H, nodes

  const NODE_COUNT_BASE = 80
  const CONNECT_DIST   = 160
  const SPEED          = 0.35

  const COLORS = ['#00ffe0', '#00ffe0', '#00ffe0', '#a8ff3e', '#ffffff']

  function resize() {
    W = c.width  = window.innerWidth
    H = c.height = window.innerHeight
  }

  function makeNode() {
    return {
      x:   Math.random() * W,
      y:   Math.random() * H,
      vx:  (Math.random() - 0.5) * SPEED,
      vy:  (Math.random() - 0.5) * SPEED,
      r:   Math.random() * 1.5 + 0.5,
      col: COLORS[Math.floor(Math.random() * COLORS.length)],
      pulse: Math.random() * Math.PI * 2,
    }
  }

  function init() {
    resize()
    const count = Math.min(NODE_COUNT_BASE, Math.floor((W * H) / 14000))
    nodes = Array.from({ length: count }, makeNode)
  }

  function draw() {
    ctx.clearRect(0, 0, W, H)

    // update
    for (const n of nodes) {
      n.x += n.vx
      n.y += n.vy
      n.pulse += 0.02
      if (n.x < 0 || n.x > W) n.vx *= -1
      if (n.y < 0 || n.y > H) n.vy *= -1
    }

    // edges
    for (let i = 0; i < nodes.length; i++) {
      for (let j = i + 1; j < nodes.length; j++) {
        const a = nodes[i], b = nodes[j]
        const dx = a.x - b.x, dy = a.y - b.y
        const dist = Math.sqrt(dx * dx + dy * dy)
        if (dist < CONNECT_DIST) {
          const alpha = (1 - dist / CONNECT_DIST) * 0.35
          ctx.beginPath()
          ctx.moveTo(a.x, a.y)
          ctx.lineTo(b.x, b.y)
          ctx.strokeStyle = `rgba(0,255,224,${alpha})`
          ctx.lineWidth = 0.6
          ctx.stroke()
        }
      }
    }

    // nodes
    for (const n of nodes) {
      const pulseR = n.r + Math.sin(n.pulse) * 0.4
      ctx.beginPath()
      ctx.arc(n.x, n.y, pulseR, 0, Math.PI * 2)
      ctx.fillStyle = n.col
      ctx.shadowBlur  = 8
      ctx.shadowColor = n.col
      ctx.fill()
      ctx.shadowBlur = 0
    }

    animId = requestAnimationFrame(draw)
  }

  init()
  draw()
  window.addEventListener('resize', init)
})

onUnmounted(() => {
  cancelAnimationFrame(animId)
  window.removeEventListener('resize', () => {})
})
</script>

<style scoped>
.particle-canvas {
  position: fixed;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  opacity: 0.55;
}
</style>
