import type { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft } from "lucide-react";

interface PageHeaderProps {
  title: ReactNode;
  backLabel?: string;
  onBack?: () => void;
  children?: ReactNode;
}

export function PageHeader({
  title,
  backLabel = "Go back",
  onBack,
  children,
}: PageHeaderProps) {
  const navigate = useNavigate();
  return (
    <div className="flex w-full flex-wrap items-center justify-between gap-4">
      <div className="flex min-w-0 items-center gap-7">
        <button
          type="button"
          data-slot="back"
          aria-label={backLabel}
          onClick={onBack ?? (() => void navigate(-1))}
          className="flex size-8 shrink-0 items-center justify-center rounded-md text-charcoal outline-none hover:bg-black/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30"
        >
          <ArrowLeft className="size-8" strokeWidth={1.5} />
        </button>
        <h1 className="min-w-0 font-jakarta text-[30px] leading-[normal] font-bold break-words text-charcoal">
          {title}
        </h1>
      </div>
      <div className="flex items-center gap-2">{children}</div>
    </div>
  );
}
