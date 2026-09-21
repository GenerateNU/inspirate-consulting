import { useEffect, useState } from 'react'

const API_BASE = 'http://localhost:8080'
const TEST_STUDENT_ID = '00000000-0000-0000-0000-000000000002'

type Extracurricular = {
  id: number
  name: string
  status: string
  type: string
  description: string
  leadership_role?: string | null
  start_date: string
  end_date?: string | null
  organization: string
}

export default function Extracurriculars() {
  const [list, setList] = useState<Extracurricular[]>([])
  const [editing, setEditing] = useState<number | null>(null)
  const [form, setForm] = useState<Partial<Extracurricular>>({})
  const [msg, setMsg] = useState<string | null>(null)

  const fetchList = async () => {
    setMsg(null)
    try {
      const res = await fetch(`${API_BASE}/extracurriculars?student_id=${TEST_STUDENT_ID}`)
      if (!res.ok) throw new Error(await res.text())
      const data = await res.json()
      const listData = data?.body ?? data?.Body ?? data
      setList(Array.isArray(listData) ? listData : [])
    } catch (err: any) {
      setMsg('Error: ' + (err.message || String(err)))
    }
  }

  useEffect(() => { fetchList() }, [])

  const startEdit = (item: Extracurricular) => {
    // Normalize ISO datetimes (e.g. "2026-09-21T00:00:00Z") to yyyy-MM-dd
    const isoToYMD = (s?: string | null) => {
      if (!s) return ''
      const t = s.indexOf('T')
      if (t > 0) return s.slice(0, t)
      return s.length >= 10 ? s.slice(0, 10) : s
    }
    setEditing(item.id)
    setForm({
      name: item.name,
      status: item.status,
      type: item.type,
      description: item.description,
      leadership_role: item.leadership_role ?? null,
      start_date: isoToYMD(item.start_date) as any,
      end_date: isoToYMD(item.end_date ?? undefined) as any,
      organization: item.organization,
    })
  }

  const cancelEdit = () => { setEditing(null); setForm({}) }

  const saveEdit = async () => {
    if (!editing) return
    setMsg(null)
    try {
      const isoToYMD = (s?: string | null | number) => {
        if (!s) return null
        const str = String(s)
        const t = str.indexOf('T')
        if (t > 0) return str.slice(0, t)
        return str.length >= 10 ? str.slice(0, 10) : str
      }
      const payload = {
        student_id: TEST_STUDENT_ID,
        name: form.name,
        status: form.status,
        type: form.type,
        description: form.description,
        leadership_role: form.leadership_role ?? null,
        start_date: isoToYMD(form.start_date),
        end_date: form.end_date ? isoToYMD(form.end_date) : null,
        organization: form.organization,
      }
      const res = await fetch(`${API_BASE}/extracurriculars/` + editing, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) throw new Error(await res.text())
      await fetchList()
      setEditing(null)
    } catch (err: any) {
      setMsg('Error: ' + (err.message || String(err)))
    }
  }

  const updateField = (k: keyof Extracurricular, v: any) => setForm(prev => ({ ...prev, [k]: v }))

  return (
    <div>
      <h2>Extracurriculars</h2>
      {msg && <p>{msg}</p>}
      <button onClick={fetchList}>Refresh</button>
      <ul>
        {list.map(item => (
          <li key={item.id} style={{ marginBottom: 12 }}>
            {editing === item.id ? (
              <div>
                <div>
                  <label>Name: <input value={form.name ?? ''} onChange={e => updateField('name', e.target.value)} /></label>
                </div>
                <div>
                  <label>Status: <select value={form.status ?? ''} onChange={e => updateField('status', e.target.value)}>
                    <option value="doing">doing</option>
                    <option value="have_done">have_done</option>
                  </select></label>
                </div>
                <div>
                  <label>Type: <select value={form.type ?? ''} onChange={e => updateField('type', e.target.value)}>
                    <option value="maintenance">maintenance</option>
                    <option value="investment">investment</option>
                  </select></label>
                </div>
                <div>
                  <label>Description:<br />
                    <textarea value={form.description ?? ''} onChange={e => updateField('description', e.target.value)} />
                  </label>
                </div>
                <div>
                  <label>Leadership: <input value={form.leadership_role ?? ''} onChange={e => updateField('leadership_role', e.target.value)} /></label>
                </div>
                <div>
                  <label>Start date: <input type="date" value={form.start_date ?? ''} onChange={e => updateField('start_date', e.target.value)} /></label>
                </div>
                <div>
                  <label>End date: <input type="date" value={form.end_date ?? ''} onChange={e => updateField('end_date', e.target.value)} /></label>
                </div>
                <div>
                  <label>Organization: <input value={form.organization ?? ''} onChange={e => updateField('organization', e.target.value)} /></label>
                </div>
                <div>
                  <button onClick={saveEdit}>Save</button>
                  <button onClick={cancelEdit}>Cancel</button>
                </div>
              </div>
            ) : (
              <div>
                <strong>{item.name}</strong> — {item.organization} — {item.status}
                <div>
                  <button onClick={() => startEdit(item)}>Edit</button>
                </div>
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}
