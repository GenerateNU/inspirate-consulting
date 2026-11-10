import { useListVideos } from "../api/endpoints";
import PageLinks, { type PageLink } from "@/components/page-links";

// Temporary links to this page's sub-pages (see the Figma sitemap).
const pageLinks: PageLink[] = [
  { to: "/student-media-list", label: "My Videos", role: "Student" },
  { to: "/counselor-media-list", label: "All Videos", role: "Counselor" },
  { to: "/upload-video", label: "Upload Video", role: "Counselor" },
  { to: "/upload-media", label: "Upload Media", role: "Counselor" },
];

export default function VideoList() {
  const { data, error, isLoading, refetch } = useListVideos();

  if (isLoading) {
    return (
      <>
        <PageLinks links={pageLinks} />
        <div className="form-container">
          <h2>Videos</h2>
          <p>Loading...</p>
        </div>
      </>
    );
  }

  if (error) {
    return (
      <>
        <PageLinks links={pageLinks} />
        <div className="form-container">
          <h2>Videos </h2>
          <p>Something went wrong. </p>
        </div>
      </>
    );
  }

  const videos = data ?? [];

  return (
    <>
      <PageLinks links={pageLinks} />
      <div className="form-container">
        <h2>Videos</h2>
        {videos.length === 0 && <p>No videos yet.</p>}
        <ul>
          {videos.map((video) => (
            <li key={video.s3_key}>
              <p>{video.s3_key}</p>
              <video
                src={video.download_url}
                controls
                style={{ maxWidth: "400px" }}
                onError={async () => {
                  // The presigned URL likely expired —> refetch to get a fresh one.
                  // just calling the same list endpoint again gives a new URL
                  await refetch();
                }}
              />
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}
