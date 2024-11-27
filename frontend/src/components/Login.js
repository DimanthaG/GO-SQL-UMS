import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './styles.css';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [forgotPasswordEmail, setForgotPasswordEmail] = useState('');
  const [showForgotPassword, setShowForgotPassword] = useState(false);
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        navigate('/video');
      } else {
        setMessage('Login failed!');
      }
    } catch (error) {
      setMessage('An error occurred.');
    }
  };

  const handleForgotPassword = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('/forgot-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email: forgotPasswordEmail }),
      });

      if (response.ok) {
        setMessage('A reset token has been sent to your email.');
      } else {
        setMessage('Failed to send reset token. Please try again.');
      }
    } catch (error) {
      setMessage('An error occurred.');
    }
  };

  return (
    <div className="login-container">
      <h2>Login</h2>
      {!showForgotPassword ? (
        <form onSubmit={handleLogin}>
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
          <button type="submit">Login</button>
          <p onClick={() => setShowForgotPassword(true)} className="forgot-password-link">
            Forgot Password?
          </p>
          <div>
          <p>Don't have an account?</p>
        <button onClick={() => navigate('/signup')}>Go to Signup</button>
      </div>
        </form>
      ) : (
        <form onSubmit={handleForgotPassword}>
          <input
            type="email"
            placeholder="Enter your email"
            value={forgotPasswordEmail}
            onChange={(e) => setForgotPasswordEmail(e.target.value)}
            required
          />
          <button type="submit">Send Reset Token</button>
          <p onClick={() => setShowForgotPassword(false)} className="forgot-password-link">
            Back to Login
          </p>
        </form>
      )}
      {message && <p>{message}</p>}
    </div>
  );
};

export default Login;
