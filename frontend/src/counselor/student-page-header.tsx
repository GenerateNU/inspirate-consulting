import { useNavigate } from "react-router-dom";

import arrowLeftIcon from "@/assets/college/arrow-left.svg";

interface StudentPageHeaderProps {
  // e.g. "Anna's College List"
  title: string;
  // where the back arrow goes. A fixed path, not the previous page, so going
  // back can't loop between a list and its items
  backTo: string;
}

// header for a counselor's "view all" pages about one student: a back arrow and the title
export default function StudentPageHeader({
  title,
  backTo,
}: StudentPageHeaderProps) {
  const navigate = useNavigate();

  return (
    <div className="flex w-full items-start gap-9">
      <button
        type="button"
        data-slot="icon-button"
        aria-label="Back"
        onClick={() => navigate(backTo)}
        className="shrink-0 rounded-md hover:bg-black/5"
      >
        <img src={arrowLeftIcon} alt="" width={32} height={32} />
      </button>
      <h1 className="min-w-0 flex-1 truncate font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
        {title}
      </h1>
    </div>
  );
}
