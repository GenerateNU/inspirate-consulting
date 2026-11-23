import type { FormEvent } from "react";
import { CollegePreferenceFields } from "@/profile/components/college-preference-fields";
import { DiscardChangesDialog } from "@/profile/components/discard-changes-dialog";
import { ErrorAlert } from "@/profile/components/error-alert";
import { NotificationFields } from "@/profile/components/notification-fields";
import { PageHeader } from "@/profile/components/page-header";
import { PageShell } from "@/profile/components/page-shell";
import { PersonalInfoFields } from "@/profile/components/personal-info-fields";
import { SaveButton } from "@/profile/components/save-button";
import { useLeaveGuard } from "@/profile/use-leave-guard";
import { useProfileForm, type ProfileSource } from "@/profile/use-profile-form";

interface ProfileFormProps extends ProfileSource {
  onCancel: () => void;
  onSaved: () => void;
}

export function ProfileForm({
  onCancel,
  onSaved,
  ...source
}: ProfileFormProps) {
  const form = useProfileForm(source);
  const leaveGuard = useLeaveGuard(form.isDirty);
  const { user, student } = source;

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    void form.save();
  };

  return (
    <PageShell>
      <form className="contents" onSubmit={handleSubmit}>
        <PageHeader
          title="Edit Profile"
          backLabel="Back to profile"
          onBack={() => leaveGuard.guard(onCancel)}
        >
          <ErrorAlert
            error={form.error}
            attempt={form.failedSaves}
            label="Couldn't save changes"
          />
          <SaveButton
            state={form.saveState}
            disabled={!form.isDirty || form.hasErrors}
            onSavedShown={onSaved}
          />
        </PageHeader>

        <div className="flex w-full flex-col gap-6 sm:pl-15">
          <PersonalInfoFields
            user={user}
            isStudent={!!student}
            values={form.values}
            errors={form.errors}
            setValue={form.setValue}
          />
          {student && (
            <CollegePreferenceFields
              values={form.values}
              setValue={form.setValue}
            />
          )}
          <NotificationFields
            notify={form.notify}
            setNotification={form.setNotification}
            daysError={form.errors.daysBeforeDue}
          />
        </div>
      </form>

      <DiscardChangesDialog {...leaveGuard.dialog} />
    </PageShell>
  );
}
