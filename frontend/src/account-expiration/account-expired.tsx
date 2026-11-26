import { useState } from "react";
import { useAuthContext } from "../contexts/use-auth-context";
import { ErrorMessage } from "@/components/ui/error-message";

export default function AccountExpired() {
  const { logout } = useAuthContext();
  const [error, setError] = useState<unknown>(null);

  const handleLogout = async () => {
    setError(null);
    try {
      await logout();
    } catch (err) {
      setError(err);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md rounded-lg bg-white p-8 text-center shadow-sm">
        <h1 className="text-2xl font-semibold text-gray-900">
          Account access expired
        </h1>
        <p className="mt-4 text-gray-600">
          Access to your account has expired. If you believe this is an error,
          please contact support.
        </p>
        <button
          onClick={handleLogout}
          className="mt-6 rounded-md bg-blue-700 px-4 py-2 text-sm font-medium text-white hover:bg-blue-800"
        >
          Log out
        </button>
        <ErrorMessage
          error={error}
          prefix="Couldn't log out"
          className="mt-4"
        />
      </div>
    </div>
  );
}
