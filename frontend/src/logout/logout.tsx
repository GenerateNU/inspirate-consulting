import { useState } from "react";
import { useAuthContext } from "../contexts/use-auth-context.tsx";
import { ErrorMessage } from "@/components/ui/error-message";

export default function LogOut() {
  const [status, setStatus] = useState("");
  const [error, setError] = useState<unknown>(null);

  const { logout } = useAuthContext();

  const handleSubmit = async () => {
    setError(null);
    try {
      await logout();
      setStatus(`Logged out`);
    } catch (err) {
      setError(err);
    }
  };

  return (
    <div>
      <button onClick={handleSubmit}>Log Out</button>
      {status && <p>{status}</p>}
      <ErrorMessage error={error} prefix="Couldn't log out" />
    </div>
  );
}
