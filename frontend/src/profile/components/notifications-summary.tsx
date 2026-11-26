import type { NotificationPreferences } from "@/api/models";
import { Panel } from "@/profile/components/panel";
import { ReadOnlyField } from "@/profile/components/read-only-field";

interface NotificationsSummaryProps {
  notifications: NotificationPreferences;
}

const onOff = (on: boolean) => (on ? "On" : "Off");

const daysLabel = (days: number) =>
  days === 0 ? "on the due date" : `${days} day${days === 1 ? "" : "s"} before`;

export function NotificationsSummary({
  notifications,
}: NotificationsSummaryProps) {
  const emailOn = notifications.email_enabled;
  const sub = (on: boolean) => (emailOn ? onOff(on) : "Off (email is off)");
  const dueDate =
    emailOn && notifications.due_date_notifications_enabled
      ? `On, ${daysLabel(notifications.days_before_due)}`
      : sub(notifications.due_date_notifications_enabled);

  return (
    <Panel title="Notifications" readOnly>
      <ReadOnlyField label="Email notifications" value={onOff(emailOn)} />
      <ReadOnlyField
        label="Weekly summary"
        value={sub(notifications.weekly_summary_enabled)}
      />
      <ReadOnlyField label="Due date reminders" value={dueDate} />
      <ReadOnlyField
        label="Past due alerts"
        value={sub(notifications.notify_past_due)}
      />
    </Panel>
  );
}
