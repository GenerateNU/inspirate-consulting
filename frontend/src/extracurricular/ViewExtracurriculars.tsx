import { Link, NavLink, useParams } from "react-router-dom";
import { errorMessage } from "../api/client";
import { useListExtracurriculars } from "../api/endpoints";
import {
  CATEGORY_LABELS,
  STATUS_LABELS,
  toDateInput,
  type ExtracurricularFormValues,
} from "./ExtracurricularShared";
import { PencilIcon } from "./Icons";
import "./Extracurriculars.css";

export default function ViewExtracurricular() {
  const { id } = useParams();
  const { data, isLoading, error } = useListExtracurriculars();
  const list = data ?? [];
  const item = list.find((entry) => String(entry.id) === id);

  return (
    <div className="ec-split">
      <aside className="ec-sidebar">
        <h2 className="ec-title">Extracurriculars</h2>
        <Link to="/extracurriculars/new" className="ec-link-button">
          Add extracurricular
        </Link>
        <ul className="ec-list">
          {list.map((entry) => (
            <li key={entry.id}>
              <NavLink to={`/extracurriculars/${entry.id}`} className="ec-card">
                <span className="ec-card-text">
                  <span className="ec-card-title">{entry.name}</span>
                  {entry.leadership_role && (
                    <span className="ec-card-subtitle">{entry.leadership_role}</span>
                  )}
                </span>
              </NavLink>
            </li>
          ))}
        </ul>
      </aside>

      <section className="ec-detail">
        {isLoading && <p>Loading...</p>}
        {error && <p>Error: {errorMessage(error)}</p>}
        {!isLoading && !error && !item && (
          <p>
            This extracurricular doesn't exist or was removed.{" "}
            <Link to="/extracurriculars">Back to extracurriculars</Link>
          </p>
        )}

        {item && (
          <>
            <div className="ec-header">
              <Link to="/extracurriculars" className="ec-back" aria-label="Back to extracurriculars">
                ←
              </Link>
              <h2 className="ec-title">{item.name} - View Only</h2>
              <Link
                to={`/extracurriculars/${item.id}/edit`}
                className="ec-back"
                aria-label={`Edit ${item.name}`}
              >
                <PencilIcon />
              </Link>
            </div>

            <div className="ec-section">
              <h3>Details</h3>
              <dl className="ec-details">
                <div>
                  <dt>Position / role</dt>
                  <dd>{item.leadership_role || "Not set"}</dd>
                </div>
                <div>
                  <dt>Organization</dt>
                  <dd>{item.organization}</dd>
                </div>
                <div>
                  <dt>Category</dt>
                  <dd>
                    {CATEGORY_LABELS[item.type as ExtracurricularFormValues["type"]] ??
                      item.type}
                  </dd>
                </div>
                <div>
                  <dt>Status</dt>
                  <dd>
                    {STATUS_LABELS[item.status as ExtracurricularFormValues["status"]] ??
                      item.status}
                  </dd>
                </div>
                <div>
                  <dt>Dates</dt>
                  <dd>
                    {toDateInput(item.start_date)} to{" "}
                    {item.end_date ? toDateInput(item.end_date) : "present"}
                  </dd>
                </div>
              </dl>
            </div>

            <div className="ec-section">
              <h3>Application description</h3>
              <p>{item.description}</p>
            </div>
          </>
        )}
      </section>
    </div>
  );
}
