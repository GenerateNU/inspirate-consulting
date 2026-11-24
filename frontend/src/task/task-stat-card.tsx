import type { ReactNode } from "react";

type TaskStatCardProps = {
  label: string;
  value: ReactNode;
  // small text after the value, e.g. "/10 Completed"
  valueSuffix?: string;
  // 0-100; shows a progress bar when set
  progress?: number;
  caption?: string;
};

// The hi-fi stat card at the top of the task dashboard
export default function TaskStatCard({
  label,
  value,
  valueSuffix,
  progress,
  caption,
}: TaskStatCardProps) {
  return (
    <div className="flex w-64 flex-col gap-2 rounded-md border border-cream bg-white px-6 py-5 shadow-card">
      <p className="font-host text-xs leading-[normal] text-light-black">
        {label}
      </p>
      <p className="flex items-baseline gap-1 font-jakarta text-2xl leading-[normal] font-semibold text-charcoal">
        {value}
        {valueSuffix && (
          <span className="font-host text-xs font-normal text-warm-grey-dark">
            {valueSuffix}
          </span>
        )}
      </p>
      {progress !== undefined && (
        <div
          role="progressbar"
          aria-valuenow={progress}
          aria-valuemin={0}
          aria-valuemax={100}
          aria-label={label}
          className="h-1 w-40 overflow-hidden rounded-full bg-cream"
        >
          <div
            className="h-full rounded-full bg-light-black transition-[width] duration-500"
            style={{ width: `${progress}%` }}
          />
        </div>
      )}
      {caption && (
        <p className="font-host text-[10px] leading-[normal] text-warm-grey-dark">
          {caption}
        </p>
      )}
    </div>
  );
}
