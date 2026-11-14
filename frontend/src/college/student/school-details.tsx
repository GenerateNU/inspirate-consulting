import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { errorMessage } from "../../api/client";
import {
  useGetEssaysFromStudent,
  useGetGlobalCollege,
  useListPersonalCollegeApplications,
  useUpdatePersonalCollegeApplication,
} from "../../api/endpoints";
import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  UpdatePersonalCollegeApplicationRequestBodyInterest as Interest,
  type Essays,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../../api/models";
import arrowLeftIcon from "@/assets/college/arrow-left.svg";
import chevronDownIcon from "@/assets/college/chevron-down.svg";
import essayFileIcon from "@/assets/college/essay-file.svg";
import locationPinIcon from "@/assets/college/location-pin.svg";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuthContext } from "@/contexts/use-auth-context";
import ChoiceChips from "../utils/components/choice-chips";
import Tag, { type TagTone } from "../utils/components/tag";
import { cycleDateLabel, cycleLabel, findCycle } from "../utils/shared";
import { categoryOptions, interestOptions } from "./application-labels";

// the fields the student edits here; each change saves right away
interface Plan {
  application_type: ApplicationType;
  category: Category;
  interest: Interest | null;
}

const essayStatusTag: Record<string, { label: string; tone: TagTone }> = {
  Draft: { label: "In progress", tone: "yellow" },
  Review: { label: "In review", tone: "blue" },
  Submitted: { label: "Completed", tone: "green" },
};

function PageFrame({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-full w-full flex-col gap-6 bg-off-white p-12 max-lg:px-6">
      {children}
    </div>
  );
}

function CycleDropdown({
  college,
  value,
  disabled,
  onChange,
}: {
  college: GlobalCollege | undefined;
  value: string;
  disabled: boolean;
  onChange: (cycle: ApplicationType) => void;
}) {
  const offered = college?.cycles ?? [];
  const current = findCycle(college, value);

  return (
    <div className="flex w-[408px] max-w-full flex-col gap-2">
      <span
        id="application-cycle-label"
        className="font-host text-sm leading-5 text-charcoal"
      >
        Application Cycle
      </span>
      <DropdownMenu>
        <DropdownMenuTrigger asChild disabled={disabled}>
          <button
            type="button"
            data-slot="cycle-trigger"
            aria-labelledby="application-cycle-label"
            className="flex w-full items-center justify-center gap-2.5 rounded-lg border border-greyish-brown bg-white px-4 py-3 text-left font-host text-sm leading-5 drop-shadow-[0px_1px_1.5px_rgba(22,25,29,0.1)] hover:bg-off-white disabled:cursor-not-allowed disabled:opacity-60"
          >
            <span className="whitespace-nowrap text-charcoal">
              {cycleLabel(value)}
            </span>
            <span className="min-w-0 flex-1 truncate text-greyish-brown">
              {current ? cycleDateLabel(current) : "No longer offered"}
            </span>
            <img src={chevronDownIcon} alt="" width={24} height={24} />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          align="start"
          className="w-(--radix-dropdown-menu-trigger-width) font-host"
        >
          <DropdownMenuRadioGroup
            value={value}
            onValueChange={(cycle) => onChange(cycle as ApplicationType)}
          >
            {offered.map((cycle) => (
              <DropdownMenuRadioItem key={cycle.cycle} value={cycle.cycle}>
                <span className="text-charcoal">{cycleLabel(cycle.cycle)}</span>
                <span className="text-warm-grey">{cycleDateLabel(cycle)}</span>
              </DropdownMenuRadioItem>
            ))}
          </DropdownMenuRadioGroup>
        </DropdownMenuContent>
      </DropdownMenu>
      {!current && (
        <p role="status" className="font-host text-xs text-danger">
          {college?.school_name ?? "This school"} no longer offers{" "}
          {cycleLabel(value)}. Pick another cycle.
        </p>
      )}
    </div>
  );
}

function EssayCard({ essay }: { essay: Essays }) {
  const status = essayStatusTag[essay.status];

  return (
    <li className="relative flex min-h-[168px] w-full items-center gap-4 overflow-clip rounded-lg border border-cream bg-white py-6 pr-6 pl-8">
      <div
        aria-hidden="true"
        className="absolute top-[-1px] bottom-[-1px] left-[-1px] w-1.5 bg-blue"
      />
      <div className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-lightest-blue">
        <span className="flex size-[22px] items-center justify-center">
          <img
            src={essayFileIcon}
            alt=""
            width={16.6667}
            height={20.3333}
            className="max-w-none"
          />
        </span>
      </div>
      <div className="flex min-w-0 flex-1 flex-col gap-2">
        <div className="flex min-w-0 items-center gap-5">
          <p className="max-w-[560px] truncate font-jakarta text-xl leading-[normal] font-medium text-charcoal">
            {essay.type}
          </p>
          {status && <Tag tone={status.tone}>{status.label}</Tag>}
        </div>
        {essay.link_to_content && (
          <a
            href={essay.link_to_content}
            target="_blank"
            rel="noreferrer"
            className="w-fit font-host text-xs leading-4 text-dark-greyish-brown underline"
          >
            ↗ Open in Google Docs
          </a>
        )}
      </div>
    </li>
  );
}

