import { useState } from "react";

import {
  updateExtracurricular,
  useListExtracurriculars,
} from "../api/endpoints/extracurriculars/extracurriculars";
import {
  UpdateExtracurricularRequestStatus,
  UpdateExtracurricularRequestType,
  type Extracurricular,
  type UpdateExtracurricularRequest,
} from "../api/models";

// Normalize ISO datetimes (e.g. "2026-09-21T00:00:00Z") to yyyy-MM-dd for date inputs
const toDateInput = (value?: string | null) =>
  value ? value.slice(0, 10) : "";

export default function Extracurriculars() {
  const { data, error, isLoading, mutate } = useListExtracurriculars<Error>();
  const [editing, setEditing] = useState<number | null>(null);
  const [form, setForm] = useState<UpdateExtracurricularRequest>({});
  const [msg, setMsg] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const list: Extracurricular[] = data?.status === 200 ? (data.data ?? []) : [];
  const listError = error
    ? error.message
    : data && data.status !== 200
      ? JSON.stringify(data.data)
      : null;

  const startEdit = (item: Extracurricular) => {
    setMsg(null);
    setEditing(item.id);
    setForm({
      name: item.name,
      status: item.status as UpdateExtracurricularRequestStatus,
      type: item.type as UpdateExtracurricularRequestType,
      description: item.description,
      leadership_role: item.leadership_role ?? undefined,
      start_date: toDateInput(item.start_date),
      end_date: toDateInput(item.end_date) || undefined,
      organization: item.organization,
    });
  };

  const cancelEdit = () => {
    setEditing(null);
    setForm({});
  };

  const saveEdit = async () => {
    if (editing === null) return;
    setMsg(null);
    setSaving(true);

    // Optional fields are sent only when filled in
    const payload: UpdateExtracurricularRequest = {
      ...form,
      leadership_role: form.leadership_role || undefined,
      end_date: form.end_date || undefined,
    };

    try {
      const res = await updateExtracurricular(editing, payload);
      if (res.status === 200) {
        await mutate();
        cancelEdit();
      } else {
        setMsg("Error: " + JSON.stringify(res.data));
      }
    } catch (err) {
      setMsg("Error: " + (err instanceof Error ? err.message : String(err)));
    } finally {
      setSaving(false);
    }
  };

  const updateField = <K extends keyof UpdateExtracurricularRequest>(
    key: K,
    value: UpdateExtracurricularRequest[K],
  ) => setForm((prev) => ({ ...prev, [key]: value }));

  return (
    <div>
      <h2>Extracurriculars</h2>
      {listError && <p>Error: {listError}</p>}
      {msg && <p>{msg}</p>}
      <button onClick={() => mutate()}>Refresh</button>
      {isLoading && <p>Loading...</p>}
      <ul>
        {list.map((item) => (
          <li key={item.id} style={{ marginBottom: 12 }}>
            {editing === item.id ? (
              <div>
                <div>
                  <label>
                    Name:{" "}
                    <input
                      value={form.name ?? ""}
                      onChange={(e) => updateField("name", e.target.value)}
                    />
                  </label>
                </div>
                <div>
                  <label>
                    Status:{" "}
                    <select
                      value={form.status ?? ""}
                      onChange={(e) =>
                        updateField(
                          "status",
                          e.target.value as UpdateExtracurricularRequestStatus,
                        )
                      }
                    >
                      {Object.values(UpdateExtracurricularRequestStatus).map(
                        (status) => (
                          <option key={status} value={status}>
                            {status}
                          </option>
                        ),
                      )}
                    </select>
                  </label>
                </div>
                <div>
                  <label>
                    Type:{" "}
                    <select
                      value={form.type ?? ""}
                      onChange={(e) =>
                        updateField(
                          "type",
                          e.target.value as UpdateExtracurricularRequestType,
                        )
                      }
                    >
                      {Object.values(UpdateExtracurricularRequestType).map(
                        (type) => (
                          <option key={type} value={type}>
                            {type}
                          </option>
                        ),
                      )}
                    </select>
                  </label>
                </div>
                <div>
                  <label>
                    Description:
                    <br />
                    <textarea
                      value={form.description ?? ""}
                      onChange={(e) =>
                        updateField("description", e.target.value)
                      }
                    />
                  </label>
                </div>
                <div>
                  <label>
                    Leadership:{" "}
                    <input
                      value={form.leadership_role ?? ""}
                      onChange={(e) =>
                        updateField("leadership_role", e.target.value)
                      }
                    />
                  </label>
                </div>
                <div>
                  <label>
                    Start date:{" "}
                    <input
                      type="date"
                      value={form.start_date ?? ""}
                      onChange={(e) =>
                        updateField("start_date", e.target.value)
                      }
                    />
                  </label>
                </div>
                <div>
                  <label>
                    End date:{" "}
                    <input
                      type="date"
                      value={form.end_date ?? ""}
                      onChange={(e) => updateField("end_date", e.target.value)}
                    />
                  </label>
                </div>
                <div>
                  <label>
                    Organization:{" "}
                    <input
                      value={form.organization ?? ""}
                      onChange={(e) =>
                        updateField("organization", e.target.value)
                      }
                    />
                  </label>
                </div>
                <div>
                  <button onClick={saveEdit} disabled={saving}>
                    {saving ? "Saving..." : "Save"}
                  </button>
                  <button onClick={cancelEdit} disabled={saving}>
                    Cancel
                  </button>
                </div>
              </div>
            ) : (
              <div>
                <strong>{item.name}</strong> — {item.organization} —{" "}
                {item.status}
                <div>
                  <button onClick={() => startEdit(item)}>Edit</button>
                </div>
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
