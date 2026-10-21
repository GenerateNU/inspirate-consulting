import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { uploadVideo } from "../api/endpoints";
import { errorMessage } from "../api/client";

export default function UploadVideo() {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState("");

  const { mutate, isPending, error } = useMutation({
    mutationFn: async (file: File) => {
      // 1. get a presigned upload URL from backend
      // will need s3 key to store in database
      // const { s3_key, upload_url } = await uploadVideo(...);
      const { upload_url } = await uploadVideo({
        original_filename: file.name,
      });

      // 2. PUT the file to the fetched presigned URL (S3)
      const uploadResponse = await fetch(upload_url, {
        method: "PUT",
        body: file,
        headers: { "Content-Type": file.type },
      });
      if (!uploadResponse.ok) throw new Error("Upload to S3 failed.");
    },
    onSuccess: (_, file) => setStatus(`Uploaded: ${file.name}`),
  });

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
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
      <h2>Upload a Video</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Video File:</label>
          <input
            type="file"
            accept="video/*"
            onChange={(e) => setFile(e.target.files?.[0] ?? null)}
            required
          />
        </div>
        <button type="submit" disabled={isPending}>
          {isPending ? "Uploading..." : "Upload"}
        </button>
      </form>
      {status && <p>{status}</p>}
      {error && <p>Error: {errorMessage(error)}</p>}
    </div>
  );
}
