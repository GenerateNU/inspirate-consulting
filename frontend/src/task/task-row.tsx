import type { TodoItem } from '../api/models';

type TaskRowProps = {
  task: TodoItem;
  title: string;
  description?: string;
  onToggleComplete: (task: TodoItem) => void;
  onOpen: (task: TodoItem) => void;
};

export default function TaskRow({
  task,
  title,
  description,
  onToggleComplete,
  onOpen,
}: TaskRowProps) {
  const completed = task.completed_at != null;

  return (
    <div
      className={`flex items-center gap-4 rounded-lg px-4 py-3 ${
        completed ? 'bg-gray-200' : 'bg-white'
      }`}
    >
      <input
        type="checkbox"
        checked={completed}
        onChange={() => onToggleComplete(task)}
        aria-label={`Mark "${title}" as ${completed ? 'not completed' : 'completed'}`}
        className="h-4 w-4 shrink-0"
      />

      <button
        type="button"
        onClick={() => onOpen(task)}
        className="flex flex-1 items-center justify-between text-left"
      >
        <div>
          <p className="text-sm font-medium">{title}</p>
          {description && <p className="text-xs text-gray-500">{description}</p>}
        </div>

        <span className="text-2xl" aria-hidden="true">
          →
        </span>
      </button>
    </div>
  );
}