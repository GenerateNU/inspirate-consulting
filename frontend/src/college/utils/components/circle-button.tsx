import type { ReactNode } from "react";

import { circleClassName } from "../shared";

// round icon button on the right of a card (e.g. Plus to add or X to delete)
interface CircleButtonProps {
  label: string;
  title?: string;
  disabled?: boolean;
  onClick: () => void;
  children: ReactNode;
}

export default function CircleButton({
  label,
  title,
  disabled,
  onClick,
  children,
}: CircleButtonProps) {
  return (
    <button
      type="button"
      aria-label={label}
      title={title}
      disabled={disabled}
      onClick={onClick}
      className={`${circleClassName} hover:bg-gray-100 disabled:cursor-not-allowed disabled:hover:bg-transparent`}
    >
      {children}
    </button>
  );
}
