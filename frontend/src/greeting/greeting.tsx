import { useState } from 'react';
import { getGreeting } from '../api/endpoints';

export default function Greeting() {
    const [name, setName] = useState('');
    const [greeting, setGreeting] = useState('');

    const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setName(event.target.value);
    }

    const handleSubmit = async (event: React.SubmitEvent<HTMLFormElement>) => {
        // Prevent default form submission
        event.preventDefault();
        try {
            const result = await getGreeting(name);
            // getGreeting returns 'CreateGreetingOutputBody | ErrorModel' so first, check if there is any
            // error (no data or message), then, TS will allow us to access data.message after we are sure it exists
            if (!result.data || !('message' in result.data)) {
                setGreeting('Something went wrong.');
                return;
            }

            setGreeting(result.data.message);
  
        } catch (err) {
            console.error(err);
            setGreeting('Something went wrong.');
        }
    }

    return (
        <div className="form-container">
            <h2>Input Form</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>Name:</label>
                    <input
                        type="text"
                        value={name}
                        onChange={handleNameChange}
                    />
                </div>
                <button type="submit">Submit</button>
            </form>
            {greeting && <p>{greeting}</p>}
        </div>
    )
};
