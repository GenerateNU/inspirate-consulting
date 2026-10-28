import { useState, type ReactNode } from "react";
import { useQueries } from "@tanstack/react-query";
import {
  ChevronDown,
  ChevronUp,
  Clock,
  Folder,
  FolderOpen,
  GraduationCap,
} from "lucide-react";

import type { ApiError } from "../api/client";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  getGetEssayReviewStatusQueryOptions,
  getGetGlobalCollegeQueryOptions,
  useGetGlobalCollege,
} from "../api/endpoints";
import type {
  EssayGroups,
  EssayReviewStatus,
  Essays,
  GlobalCollege,
} from "../api/models";

export type ReviewMap = Record<string, EssayReviewStatus | undefined>;

export function useReviewStatuses(essays: Essays[]) {
  const results = useQueries({
    queries: essays.map((essay) =>
      getGetEssayReviewStatusQueryOptions(essay.id),
    ),
  });

  const reviews: ReviewMap = {};
  essays.forEach((essay, i) => {
    reviews[essay.id] = results[i]?.data;
  });
  const error = results.find(
    (r) => r.error && (r.error as ApiError).status !== 404,
  )?.error;

  return { reviews, error };
}

// Review state comes only from the review ledger; essays.status is not used
export const isReviewOpen = (review?: EssayReviewStatus) =>
  review?.status === "open";

export const isReviewCompleted = (review?: EssayReviewStatus) =>
  review?.status === "completed";

export function reviewStats(essays: Essays[], reviews: ReviewMap) {
  const open = essays.filter((e) => isReviewOpen(reviews[e.id])).length;
  const completed = essays.filter((e) =>
    isReviewCompleted(reviews[e.id]),
  ).length;
  return {
    total: essays.length,
    open,
    completed,
    notReviewed: essays.length - open - completed,
  };
}

export function StatGrid({
  stats,
}: {
  stats: { label: string; value: number; highlight?: boolean }[];
}) {
  return (
    <dl className="grid grid-cols-2 gap-4 sm:grid-cols-4">
      {stats.map(({ label, value, highlight }) => (
        <Card key={label} className="gap-2 p-4">
          <dt
            className={`text-xs font-semibold tracking-wide uppercase ${
              highlight && value > 0
                ? "text-destructive"
                : "text-muted-foreground"
            }`}
          >
            {label}
          </dt>
          <dd className="text-3xl font-semibold">{value}</dd>
        </Card>
      ))}
    </dl>
  );
}

const applicationTypeLabels: Record<string, string> = {
  EA: "Early Action",
  ED: "Early Decision",
  RD: "Regular Decision",
};

const deadlineFor = (college: GlobalCollege, applicationType: string) =>
  ({
    EA: college.ea_deadline,
    ED: college.ed_deadline,
    RD: college.rd_deadline,
  })[applicationType];

// Deadlines are calendar dates, so format in UTC to avoid shifting a day back
const formatDate = (value: string) =>
  new Date(value).toLocaleDateString(undefined, {
    month: "short",
    day: "2-digit",
    year: "numeric",
    timeZone: "UTC",
  });

export function useCollegeNames(collegeIds: number[]) {
  const results = useQueries({
    queries: collegeIds.map((id) => getGetGlobalCollegeQueryOptions(id)),
  });
  return Object.fromEntries(
    collegeIds.map((id, i) => [
      id,
      results[i]?.data?.school_name ?? `College #${id}`,
    ]),
  );
}

function CollegeHeading({
  collegeId,
  applicationType,
}: {
  collegeId: number;
  applicationType?: string;
}) {
  const { data: college } = useGetGlobalCollege(collegeId);
  const deadline =
    college && applicationType && deadlineFor(college, applicationType);

  return (
    <div className="flex flex-wrap items-center justify-between gap-2 rounded-md bg-muted px-4 py-3">
      <h3 className="flex items-center gap-2 text-lg font-medium">
        <GraduationCap className="size-5" aria-hidden />
        {college?.school_name ?? `College #${collegeId}`}
      </h3>
      <div className="flex items-center gap-3 text-xs font-semibold tracking-wide text-muted-foreground uppercase">
        {applicationType && (
          <Badge variant="outline">
            {applicationTypeLabels[applicationType] ?? applicationType}
          </Badge>
        )}
        {deadline && <span>Due: {formatDate(deadline)}</span>}
      </div>
    </div>
  );
}

