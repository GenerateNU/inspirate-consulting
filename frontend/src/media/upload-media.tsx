import { useState } from "react";
import { errorMessage } from "../api/client";
import { useCreateMedia } from "../api/endpoints";

export default function CreateMedia() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState("");

  const { mutate, isPending, error } = useCreateMedia();

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setStatus("");
    mutate(
      {
        data: {
          title,
          description,
          length_in_mins: 10,
          school_year: null,
          s3_key: "video-bucket/dummy-video.mp4",
        },
      },
      {
        onSuccess: (media) => {
          setStatus(`Created: ${media.title}`);
          setTitle("");
          setDescription("");
        },
      },
    );
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
        <button type="submit" disabled={isPending}>
          {isPending ? "Uploading..." : "Upload Video"}
        </button>
      </form>
      {status && <p>{status}</p>}
      {error && <p>Error: {errorMessage(error)}</p>}
    </div>
  );
}
