import type { ReactNode } from "react";

import { circleClassName } from "../shared";

// same circle as CircleButton but not clickable (e.g. Check for already added)
interface CircleBadgeProps {
  label: string;
  title?: string;
  children: ReactNode;
}

export default function CircleBadge({
  label,
  title,
  children,
}: CircleBadgeProps) {
  return (
    <span
      role="img"
      aria-label={label}
      title={title}
      className={circleClassName}
    >
      {children}
    </span>
  );
}
