import React from 'react';
import { useNavigate } from 'react-router-dom';
import './styles.css';

const YouTubePage = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('userToken');
    navigate('/login');
  };

  const handleResetPassword = () => {
    navigate('/reset-password');
  };

  return (
    <div className="youtube-container">
      <h2>Woah, you logged in!</h2>
      <div style={{ textAlign: 'center', marginTop: '20px' }}>
        <iframe
          width="560"
          height="315"
          src="https://www.youtube.com/embed/SuwPaNAZlF0" // Replace with your video URL
          title="YouTube video"
          frameBorder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowFullScreen
        ></iframe>
      </div>
      <button className="logout-button" onClick={handleLogout}>
        Logout
      </button>
      <button className="reset-password-button" onClick={handleResetPassword}>
        Reset Password
      </button>
    </div>
  );
};

export default YouTubePage;
