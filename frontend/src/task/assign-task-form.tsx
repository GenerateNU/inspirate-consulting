import { useState } from "react";
import { useSearchParams } from "react-router-dom";
import AssignToSelect, { type AssignToValue } from "./assign-to-select";
import ConnectToSelect, { type ConnectToValue } from "./connect-to-select";

export default function AssignTaskForm() {
  // If we came from a student's page, preselect that student
  const [searchParams] = useSearchParams();
  const preselectedStudentId = searchParams.get("studentId");

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [assignTo, setAssignTo] = useState<AssignToValue>({
    mode: preselectedStudentId ? "individual" : "",
    grades: [],
    studentIds: preselectedStudentId ? [preselectedStudentId] : [],
  });
  const [connectTo, setConnectTo] = useState<ConnectToValue>({
    type: "",
    itemIds: [],
  });

  return (
    <div className="flex flex-col gap-4 p-10">
      <h1 className="text-2xl font-semibold">Assign New Task</h1>

      <input
        type="text"
        placeholder="Task Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        className="border px-3 py-2"
      />

      <textarea
        placeholder="Task Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        className="border px-3 py-2"
      />

      <div className="flex gap-6">
        {/* Left: due date */}
        <label className="flex w-64 flex-col text-sm">
          Select Due Date
          <input
            type="date"
            value={dueDate}
            onChange={(e) => setDueDate(e.target.value)}
            className="mt-1 border px-2 py-1"
          />
        </label>

        {/* Right: assign to / connect to */}
        <div className="flex flex-1 flex-col gap-4">
          <AssignToSelect value={assignTo} onChange={setAssignTo} />
          <ConnectToSelect value={connectTo} onChange={setConnectTo} />
        </div>
      </div>

      {/* PLACEHOLDER: submitting needs a backend endpoint that assigns a task to the selected students */}
      <div className="flex flex-col items-end gap-1">
        <button
          type="button"
          disabled
          className="cursor-not-allowed rounded bg-black px-6 py-2 text-white opacity-50"
        >
          Assign Task →
        </button>
        <p className="text-xs text-gray-500">
          Assigning tasks isn&apos;t available yet. It needs a new backend
          endpoint.
        </p>
      </div>
    </div>
  );
}
