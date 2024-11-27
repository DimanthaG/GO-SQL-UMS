import React, { useState } from 'react';
import axios from 'axios';

const Signup = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();

    // Basic form validation
    if (!email || !password) {
      setMessage('Email and password are required.');
      return;
    }

    // Make a POST request to the signup endpoint
    axios.post('http://localhost:8080/signup', { email, password })
      .then((response) => {
        setMessage('Signup successful!');
      })
      .catch((error) => {
        if (error.response) {
          // Handle server-side errors
          setMessage(`Signup failed: ${error.response.data}`);
        } else if (error.request) {
          // Handle cases where the request was made but no response was received
          setMessage('Signup failed: No response from the server.');
        } else {
          // Handle other types of errors
          setMessage(`Signup failed: ${error.message}`);
        }
      });
  };

  return (
    <div>
      <h2>Signup</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Signup</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
};

export default Signup;
