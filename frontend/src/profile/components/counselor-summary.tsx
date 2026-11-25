import type { Counselor } from "@/api/models";
import { Panel } from "@/profile/components/panel";
import { ReadOnlyField } from "@/profile/components/read-only-field";

interface CounselorSummaryProps {
  counselor: Counselor;
}

export function CounselorSummary({ counselor }: CounselorSummaryProps) {
  const link = counselor.custom_link && (
    <a
      href={counselor.custom_link}
      target="_blank"
      rel="noreferrer"
      className="text-dark-green underline underline-offset-2"
    >
      {counselor.custom_link}
    </a>
  );
  const expires =
    counselor.expires_at &&
    new Date(counselor.expires_at).toLocaleDateString(undefined, {
      dateStyle: "long",
    });

  return (
    <Panel title="Counselor Details" readOnly>
      <ReadOnlyField label="Custom Link" value={link} />
      <ReadOnlyField label="Account Expiration" value={expires || "None"} />
    </Panel>
  );
}
