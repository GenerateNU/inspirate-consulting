import { Panel } from "@/profile/components/panel";
import { SwitchRow } from "@/profile/components/switch-row";
import { TextField } from "@/profile/components/text-field";
import {
  daysEditable,
  type NotificationValues,
  type SetNotificationValue,
} from "@/profile/profile-form-values";

interface NotificationFieldsProps {
  notify: NotificationValues;
  setNotification: SetNotificationValue;
  daysError?: string;
}

export function NotificationFields({
  notify,
  setNotification,
  daysError,
}: NotificationFieldsProps) {
  const emailOff = !notify.email;
  return (
    <Panel title="Notifications">
      <SwitchRow
        label="Email notifications"
        description="Turn off to stop all notification emails."
        checked={notify.email}
        onChange={setNotification("email")}
      />
      <SwitchRow
        label="Weekly summary"
        description="A weekly email of your outstanding assignments."
        checked={notify.weeklySummary}
        onChange={setNotification("weeklySummary")}
        disabled={emailOff}
      />
      <SwitchRow
        label="Due date reminders"
        description="An email before an assignment is due."
        checked={notify.dueDate}
        onChange={setNotification("dueDate")}
        disabled={emailOff}
      />
      <TextField
        label="Days before due date"
        value={notify.daysBeforeDue}
        onChange={setNotification("daysBeforeDue")}
        error={daysError}
        disabled={!daysEditable(notify)}
        type="number"
        inputMode="numeric"
        min={0}
        max={30}
        step={1}
      />
      <SwitchRow
        label="Past due alerts"
        description="An email when an assignment becomes past due."
        checked={notify.pastDue}
        onChange={setNotification("pastDue")}
        disabled={emailOff}
      />
    </Panel>
  );
}
