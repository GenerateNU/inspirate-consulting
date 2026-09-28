import { useState } from 'react';
import { useCreateUser } from '../api/endpoints';

export default function CreateUser() {
    const [Name, setName] = useState('');
    const [Email, setEmail] = useState('');
    const [Password, setPassword] = useState('');
    const [Status, setStatus] = useState('');


     const { trigger, isMutating, error } = useCreateUser();


    const handleSubmit = async (event: React.SubmitEvent<HTMLElement>) => {
        event.preventDefault();
        try {
        const response = await trigger({
            name: Name,
            email: Email,
            password: Password,
            pfp_key: null,

        })
        if (response.status === 200) {
            setStatus(`Created user ${Name}`)
        } else {
            setStatus('Something went wrong.');
        }
        }
     catch (err){
        console.error(err);
        setStatus('Something went wrong.'); 
     }

    };

    return (
        <div className="form-container">
            
            <h2>Create Account</h2>
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>Name:</label>
                    <input
                        type="text"
                        value={Name}
                        onChange={(e) => setName(e.target.value)}
                        required
                    />
                </div>
                <div className="form-group">
                    <label>Email:</label>
                    <input
                        type="text"
                        value={Email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                </div>
                <div className="form-group">
                    <label>Password:</label>
                    <input
                        type="text"
                        value={Password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>
                <button type="submit" disabled={isMutating}>
                    {isMutating ? 'Creating User...' : 'Create User'}
                </button>
            </form>
            {Status && <p>{Status}</p>}
            {error && <p>Error: {String(error)}</p>}
        
        
        </div>


    )


}