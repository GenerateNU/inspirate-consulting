import { useState } from "react";
import { useCreateUser } from "../api/endpoints";
import { CreateUserInputBodyRole } from "../api/models";

export default function CreateUser() {
  const [Name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [Role, setRole] = useState<CreateUserInputBodyRole>(
    CreateUserInputBodyRole.student,
  );
  const [Status, setStatus] = useState("");

  const { trigger, isMutating, error } = useCreateUser<Error>();

  const handleSubmit = async (event: React.SubmitEvent<HTMLElement>) => {
    event.preventDefault();
    try {
      const response = await trigger({
        name: Name,
        email: email,
        role: Role,
        pfp_key: null,
      });
      if (response.status === 200) {
        setStatus(
          `Created user ${Name}. Temporary password: ${response.data.temp_password}`,
        );
      } else {
        setStatus("Something went wrong.");
      }
    } catch (err) {
      setStatus("Something went wrong." + err);
    }
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
            value={email}
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
        <button type="submit" disabled={isMutating}>
          {isMutating ? "Creating User..." : "Create User"}
        </button>
      </form>
      {Status && <p>{Status}</p>}
      {error && <p>Error: {String(error)}</p>}
    </div>
  );
}
