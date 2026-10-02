import { useState } from "react";

const API_BASE_URL = "http://localhost:8080";

type Essay = {
  id: string;
  student_id: string;
  type: string;
  college_id: number | null;
  link_to_content: string;
  status: "Submitted" | "Draft" | "Review" | "Archived";
};

export default function Essay() {
  const [studentId, setStudentId] = useState("");
  const [type, setType] = useState("");
  const [collegeId, setCollegeId] = useState("");
  const [linkToContent, setLinkToContent] = useState("");

  const [essayId, setEssayId] = useState("");
  const [status, setStatus] = useState<Essay["status"]>("Draft");

  const [essays, setEssays] = useState<Essay[]>([]);
  const [response, setResponse] = useState("");
  const [error, setError] = useState("");

  const getEssays = async () => {
    setError("");
    setResponse("");

    try {
      const result = await fetch(
        `${API_BASE_URL}/students/${studentId.trim()}/essays`,
      );

      const data = await result.json();

      if (!result.ok) {
        throw new Error(JSON.stringify(data));
      }

      setEssays(data.essays ?? []);
      setResponse(JSON.stringify(data, null, 2));
    } catch (err) {
      console.error(err);
      setError("Something went wrong.");
    }
  };

  const createEssay = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    setError("");
    setResponse("");

    try {
      const result = await fetch(`${API_BASE_URL}/essays`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          student_id: studentId.trim(),
          type,
          college_id: collegeId.trim() === "" ? null : Number(collegeId),
          link_to_content: linkToContent,
        }),
      });

      const data = await result.json().catch(() => ({}));

      if (!result.ok) {
        throw new Error(JSON.stringify(data));
      }

      setResponse(JSON.stringify(data, null, 2));

      await getEssays();
    } catch (err) {
      console.error(err);
      setError("Something went wrong.");
    }
  };

  const updateStatus = async () => {
    setError("");
    setResponse("");

    try {
      const result = await fetch(
        `${API_BASE_URL}/essays/${essayId.trim()}/status`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            status,
          }),
        },
      );

      const data = await result.json().catch(() => ({}));

      if (!result.ok) {
        throw new Error(JSON.stringify(data));
      }

      setResponse(JSON.stringify(data, null, 2));

      await getEssays();
    } catch (err) {
      console.error(err);
      setError("Something went wrong.");
    }
  };

  return (
    <div>
      <h1>Essays</h1>

      <hr />

      <h2>Get Essays From Student</h2>

      <input
        type="text"
        placeholder="Student ID"
        value={studentId}
        onChange={(e) => setStudentId(e.target.value)}
      />

      <button type="button" onClick={getEssays}>
        Get Essays
      </button>

      <hr />

      <h2>Create Essay</h2>

      <form onSubmit={createEssay}>
        <div>
          <label>Student ID: </label>
          <input
            type="text"
            value={studentId}
            onChange={(e) => setStudentId(e.target.value)}
          />
        </div>

        <div>
          <label>Type: </label>
          <input
            type="text"
            value={type}
            onChange={(e) => setType(e.target.value)}
          />
        </div>

        <div>
          <label>College ID: </label>
          <input
            type="text"
            value={collegeId}
            onChange={(e) => setCollegeId(e.target.value)}
          />
        </div>

        <div>
          <label>Link to Content: </label>
          <input
            type="text"
            value={linkToContent}
            onChange={(e) => setLinkToContent(e.target.value)}
          />
        </div>

        <button type="submit">Create Essay</button>
      </form>

      <hr />

      <h2>Update Essay Status</h2>

      <div>
        <label>Essay ID: </label>
        <input
          type="text"
          placeholder="Essay ID"
          value={essayId}
          onChange={(e) => setEssayId(e.target.value)}
        />
      </div>

      <div>
        <label>Status: </label>
        <select
          value={status}
          onChange={(e) => setStatus(e.target.value as Essay["status"])}
        >
          <option value="Draft">Draft</option>
          <option value="Submitted">Submitted</option>
          <option value="Review">Review</option>
          <option value="Archived">Archived</option>
        </select>
      </div>

      <button type="button" onClick={updateStatus}>
        Update Status
      </button>

      <hr />

      <h2>Essays</h2>

      {essays.length === 0 ? (
        <p>No essays found.</p>
      ) : (
        <ul>
          {essays.map((essay) => (
            <li key={essay.id}>
              <p>ID: {essay.id}</p>
              <p>Student ID: {essay.student_id}</p>
              <p>Type: {essay.type}</p>
              <p>College ID: {essay.college_id ?? "None"}</p>
              <p>Link: {essay.link_to_content}</p>
              <p>Status: {essay.status}</p>
            </li>
          ))}
        </ul>
      )}

      {error && (
        <div>
          <h3>Error</h3>
          <pre>{error}</pre>
        </div>
      )}

      {response && (
        <div>
          <h3>API Response</h3>
          <pre>{response}</pre>
        </div>
      )}
    </div>
  );
}
