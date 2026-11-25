import { Skeleton } from "@/components/ui/skeleton";
import { PageHeader } from "@/profile/components/page-header";
import { PageShell } from "@/profile/components/page-shell";
import { Panel } from "@/profile/components/panel";
import { ReadOnlyFieldSkeleton } from "@/profile/components/read-only-field";

interface ProfileSkeletonProps {
  academic?: boolean;
  own?: boolean;
}

interface SkeletonPanel {
  title: string;
  fields: number;
  avatar?: boolean;
}

const skeletonPanels = (academic: boolean, own: boolean): SkeletonPanel[] => [
  { title: "Personal Information", fields: 2, avatar: true },
  ...(academic
    ? [
        { title: "Academic Overview", fields: 3 },
        { title: "College Preferences & Specializations", fields: 3 },
      ]
    : []),
  ...(own ? [{ title: "Notifications", fields: 4 }] : []),
];

export function ProfileSkeleton({
  academic = true,
  own = true,
}: ProfileSkeletonProps) {
  return (
    <PageShell>
      <PageHeader title={<Skeleton className="h-9 w-56 bg-cream" />}>
        {own && <ActionSkeletons />}
      </PageHeader>
      <div
        role="status"
        aria-label={own ? "Loading your profile" : "Loading profile"}
        className="flex w-full flex-col gap-6 sm:pl-15"
      >
        {skeletonPanels(academic, own).map((panel) => (
          <PanelSkeleton key={panel.title} {...panel} />
        ))}
      </div>
    </PageShell>
  );
}

function ActionSkeletons() {
  return (
    <>
      <Skeleton className="h-8 w-[92px] rounded-lg bg-cream" />
      <Skeleton className="h-8 w-[98px] rounded-lg bg-cream" />
    </>
  );
}

function PanelSkeleton({ title, fields, avatar }: SkeletonPanel) {
  const aside = avatar && (
    <Skeleton className="size-24 shrink-0 rounded-full bg-cream" />
  );
  return (
    <Panel title={title} aside={aside} readOnly>
      {Array.from({ length: fields }, (_, i) => (
        <ReadOnlyFieldSkeleton key={i} />
      ))}
    </Panel>
  );
}
