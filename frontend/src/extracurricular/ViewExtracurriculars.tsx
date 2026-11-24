import type { ReactNode } from "react";
import { Link, useParams } from "react-router-dom";
import type { LucideIcon } from "lucide-react";
import {
  ArrowLeft,
  Building2,
  Calendar,
  CalendarCheck,
  CalendarDays,
  ChevronRight,
  CircleCheck,
  Goal,
  Palette,
  Paperclip,
  Pencil,
  SquareUser,
  Star,
  StickyNote,
  Trophy,
  Zap,
} from "lucide-react";
import { errorMessage } from "../api/client";
import { useGetExtracurricular } from "../api/endpoints";
import type { ExtracurricularDetail } from "../api/models";
import { cn } from "cn";
import { activityPageClass } from "./ActivityForm";
import {
  AchievementTile,
  CategoryIcon,
  NoteRow,
  SectionCard,
  Toggle,
} from "./sections";
import {
  APPLICATION_WORD_LIMIT,
  categoryLabel,
  edit,
  longDateLabel,
  statusLabel,
  wordCount,
} from "./ExtracurricularShared";

const byNewest = <T,>(items: T[] | null, key: (item: T) => string) =>
  [...(items ?? [])].sort((a, b) => key(b).localeCompare(key(a)));

// Read-only stand-in for the edit page's textareas.
const readOnlyBox =
  "rounded-lg border border-[#bdb7ad] bg-white px-4 py-3 text-sm leading-5 break-words whitespace-pre-wrap text-charcoal";

export default function ViewExtracurricular() {
  const { id } = useParams();
  const numericId = Number(id);
  const { data, isLoading, error } = useGetExtracurricular(numericId, {
    query: { enabled: Number.isInteger(numericId) },
  });

  if (isLoading) return <p className={activityPageClass}>Loading...</p>;
  if (error || !data) {
    return (
      <div className={activityPageClass}>
        <p>
          {error && error.status !== 404
            ? `Error: ${errorMessage(error)}`
            : "This extracurricular doesn't exist or was removed."}
        </p>
        <Link to="/extracurriculars" className="underline">
          Back to extracurriculars
        </Link>
      </div>
    );
  }

  return <ActivityView item={data} />;
}

