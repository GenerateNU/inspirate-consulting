import { useState } from "react";

import { errorMessage } from "../api/client";
import {
  useListExtracurriculars,
  useUpdateExtracurricular,
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
  const { data, error, isLoading, refetch } = useListExtracurriculars();
  const update = useUpdateExtracurricular();
  const [editing, setEditing] = useState<number | null>(null);
  const [form, setForm] = useState<UpdateExtracurricularRequest>({});

  const list: Extracurricular[] = data ?? [];

  const startEdit = (item: Extracurricular) => {
    update.reset();
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

  const saveEdit = () => {
    if (editing === null) return;

    // Optional fields are sent only when filled in
    const payload: UpdateExtracurricularRequest = {
      ...form,
      leadership_role: form.leadership_role || undefined,
      end_date: form.end_date || undefined,
    };

    update.mutate(
      { id: editing, data: payload },
      {
        onSuccess: async () => {
          await refetch();
          cancelEdit();
        },
      },
    );
  };

  const updateField = <K extends keyof UpdateExtracurricularRequest>(
    key: K,
    value: UpdateExtracurricularRequest[K],
  ) => setForm((prev) => ({ ...prev, [key]: value }));

  return (
    <div>
      <h2>Extracurriculars</h2>
      {error && <p>Error: {errorMessage(error)}</p>}
      {update.error && <p>Error: {errorMessage(update.error)}</p>}
      <button onClick={() => refetch()}>Refresh</button>
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
                  <button onClick={saveEdit} disabled={update.isPending}>
                    {update.isPending ? "Saving..." : "Save"}
                  </button>
                  <button onClick={cancelEdit} disabled={update.isPending}>
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
