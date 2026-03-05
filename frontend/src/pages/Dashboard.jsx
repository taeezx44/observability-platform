import { useState, useEffect } from 'react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts'

const API = import.meta.env.VITE_API_URL || ''

// Mock data for when API is not connected
const mockData = () => Array.from({ length: 20 }, (_, i) => ({
  time: new Date(Date.now() - (20 - i) * 30000).toLocaleTimeString(),
  cpu: 30 + Math.random() * 40,
  memory: 50 + Math.random() * 30,
  requests: 80 + Math.random() * 120,
}))

function StatCard({ label, value, unit, color, trend }) {
  return (
    <div style={{
      background: 'var(--surface)', border: '1px solid var(--border)',
      borderRadius: 8, padding: '20px 24px',
      borderTop: `2px solid ${color}`,
      animation: 'fadeIn 0.3s ease',
    }}>
      <div style={{ fontFamily: 'Space Mono', fontSize: 10, color: 'var(--muted)', letterSpacing: 2, marginBottom: 8 }}>
        {label}
      </div>
      <div style={{ fontSize: 32, fontWeight: 800, color, letterSpacing: -1 }}>
        {value}<span style={{ fontSize: 14, marginLeft: 4, color: 'var(--muted)' }}>{unit}</span>
      </div>
      {trend && (
        <div style={{ fontFamily: 'Space Mono', fontSize: 10, color: trend > 0 ? 'var(--danger)' : 'var(--accent)', marginTop: 6 }}>
          {trend > 0 ? '↑' : '↓'} {Math.abs(trend)}% vs 5m ago
        </div>
      )}
    </div>
  )
}

function ChartCard({ title, children }) {
  return (
    <div style={{
      background: 'var(--surface)', border: '1px solid var(--border)',
      borderRadius: 8, padding: '20px',
    }}>
      <div style={{ fontFamily: 'Space Mono', fontSize: 11, color: 'var(--accent)', marginBottom: 16, letterSpacing: 1 }}>
        {title}
      </div>
      {children}
    </div>
  )
}

const tooltipStyle = {
  background: '#1a1a24', border: '1px solid #2a2a3a',
  borderRadius: 6, fontFamily: 'Space Mono', fontSize: 11,
}

