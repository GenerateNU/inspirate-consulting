import { useListMedia } from '../api/endpoints';
import type { Media } from '../api/models';

export default function CounselorMediaList() {
  const { data, error, isLoading } = useListMedia({ limit: 10, offset: 0 });

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>Something went wrong.</p>;

  if (!data || !Array.isArray(data.data)) {
    return <p>Something went wrong.</p>;
  }

  return (
    <div>
      <h2>Counselor: All Videos</h2>
      <ul>
        {data.data.map((media: Media) => (
          <li key={media.id}>{media.title}</li>
        ))}
      </ul>
    </div>
  );
}