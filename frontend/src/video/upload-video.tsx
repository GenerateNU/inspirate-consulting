import { useState } from "react";
import { useUploadVideo } from "../api/endpoints";

export default function UploadVideo() {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState("");

  const { trigger, isMutating } = useUploadVideo();

  const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!file) {
      setStatus("Please choose a file.");
      return;
    }

    try {
      // 1. get a presigned upload URL from backend
      const presignResult = await trigger({ original_filename: file.name });
      if (presignResult.status !== 200) {
        setStatus("Could not get an upload URL.");
        return;
      }

      // will need s3 key to store in database
      // const { s3_key, upload_url } = presignResult.data;
      const { upload_url } = presignResult.data;

      // 2. PUT the file to the fetched presigned URL (S3)
      const uploadResponse = await fetch(upload_url, {
        method: "PUT",
        body: file,
        headers: { "Content-Type": file.type },
      });

      if (!uploadResponse.ok) {
        setStatus("Upload to S3 failed.");
        return;
      }

      setStatus(`Uploaded: ${file.name}`);
    } catch (err) {
      console.error(err);
      setStatus("Something went wrong.");
    }
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
        <button type="submit" disabled={isMutating}>
          {isMutating ? "Uploading..." : "Upload"}
        </button>
      </form>
      {status && <p>{status}</p>}
    </div>
  );
}
