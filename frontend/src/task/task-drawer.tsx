import { useEffect } from 'react';
import type { TodoItem } from '../api/models';

type TaskDrawerProps = {
  task: TodoItem | null;
  onClose: () => void;
  onStart?: (task: TodoItem) => void;
};

function formatDueDateTime(deadline: string | null | undefined) {
  if (!deadline) return 'No due date';
  const due = new Date(deadline);
  const date = due.toLocaleDateString('en-US', { month: 'numeric', day: 'numeric' });
  const time = due.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    timeZoneName: 'short',
  });
  return `Due ${date} at ${time}`;
}

export default function TaskDrawer({ task, onClose, onStart }: TaskDrawerProps) {
  const isOpen = task !== null;

  // Close on Escape
  useEffect(() => {
    if (!isOpen) return;
    const handleKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
    };
    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  }, [isOpen, onClose]);

  return (
    <>
      {/* Backdrop: click to close */}
      <div
        onClick={onClose}
        className={`fixed inset-0 bg-black/30 transition-opacity ${
          isOpen ? 'opacity-100' : 'pointer-events-none opacity-0'
        }`}
      />

      {/* Drawer panel */}
      <aside
        role="dialog"
        aria-modal="true"
        className={`fixed right-0 top-0 flex h-full w-96 flex-col bg-white p-8 shadow-xl transition-transform duration-300 ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
      >
        {task && (
          <>
            <div className="flex items-start justify-between">
              <h2 className="text-2xl font-semibold">{task.todo_description}</h2>
              <button type="button" onClick={onClose} aria-label="Close" className="text-2xl">
                ×
              </button>
            </div>

            <p className="mt-6 text-xs">{formatDueDateTime(task.deadline)}</p>

            {/* mt-auto pushes the button to the bottom of the drawer */}
            <button
              type="button"
              onClick={() => onStart?.(task)}
              className="mt-auto flex items-center justify-end gap-2 rounded bg-black px-6 py-3 text-lg text-white"
            >
              Start →
            </button>
          </>
        )}
      </aside>
    </>
  );
}