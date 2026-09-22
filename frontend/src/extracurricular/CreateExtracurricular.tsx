import { useState } from 'react'

const API_BASE = 'http://localhost:8080'
const TEST_STUDENT_ID = '00000000-0000-0000-0000-000000000002'

type FormState = {
  name: string
  status: 'doing' | 'have_done'
  type: 'maintenance' | 'investment'
  description: string
  leadership_role: string
  start_date: string
  end_date: string
  organization: string
}

export default function CreateExtracurricular() {
  const [form, setForm] = useState<FormState>({
    name: '',
    status: 'doing',
    type: 'maintenance',
    description: '',
    leadership_role: '',
    start_date: '',
    end_date: '',
    organization: '',
  })
  const [msg, setMsg] = useState<string | null>(null)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setForm(prev => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setMsg(null)
    const payload = {
      student_id: TEST_STUDENT_ID,
      user_id: TEST_STUDENT_ID,
      name: form.name,
      status: form.status,
      type: form.type,
      description: form.description,
      leadership_role: form.leadership_role || null,
      start_date: form.start_date,
      end_date: form.end_date || null,
      organization: form.organization,
    }

    try {
      const res = await fetch(`${API_BASE}/extracurriculars`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) throw new Error(await res.text())
      const data = await res.json()
      const created = data?.body ?? data?.Body ?? data
      setMsg('Created id: ' + (created?.id ?? 'unknown'))
    } catch (err: any) {
      setMsg('Error: ' + (err.message || String(err)))
    }
  }

  return (
    <div>
      <h2>Create Extracurricular</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Name: <input name="name" value={form.name} onChange={handleChange} /></label>
        </div>
        <div>
          <label>Status: 
            <select name="status" value={form.status} onChange={handleChange}>
              <option value="doing">doing</option>
              <option value="have_done">have_done</option>
            </select>
          </label>
        </div>
        <div>
          <label>Type: 
            <select name="type" value={form.type} onChange={handleChange}>
              <option value="maintenance">maintenance</option>
              <option value="investment">investment</option>
            </select>
          </label>
        </div>
        <div>
          <label>Description:<br />
            <textarea name="description" value={form.description} onChange={handleChange} />
          </label>
        </div>
        <div>
          <label>Leadership role: <input name="leadership_role" value={form.leadership_role} onChange={handleChange} /></label>
        </div>
        <div>
          <label>Start date: <input type="date" name="start_date" value={form.start_date} onChange={handleChange} /></label>
        </div>
        <div>
          <label>End date: <input type="date" name="end_date" value={form.end_date} onChange={handleChange} /></label>
        </div>
        <div>
          <label>Organization: <input name="organization" value={form.organization} onChange={handleChange} /></label>
        </div>
        <div>
          <button type="submit">Create</button>
        </div>
      </form>
      {msg && <p>{msg}</p>}
    </div>
  )
}
