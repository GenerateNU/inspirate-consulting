import type { User } from "@/api/models";
import UserAvatar from "@/components/user-avatar";
import { EmailValue } from "@/profile/components/email-field";
import { Panel } from "@/profile/components/panel";
import { ReadOnlyField } from "@/profile/components/read-only-field";

export interface ProfileUser extends Pick<
  User,
  "id" | "name" | "email" | "pfp_key"
> {
  pfp_url?: string | null;
}

interface PersonalInfoSummaryProps {
  user: ProfileUser;
  own: boolean;
}

export function PersonalInfoSummary({ user, own }: PersonalInfoSummaryProps) {
  const email =
    own || !user.email ? (
      user.email
    ) : (
      <EmailValue email={user.email} userId={user.id} />
    );

  return (
    <Panel
      title="Personal Information"
      readOnly
      aside={
        <UserAvatar
          name={user.name}
          userId={user.id}
          pfpKey={user.pfp_key}
          src={user.pfp_url}
          className="size-24 shrink-0"
          fallbackClassName="bg-mist text-3xl text-neutral-600"
        />
      }
    >
      <ReadOnlyField label="Name" value={user.name} />
      <ReadOnlyField label="Email" value={email} />
    </Panel>
  );
}