function SchoolDetailsBody({ app }: { app: PersonalCollegeApplication }) {
  const { studentID } = useAuthContext();
  const collegeQuery = useGetGlobalCollege(app.global_college_id);
  const essaysQuery = useGetEssaysFromStudent(studentID ?? "", {
    query: { enabled: !!studentID },
  });
  const { refetch } = useListPersonalCollegeApplications();
  const update = useUpdatePersonalCollegeApplication();

  // shown right away while the save is in flight; cleared once the list refetches
  const [draft, setDraft] = useState<Partial<Plan>>({});

  const plan: Plan = {
    application_type: app.application_type as ApplicationType,
    category: app.category as Category,
    interest: app.interest,
    ...draft,
  };

  const save = (change: Partial<Plan>) => {
    const next = { ...plan, ...change };
    setDraft(next);
    update.mutate(
      {
        id: app.id,
        data: {
          global_college_id: app.global_college_id,
          application_type: next.application_type,
          category: next.category,
          // leaving interest out clears it
          interest: next.interest ?? undefined,
          application_status: app.application_status,
        },
      },
      {
        onSuccess: async () => {
          await refetch();
          setDraft({});
        },
        onError: () => setDraft({}),
      },
    );
  };

  const college = collegeQuery.data;
  const essays = (essaysQuery.data?.essays ?? []).filter(
    (essay) =>
      essay.college_id === app.global_college_id && essay.status !== "Archived",
  );

  return (
    <>
      <section className="flex flex-col gap-6 rounded-lg border border-cream bg-white p-6 shadow-card">
        <div className="flex flex-col gap-1 leading-[normal]">
          <h2 className="font-jakarta text-xl font-medium text-charcoal">
            Your application plan
          </h2>
          <p className="font-host text-xs text-dark-greyish-brown">
            Choose how you&apos;re applying and how this school fits on your
            list.
          </p>
        </div>

        {collegeQuery.isLoading ? (
          <p className="font-host text-sm">Loading...</p>
        ) : (
          <CycleDropdown
            college={college}
            value={plan.application_type}
            disabled={update.isPending}
            onChange={(cycle) => save({ application_type: cycle })}
          />
        )}

        <div className="flex flex-wrap items-start gap-x-[65px] gap-y-6">
          <div className="flex flex-col gap-2">
            <span className="font-host text-sm leading-5 text-charcoal">
              Classification
            </span>
            <ChoiceChips
              label="Classification"
              options={categoryOptions}
              value={plan.category}
              disabled={update.isPending}
              onChange={(category) => category && save({ category })}
              chipClassName="min-w-[60px]"
            />
          </div>
          <div className="flex flex-col gap-2">
            <span className="font-host text-sm leading-5 text-charcoal">
              Interest
            </span>
            <ChoiceChips
              label="Interest"
              options={interestOptions}
              value={plan.interest}
              allowDeselect
              disabled={update.isPending}
              onChange={(interest) => save({ interest })}
            />
          </div>
        </div>

        {update.error && (
          <p role="alert" className="font-host text-sm text-danger">
            Couldn&apos;t save: {errorMessage(update.error)}
          </p>
        )}
      </section>

      <section className="flex flex-col gap-4">
        <h2 className="pt-2.5 font-jakarta text-xl leading-[normal] font-medium text-black">
          Essay Prompts
        </h2>
        {essaysQuery.isLoading ? (
          <p className="font-host text-sm">Loading...</p>
        ) : essaysQuery.error ? (
          <p className="font-host text-sm">
            Error: {errorMessage(essaysQuery.error)}
          </p>
        ) : essays.length === 0 ? (
          <p className="font-host text-sm text-dark-greyish-brown">
            No essays for this school yet.
          </p>
        ) : (
          <ul className="flex flex-col gap-4">
            {essays.map((essay) => (
              <EssayCard key={essay.id} essay={essay} />
            ))}
          </ul>
        )}
      </section>
    </>
  );
}

// one school on the student's list: how they're applying, and its essays
export default function SchoolDetails() {
  const id = Number(useParams().id);
  const navigate = useNavigate();
  const isStudent = useAuthContext().role === "student";
  const applicationsQuery = useListPersonalCollegeApplications(undefined, {
    query: { enabled: isStudent },
  });
  const app = applicationsQuery.data?.find(
    (application) => application.id === id,
  );
  const collegeQuery = useGetGlobalCollege(app?.global_college_id ?? 0, {
    query: { enabled: !!app },
  });
  const college = collegeQuery.data;

  let body;
  if (!isStudent) {
    body = (
      <p className="font-host text-sm">
        School Details is for students. Log in as a student to see it.
      </p>
    );
  } else if (applicationsQuery.isLoading) {
    body = <p className="font-host text-sm">Loading...</p>;
  } else if (applicationsQuery.error) {
    body = (
      <p className="font-host text-sm">
        Error: {errorMessage(applicationsQuery.error)}
      </p>
    );
  } else if (!app) {
    body = (
      <p className="font-host text-sm">This school isn&apos;t on your list.</p>
    );
  } else {
    body = <SchoolDetailsBody key={app.id} app={app} />;
  }

  return (
    <PageFrame>
      <div className="flex items-start gap-[35px]">
        <button
          type="button"
          data-slot="icon-button"
          aria-label="Back to My List"
          onClick={() => navigate("/college-applications")}
          className="shrink-0 rounded-md hover:bg-black/5"
        >
          <img src={arrowLeftIcon} alt="" width={32} height={32} />
        </button>
        <div className="flex flex-col gap-[3px]">
          <h1 className="font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
            {college?.school_name ?? (app ? "Loading..." : "School Details")}
          </h1>
          {college && (
            <p className="flex h-4 items-center gap-px font-host text-[13px] leading-4 tracking-[0.5px] text-warm-grey">
              <img src={locationPinIcon} alt="" width={16} height={16} />
              {college.school_location}
            </p>
          )}
        </div>
      </div>
      <div className="flex flex-col gap-6 pl-[66px] max-lg:pl-0">{body}</div>
    </PageFrame>
  );
}
