import { useState } from "react";
import { errorMessage } from "../api/client";
import { useCreateUser } from "../api/endpoints";
import { CreateUserInputBodyRole } from "../api/models";

export default function CreateUser() {
  const [Name, setName] = useState("");
  const [Email, setEmail] = useState("");
  const [Role, setRole] = useState<CreateUserInputBodyRole>(
    CreateUserInputBodyRole.student,
  );
  const [Status, setStatus] = useState("");

  const { mutate, isPending, error } = useCreateUser();

  const handleSubmit = (event: React.SubmitEvent<HTMLElement>) => {
    event.preventDefault();
    setStatus("");
    mutate(
      { data: { name: Name, email: Email, role: Role, pfp_key: null } },
      {
        onSuccess: () =>
          setStatus(
            `Created user ${Name}. An email with sign-in instructions has been sent to ${Email}.`,
          ),
      },
    );
  };

  return (
    <div className="form-container">
      <h2>Create Account</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Name:</label>
          <input
            type="text"
            value={Name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Email:</label>
          <input
            type="text"
            value={Email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Role:</label>
          <select
            value={Role}
            onChange={(e) => setRole(e.target.value as CreateUserInputBodyRole)}
          >
            <option value={CreateUserInputBodyRole.student}>Student</option>
            <option value={CreateUserInputBodyRole.counselor}>Counselor</option>
          </select>
        </div>
        <button type="submit" disabled={isPending}>
          {isPending ? "Creating User..." : "Create User"}
        </button>
      </form>
      {Status && <p>{Status}</p>}
      {error && <p>Error: {errorMessage(error)}</p>}
    </div>
  );
}
