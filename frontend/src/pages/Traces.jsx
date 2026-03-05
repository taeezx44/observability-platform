import { useState, useEffect } from 'react'

const API = import.meta.env.VITE_API_URL || ''

// Generate mock trace data
const mockTraces = () => [
  { trace_id: 'abc123def456', service: 'api',        operation: 'POST /checkout',    span_count: 8, duration_ms: 843 },
  { trace_id: 'bcd234efg567', service: 'api',        operation: 'GET /products',     span_count: 3, duration_ms: 234 },
  { trace_id: 'cde345fgh678', service: 'collector',  operation: 'scrapeTarget',      span_count: 2, duration_ms: 567 },
  { trace_id: 'def456ghi789', service: 'api',        operation: 'PUT /cart',         span_count: 5, duration_ms: 1203 },
  { trace_id: 'efg567hij890', service: 'alerting',   operation: 'evaluateRule',      span_count: 2, duration_ms: 89 },
]

const mockWaterfall = (traceId) => {
  const spans = [
    { span: { trace_id: traceId, span_id: '1', parent_id: '',  service: 'api',      operation: 'HTTP POST /checkout', status: 'OK',    tags: { 'http.method': 'POST' } }, depth: 0, start_offset_ms: 0,   duration_ms: 843 },
    { span: { trace_id: traceId, span_id: '2', parent_id: '1', service: 'api',      operation: 'validateCart',        status: 'OK',    tags: {} },                         depth: 1, start_offset_ms: 5,   duration_ms: 12 },
    { span: { trace_id: traceId, span_id: '3', parent_id: '1', service: 'auth',     operation: 'verifyToken',         status: 'OK',    tags: {} },                         depth: 1, start_offset_ms: 18,  duration_ms: 45 },
    { span: { trace_id: traceId, span_id: '4', parent_id: '1', service: 'payment',  operation: 'chargeCard',          status: 'OK',    tags: {} },                         depth: 1, start_offset_ms: 64,  duration_ms: 620 },
    { span: { trace_id: traceId, span_id: '5', parent_id: '4', service: 'postgres', operation: 'INSERT transactions', status: 'OK',    tags: { 'db.type': 'postgres' } },  depth: 2, start_offset_ms: 80,  duration_ms: 580 },
    { span: { trace_id: traceId, span_id: '6', parent_id: '1', service: 'email',    operation: 'sendReceipt',         status: 'OK',    tags: {} },                         depth: 1, start_offset_ms: 690, duration_ms: 120 },
    { span: { trace_id: traceId, span_id: '7', parent_id: '1', service: 'redis',    operation: 'SET cart:cache',      status: 'OK',    tags: {} },                         depth: 1, start_offset_ms: 815, duration_ms: 8 },
  ]
  return spans
}

function TraceRow({ trace, onClick, selected }) {
  const slow = trace.duration_ms > 500
  return (
    <div
      onClick={() => onClick(trace.trace_id)}
      style={{
        display: 'flex', alignItems: 'center', gap: 16,
        padding: '12px 16px', cursor: 'pointer',
        borderBottom: '1px solid rgba(255,255,255,0.04)',
        background: selected ? 'rgba(0,255,136,0.05)' : 'transparent',
        transition: 'background 0.15s',
      }}
    >
      <span style={{ fontFamily: 'Space Mono', fontSize: 11, color: '#a78bfa', width: 100, flexShrink: 0 }}>
        {trace.service}
      </span>
      <span style={{ flex: 1, fontSize: 13, fontWeight: 600 }}>
        {trace.operation}
      </span>
      <span style={{ fontFamily: 'Space Mono', fontSize: 10, color: 'var(--muted)' }}>
        {trace.span_count} spans
      </span>
      <span style={{
        fontFamily: 'Space Mono', fontSize: 11,
        color: slow ? 'var(--danger)' : 'var(--accent)',
        fontWeight: 700, minWidth: 64, textAlign: 'right',
      }}>
        {trace.duration_ms}ms
      </span>
    </div>
  )
}

