import { Link, NavLink, useParams } from "react-router-dom";
import { errorMessage } from "../api/client";
import { useListExtracurriculars } from "../api/endpoints";
import {
  categoryLabel,
  statusLabel,
  toDateInput,
  ec,
} from "./ExtracurricularShared";
import { Pencil } from "lucide-react";

export default function ViewExtracurricular() {
  const { id } = useParams();
  const { data, isLoading, error } = useListExtracurriculars();
  const list = data ?? [];
  const item = list.find((entry) => String(entry.id) === id);

  return (
    <div className="grid items-start gap-6 p-6 md:grid-cols-[minmax(14rem,18rem)_1fr]">
      <aside className="border border-current p-5">
        <h2 className="mb-3 text-[1.75rem]">Extracurriculars</h2>
        <Link to="/extracurriculars/new" className={ec.linkButton}>
          Add extracurricular
        </Link>
        <ul className={`${ec.list} mt-4`}>
          {list.map((entry) => (
            <li key={entry.id}>
              <NavLink
                to={`/extracurriculars/${entry.id}`}
                className={({ isActive }) =>
                  isActive ? `${ec.card} bg-gray-500/25` : ec.card
                }
              >
                <span className={ec.cardText}>
                  <span className={ec.cardTitle}>{entry.name}</span>
                  {entry.leadership_role && (
                    <span className={ec.cardSubtitle}>
                      {entry.leadership_role}
                    </span>
                  )}
                </span>
              </NavLink>
            </li>
          ))}
        </ul>
      </aside>

      <section className="border border-current p-5">
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
            <div className={ec.header}>
              <Link
                to="/extracurriculars"
                className={ec.back}
                aria-label="Back to extracurriculars"
              >
                ←
              </Link>
              <h2 className={ec.title}>{item.name}</h2>
              <Link
                to={`/extracurriculars/${item.id}/edit`}
                className={ec.back}
                aria-label={`Edit ${item.name}`}
              >
                <Pencil size={20} aria-hidden="true" />
              </Link>
            </div>

            <div className={ec.section}>
              <h3 className={ec.sectionHeading}>Details</h3>
              <dl className="grid gap-4 md:grid-cols-3">
                <div>
                  <dt className="text-xs">Position / role</dt>
                  <dd className="font-semibold">
                    {item.leadership_role || "Not set"}
                  </dd>
                </div>
                <div>
                  <dt className="text-xs">Organization</dt>
                  <dd className="font-semibold">{item.organization}</dd>
                </div>
                <div>
                  <dt className="text-xs">Category</dt>
                  <dd className="font-semibold">{categoryLabel(item.type)}</dd>
                </div>
                <div>
                  <dt className="text-xs">Status</dt>
                  <dd className="font-semibold">{statusLabel(item.status)}</dd>
                </div>
                <div>
                  <dt className="text-xs">Dates</dt>
                  <dd className="font-semibold">
                    {toDateInput(item.start_date)} to{" "}
                    {item.end_date ? toDateInput(item.end_date) : "present"}
                  </dd>
                </div>
              </dl>
            </div>

            <div className={ec.section}>
              <h3 className={ec.sectionHeading}>Application description</h3>
              <p>{item.description}</p>
            </div>
          </>
        )}
      </section>
    </div>
  );
}
