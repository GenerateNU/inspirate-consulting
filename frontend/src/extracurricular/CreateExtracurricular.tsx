import { useState } from "react";

import { useCreateExtracurricular } from "../api/endpoints/extracurriculars/extracurriculars";
import {
  CreateExtracurricularRequestStatus,
  CreateExtracurricularRequestType,
  type CreateExtracurricularRequest,
} from "../api/models";

// TODO: replace with the logged-in student's id once auth is wired up
const TEST_STUDENT_ID = "00000000-0000-0000-0000-000000000002";

const emptyForm: CreateExtracurricularRequest = {
  name: "",
  status: CreateExtracurricularRequestStatus.doing,
  type: CreateExtracurricularRequestType.maintenance,
  description: "",
  leadership_role: undefined,
  start_date: "",
  end_date: undefined,
  organization: "",
};

export default function CreateExtracurricular() {
  const [form, setForm] = useState<CreateExtracurricularRequest>(emptyForm);
  const [msg, setMsg] = useState<string | null>(null);
  const { trigger, isMutating } = useCreateExtracurricular(TEST_STUDENT_ID);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >,
  ) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMsg(null);

    // Optional fields are sent only when filled in
    const payload: CreateExtracurricularRequest = {
      ...form,
      leadership_role: form.leadership_role || undefined,
      end_date: form.end_date || undefined,
    };

    try {
      const res = await trigger(payload);
      if (res.status === 200) {
        setMsg("Created id: " + res.data.id);
        setForm(emptyForm);
      } else {
        setMsg("Error: " + JSON.stringify(res.data));
      }
    } catch (err) {
      setMsg("Error: " + (err instanceof Error ? err.message : String(err)));
    }
  };

  return (
    <div>
      <h2>Create Extracurricular</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>
            Name:{" "}
            <input name="name" value={form.name} onChange={handleChange} />
          </label>
        </div>
        <div>
          <label>
            Status:
            <select name="status" value={form.status} onChange={handleChange}>
              {Object.values(CreateExtracurricularRequestStatus).map(
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
            Type:
            <select name="type" value={form.type} onChange={handleChange}>
              {Object.values(CreateExtracurricularRequestType).map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          </label>
        </div>
        <div>
          <label>
            Description:
            <br />
            <textarea
              name="description"
              value={form.description}
              onChange={handleChange}
            />
          </label>
        </div>
        <div>
          <label>
            Leadership role:{" "}
            <input
              name="leadership_role"
              value={form.leadership_role ?? ""}
              onChange={handleChange}
            />
          </label>
        </div>
        <div>
          <label>
            Start date:{" "}
            <input
              type="date"
              name="start_date"
              value={form.start_date}
              onChange={handleChange}
            />
          </label>
        </div>
        <div>
          <label>
            End date:{" "}
            <input
              type="date"
              name="end_date"
              value={form.end_date ?? ""}
              onChange={handleChange}
            />
          </label>
        </div>
        <div>
          <label>
            Organization:{" "}
            <input
              name="organization"
              value={form.organization}
              onChange={handleChange}
            />
          </label>
        </div>
        <div>
          <button type="submit" disabled={isMutating}>
            {isMutating ? "Creating..." : "Create"}
          </button>
        </div>
      </form>
      {msg && <p>{msg}</p>}
    </div>
  );
}
