import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { ExternalLink, FilePlus, FolderPlus, Send } from "lucide-react";

import { errorMessage } from "../api/client";
import {
  getGetEssayGroupsFromStudentQueryKey,
  getGetEssayReviewStatusQueryKey,
  getGetEssaysFromStudentQueryKey,
  getGetStudentQueryKey,
  useCreateEssay,
  useCreateEssayGroup,
  useGetEssayGroupsFromStudent,
  useGetEssaysFromStudent,
  useGetStudent,
  useListPersonalCollegeApplications,
  useRequestEssayReview,
} from "../api/endpoints";
import type { Essays } from "../api/models";
import { useAuthContext } from "../contexts/use-auth-context";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  NativeSelect,
  NativeSelectOption,
} from "@/components/ui/native-select";
import {
  EssaySections,
  RequestedAt,
  StatGrid,
  isReviewCompleted,
  isReviewOpen,
  reviewStats,
  useCollegeNames,
  useReviewStatuses,
} from "./essay-sections";

export default function Essay() {
  const { studentID } = useAuthContext();
  const studentId = studentID ?? "";
  const queryClient = useQueryClient();

  const [targetFilter, setTargetFilter] = useState("all");
  const [showEssayForm, setShowEssayForm] = useState(false);
  const [showGroupForm, setShowGroupForm] = useState(false);

  const essaysQuery = useGetEssaysFromStudent(studentId, {
    query: { enabled: !!studentId },
  });
  const groupsQuery = useGetEssayGroupsFromStudent(studentId, {
    query: { enabled: !!studentId },
  });
  const studentQuery = useGetStudent(studentId, {
    query: { enabled: !!studentId },
  });
  const applicationsQuery = useListPersonalCollegeApplications({
    query: { enabled: !!studentId },
  });
  const requestReview = useRequestEssayReview();

  const essays = essaysQuery.data?.essays ?? [];
  const groups = groupsQuery.data?.essay_groups ?? [];
  const applications = applicationsQuery.data ?? [];
  const { reviews, error: reviewsError } = useReviewStatuses(essays);
  const applicationCollegeIds = applications.map((a) => a.global_college_id);
  const targetIds = [
    ...new Set(
      essays.flatMap((e) => (e.college_id == null ? [] : [e.college_id])),
    ),
  ];
  const collegeNames = useCollegeNames([
    ...new Set([...targetIds, ...applicationCollegeIds]),
  ]);

  if (!studentId) {
    return <p className="p-8">Log in as a student to see your essays.</p>;
  }

  const error =
    essaysQuery.error ||
    groupsQuery.error ||
    reviewsError ||
    requestReview.error;

  const applicationTypes = Object.fromEntries(
    applications.map((a) => [a.global_college_id, a.application_type]),
  );
  const visibleEssays = essays.filter((e) =>
    targetFilter === "all"
      ? true
      : targetFilter === "general"
        ? e.college_id == null
        : e.college_id === Number(targetFilter),
  );

  const handleRequestReview = (essay: Essays) => {
    requestReview.mutate(
      { data: { essay_id: essay.id, amount: 1 } },
      {
        onSuccess: () => {
          void queryClient.invalidateQueries({
            queryKey: getGetEssayReviewStatusQueryKey(essay.id),
          });
          void queryClient.invalidateQueries({
            queryKey: getGetStudentQueryKey(studentId),
          });
        },
      },
    );
  };

  const balance = studentQuery.data?.review_balance;
  const stats = reviewStats(visibleEssays, reviews);

  return (
    <div className="flex-1 bg-muted/40 px-4 py-8">
      <div className="mx-auto flex max-w-5xl flex-col gap-6">
        <Card className="gap-6 p-6 sm:p-8">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold tracking-tight uppercase">
                My Essay Workspace
              </h1>
              {balance != null && (
                <p className="mt-1 text-sm text-muted-foreground">
                  Review credits remaining: {balance}
                </p>
              )}
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <NativeSelect
                aria-label="Filter by target"
                value={targetFilter}
                onChange={(e) => setTargetFilter(e.target.value)}
              >
                <NativeSelectOption value="all">All targets</NativeSelectOption>
                <NativeSelectOption value="general">
                  No specific college
                </NativeSelectOption>
                {targetIds.map((id) => (
                  <NativeSelectOption key={id} value={id}>
                    {collegeNames[id]}
                  </NativeSelectOption>
                ))}
              </NativeSelect>
              <Button
                variant="secondary"
                onClick={() => setShowGroupForm((s) => !s)}
              >
                <FolderPlus className="size-4" aria-hidden />
                New section
              </Button>
              <Button onClick={() => setShowEssayForm((s) => !s)}>
                <FilePlus className="size-4" aria-hidden />
                Add new essay
              </Button>
            </div>
          </div>

          {showGroupForm && (
            <NewGroupForm
              studentId={studentId}
              onDone={() => setShowGroupForm(false)}
            />
          )}
          {showEssayForm && (
            <NewEssayForm
              studentId={studentId}
              groups={groups}
              colleges={applicationCollegeIds.map((id) => ({
                id,
                name: collegeNames[id],
              }))}
              onDone={() => setShowEssayForm(false)}
            />
          )}

          <StatGrid
            stats={[
              { label: "Total essays", value: stats.total },
              { label: "Counselor queue", value: stats.open },
              { label: "Reviewed", value: stats.completed },
              { label: "Not reviewed", value: stats.notReviewed },
            ]}
          />
        </Card>

        {error && (
          <p
            role="alert"
            className="rounded-md bg-destructive/10 p-4 text-destructive"
          >
            {errorMessage(error)}
          </p>
        )}

        {essaysQuery.isLoading || groupsQuery.isLoading ? (
          <p>Loading...</p>
        ) : visibleEssays.length === 0 ? (
          <Card className="p-8 text-center text-muted-foreground">
            {essays.length === 0
              ? "No essays yet. Add one to get started."
              : "No essays match this filter."}
          </Card>
        ) : (
          <EssaySections
            essays={visibleEssays}
            groups={groups}
            applicationTypes={applicationTypes}
            renderBadges={(essay) => {
              const review = reviews[essay.id];
              return (
                <>
                  {!isReviewOpen(review) && !isReviewCompleted(review) && (
                    <Badge variant="secondary">Not reviewed</Badge>
                  )}
                  {isReviewOpen(review) && (
                    <>
                      <Badge variant="secondary">Under review</Badge>
                      <RequestedAt review={review!} />
                    </>
                  )}
                  {isReviewCompleted(review) && (
                    <Badge variant="outline">Reviewed</Badge>
                  )}
                </>
              );
            }}
            renderActions={(essay) => {
              const canRequest = !isReviewOpen(reviews[essay.id]);
              return (
                <>
                  {canRequest && (
                    <Button
                      variant="secondary"
                      disabled={requestReview.isPending || balance === 0}
                      title={balance === 0 ? "No review credits left" : ""}
                      onClick={() => handleRequestReview(essay)}
                    >
                      <Send className="size-4" aria-hidden />
                      Request review
                    </Button>
                  )}
                  <Button asChild>
                    <a
                      href={essay.link_to_content}
                      target="_blank"
                      rel="noreferrer"
                    >
                      Open essay
                      <ExternalLink className="size-4" aria-hidden />
                    </a>
                  </Button>
                </>
              );
            }}
          />
        )}
      </div>
    </div>
  );
}

