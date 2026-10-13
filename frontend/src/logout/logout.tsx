import { useState } from "react";
import { useAuthContext } from "../contexts/use-auth-context.tsx";

export default function LogOut() {
  const [Status, setStatus] = useState("");

  const { logout } = useAuthContext();

  const handleSubmit = async () => {
    try {
      const response = await logout();
      if (response.status === 204) {
        setStatus(`Logged out`);
      } else {
        setStatus("Something went wrong.");
      }
    } catch (err) {
      console.error(err);
      setStatus("Something went wrong.");
    }
  };

  return (
    <div>
      <button onClick={handleSubmit}>Log Out</button>
      {Status && <p>{Status}</p>}
    </div>
  );
}