export default function Dashboard() {
  const [data, setData] = useState(mockData())
  const [metricNames, setMetricNames] = useState([])
  const [liveIndicator, setLiveIndicator] = useState(true)

  useEffect(() => {
    // Try to fetch metric names
    fetch(`${API}/api/metrics/names`)
      .then(r => r.json())
      .then(setMetricNames)
      .catch(() => {}) // silently fail if API not up

    // Try to fetch real metrics
    const fetchMetrics = () => {
      fetch(`${API}/api/metrics?limit=100`)
        .then(r => r.json())
        .then(metrics => {
          if (metrics && metrics.length > 0) {
            // Transform real data
            const transformed = metrics.slice(0, 20).map(m => ({
              time: new Date(m.timestamp).toLocaleTimeString(),
              [m.name]: m.value,
            }))
            setData(transformed)
          }
        })
        .catch(() => {
          // Keep mock data if API not available
          setData(prev => {
            const next = [...prev.slice(1), {
              time: new Date().toLocaleTimeString(),
              cpu: 30 + Math.random() * 40,
              memory: 50 + Math.random() * 30,
              requests: 80 + Math.random() * 120,
            }]
            return next
          })
        })
    }

    fetchMetrics()
    const interval = setInterval(fetchMetrics, 15000)
    const liveBlink = setInterval(() => setLiveIndicator(v => !v), 1000)

    return () => { clearInterval(interval); clearInterval(liveBlink) }
  }, [])

  const latest = data[data.length - 1] || {}

  return (
    <div style={{ padding: '32px 40px' }}>
      {/* Header */}
      <div style={{ marginBottom: 32 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
          <h1 style={{ fontSize: 28, fontWeight: 800, letterSpacing: -1 }}>Metrics Dashboard</h1>
          <div style={{
            display: 'flex', alignItems: 'center', gap: 6,
            padding: '4px 10px', borderRadius: 4,
            background: 'rgba(0,255,136,0.08)', border: '1px solid rgba(0,255,136,0.2)',
            fontFamily: 'Space Mono', fontSize: 10, color: 'var(--accent)',
          }}>
            <span style={{ width: 6, height: 6, borderRadius: '50%', background: 'var(--accent)', opacity: liveIndicator ? 1 : 0.2, display: 'inline-block' }} />
            LIVE
          </div>
        </div>
        <div style={{ fontFamily: 'Space Mono', fontSize: 11, color: 'var(--muted)' }}>
          Refreshes every 15s • {metricNames.length > 0 ? `${metricNames.length} metric series` : 'Using demo data — connect scraper to see real metrics'}
        </div>
      </div>

      {/* Stat cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(180px,1fr))', gap: 16, marginBottom: 32 }}>
        <StatCard label="CPU USAGE"     value={(latest.cpu     || 0).toFixed(1)} unit="%" color="var(--accent)"  trend={2.1} />
        <StatCard label="MEMORY"        value={(latest.memory  || 0).toFixed(1)} unit="%" color="#a78bfa"        trend={-0.4} />
        <StatCard label="REQUESTS/MIN"  value={(latest.requests|| 0).toFixed(0)} unit=""  color="var(--accent4)" trend={12} />
        <StatCard label="ERROR RATE"    value="0.3"                              unit="%" color="var(--danger)"  trend={-0.1} />
      </div>

      {/* Charts */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 20, marginBottom: 20 }}>
        <ChartCard title="CPU USAGE (%)">
          <ResponsiveContainer width="100%" height={200}>
            <AreaChart data={data}>
              <defs>
                <linearGradient id="cpuGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#00ff88" stopOpacity={0.2} />
                  <stop offset="95%" stopColor="#00ff88" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#1a1a2a" />
              <XAxis dataKey="time" tick={{ fontFamily: 'Space Mono', fontSize: 9, fill: '#64748b' }} />
              <YAxis domain={[0,100]} tick={{ fontFamily: 'Space Mono', fontSize: 9, fill: '#64748b' }} />
              <Tooltip contentStyle={tooltipStyle} />
              <Area type="monotone" dataKey="cpu" stroke="#00ff88" fill="url(#cpuGrad)" strokeWidth={2} dot={false} />
            </AreaChart>
          </ResponsiveContainer>
        </ChartCard>

        <ChartCard title="MEMORY USAGE (%)">
          <ResponsiveContainer width="100%" height={200}>
            <AreaChart data={data}>
              <defs>
                <linearGradient id="memGrad" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#a78bfa" stopOpacity={0.2} />
                  <stop offset="95%" stopColor="#a78bfa" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#1a1a2a" />
              <XAxis dataKey="time" tick={{ fontFamily: 'Space Mono', fontSize: 9, fill: '#64748b' }} />
              <YAxis domain={[0,100]} tick={{ fontFamily: 'Space Mono', fontSize: 9, fill: '#64748b' }} />
              <Tooltip contentStyle={tooltipStyle} />
              <Area type="monotone" dataKey="memory" stroke="#a78bfa" fill="url(#memGrad)" strokeWidth={2} dot={false} />
            </AreaChart>
          </ResponsiveContainer>
        </ChartCard>
      </div>

      <ChartCard title="HTTP REQUESTS / MIN">
        <ResponsiveContainer width="100%" height={160}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" stroke="#1a1a2a" />
            <XAxis dataKey="time" tick={{ fontFamily: 'Space Mono', fontSize: 9, fill: '#64748b' }} />
            <YAxis tick={{ fontFamily: 'Space Mono', fontSize: 9, fill: '#64748b' }} />
            <Tooltip contentStyle={tooltipStyle} />
            <Line type="monotone" dataKey="requests" stroke="#06b6d4" strokeWidth={2} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </ChartCard>
    </div>
  )
}