function WaterfallBar({ item, totalMs }) {
  const startPct = (item.start_offset_ms / totalMs) * 100
  const widthPct = Math.max((item.duration_ms / totalMs) * 100, 0.5)
  const isSlow = item.duration_ms > 300
  const isError = item.span.status === 'ERROR'

  const color = isError ? '#ef4444' : isSlow ? '#f59e0b' : '#00ff88'

  return (
    <div style={{
      display: 'grid', gridTemplateColumns: '200px 1fr',
      borderBottom: '1px solid rgba(255,255,255,0.04)',
      minHeight: 32,
    }}>
      {/* Label */}
      <div style={{
        display: 'flex', alignItems: 'center', gap: 8,
        padding: '6px 12px', paddingLeft: 12 + item.depth * 16,
        borderRight: '1px solid var(--border)',
      }}>
        <span style={{ fontFamily: 'Space Mono', fontSize: 9, color: '#a78bfa', flexShrink: 0 }}>
          {item.span.service}
        </span>
        <span style={{ fontFamily: 'Space Mono', fontSize: 9, color: 'var(--muted)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
          {item.span.operation}
        </span>
      </div>

      {/* Bar */}
      <div style={{ position: 'relative', background: '#0d0d14', display: 'flex', alignItems: 'center' }}>
        <div style={{
          position: 'absolute',
          left: `${startPct}%`,
          width: `${widthPct}%`,
          height: 16, borderRadius: 2,
          background: color,
          opacity: 0.85,
          transition: 'opacity 0.2s',
          minWidth: 4,
        }} />
        <span style={{
          position: 'absolute',
          left: `calc(${startPct}% + ${widthPct}% + 6px)`,
          fontFamily: 'Space Mono', fontSize: 9,
          color: isSlow ? color : 'var(--muted)',
          whiteSpace: 'nowrap',
        }}>
          {item.duration_ms}ms
        </span>
      </div>
    </div>
  )
}

export default function Traces() {
  const [traces, setTraces] = useState(mockTraces())
  const [selectedTrace, setSelectedTrace] = useState(null)
  const [waterfall, setWaterfall] = useState([])

  useEffect(() => {
    fetch(`${API}/api/traces?min_ms=100&limit=50`)
      .then(r => r.json())
      .then(data => { if (data?.length) setTraces(data) })
      .catch(() => {})
  }, [])

  const handleSelect = (traceId) => {
    setSelectedTrace(traceId)

    fetch(`${API}/api/traces/${traceId}`)
      .then(r => r.json())
      .then(data => { if (data?.length) setWaterfall(data) })
      .catch(() => {
        setWaterfall(mockWaterfall(traceId))
      })
  }

  const totalMs = waterfall.length
    ? Math.max(...waterfall.map(s => s.start_offset_ms + s.duration_ms))
    : 1

  return (
    <div style={{ padding: '32px 40px' }}>
      <div style={{ marginBottom: 24 }}>
        <h1 style={{ fontSize: 28, fontWeight: 800, letterSpacing: -1, marginBottom: 8 }}>Distributed Traces</h1>
        <div style={{ fontFamily: 'Space Mono', fontSize: 11, color: 'var(--muted)' }}>
          Slow traces ({'>'}500ms) • Click a trace to see waterfall
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1.4fr', gap: 20, height: 'calc(100vh - 220px)' }}>
        {/* Trace list */}
        <div style={{
          background: 'var(--surface)', border: '1px solid var(--border)',
          borderRadius: 8, overflow: 'auto',
        }}>
          <div style={{
            padding: '12px 16px', borderBottom: '1px solid var(--border)',
            fontFamily: 'Space Mono', fontSize: 10, color: 'var(--accent)',
            letterSpacing: 1,
          }}>
            RECENT SLOW TRACES
          </div>
          {traces.map(t => (
            <TraceRow
              key={t.trace_id}
              trace={t}
              onClick={handleSelect}
              selected={selectedTrace === t.trace_id}
            />
          ))}
        </div>

        {/* Waterfall */}
        <div style={{
          background: 'var(--surface)', border: '1px solid var(--border)',
          borderRadius: 8, overflow: 'auto',
        }}>
          <div style={{
            padding: '12px 16px', borderBottom: '1px solid var(--border)',
            fontFamily: 'Space Mono', fontSize: 10, color: 'var(--accent4)',
            letterSpacing: 1, display: 'flex', justifyContent: 'space-between',
          }}>
            <span>TRACE WATERFALL</span>
            {selectedTrace && (
              <span style={{ color: 'var(--muted)', fontSize: 9 }}>
                {selectedTrace.slice(0, 12)}...
              </span>
            )}
          </div>

          {waterfall.length === 0 ? (
            <div style={{
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              height: 200, fontFamily: 'Space Mono', fontSize: 11, color: 'var(--muted)',
            }}>
              ← Select a trace to view waterfall
            </div>
          ) : (
            waterfall.map((item, i) => (
              <WaterfallBar key={i} item={item} totalMs={totalMs} />
            ))
          )}
        </div>
      </div>
    </div>
  )
}
