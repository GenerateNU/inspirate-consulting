import { useState } from 'react';
import { createGreeting } from '../api/endpoints';

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
            const result = await createGreeting({ name });
            // createGreeting returns 'GreetingMessageBody | ErrorModel' so check if greeting is present
            if (!result.data || !('greeting' in result.data)) {
                setGreeting('Something went wrong.');
                return;
            }

            setGreeting(result.data.greeting);
  
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
