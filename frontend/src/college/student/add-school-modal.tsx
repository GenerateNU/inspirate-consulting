import { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Check } from "lucide-react";

import { errorMessage } from "../../api/client";
import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  UpdatePersonalCollegeApplicationRequestBodyInterest as Interest,
  type GlobalCollege,
} from "../../api/models";
import chatCircleIcon from "@/assets/college/chat-circle.svg";
import closeIcon from "@/assets/college/close-x.svg";
import plusIcon from "@/assets/college/plus.svg";
import searchIcon from "@/assets/college/search.svg";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import ChoiceChips from "../utils/components/choice-chips";
import Tag from "../utils/components/tag";
import { cycleLabel } from "../utils/shared";
import { useCollegeSearch } from "../utils/use-college-search";
import { categoryOptions, interestOptions } from "./application-labels";
import type { NewSchool } from "./use-my-college-list";

// what the student has picked on one search result
interface Selection {
  cycle: ApplicationType | null;
  category: Category | null;
  interest: Interest | null;
  // marked with + to be added on Save Changes
  picked: boolean;
}

const emptySelection: Selection = {
  cycle: null,
  category: null,
  interest: null,
  picked: false,
};

// a school can be added once it has a cycle and a classification; interest is optional
const isComplete = (selection: Selection) =>
  selection.cycle != null && selection.category != null;

interface AddSchoolModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  addedCollegeIds: Set<number>;
  adding: boolean;
  error: unknown;
  // resolves with the schools that failed to save
  onSave: (schools: NewSchool[]) => Promise<NewSchool[]>;
}

