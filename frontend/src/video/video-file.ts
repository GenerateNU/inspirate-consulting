import {
  PresignUploadRequestBodyContentType,
  type PresignUploadRequestBody,
} from "../api/models";

const allowedVideoTypes: string[] = Object.values(
  PresignUploadRequestBodyContentType,
);

export const videoAccept = allowedVideoTypes.join(",");

// The presigned URL is signed for this exact type and size, so the PUT must send the same file
export function presignVideoBody(file: File): PresignUploadRequestBody {
  if (!allowedVideoTypes.includes(file.type)) {
    throw new Error("Only MP4, MOV, or WebM videos can be uploaded.");
  }
  return {
    original_filename: file.name,
    content_type: file.type as PresignUploadRequestBodyContentType,
    content_length: file.size,
  };
}
