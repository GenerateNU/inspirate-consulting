import { useListVideos } from "../api/endpoints";

export default function VideoList() {
  const { data: response, error, isLoading, mutate } = useListVideos();

  if (isLoading) {
    return (
      <div className="form-container">
        <h2>Videos</h2>
        <p>Loading...</p>
      </div>
    );
  }

  if (error || !response || response.status !== 200) {
    return (
      <div className="form-container">
        <h2>Videos</h2>
        <p>Something went wrong.</p>
      </div>
    );
  }

  const videos = response.data ?? [];

  return (
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
              onError={() => {
                // The presigned URL likely expired —> refetch to get a fresh one.
                // just calling the same list endpoint again gives a new URL
                mutate();
              }}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}