export function RequestedAt({ review }: { review: EssayReviewStatus }) {
  return (
    <span className="flex items-center gap-1 text-xs text-muted-foreground">
      <Clock className="size-3.5" aria-hidden />
      Requested {new Date(review.requested_at).toLocaleString()}
    </span>
  );
}

function EssayCard({
  essay,
  badges,
  actions,
}: {
  essay: Essays;
  badges: ReactNode;
  actions: ReactNode;
}) {
  return (
    <li>
      <Card className="flex-col gap-4 p-5 shadow-none sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">{badges}</div>
          <h4 className="mt-2 text-lg font-medium">{essay.type}</h4>
        </div>
        <div className="flex shrink-0 flex-wrap items-center gap-2">
          {actions}
        </div>
      </Card>
    </li>
  );
}

function Section({
  title,
  count,
  children,
}: {
  title: string;
  count: number;
  children: ReactNode;
}) {
  const [open, setOpen] = useState(true);

  return (
    <Collapsible open={open} onOpenChange={setOpen} asChild>
      <Card className="gap-0 overflow-hidden py-0">
        <CollapsibleTrigger className="flex w-full cursor-pointer items-center justify-between bg-muted px-6 py-4 text-left hover:bg-accent">
          <span className="flex items-center gap-3 text-lg font-medium tracking-wide uppercase">
            {open ? (
              <FolderOpen className="size-5" aria-hidden />
            ) : (
              <Folder className="size-5" aria-hidden />
            )}
            {title}
            <Badge variant="outline" className="bg-background">
              {count} {count === 1 ? "essay" : "essays"}
            </Badge>
          </span>
          {open ? (
            <ChevronUp className="size-5" aria-hidden />
          ) : (
            <ChevronDown className="size-5" aria-hidden />
          )}
        </CollapsibleTrigger>
        <CollapsibleContent className="flex flex-col gap-6 p-6">
          {children}
        </CollapsibleContent>
      </Card>
    </Collapsible>
  );
}

// Essays without a college come first, then one block per target college
function SectionBody({
  essays,
  applicationTypes,
  renderCard,
}: {
  essays: Essays[];
  applicationTypes?: Record<number, string>;
  renderCard: (essay: Essays) => ReactNode;
}) {
  const general = essays.filter((e) => e.college_id == null);
  const byCollege = new Map<number, Essays[]>();
  for (const essay of essays) {
    if (essay.college_id == null) continue;
    byCollege.set(essay.college_id, [
      ...(byCollege.get(essay.college_id) ?? []),
      essay,
    ]);
  }

  return (
    <>
      {general.length > 0 && (
        <ul className="flex flex-col gap-3">{general.map(renderCard)}</ul>
      )}
      {[...byCollege].map(([collegeId, collegeEssays]) => (
        <div key={collegeId} className="flex flex-col gap-3">
          <CollegeHeading
            collegeId={collegeId}
            applicationType={applicationTypes?.[collegeId]}
          />
          <ul className="flex flex-col gap-3">
            {collegeEssays.map(renderCard)}
          </ul>
        </div>
      ))}
    </>
  );
}

export function EssaySections({
  essays,
  groups,
  applicationTypes,
  renderBadges,
  renderActions,
}: {
  essays: Essays[];
  groups: EssayGroups[];
  applicationTypes?: Record<number, string>;
  renderBadges: (essay: Essays) => ReactNode;
  renderActions: (essay: Essays) => ReactNode;
}) {
  const groupIds = new Set(groups.map((g) => g.id));
  const sections = groups
    .map((group) => ({
      id: group.id,
      name: group.name,
      essays: essays.filter((e) => e.essay_group_id === group.id),
    }))
    .filter((section) => section.essays.length > 0);

  const ungrouped = essays.filter(
    (e) => !e.essay_group_id || !groupIds.has(e.essay_group_id),
  );
  if (ungrouped.length > 0) {
    sections.push({ id: "ungrouped", name: "Other essays", essays: ungrouped });
  }

  const renderCard = (essay: Essays) => (
    <EssayCard
      key={essay.id}
      essay={essay}
      badges={renderBadges(essay)}
      actions={renderActions(essay)}
    />
  );

  return (
    <div className="flex flex-col gap-6">
      {sections.map((section) => (
        <Section
          key={section.id}
          title={section.name}
          count={section.essays.length}
        >
          <SectionBody
            essays={section.essays}
            applicationTypes={applicationTypes}
            renderCard={renderCard}
          />
        </Section>
      ))}
    </div>
  );
}
