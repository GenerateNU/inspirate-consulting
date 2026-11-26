import type { NotificationPreferences, Student } from "@/api/models";
import { InviteUserCard } from "@/profile/components/invite-user-card";
import { LogOutButton } from "@/profile/components/log-out-button";
import { NotificationsSummary } from "@/profile/components/notifications-summary";
import { PageHeader } from "@/profile/components/page-header";
import { PageShell } from "@/profile/components/page-shell";
import {
  PersonalInfoSummary,
  type ProfileUser,
} from "@/profile/components/personal-info-summary";
import { StudentSummary } from "@/profile/components/student-summary";

export interface OwnProfile {
  notifications: NotificationPreferences;
  onEdit: () => void;
  canInvite?: boolean;
}

interface ProfileViewProps {
  user: ProfileUser;
  student?: Student;
  own?: OwnProfile;
}

export function ProfileView({ user, student, own }: ProfileViewProps) {
  return (
    <PageShell>
      <PageHeader title={user.name}>
        {own && <OwnProfileActions onEdit={own.onEdit} />}
      </PageHeader>

      <div className="flex w-full flex-col gap-6 sm:pl-15">
        <PersonalInfoSummary user={user} own={!!own} />
        {student && <StudentSummary student={student} />}
        {own && <NotificationsSummary notifications={own.notifications} />}
        {own?.canInvite && <InviteUserCard />}
      </div>
    </PageShell>
  );
}

const editButtonClass =
  "flex h-8 items-center rounded-lg bg-dark-green px-4 font-host text-xs leading-[normal] whitespace-nowrap text-off-white outline-none transition-colors hover:bg-dark-green/90 focus-visible:ring-[3px] focus-visible:ring-dark-green/30";

interface OwnProfileActionsProps {
  onEdit: () => void;
}

function OwnProfileActions({ onEdit }: OwnProfileActionsProps) {
  return (
    <>
      <LogOutButton />
      <button
        type="button"
        data-slot="edit-profile"
        onClick={onEdit}
        className={editButtonClass}
      >
        Edit Profile
      </button>
    </>
  );
}
