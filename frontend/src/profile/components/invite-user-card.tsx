import { Link } from "react-router-dom";
import { PersonStanding } from "lucide-react";

const INVITE_USER_PATH = "/create-account";

export function InviteUserCard() {
  return (
    <section
      aria-label="Invite a user"
      className="flex w-full flex-wrap items-center justify-between gap-4 rounded-md border border-cream bg-white p-6 shadow-card"
    >
      <div className="flex items-center gap-4">
        <PersonStanding
          aria-hidden="true"
          className="size-11 shrink-0 text-charcoal"
          strokeWidth={1.5}
        />
        <div className="flex flex-col">
          <p className="font-host text-sm leading-5 text-charcoal">
            Don&apos;t see a student in the directory?
          </p>
          <p className="font-host text-xs leading-4 text-dark-greyish-brown">
            Invite a user to the Inspirate Consulting Hub!
          </p>
        </div>
      </div>
      <Link
        to={INVITE_USER_PATH}
        data-slot="send-invite"
        className="flex h-9 items-center rounded-lg border border-greyish-brown bg-white px-4 font-host text-xs leading-[normal] whitespace-nowrap text-charcoal outline-none transition-colors hover:bg-black/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30"
      >
        Send Invite
      </Link>
    </section>
  );
}
