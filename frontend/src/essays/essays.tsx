import { useState } from "react";

import { errorMessage } from "../api/client";
import {
  useCreateEssay,
  useGetEssaysFromStudent,
  useUpdateEssayStatus,
} from "../api/endpoints";
import { UpdateStatusBodyStatus } from "../api/models";

export default function Essay() {
  const [studentId, setStudentId] = useState("");
  const [submittedStudentId, setSubmittedStudentId] = useState("");
  const [type, setType] = useState("");
  const [collegeId, setCollegeId] = useState("");
  const [linkToContent, setLinkToContent] = useState("");

  const [essayId, setEssayId] = useState("");
  const [status, setStatus] = useState<UpdateStatusBodyStatus>("Draft");

  const [response, setResponse] = useState("");

  const essaysQuery = useGetEssaysFromStudent(submittedStudentId, {
    query: { enabled: !!submittedStudentId },
  });
  const create = useCreateEssay();
  const update = useUpdateEssayStatus();

  const essays = (submittedStudentId && essaysQuery.data?.essays) || [];
  const error =
    (submittedStudentId && essaysQuery.error) || create.error || update.error;

  const refreshEssays = () => {
    if (submittedStudentId) void essaysQuery.refetch();
  };

  const getEssays = () => {
    setResponse("");
    create.reset();
    update.reset();
    const trimmed = studentId.trim();
    if (trimmed === submittedStudentId) refreshEssays();
    else setSubmittedStudentId(trimmed);
  };

  const createEssay = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setResponse("");
    update.reset();
    create.mutate(
      {
        data: {
          type,
          college_id: collegeId.trim() === "" ? null : Number(collegeId),
          link_to_content: linkToContent,
        },
      },
      {
        onSuccess: () => {
          setResponse("Essay created.");
          refreshEssays();
        },
      },
    );
  };

  const updateStatus = () => {
    setResponse("");
    create.reset();
    update.mutate(
      { essayId: essayId.trim(), data: { status } },
      {
        onSuccess: (essay) => {
          setResponse(JSON.stringify(essay, null, 2));
          refreshEssays();
        },
      },
    );
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

        <button type="submit" disabled={create.isPending}>
          Create Essay
        </button>
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
          onChange={(e) => setStatus(e.target.value as UpdateStatusBodyStatus)}
        >
          {Object.values(UpdateStatusBodyStatus).map((value) => (
            <option key={value} value={value}>
              {value}
            </option>
          ))}
        </select>
      </div>

      <button type="button" onClick={updateStatus} disabled={update.isPending}>
        Update Status
      </button>

      <hr />

      <h2>Essays</h2>

      {essaysQuery.isLoading ? (
        <p>Loading...</p>
      ) : essays.length === 0 ? (
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
          <pre>{errorMessage(error)}</pre>
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
