import { useState } from 'react'
import Dashboard from './pages/Dashboard.jsx'
import Logs from './pages/Logs.jsx'
import Traces from './pages/Traces.jsx'

const NAV = [
  { id: 'dashboard', label: '📊 Metrics',  color: '#00ff88' },
  { id: 'logs',      label: '📝 Logs',     color: '#a78bfa' },
  { id: 'traces',    label: '🔍 Traces',   color: '#06b6d4' },
]

export default function App() {
  const [page, setPage] = useState('dashboard')

  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      {/* Sidebar */}
      <nav style={{
        width: 220, flexShrink: 0,
        background: 'var(--surface)',
        borderRight: '1px solid var(--border)',
        display: 'flex', flexDirection: 'column',
        padding: '24px 0',
        position: 'sticky', top: 0, height: '100vh',
      }}>
        <div style={{ padding: '0 20px 24px', borderBottom: '1px solid var(--border)' }}>
          <div style={{ fontFamily: 'Space Mono', fontSize: 10, color: 'var(--accent)', letterSpacing: 2, marginBottom: 6 }}>
            ◉ LIVE
          </div>
          <div style={{ fontSize: 16, fontWeight: 800, letterSpacing: -0.5 }}>Observability</div>
          <div style={{ fontFamily: 'Space Mono', fontSize: 10, color: 'var(--muted)', marginTop: 2 }}>self-hosted platform</div>
        </div>

        <div style={{ padding: '16px 12px', flex: 1 }}>
          {NAV.map(n => (
            <button
              key={n.id}
              onClick={() => setPage(n.id)}
              style={{
                display: 'block', width: '100%',
                padding: '10px 12px', marginBottom: 4,
                borderRadius: 6, border: 'none',
                background: page === n.id ? `rgba(${n.id === 'dashboard' ? '0,255,136' : n.id === 'logs' ? '167,139,250' : '6,182,212'},0.1)` : 'transparent',
                color: page === n.id ? n.color : 'var(--muted)',
                fontFamily: 'Space Mono', fontSize: 12,
                cursor: 'pointer', textAlign: 'left',
                transition: 'all 0.15s',
                outline: page === n.id ? `1px solid ${n.color}33` : '1px solid transparent',
              }}
            >
              {n.label}
            </button>
          ))}
        </div>

        <div style={{ padding: '16px 20px', borderTop: '1px solid var(--border)' }}>
          <div style={{ fontFamily: 'Space Mono', fontSize: 9, color: 'var(--muted)', lineHeight: 1.6 }}>
            Go + ClickHouse<br />React + Recharts
          </div>
        </div>
      </nav>

      {/* Main content */}
      <main style={{ flex: 1, overflow: 'auto' }}>
        {page === 'dashboard' && <Dashboard />}
        {page === 'logs'      && <Logs />}
        {page === 'traces'    && <Traces />}
      </main>
    </div>
  )
}
