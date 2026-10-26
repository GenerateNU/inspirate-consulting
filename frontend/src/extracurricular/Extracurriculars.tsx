import { Link } from "react-router-dom";
import { errorMessage } from "../api/client";
import { useListExtracurriculars } from "../api/endpoints";
import { cardSubtitle } from "./ExtracurricularShared";
import { EyeIcon } from "./Icons";
import "./Extracurriculars.css";

export default function Extracurriculars() {
  const { data, isLoading, error } = useListExtracurriculars();
  const list = data ?? [];

  return (
    <div className="ec-page">
      <div className="ec-header">
        <h2 className="ec-title">Extracurriculars</h2>
        <Link
          to="/extracurriculars/new"
          className="ec-link-button ec-header-push"
        >
          Add extracurricular
        </Link>
      </div>

      <div className="ec-box">
        {isLoading && <p>Loading...</p>}
        {error && <p>Error: {errorMessage(error)}</p>}
        {!isLoading && !error && list.length === 0 && (
          <p>
            No extracurriculars yet.{" "}
            <Link to="/extracurriculars/new">Add your first one</Link>.
          </p>
        )}

        <ul className="ec-list">
          {list.map((item) => (
            <li key={item.id}>
              <Link to={`/extracurriculars/${item.id}`} className="ec-card">
                <span className="ec-card-text">
                  <span className="ec-card-title">{item.name}</span>
                  <span className="ec-card-subtitle">{cardSubtitle(item)}</span>
                </span>
                <span className="ec-icon" aria-label={`View ${item.name}`}>
                  <EyeIcon />
                </span>
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
