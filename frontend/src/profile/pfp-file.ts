import Compressor from "compressorjs";
import { presignPfpUpload, updatePfp } from "../api/endpoints";
import { seedPfpCache } from "@/lib/pfp-cache";
import {
  PresignPfpUploadRequestBodyContentType,
  type PresignPfpUploadRequestBody,
} from "../api/models";

const allowedPfpTypes: string[] = Object.values(
  PresignPfpUploadRequestBodyContentType,
);

// Matches the backend's 5 MiB limit on PUT /user/pfp
const maxPfpBytes = 5 * 1024 * 1024;

export const pfpAccept = allowedPfpTypes.join(",");

// Best effort: shrink to a 512px JPEG, falling back to the original if the browser can't
export function compressPfp(file: File): Promise<File> {
  if (!allowedPfpTypes.includes(file.type)) {
    return Promise.reject(new Error("Choose a JPEG, PNG, or WebP image."));
  }
  return new Promise((resolve) => {
    new Compressor(file, {
      maxWidth: 512,
      maxHeight: 512,
      quality: 0.8,
      mimeType: "image/jpeg",
      success: (result) =>
        resolve(new File([result], "pfp.jpg", { type: result.type })),
      error: () => resolve(file),
    });
  });
}

// The presigned URL is signed for this exact type and size, so the PUT must send the same file
export function presignPfpBody(file: File): PresignPfpUploadRequestBody {
  if (file.size > maxPfpBytes) {
    throw new Error("Image must be 5 MB or smaller.");
  }
  return {
    content_type: file.type as PresignPfpUploadRequestBodyContentType,
    content_length: file.size,
  };
}

export async function uploadPfp(file: File, userId: string) {
  const { s3_key, upload_url } = await presignPfpUpload(presignPfpBody(file));
  const uploadResponse = await fetch(upload_url, {
    method: "PUT",
    body: file,
    headers: { "Content-Type": file.type },
  });
  if (!uploadResponse.ok) throw new Error("Upload to S3 failed.");
  const pfp = await updatePfp({ s3_key });
  await seedPfpCache(userId, pfp.s3_key, file);
  return pfp;
}
