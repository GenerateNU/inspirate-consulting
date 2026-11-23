import { uploadVideo } from "../../../api/endpoints";
import { presignVideoBody } from "../../../video/video-file";

// reads a video file's duration in the browser, rounded to whole minutes
export function getVideoDurationMins(file: File): Promise<number> {
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

// gets a presigned URL, PUTs the file to S3, and returns the new S3 key
export async function uploadToS3(file: File) {
  const { s3_key, upload_url } = await uploadVideo(presignVideoBody(file));
  const response = await fetch(upload_url, {
    method: "PUT",
    body: file,
    headers: { "Content-Type": file.type },
  });
  if (!response.ok) throw new Error("Upload to S3 failed.");
  return s3_key;
}