function NewEssayForm({
  studentId,
  groups,
  colleges,
  onDone,
}: {
  studentId: string;
  groups: { id: string; name: string }[];
  colleges: { id: number; name: string }[];
  onDone: () => void;
}) {
  const queryClient = useQueryClient();
  const create = useCreateEssay();
  const [type, setType] = useState("");
  const [link, setLink] = useState("");
  const [groupId, setGroupId] = useState("");
  const [collegeId, setCollegeId] = useState("");

  const submit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    create.mutate(
      {
        data: {
          type: type.trim(),
          link_to_content: link.trim(),
          college_id: collegeId ? Number(collegeId) : null,
          essay_group_id: groupId || undefined,
        },
      },
      {
        onSuccess: () => {
          void queryClient.invalidateQueries({
            queryKey: getGetEssaysFromStudentQueryKey(studentId),
          });
          onDone();
        },
      },
    );
  };

  return (
    <form
      onSubmit={submit}
      className="grid gap-4 rounded-md border p-4 sm:grid-cols-2"
    >
      <div className="flex flex-col gap-2">
        <Label htmlFor="essay-name">Essay name</Label>
        <Input
          id="essay-name"
          required
          placeholder="Personal Statement"
          value={type}
          onChange={(e) => setType(e.target.value)}
        />
      </div>
      <div className="flex flex-col gap-2">
        <Label htmlFor="essay-link">Google Doc link</Label>
        <Input
          id="essay-link"
          required
          type="url"
          placeholder="https://docs.google.com/..."
          value={link}
          onChange={(e) => setLink(e.target.value)}
        />
      </div>
      <div className="flex flex-col gap-2">
        <Label htmlFor="essay-section">Section</Label>
        <NativeSelect
          id="essay-section"
          className="w-full"
          value={groupId}
          onChange={(e) => setGroupId(e.target.value)}
        >
          <NativeSelectOption value="">None</NativeSelectOption>
          {groups.map((g) => (
            <NativeSelectOption key={g.id} value={g.id}>
              {g.name}
            </NativeSelectOption>
          ))}
        </NativeSelect>
      </div>
      <div className="flex flex-col gap-2">
        <Label htmlFor="essay-college">Target college</Label>
        <NativeSelect
          id="essay-college"
          className="w-full"
          value={collegeId}
          onChange={(e) => setCollegeId(e.target.value)}
        >
          <NativeSelectOption value="">None</NativeSelectOption>
          {colleges.map((c) => (
            <NativeSelectOption key={c.id} value={c.id}>
              {c.name}
            </NativeSelectOption>
          ))}
        </NativeSelect>
      </div>
      {create.error && (
        <p role="alert" className="text-sm text-destructive sm:col-span-2">
          {errorMessage(create.error)}
        </p>
      )}
      <div className="flex justify-end gap-2 sm:col-span-2">
        <Button variant="ghost" onClick={onDone}>
          Cancel
        </Button>
        <Button type="submit" disabled={create.isPending}>
          Add essay
        </Button>
      </div>
    </form>
  );
}

