import { Link } from "react-router-dom";
import { errorMessage } from "../api/client";
import { useListExtracurriculars } from "../api/endpoints";
import { cardSubtitle, ec } from "./ExtracurricularShared";
import { Eye } from "lucide-react";

export default function Extracurriculars() {
  const { data, isLoading, error } = useListExtracurriculars();
  const list = data ?? [];

  return (
    <div className={ec.page}>
      <div className={ec.header}>
        <h2 className={ec.title}>Extracurriculars</h2>
        <Link to="/extracurriculars/new" className={`${ec.linkButton} ml-auto`}>
          Add extracurricular
        </Link>
      </div>

      <div className={ec.box}>
        {isLoading && <p>Loading...</p>}
        {error && <p>Error: {errorMessage(error)}</p>}
        {!isLoading && !error && list.length === 0 && (
          <p>
            No extracurriculars yet.{" "}
            <Link to="/extracurriculars/new">Add your first one</Link>.
          </p>
        )}

        <ul className={ec.list}>
          {list.map((item) => (
            <li key={item.id}>
              <Link to={`/extracurriculars/${item.id}`} className={ec.card}>
                <span className={ec.cardText}>
                  <span className={ec.cardTitle}>{item.name}</span>
                  <span className={ec.cardSubtitle}>{cardSubtitle(item)}</span>
                </span>
                <span
                  className="inline-flex size-8 shrink-0 items-center justify-center rounded-full border border-current"
                  aria-label={`View ${item.name}`}
                >
                  <Eye size={16} aria-hidden="true" />
                </span>
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
