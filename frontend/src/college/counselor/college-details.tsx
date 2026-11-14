import { useState, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";

import { errorMessage, type ApiError } from "../../api/client";
import type { GlobalCollege } from "../../api/models";
import arrowLeftIcon from "@/assets/college/arrow-left.svg";
import { cn } from "@/lib/utils";
import { useCollegeFromRoute } from "../utils/use-college-from-route";
import { cycleDateLabel, cycleLabel } from "../utils/shared";
import CyclesEditor, { editorCycleOrder } from "./cycles-editor";
import { useCollegeEditForm } from "./use-college-edit-form";

const primaryButtonClassName =
  "rounded-lg bg-dark-green px-4 py-2 font-host text-xs leading-[normal] text-white hover:bg-dark-green/90 disabled:cursor-not-allowed disabled:opacity-50";
const secondaryButtonClassName =
  "rounded-[7px] border border-greyish-brown bg-off-white px-[21px] py-[7px] font-host text-xs leading-[18px] text-dark-greyish-brown hover:bg-white disabled:cursor-not-allowed disabled:opacity-50";

// white card with a title (and optional subtitle) used for each section
function Panel({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle?: string;
  children: ReactNode;
}) {
  return (
    <section className="flex w-full flex-col gap-6 overflow-clip rounded-lg border border-cream bg-white p-6 shadow-card">
      <div className="flex flex-col gap-1 leading-[normal]">
        <h2 className="font-jakarta text-xl font-medium text-charcoal">
          {title}
        </h2>
        {subtitle && (
          <p className="font-host text-xs text-dark-greyish-brown">
            {subtitle}
          </p>
        )}
      </div>
      {children}
    </section>
  );
}

// header row: back arrow, school name, and the page's buttons on the right
function Header({ title, actions }: { title: string; actions?: ReactNode }) {
  const navigate = useNavigate();
  return (
    <div className="flex w-full items-start gap-9">
      <button
        type="button"
        data-slot="icon-button"
        aria-label="Back to colleges"
        onClick={() => navigate("/counselor-college-list")}
        className="shrink-0 rounded-md hover:bg-black/5"
      >
        <img src={arrowLeftIcon} alt="" width={32} height={32} />
      </button>
      <h1 className="min-w-0 flex-1 truncate font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
        {title}
      </h1>
      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </div>
  );
}

function SchoolView({ college }: { college: GlobalCollege }) {
  const offered = editorCycleOrder.flatMap((cycle) => {
    const saved = college.cycles?.find((c) => c.cycle === cycle);
    return saved ? [saved] : [];
  });

  return (
    <>
      <Panel title="School information">
        <div className="flex gap-6 font-host max-md:flex-col">
          {[
            ["School name", college.school_name],
            ["Location", college.school_location],
          ].map(([label, value]) => (
            <div key={label} className="flex min-w-0 flex-1 flex-col gap-1">
              <span className="text-xs leading-[normal] text-dark-greyish-brown">
                {label}
              </span>
              <span className="text-sm leading-5 text-charcoal">{value}</span>
            </div>
          ))}
        </div>
      </Panel>

      <Panel title="Application cycles & deadlines">
        {offered.length === 0 ? (
          <p className="font-host text-sm text-dark-greyish-brown">
            No cycles yet. Choose Edit to add the cycles this school offers.
          </p>
        ) : (
          <table className="w-full font-host">
            <thead>
              <tr className="text-left text-xs leading-[normal] text-dark-greyish-brown">
                <th className="w-[304px] pb-2 font-normal">Cycle</th>
                <th className="pb-2 font-normal">Deadline</th>
              </tr>
            </thead>
            <tbody>
              {offered.map((cycle) => (
                <tr
                  key={cycle.cycle}
                  className="border-t border-cream text-sm leading-5 text-charcoal"
                >
                  <td className="py-3 pr-6">{cycleLabel(cycle.cycle)}</td>
                  <td className="py-3">{cycleDateLabel(cycle) || "Open"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Panel>
    </>
  );
}

function TextField({
  label,
  value,
  disabled,
  onChange,
}: {
  label: string;
  value: string;
  disabled: boolean;
  onChange: (value: string) => void;
}) {
  return (
    <label className="flex min-w-0 flex-1 flex-col gap-2">
      <span className="font-host text-sm leading-5 text-charcoal">{label}</span>
      <input
        type="text"
        data-slot="text-input"
        value={value}
        required
        disabled={disabled}
        onChange={(e) => onChange(e.target.value)}
        className="w-full rounded-lg border border-greyish-brown bg-white px-4 py-3 font-host text-sm leading-5 text-charcoal outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-60"
      />
    </label>
  );
}

function SchoolEditor({
  college,
  onSaved,
  onExit,
}: {
  college: GlobalCollege;
  onSaved: () => Promise<unknown>;
  onExit: () => void;
}) {
  const form = useCollegeEditForm(college, onSaved);

  return (
    <form
      className="contents"
      onSubmit={(e) => {
        e.preventDefault();
        form.save(onExit);
      }}
    >
      <Header
        title={college.school_name}
        actions={
          <>
            <button
              type="button"
              data-slot="secondary-action"
              disabled={form.saving}
              onClick={() => {
                form.discard();
                onExit();
              }}
              className={secondaryButtonClassName}
            >
              Discard
            </button>
            <button
              type="submit"
              data-slot="primary-action"
              disabled={form.saving}
              className={primaryButtonClassName}
            >
              {form.saving ? "Saving..." : "Save Changes"}
            </button>
          </>
        }
      />

      <div className="flex w-full flex-col gap-5 pl-[67px] max-lg:pl-0">
        {form.isDuplicateWarning ? (
          // the backend's 409 message is developer-facing, so show our own wording
          <div
            role="alert"
            className="flex flex-wrap items-center justify-between gap-3 rounded-lg bg-lightest-yellow p-4 font-host text-sm text-dark-yellow"
          >
            &quot;{form.name}&quot; looks very similar to another college on the
            list. Check that it isn&apos;t a duplicate before saving.
            <button
              type="button"
              data-slot="secondary-action"
              disabled={form.saving}
              onClick={() => form.save(onExit, true)}
              className={secondaryButtonClassName}
            >
              Save anyway
            </button>
          </div>
        ) : (
          (form.validationError || form.error) && (
            <p role="alert" className="font-host text-sm text-danger">
              {form.validationError ?? `Error: ${errorMessage(form.error)}`}
            </p>
          )
        )}

        <Panel
          title="School information"
          subtitle="The name and location shown on every student’s college list."
        >
          <div className="flex gap-6 max-md:flex-col">
            <TextField
              label="School name"
              value={form.name}
              disabled={form.saving}
              onChange={form.setName}
            />
            <TextField
              label="Location"
              value={form.location}
              disabled={form.saving}
              onChange={form.setLocation}
            />
          </div>
        </Panel>

        <Panel
          title="Application cycles & deadlines"
          subtitle="Turn on the cycles this school offers, then set the deadline for each one."
        >
          <CyclesEditor
            draft={form.cycles}
            disabled={form.saving}
            onChange={form.setCycles}
          />
        </Panel>
      </div>
    </form>
  );
}

// a global college for counselors: read-only by default, Edit switches to the form
export default function CounselorCollegeDetails() {
  const { data: college, isLoading, error, refetch } = useCollegeFromRoute();
  const [editing, setEditing] = useState(false);

  let content: ReactNode;
  if (isLoading) {
    content = (
      <>
        <Header title="" />
        <p className="pl-[67px] font-host text-sm">Loading...</p>
      </>
    );
  } else if (error || !college) {
    content = (
      <>
        <Header title="School Details" />
        <p className="pl-[67px] font-host text-sm">
          {(error as ApiError | null)?.status === 404
            ? "College not found."
            : `Error: ${errorMessage(error)}`}
        </p>
      </>
    );
  } else if (editing) {
    content = (
      <SchoolEditor
        key={college.updated_at}
        college={college}
        onSaved={refetch}
        onExit={() => setEditing(false)}
      />
    );
  } else {
    content = (
      <>
        <Header
          title={college.school_name}
          actions={
            <button
              type="button"
              data-slot="primary-action"
              onClick={() => setEditing(true)}
              className={primaryButtonClassName}
            >
              Edit
            </button>
          }
        />
        <div className="flex w-full flex-col gap-8 pl-[67px] max-lg:pl-0">
          <SchoolView college={college} />
        </div>
      </>
    );
  }

  return (
    <div
      className={cn(
        "flex min-h-full w-full flex-col bg-off-white py-12 pr-20 pl-12 max-lg:px-6",
        editing ? "gap-[47px]" : "gap-5",
      )}
    >
      {content}
    </div>
  );
}
