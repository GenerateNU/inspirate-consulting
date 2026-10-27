type TaskStatCardProps = {
  label: string;
  value: string;
  caption: string;
};

export default function TaskStatCard({ label, value, caption }: TaskStatCardProps) {
  return (
    <div className="flex w-36 flex-col rounded-md bg-white p-2">
      <p className="text-xs text-gray-600">{label}</p>
      <p className="text-2xl font-semibold">{value}</p>
      <p className="text-xs text-gray-600">{caption}</p>
    </div>
  );
}

