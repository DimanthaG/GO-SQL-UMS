import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Login.css';

const Signup = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const handleSignup = async (e) => {
    e.preventDefault();

    // Frontend validation for password length
    if (password.length < 8) {
      setMessage('Password must be at least 8 characters long.');
      return;
    }

    try {
      const response = await fetch('/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        setMessage('Signup successful!');
        setTimeout(() => navigate('/login'), 2000); // Redirect to login page after 2 seconds
      } else if (response.status === 409) {
        setMessage('Email already exists. Please use a different email.');
      } else if (response.status === 400) {
        setMessage('Invalid input. Ensure your email and password are correct.');
      } else {
        setMessage('Signup failed. Please try again later.');
      }
    } catch (error) {
      setMessage('An error occurred. Please try again.');
    }
  };

  return (
    <div className="signup-container">
      <h2>Signup</h2>
      <form onSubmit={handleSignup}>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password (min 8 characters)"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Signup</button>
      </form>
      {message && <p>{message}</p>}
      <div>
        <p>Already have an account?</p>
        <button onClick={() => navigate('/')}>Go to Login</button>
      </div>
    </div>
  );
};

export default Signup;
