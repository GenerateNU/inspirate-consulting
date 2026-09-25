import { useGetTodoItems, useUpdateTodoItemCompleted } from "../api/endpoints"
import type { TodoItem } from "../api/models";

function TaskItem({ task, onToggled }: { task: TodoItem; onToggled: () => void }){
    const {trigger} = useUpdateTodoItemCompleted(task.id);

    const handleToggle = async () =>{
        try{
            await trigger({ completed: task.completed_at == null });
            onToggled();
        } catch(err){
            console.error(err)
        }
    };

    return(
        <li>
            <input
                type="checkbox"
                checked={task.completed_at !== null}
                onChange={handleToggle}
            />
            {task.todo_description}
        </li>
    );

}

export default function TaskList() {
    const { data, error, isLoading, mutate } = useGetTodoItems();

    if (isLoading) return <p>Loading...</p>;
    if (error) return <p>Something went wrong.</p>;

    if (!data || !Array.isArray(data.data)) {
        return <p>Something went wrong.</p>;
    }

    return (
        <div>
            <h2>To-Do List</h2>
            <ul>
                {data?.data.map((task) => (
                    <TaskItem key={task.id} task={task} onToggled={mutate} />
                ))}
            </ul>
        </div>
    );
}