// search for schools, pick the cycle -> classification -> interest for each,
// mark them with +, then Save Changes adds them all to My List
export default function AddSchoolModal({
  open,
  onOpenChange,
  addedCollegeIds,
  adding,
  error,
  onSave,
}: AddSchoolModalProps) {
  const search = useCollegeSearch();
  const navigate = useNavigate();
  const [selections, setSelections] = useState<Record<number, Selection>>({});
  const searchRef = useRef<HTMLInputElement>(null);

  const selectionFor = (id: number) => selections[id] ?? emptySelection;
  const select = (id: number, change: Partial<Selection>) =>
    setSelections((prev) => {
      const next = { ...(prev[id] ?? emptySelection), ...change };
      // un-picking an earlier step clears the steps after it, which also unmarks the school
      if (next.cycle == null) next.category = null;
      if (next.category == null) next.interest = null;
      if (!isComplete(next)) next.picked = false;
      return { ...prev, [id]: next };
    });

  const pickedSchools: NewSchool[] = Object.entries(selections).flatMap(
    ([id, selection]) =>
      selection.picked && isComplete(selection)
        ? [
            {
              globalCollegeId: Number(id),
              applicationType: selection.cycle!,
              category: selection.category!,
              interest: selection.interest,
            },
          ]
        : [],
  );

  const close = (nextOpen: boolean) => {
    if (!nextOpen) setSelections({});
    onOpenChange(nextOpen);
  };

  const save = async () => {
    const failed = await onSave(pickedSchools);
    if (failed.length === 0) {
      close(false);
      return;
    }
    // keep only the schools that didn't save, still marked, so the student can retry
    const failedIds = new Set(failed.map((school) => school.globalCollegeId));
    setSelections((prev) =>
      Object.fromEntries(
        Object.entries(prev).filter(([id]) => failedIds.has(Number(id))),
      ),
    );
  };

  const renderResult = (college: GlobalCollege) => {
    const selection = selectionFor(college.id);
    const { cycle, category, interest, picked } = selection;
    const offered = college.cycles ?? [];
    const added = addedCollegeIds.has(college.id);
    const complete = isComplete(selection);

    return (
      <li key={college.id} className="flex flex-col gap-2.5">
        <div className="flex items-center justify-between gap-4 px-5 py-2.5">
          <div className="flex min-w-0 flex-1 flex-col gap-3">
            <div className="flex flex-col gap-1">
              <p className="font-jakarta text-xl leading-[normal] font-medium text-charcoal">
                {college.school_name}
              </p>
              <p className="font-host text-xs leading-[normal] text-charcoal/50">
                {college.school_location}
              </p>
            </div>

            {added ? null : offered.length === 0 ? (
              <p className="font-host text-xs text-warm-grey-dark">
                No application cycles yet. Ask your counselor to add them.
              </p>
            ) : (
              <>
                <ChoiceChips
                  label={`Application cycle for ${college.school_name}`}
                  options={offered.map(({ cycle: value }) => ({
                    value: value as ApplicationType,
                    label: cycleLabel(value),
                  }))}
                  value={cycle}
                  allowDeselect
                  disabled={adding}
                  onChange={(value) => select(college.id, { cycle: value })}
                />
                {cycle && (
                  <ChoiceChips
                    label={`Classification for ${college.school_name}`}
                    options={categoryOptions}
                    value={category}
                    allowDeselect
                    disabled={adding}
                    onChange={(value) =>
                      select(college.id, { category: value })
                    }
                    chipClassName="min-w-[60px]"
                  />
                )}
                {category && (
                  <ChoiceChips
                    label={`Interest in ${college.school_name}`}
                    options={interestOptions}
                    value={interest}
                    allowDeselect
                    disabled={adding}
                    onChange={(value) =>
                      select(college.id, { interest: value })
                    }
                  />
                )}
              </>
            )}
          </div>

          {added ? (
            <Tag tone="green">On your list</Tag>
          ) : (
            <button
              type="button"
              data-slot="icon-button"
              aria-pressed={picked}
              aria-label={
                picked
                  ? `Don't add ${college.school_name}`
                  : `Add ${college.school_name} to your list`
              }
              title={
                complete
                  ? picked
                    ? "Will be added when you save"
                    : "Add to your list"
                  : "Pick a cycle and classification first"
              }
              disabled={!complete || adding}
              onClick={() => select(college.id, { picked: !picked })}
              className={cn(
                "flex size-8 shrink-0 items-center justify-center rounded-full border disabled:cursor-not-allowed disabled:opacity-40",
                picked
                  ? "border-dark-green bg-dark-green text-white hover:bg-dark-green/90"
                  : "border-greyish-brown bg-off-white hover:bg-lightest-green disabled:hover:bg-off-white",
              )}
            >
              {picked ? (
                <Check aria-hidden="true" className="size-4" />
              ) : (
                <img src={plusIcon} alt="" width={16} height={16} />
              )}
            </button>
          )}
        </div>
        <div aria-hidden="true" className="h-[2px] w-full bg-off-white" />
      </li>
    );
  };

  let results;
  if (search.isLoading || (search.searching && search.searchLoading)) {
    results = <p className="py-6 text-center font-host text-sm">Loading...</p>;
  } else if (search.error || (search.searching && search.searchError)) {
    results = (
      <p className="py-6 text-center font-host text-sm">
        Error: {errorMessage(search.error ?? search.searchError)}
      </p>
    );
  } else if (search.colleges.length === 0) {
    results = (
      <p className="py-6 text-center font-host text-sm text-warm-grey-dark">
        {search.searching ? "No matching schools." : "No schools yet."}
      </p>
    );
  } else {
    results = (
      <ul className="flex flex-col gap-2.5">
        {search.colleges.map(renderResult)}
      </ul>
    );
  }

  return (
    <Dialog open={open} onOpenChange={close}>
      <DialogContent
        showCloseButton={false}
        // start in the search box instead of on the close button
        onOpenAutoFocus={(e) => {
          e.preventDefault();
          searchRef.current?.focus();
        }}
        className="flex w-[712px] flex-col items-center gap-2 rounded-lg border-0 bg-off-white px-6 pt-5 pb-[25px] drop-shadow-[0px_1px_1.5px_rgba(22,25,29,0.1)] sm:max-w-[712px]"
      >
        <div className="flex w-full justify-end">
          <DialogClose asChild>
            <button
              type="button"
              data-slot="icon-button"
              aria-label="Close"
              className="rounded-sm hover:opacity-70"
            >
              <img src={closeIcon} alt="" width={20} height={20} />
            </button>
          </DialogClose>
        </div>
        <DialogTitle className="font-jakarta text-xl leading-[normal] font-medium text-black">
          Add School
        </DialogTitle>
        <DialogDescription className="sr-only">
          Search for schools, pick how you&apos;re applying to each, mark them
          with the plus button, then save to add them to your list.
        </DialogDescription>

        <div className="flex h-[500px] max-h-[55vh] w-full flex-col gap-2.5 px-[30px] py-2.5">
          <label className="flex h-10 w-full shrink-0 items-center gap-2 rounded-full border border-greyish-brown bg-white px-3">
            <span className="flex size-4 shrink-0 items-center justify-center">
              <img src={searchIcon} alt="" width={12.6667} height={12.6667} />
            </span>
            <input
              type="text"
              ref={searchRef}
              data-slot="search-input"
              value={search.query}
              onChange={(e) => search.setQuery(e.target.value)}
              placeholder="Search by name or location"
              aria-label="Search schools"
              className="min-w-0 flex-1 bg-transparent font-host text-sm text-charcoal outline-none placeholder:text-warm-grey"
            />
          </label>

          {error != null && (
            <p role="alert" className="font-host text-sm text-danger">
              Couldn&apos;t add every school: {errorMessage(error)}
            </p>
          )}

          <div className="min-h-0 flex-1 overflow-y-auto rounded-lg bg-white py-5 shadow-card">
            {results}
          </div>
        </div>

        <div className="flex flex-col items-center gap-[3px]">
          <p className="font-host text-xs leading-[normal] text-warm-grey-dark">
            Can&apos;t find your school?
          </p>
          {/* TODO: open the chat with this student's counselor directly once
              /user/me returns the counselor's user id */}
          <button
            type="button"
            data-slot="chat-counselor"
            onClick={() => navigate("/chats")}
            className="flex h-[29px] items-center justify-center gap-[3px] rounded-[10px] border border-white bg-light-black px-4 py-2 font-host text-[13px] leading-[normal] text-white hover:bg-light-black/85"
          >
            <img src={chatCircleIcon} alt="" width={15} height={15} />
            Chat with your counselor
          </button>
        </div>

        <div className="flex w-full justify-end px-[27px] py-5">
          <button
            type="button"
            data-slot="primary-action"
            disabled={pickedSchools.length === 0 || adding}
            onClick={save}
            className="rounded-lg bg-dark-green px-4 py-2 font-host text-xs leading-[normal] text-white hover:bg-dark-green/90 disabled:cursor-not-allowed disabled:opacity-40"
          >
            {adding ? "Saving..." : "Save Changes"}
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
