import { useGetMediaAccess } from '../api/endpoints';

export default function StudentMediaList() {
  const { data, error, isLoading } = useGetMediaAccess();

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>Something went wrong.</p>;

  if (!data || !Array.isArray(data.data)) {
    return <p>Something went wrong.</p>;
  }

  return (
    <div>
      <h2>Student: My Videos</h2>
      <ul>
        {data?.data.map((media) => (
          <li key={media.id}>{media.title}</li>
        ))}
      </ul>
    </div>
  );
}