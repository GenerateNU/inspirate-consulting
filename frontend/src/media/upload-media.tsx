import { useState } from "react";
import { errorMessage } from "../api/client";
import { useMutation } from "@tanstack/react-query";
import { createMedia, uploadVideo } from "../api/endpoints";

// reads a video file's duration in the browser, rounded to whole minutes
function getVideoDurationMins(file: File): Promise<number> {
  return new Promise((resolve) => {
    const url = URL.createObjectURL(file);
    const video = document.createElement("video");
    video.preload = "metadata";
    video.onloadedmetadata = () => {
      URL.revokeObjectURL(url);
      resolve(Math.max(1, Math.round(video.duration / 60)));
    };
    video.onerror = () => {
      URL.revokeObjectURL(url);
      resolve(0);
    };
    video.src = url;
  });
}

export default function CreateMedia() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [fileInputKey, setFileInputKey] = useState(0);

  const { mutate, isPending, error } = useMutation({
    mutationFn: async (file: File) => {
      // get a unique S3 key and a presigned upload URL from the backend
      const { s3_key, upload_url } = await uploadVideo({
        original_filename: file.name,
      });

      // PUT the file to the presigned URL (S3)
      const uploadResponse = await fetch(upload_url, {
        method: "PUT",
        body: file,
        headers: { "Content-Type": file.type },
      });
      if (!uploadResponse.ok) throw new Error("Upload to S3 failed.");

      // save the media row with the real S3 key
      const lengthInMins = await getVideoDurationMins(file);
      return createMedia({
        title,
        description,
        length_in_mins: lengthInMins,
        school_year: null,
        s3_key,
      });
    },
    onSuccess: (media) => {
      setStatus(`Created: ${media.title}`);
      setTitle("");
      setDescription("");
      setFile(null);
      setFileInputKey((k) => k + 1);
    },
  });

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!file) {
      setStatus("Please choose a file.");
      return;
    }
    setStatus("");
    mutate(file);
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
        <div className="form-group">
          <label>Video File:</label>
          <input
            key={fileInputKey}
            type="file"
            accept="video/*"
            onChange={(e) => setFile(e.target.files?.[0] ?? null)}
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
