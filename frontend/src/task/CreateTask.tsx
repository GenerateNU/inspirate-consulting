import { useState } from "react";
import { errorMessage } from "../api/client";
import { useCreateTodoItem } from "../api/endpoints";
import type { CreateTodoItemRequestBody } from "../api/models";
import {
  DEFAULT_DUE_TIME,
  LINK_TYPE_LABELS,
  useTodoLinkOptions,
  type LinkType,
} from "./todoShared";

export default function CreateTask() {
  const [description, setDescription] = useState("");
  const [deadline, setDeadline] = useState("");
  const [deadlineTime, setDeadlineTime] = useState("");
  const [linkType, setLinkType] = useState<LinkType | "">("");
  const [linkId, setLinkId] = useState("");
  const [status, setStatus] = useState("");

  const { options } = useTodoLinkOptions();
  const { mutate, isPending, error } = useCreateTodoItem();

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus("");

    const body: CreateTodoItemRequestBody = { todo_description: description };
    if (deadline)
      body.deadline = new Date(
        `${deadline}T${deadlineTime || DEFAULT_DUE_TIME}`,
      ).toISOString();
    if (linkType === "essay") body.essay_id = linkId;
    if (linkType === "media") body.media_id = linkId;
    if (linkType === "college") body.global_college_id = Number(linkId);

    mutate(
      { data: body },
      {
        onSuccess: (todo) => {
          setStatus(`Created: ${todo.todo_description}`);
          setDescription("");
          setDeadline("");
          setDeadlineTime("");
          setLinkType("");
          setLinkId("");
        },
      },
    );
  };

  return (
    <div className="form-container">
      <h2>Create Task</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Description:</label>
          <input
            type="text"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Due date (optional):</label>
          <input
            type="date"
            value={deadline}
            onChange={(e) => setDeadline(e.target.value)}
          />{" "}
          <input
            type="time"
            aria-label="Due time (optional)"
            value={deadlineTime}
            onChange={(e) => setDeadlineTime(e.target.value)}
            disabled={!deadline}
          />
        </div>
        <div className="form-group">
          <label>Link to (optional):</label>
          <select
            value={linkType}
            onChange={(e) => {
              setLinkType(e.target.value as LinkType | "");
              setLinkId("");
            }}
          >
            <option value="">Nothing</option>
            {(Object.keys(LINK_TYPE_LABELS) as LinkType[]).map((type) => (
              <option key={type} value={type}>
                {LINK_TYPE_LABELS[type]}
              </option>
            ))}
          </select>{" "}
          {linkType && (
            <select
              value={linkId}
              onChange={(e) => setLinkId(e.target.value)}
              required
            >
              <option value="">
                Choose {LINK_TYPE_LABELS[linkType].toLowerCase()}...
              </option>
              {options[linkType].map((o) => (
                <option key={o.id} value={o.id}>
                  {o.label}
                </option>
              ))}
            </select>
          )}
        </div>
        <button type="submit" disabled={isPending}>
          {isPending ? "Creating Task..." : "Create Task"}
        </button>
      </form>
      {status && <p>{status}</p>}
      {error && <p>Error: {errorMessage(error)}</p>}
    </div>
  );
}
