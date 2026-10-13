import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useLoginUser } from "../api/endpoints";

export default function Login() {
  const [Email, setEmail] = useState("");
  const [Password, setPassword] = useState("");
  const [Status, setStatus] = useState("");
  const navigate = useNavigate();

  const { trigger, isMutating, error } = useLoginUser<Error>();

  const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    try {
      const response = await trigger({ email: Email, password: Password });
      if (response.status === 200) {
        // On success navigate to root
        void navigate("/");
      } else {
        setStatus(response.data.detail ?? "Invalid email or password.");
      }
    } catch (err) {
      setStatus("Something went wrong." + err);
    }
  };

  return (
    <div className="form-container">
      <h2>Log In</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Email:</label>
          <input
            type="email"
            value={Email}
            onChange={(e) => setEmail(e.target.value)}
            autoComplete="email"
            required
          />
        </div>
        <div className="form-group">
          <label>Password:</label>
          <input
            type="password"
            value={Password}
            onChange={(e) => setPassword(e.target.value)}
            autoComplete="current-password"
            required
          />
        </div>
        <button type="submit" disabled={isMutating}>
          {isMutating ? "Logging In..." : "Log In"}
        </button>
      </form>
      {Status && <p>{Status}</p>}
      {error && <p>Error: {String(error)}</p>}
    </div>
  );
}