// Figma: Extracurricular - View Activity - Student View.
function ActivityView({ item }: { item: ExtracurricularDetail }) {
  const achievements = byNewest(item.achievements, (a) => a.entry_date);
  const notes = byNewest(item.notes, (n) => n.created_at);
  const description = item.application_description ?? "";

  return (
    <div className={activityPageClass}>
      <nav
        aria-label="Breadcrumb"
        className="flex items-center gap-2 text-xs text-[#303436]"
      >
        <Link
          to="/extracurriculars"
          className="flex items-center gap-2 hover:underline"
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
          Extracurriculars
        </Link>
        <ChevronRight className="size-4" aria-hidden="true" />
        <span aria-current="page">{item.name}</span>
      </nav>

      <div className="flex min-w-0 items-center gap-4">
        <div className="flex min-w-0 items-center gap-6">
          <CategoryIcon icon={Palette} />
          <h1 className={cn(edit.pageTitle, "truncate")}>{item.name}</h1>
        </div>
        <Link
          to={`/extracurriculars/${item.id}/edit`}
          className="inline-flex size-8 shrink-0 items-center justify-center rounded-md text-charcoal hover:bg-accent focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-none"
          aria-label={`Edit ${item.name}`}
        >
          <Pencil className="size-6" aria-hidden="true" />
        </Link>
      </div>

      <h2 className={edit.cardTitle}>Details</h2>
      <dl className="grid gap-4 sm:grid-cols-2 md:grid-cols-3">
        <DetailCard icon={Star} label="Activity / Sport Name">
          {item.name}
        </DetailCard>
        <DetailCard icon={SquareUser} label="Leadership Position">
          {item.leadership_role}
        </DetailCard>
        <DetailCard icon={Building2} label="Organization">
          {item.organization}
        </DetailCard>
        <DetailCard icon={Zap} label="Hours Per Week">
          {item.hours_per_week}
        </DetailCard>
        <DetailCard icon={Star} label="Activity Category">
          {categoryLabel(item.type)}
        </DetailCard>
        <DetailCard icon={Calendar} label="Weeks Per Year">
          {item.weeks_per_year}
        </DetailCard>
        <DetailCard icon={CircleCheck} label="Status">
          {statusLabel(item.status)}
        </DetailCard>
        <DetailCard icon={CalendarDays} label="Start Date">
          {longDateLabel(item.start_date, { calendarDate: true })}
        </DetailCard>
        <DetailCard icon={CalendarCheck} label="End Date">
          {item.end_date
            ? longDateLabel(item.end_date, { calendarDate: true })
            : item.status === "doing"
              ? "Present"
              : undefined}
        </DetailCard>
      </dl>

      <SectionCard icon={Trophy} title="Achievements & Awards">
        {achievements.length > 0 ? (
          <ul className="flex flex-wrap gap-4">
            {achievements.map((achievement) => (
              <AchievementTile key={achievement.id} achievement={achievement} />
            ))}
          </ul>
        ) : (
          <EmptyText>No achievements or awards yet.</EmptyText>
        )}
      </SectionCard>

      <div className="grid gap-4 md:grid-cols-2">
        <SectionCard icon={StickyNote} title="Notes and Scratch Work">
          {notes.length > 0 ? (
            <ul className="flex flex-col gap-3">
              {notes.map((note) => (
                <NoteRow key={note.id} note={note} />
              ))}
            </ul>
          ) : (
            <EmptyText>No notes yet.</EmptyText>
          )}
        </SectionCard>
        <SectionCard icon={Goal} title="Goals">
          <p className={cn(readOnlyBox, "min-h-[100px] flex-1")}>
            {item.goals || <EmptyText>No goals yet.</EmptyText>}
          </p>
        </SectionCard>
      </div>

      <SectionCard icon={Paperclip} title="Application Information">
        <div className="flex flex-col gap-2">
          <p className={edit.label}>
            Application description - up to {APPLICATION_WORD_LIMIT} words
            (Common App)
          </p>
          <p className={cn(readOnlyBox, "min-h-[176px]")}>
            {description || (
              <EmptyText>No application description yet.</EmptyText>
            )}
          </p>
          <p className="text-right text-xs text-[#303436]">
            {wordCount(description)} / {APPLICATION_WORD_LIMIT} words
          </p>
        </div>
        <hr className="border-cream" />
        <div className="flex items-center justify-between gap-4">
          <span
            id="ec-view-ready-label"
            className="text-base leading-5 text-charcoal"
          >
            Ready For College Application?
          </span>
          <Toggle
            labelledBy="ec-view-ready-label"
            checked={item.ready_for_application}
          />
        </div>
      </SectionCard>
    </div>
  );
}

function DetailCard({
  icon: Icon,
  label,
  children,
}: {
  icon: LucideIcon;
  label: string;
  children: ReactNode;
}) {
  const empty = children === undefined || children === null || children === "";
  return (
    <div className="flex items-center gap-3 rounded-lg border border-cream bg-white p-4">
      <span
        aria-hidden="true"
        className="flex size-8 shrink-0 items-center justify-center rounded bg-[#f8efe4] text-[#786246]"
      >
        <Icon className="size-5" strokeWidth={1.75} />
      </span>
      <div className="flex min-w-0 flex-col gap-1 text-charcoal">
        <dt className="text-xs">{label}</dt>
        <dd
          className={cn(
            "text-base leading-5 break-words",
            empty && "text-[#5e5c56]",
          )}
        >
          {empty ? "Not set" : children}
        </dd>
      </div>
    </div>
  );
}

function EmptyText({ children }: { children: ReactNode }) {
  return <span className="text-sm text-[#5e5c56]">{children}</span>;
}
