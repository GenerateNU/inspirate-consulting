import { Skeleton } from "@/components/ui/skeleton";

// three placeholder cards while a Modules list loads
export default function ModuleListSkeleton() {
  return (
    <div className="flex w-full flex-col gap-4">
      {[0, 1, 2].map((i) => (
        <Skeleton key={i} className="h-[115px] w-full rounded-lg" />
      ))}
    </div>
  );
}
