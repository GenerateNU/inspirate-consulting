import { useState } from "react";
import { useAuthContext } from "../contexts/use-auth-context";

export default function AccountExpired() {
  const { logout } = useAuthContext();
  const [status, setStatus] = useState("");

  const handleLogout = async () => {
    try {
      await logout();
    } catch (err) {
      console.error(err);
      setStatus("Something went wrong.");
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
        {status && <p className="mt-4 text-sm text-red-600">{status}</p>}
      </div>
    </div>
  );
}
