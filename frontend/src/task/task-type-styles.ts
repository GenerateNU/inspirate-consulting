import {
  FileText,
  Palette,
  School,
  SquarePlay,
  type LucideIcon,
} from "lucide-react";

import type { TodoItem, TodoItemLinkType } from "../api/models";
import type { TagTone } from "../college/utils/components/tag";

export type TaskType = NonNullable<TodoItemLinkType>;

export interface TaskTypeStyle {
  label: string;
  icon: LucideIcon;
  // left edge of the card
  accentClassName: string;
  // the square behind the icon
  iconClassName: string;
  tagTone: TagTone;
  tagClassName?: string;
}

// colors per task type from the hi-fi, using the same icons as the sidebar
export const TASK_TYPE_STYLES: Record<TaskType, TaskTypeStyle> = {
  college: {
    label: "College",
    icon: School,
    accentClassName: "bg-mint",
    iconClassName: "bg-lightest-green text-dark-green",
    tagTone: "green",
  },
  essay: {
    label: "Essay",
    icon: FileText,
    accentClassName: "bg-blue",
    iconClassName: "bg-lightest-blue text-dark-blue",
    tagTone: "blue",
  },
  extracurricular: {
    label: "Extracurricular",
    icon: Palette,
    accentClassName: "bg-tan",
    iconClassName: "bg-lightest-yellow text-dark-yellow",
    tagTone: "yellow",
  },
  module: {
    label: "Module",
    icon: SquarePlay,
    accentClassName: "bg-greyish-brown",
    iconClassName: "bg-cream text-dark-greyish-brown",
    tagTone: "grey",
    tagClassName: "border-none bg-cream",
  },
};

// tasks without a link get a plain grey edge and no icon or tag
export const UNLINKED_ACCENT_CLASS_NAME = "bg-warm-grey";

export const taskTypeStyle = (task: TodoItem) =>
  task.link_type ? TASK_TYPE_STYLES[task.link_type] : null;
