import { useState } from "react";
import { errorMessage } from "../api/client";
import { useCreateGreeting } from "../api/endpoints";

export default function Greeting() {
  const [name, setName] = useState("");
  const [greeting, setGreeting] = useState("");
  const { mutate, error } = useCreateGreeting();

  const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setName(event.target.value);
  };

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    // Prevent default form submission
    event.preventDefault();
    setGreeting("");
    mutate(
      { data: { name } },
      { onSuccess: (result) => setGreeting(result.greeting) },
    );
  };

  return (
    <div className="form-container">
      <h2>Input Form</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Name:</label>
          <input type="text" value={name} onChange={handleNameChange} />
        </div>
        <button type="submit">Submit</button>
      </form>
      {greeting && <p>{greeting}</p>}
      {error && <p>Error: {errorMessage(error)}</p>}
    </div>
  );
}