function NewGroupForm({
  studentId,
  onDone,
}: {
  studentId: string;
  onDone: () => void;
}) {
  const queryClient = useQueryClient();
  const create = useCreateEssayGroup();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  const submit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    create.mutate(
      {
        data: {
          name: name.trim(),
          description: description.trim() || null,
          student_id: studentId,
        },
      },
      {
        onSuccess: () => {
          void queryClient.invalidateQueries({
            queryKey: getGetEssayGroupsFromStudentQueryKey(studentId),
          });
          onDone();
        },
      },
    );
  };

  return (
    <form
      onSubmit={submit}
      className="grid gap-4 rounded-md border p-4 sm:grid-cols-2"
    >
      <div className="flex flex-col gap-2">
        <Label htmlFor="section-name">Section name</Label>
        <Input
          id="section-name"
          required
          placeholder="Common App Essays"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </div>
      <div className="flex flex-col gap-2">
        <Label htmlFor="section-description">Description (optional)</Label>
        <Input
          id="section-description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
      </div>
      {create.error && (
        <p role="alert" className="text-sm text-destructive sm:col-span-2">
          {errorMessage(create.error)}
        </p>
      )}
      <div className="flex justify-end gap-2 sm:col-span-2">
        <Button variant="ghost" onClick={onDone}>
          Cancel
        </Button>
        <Button type="submit" disabled={create.isPending}>
          Create section
        </Button>
      </div>
    </form>
  );
}
