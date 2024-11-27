import React, { useState } from 'react';
import './styles.css';

const ResetPasswordPage = () => {
  const [email, setEmail] = useState('');
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [message, setMessage] = useState('');
  const [isVerified, setIsVerified] = useState(false); // Tracks if email is verified

  const handleVerifyEmail = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('/verify-email', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      if (response.ok) {
        setMessage('Email verified. You can now change your password.');
        setIsVerified(true); // Allow password change after email verification
      } else {
        setMessage('Email not found. Please try again.');
      }
    } catch (error) {
      setMessage('An error occurred. Please try again.');
    }
  };

  const handleChangePassword = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('/reset-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, currentPassword, newPassword }),
      });

      if (response.ok) {
        setMessage('Password changed successfully.');
      } else {
        setMessage('Failed to change password. Please check your credentials.');
      }
    } catch (error) {
      setMessage('An error occurred. Please try again.');
    }
  };

  return (
    <div className="reset-password-container">
      <h2>Reset Password</h2>
      {!isVerified ? (
        <form onSubmit={handleVerifyEmail}>
          <input
            type="email"
            placeholder="Enter your email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <button type="submit">Verify Email</button>
        </form>
      ) : (
        <form onSubmit={handleChangePassword}>
          <input
            type="password"
            placeholder="Enter current password"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Enter new password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
          />
          <button type="submit">Change Password</button>
        </form>
      )}
      {message && <p>{message}</p>}
    </div>
  );
};

export default ResetPasswordPage;
