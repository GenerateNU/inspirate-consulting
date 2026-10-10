import { useState } from "react";
import { useCreateMedia } from "../api/endpoints";

export default function CreateMedia() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState("");

  const { trigger, isMutating, error } = useCreateMedia<Error>();

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    try {
      const response = await trigger({
        title,
        description,
        length_in_mins: 10,
        school_year: null,
        s3_key: "video-bucket/dummy-video.mp4",
      });

      if (response.status === 200) {
        setStatus(`Created: ${response.data.title}`);
        setTitle("");
        setDescription("");
      } else {
        setStatus("Something went wrong.");
      }
    } catch (err) {
      console.error(err);
      setStatus("Something went wrong.");
    }
  };

  return (
    <div className="form-container">
      <h2>Upload Video</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Title:</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Description:</label>
          <input
            type="text"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
          />
        </div>
        <button type="submit" disabled={isMutating}>
          {isMutating ? "Uploading..." : "Upload Video"}
        </button>
      </form>
      {status && <p>{status}</p>}
      {error && <p>Error: {error.message}</p>}
    </div>
  );
}
