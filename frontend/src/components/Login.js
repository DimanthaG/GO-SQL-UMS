import React, { useState } from 'react';
import axios from 'axios';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!email || !password) {
      setMessage('Email and password are required.');
      return;
    }

    try {
      const response = await axios.post('http://localhost:8080/login', {
        email,
        password,
      });
      setMessage('Login successful!');
    } catch (error) {
      if (error.response) {
        // Backend returned an error response
        setMessage(`Login failed: ${error.response.data}`);
      } else if (error.request) {
        // Request was sent but no response
        setMessage('Login failed: No response from the server.');
      } else {
        // Something else happened
        setMessage(`Login failed: ${error.message}`);
      }
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button type="submit">Login</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
};

export default Login;
