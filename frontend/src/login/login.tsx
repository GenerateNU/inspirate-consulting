import { useState } from "react";
import { useAuthContext } from "../contexts/use-auth-context";

export default function LogIn() {
  const [email, setEmail] = useState("");
  const [Password, setPassword] = useState("");
  const [Status, setStatus] = useState("");

  const { login } = useAuthContext();

  const handleSubmit = async (event: React.SubmitEvent<HTMLElement>) => {
    event.preventDefault();
    try {
      {
        const response = await login(email, Password);
        if (response.status != 200) {
          setStatus("Something went wrong");
        } else {
          setStatus("Logged in");
        }
      }
    } catch (err) {
      console.error(err);
      setStatus("Something went wrong");
    }
  };

  return (
    <div className="form-container">
      <h2>Log in</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Email:</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Password:</label>
          <input
            type="password"
            value={Password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">Log In</button>
      </form>
      {Status && <p>{Status}</p>}
    </div>
  );
}
