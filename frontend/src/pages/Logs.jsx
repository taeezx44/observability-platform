import { useState, useEffect, useRef } from 'react'

const API = import.meta.env.VITE_API_URL || ''

const LEVELS = ['', 'ERROR', 'WARN', 'INFO', 'DEBUG']

const LEVEL_COLOR = {
  ERROR:   '#ef4444',
  WARN:    '#f59e0b',
  INFO:    '#00ff88',
  DEBUG:   '#64748b',
  UNKNOWN: '#64748b',
}

// Mock logs for demo
const mockLogs = () => {
  const services = ['api', 'collector', 'alerting']
  const levels = ['INFO', 'INFO', 'INFO', 'WARN', 'ERROR']
  const messages = [
    'Request completed in 23ms',
    'Scraping target http://app:8080/metrics',
    'BatchInsert: 42 metrics written to ClickHouse',
    'ClickHouse query took 450ms, consider adding index',
    'Failed to connect to kafka: connection refused',
    'Alert rule HighCPU: current value 87.2 > threshold 85',
    'WebSocket client connected from 127.0.0.1',
    'Metrics scraped: 156 samples from 3 targets',
  ]
  return Array.from({ length: 30 }, (_, i) => ({
    timestamp: new Date(Date.now() - i * 4000).toISOString(),
    level: levels[Math.floor(Math.random() * levels.length)],
    service: services[Math.floor(Math.random() * services.length)],
    message: messages[Math.floor(Math.random() * messages.length)],
    fields: {},
  }))
}

export default function Logs() {
  const [logs, setLogs] = useState(mockLogs())
  const [search, setSearch] = useState('')
  const [level, setLevel] = useState('')
  const [service, setService] = useState('')
  const [autoScroll, setAutoScroll] = useState(true)
  const bottomRef = useRef(null)

  useEffect(() => {
    const fetchLogs = () => {
      const params = new URLSearchParams({ limit: 200 })
      if (search)  params.set('search', search)
      if (level)   params.set('level', level)
      if (service) params.set('service', service)

      fetch(`${API}/api/logs?${params}`)
        .then(r => r.json())
        .then(data => { if (data?.length) setLogs(data) })
        .catch(() => {
          // Keep adding mock logs
          setLogs(prev => {
            const msgs = ['Request completed', 'Scraping metrics', 'Alert evaluated', 'WebSocket ping']
            const lvls = ['INFO', 'INFO', 'WARN', 'ERROR']
            const svcs = ['api', 'collector', 'alerting']
            const newLog = {
              timestamp: new Date().toISOString(),
              level: lvls[Math.floor(Math.random() * lvls.length)],
              service: svcs[Math.floor(Math.random() * svcs.length)],
              message: msgs[Math.floor(Math.random() * msgs.length)],
              fields: {},
            }
            return [newLog, ...prev.slice(0, 299)]
          })
        })
    }

    fetchLogs()
    const interval = setInterval(fetchLogs, 5000)
    return () => clearInterval(interval)
  }, [search, level, service])

  useEffect(() => {
    if (autoScroll && bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [logs, autoScroll])

  const filtered = logs.filter(l => {
    if (level && l.level !== level) return false
    if (service && l.service !== service) return false
    if (search && !l.message.toLowerCase().includes(search.toLowerCase())) return false
    return true
  })

  return (
    <div style={{ padding: '32px 40px' }}>
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ fontSize: 28, fontWeight: 800, letterSpacing: -1, marginBottom: 8 }}>Log Explorer</h1>
        <div style={{ fontFamily: 'Space Mono', fontSize: 11, color: 'var(--muted)' }}>
          Live log stream • {filtered.length} entries
        </div>
      </div>

      {/* Filters */}
      <div style={{ display: 'flex', gap: 12, marginBottom: 20, flexWrap: 'wrap' }}>
        <input
          placeholder="Search messages..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          style={{
            flex: 1, minWidth: 200,
            background: 'var(--surface)', border: '1px solid var(--border)',
            borderRadius: 6, padding: '8px 14px',
            color: 'var(--text)', fontFamily: 'Space Mono', fontSize: 12,
            outline: 'none',
          }}
        />
        <select
          value={level}
          onChange={e => setLevel(e.target.value)}
          style={{
            background: 'var(--surface)', border: '1px solid var(--border)',
            borderRadius: 6, padding: '8px 14px',
            color: level ? LEVEL_COLOR[level] || 'var(--text)' : 'var(--muted)',
            fontFamily: 'Space Mono', fontSize: 12, cursor: 'pointer',
          }}
        >
          {LEVELS.map(l => <option key={l} value={l}>{l || 'All levels'}</option>)}
        </select>
        <label style={{ display: 'flex', alignItems: 'center', gap: 8, fontFamily: 'Space Mono', fontSize: 11, color: 'var(--muted)', cursor: 'pointer' }}>
          <input type="checkbox" checked={autoScroll} onChange={e => setAutoScroll(e.target.checked)} />
          Auto-scroll
        </label>
      </div>

      {/* Log stream */}
      <div style={{
        background: 'var(--surface)', border: '1px solid var(--border)',
        borderRadius: 8, height: 'calc(100vh - 280px)',
        overflow: 'auto', fontFamily: 'Space Mono', fontSize: 12,
      }}>
        {filtered.map((log, i) => (
          <div
            key={i}
            style={{
              display: 'flex', gap: 16, padding: '8px 16px',
              borderBottom: '1px solid rgba(255,255,255,0.03)',
              alignItems: 'flex-start',
              background: log.level === 'ERROR' ? 'rgba(239,68,68,0.04)' : 'transparent',
              animation: i === 0 ? 'fadeIn 0.2s ease' : 'none',
            }}
          >
            <span style={{ color: 'var(--muted)', whiteSpace: 'nowrap', flexShrink: 0, fontSize: 10 }}>
              {new Date(log.timestamp).toLocaleTimeString()}
            </span>
            <span style={{
              color: LEVEL_COLOR[log.level] || 'var(--muted)',
              fontWeight: 700, minWidth: 48, flexShrink: 0, fontSize: 10,
            }}>
              {log.level}
            </span>
            <span style={{ color: '#a78bfa', flexShrink: 0, minWidth: 80, fontSize: 10 }}>
              {log.service}
            </span>
            <span style={{ color: 'var(--text)', lineHeight: 1.4 }}>
              {log.message}
            </span>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
    </div>
  )
}
