import { useState, type ReactNode } from "react";

import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../../api/models";
import { PencilIcon, PinIcon } from "../shared";
import {
  applicationTypes,
  categories,
  categoryClassName,
  categoryLabel,
  deadlineClassName,
  typeLabel,
} from "./application-labels";

export const smallButtonClassName =
  "text-xs text-gray-500 hover:text-black hover:underline disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:no-underline";

function Tag({
  className,
  children,
}: {
  className: string;
  children: ReactNode;
}) {
  return (
    <span className={`rounded-full border px-5 py-0.5 text-sm ${className}`}>
      {children}
    </span>
  );
}

// one application on My List
// the page owns which card is being edited and all the mutations, the card holds the fields
export default function ApplicationCard({
  app,
  college,
  collegeName,
  busy,
  saving,
  isEditing,
  isDragging,
  draggable,
  rankControls,
  onDragStart,
  onDragEnd,
  onDrop,
  onStartEdit,
  onCancelEdit,
  onSave,
  onRemove,
}: {
  app: PersonalCollegeApplication;
  college: GlobalCollege | undefined;
  collegeName: string;
  busy: boolean;
  saving: boolean;
  isEditing: boolean;
  isDragging: boolean;
  draggable: boolean;
  rankControls: ReactNode;
  onDragStart: () => void;
  onDragEnd: () => void;
  onDrop: () => void;
  onStartEdit: () => void;
  onCancelEdit: () => void;
  onSave: (applicationType: ApplicationType, category: Category) => void;
  onRemove: () => void;
}) {
  const [editType, setEditType] = useState<ApplicationType>(ApplicationType.RD);
  const [editCategory, setEditCategory] = useState<Category>(Category.target);

  const startEditing = () => {
    setEditType(app.application_type as ApplicationType);
    setEditCategory(app.category as Category);
    onStartEdit();
  };

  return (
    <li
      draggable={draggable}
      onDragStart={(e) => {
        e.dataTransfer.setData("text/plain", String(app.id));
        onDragStart();
      }}
      onDragEnd={onDragEnd}
      onDragOver={(e) => e.preventDefault()}
      onDrop={(e) => {
        e.preventDefault();
        onDrop();
      }}
      className={`flex items-start gap-4 border border-black px-6 py-8 ${
        draggable ? "cursor-grab" : ""
      } ${isDragging ? "opacity-50" : ""}`}
    >
      <span className="mt-1 w-8 shrink-0 text-xl text-gray-400">
        {app.rank ?? "–"}
      </span>

      <div className="min-w-0 flex-1">
        {/* TODO: link the name to the application's page once it exists */}
        <h3 className="text-3xl font-medium">{collegeName}</h3>
        {college && (
          <p className="mt-1 flex items-center gap-1 text-xs text-gray-400">
            <PinIcon />
            {college.school_location}
          </p>
        )}

        {isEditing ? (
          <div className="mt-3 flex flex-wrap items-center gap-2">
            <select
              aria-label="Application type"
              value={editType}
              onChange={(e) => setEditType(e.target.value as ApplicationType)}
              className={`rounded-full border px-3 py-0.5 text-sm ${deadlineClassName}`}
            >
              {applicationTypes.map((type) => (
                <option key={type} value={type}>
                  {typeLabel(type)}
                </option>
              ))}
            </select>
            <select
              aria-label="Category"
              value={editCategory}
              onChange={(e) => setEditCategory(e.target.value as Category)}
              className={`rounded-full border px-3 py-0.5 text-sm ${categoryClassName[editCategory]}`}
            >
              {categories.map((category) => (
                <option key={category} value={category}>
                  {categoryLabel(category)}
                </option>
              ))}
            </select>
          </div>
        ) : (
          <div className="mt-3 flex flex-wrap items-center gap-2">
            <Tag className={deadlineClassName}>
              {typeLabel(app.application_type)}
            </Tag>
            <Tag className={categoryClassName[app.category] ?? ""}>
              {categoryLabel(app.category)}
            </Tag>
          </div>
        )}
      </div>

      <div className="flex shrink-0 flex-col items-end gap-2">
        {isEditing ? (
          <div className="flex items-center gap-3">
            <button
              type="button"
              disabled={busy}
              onClick={() => onSave(editType, editCategory)}
              className={smallButtonClassName}
            >
              {saving ? "Saving..." : "Save"}
            </button>
            <button
              type="button"
              onClick={onCancelEdit}
              className={smallButtonClassName}
            >
              Cancel
            </button>
            <button
              type="button"
              disabled={busy}
              onClick={onRemove}
              className="text-xs text-red-600 hover:underline disabled:opacity-40"
            >
              Remove
            </button>
          </div>
        ) : (
          <>
            <button
              type="button"
              disabled={busy}
              onClick={startEditing}
              className={`flex items-center gap-1 ${smallButtonClassName}`}
            >
              <PencilIcon />
              Edit
            </button>
            {rankControls}
          </>
        )}
      </div>
    </li>
  );
}
