import {useState} from 'react'
import { useCreateTodoItem } from '../api/endpoints';

export default function CreateTask(){
    const[description, setDescription] = useState('');
    const[status, setStatus] = useState('');

    const{trigger, isMutating, error} = useCreateTodoItem();

    const handleSubmit = async(event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        
        try {
            const response = await trigger({
                todo_description: description,
            });

            if(response.status === 200){
                setStatus(`Created: ${response.data.todo_description}`)
                setDescription('');
            }
            else{
                setStatus('Something went wrong')

            }
        } catch (err) {
            console.error(err);
            setStatus('Something went wrong.');
        }

    };

    return(
        <div className = "form-container">
            <h2>Create Task</h2>
            <form onSubmit = {handleSubmit}>
                <div className = "form-group">
                    <label>Description:</label>
                    <input
                        type = "text"
                        value = {description}
                        onChange = {(e) => setDescription(e.target.value)}
                        required
                    />
                </div>
                <button type ="submit" disabled = {isMutating}>
                    {isMutating ? 'Creating Task...' : 'Create Task'}
                </button>
            </form>
            {status && <p>{status}</p>}
            {error && <p>Error: {String(error)}</p>}
        </div>
    )

}