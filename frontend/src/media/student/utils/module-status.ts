import type { AccessibleMedia } from "../../../api/models";

export type ModuleStatus = "notStarted" | "inProgress" | "completed";

// started_at is set when the student first opens the module,
// completed_at when they press Mark as Complete
export function getModuleStatus(module: AccessibleMedia): ModuleStatus {
  if (module.completed_at) return "completed";
  if (module.started_at) return "inProgress";
  return "notStarted";
}
