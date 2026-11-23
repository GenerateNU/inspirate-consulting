import { useState } from "react";
import { LogOut } from "lucide-react";
import { useAuthContext } from "@/contexts/use-auth-context";
import { ErrorAlert } from "@/profile/components/error-alert";

const logOutButtonClass =
  "flex h-8 items-center gap-1.5 rounded-lg border border-dark-green bg-white px-4 font-host text-xs leading-[normal] whitespace-nowrap text-dark-green outline-none transition-colors hover:bg-dark-green/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed disabled:opacity-60";

export function LogOutButton() {
  const { logout } = useAuthContext();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<unknown>(null);
  const [failures, setFailures] = useState(0);

  const handleLogOut = async () => {
    setPending(true);
    setError(null);
    try {
      await logout();
    } catch (err) {
      setError(err);
      setFailures((n) => n + 1);
      setPending(false);
    }
  };

  return (
    <>
      <ErrorAlert error={error} attempt={failures} label="Couldn't log out" />
      <button
        type="button"
        data-slot="log-out"
        onClick={() => void handleLogOut()}
        disabled={pending}
        aria-busy={pending}
        className={logOutButtonClass}
      >
        <LogOut aria-hidden="true" className="size-3.5" strokeWidth={1.75} />
        {pending ? "Logging out…" : "Log Out"}
      </button>
    </>
  );
}
