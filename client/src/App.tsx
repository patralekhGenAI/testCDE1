import { useState } from 'react'

export default function App() {
  const [input, setInput] = useState('')
  const [message, setMessage] = useState('Ready.')

  // NOTE: Backend not wired yet — this will 404 until we add the Go API next step.
  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setMessage('Submitting to /api/standardize (backend coming next)…')
    try {
      const res = await fetch('/api/standardize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ address: input })
      })
      if (!res.ok) {
        setMessage(`Waiting on backend… (HTTP ${res.status})`)
        return
      }
      const data = await res.json()
      setMessage(JSON.stringify(data))
    } catch (err: any) {
      setMessage(err?.message || 'Network error')
    }
  }

  return (
    <div style={{ maxWidth: 640, margin: '3rem auto', fontFamily: 'sans-serif' }}>
      <h1>Address Standardizer (Frontend)</h1>
      <form onSubmit={handleSubmit}>
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Enter any address"
          style={{ width: '100%', padding: '0.75rem' }}
        />
        <button disabled={!input} style={{ marginTop: '1rem' }}>Submit</button>
      </form>
      <p style={{ marginTop: '1rem' }}>{message}</p>
    </div>
  )
}
