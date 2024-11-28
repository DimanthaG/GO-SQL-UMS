import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './styles.css';

const YouTubePage = () => {
  const navigate = useNavigate();

  // Check authentication status when the component mounts
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await fetch('/video', {
          method: 'GET',
          credentials: 'include', // Include cookies in the request
        });

        if (!response.ok) {
          // If the user is not authenticated, redirect to login
          navigate('/');
        }
      } catch (error) {
        console.error('Error checking authentication:', error);
        navigate('/'); // Redirect to login on error
      }
    };

    checkAuth();
  }, [navigate]);

  // Handle logout
  const handleLogout = async () => {
    try {
      const response = await fetch('/logout', {
        method: 'POST',
        credentials: 'include', // Include cookies in the request
      });

      if (response.ok) {
        // Clear any tokens stored in local storage (if applicable)
        localStorage.removeItem('userToken');

        // Navigate to the login page
        navigate('/');
      } else {
        alert('Failed to log out. Please try again.');
      }
    } catch (error) {
      console.error('An error occurred during logout:', error);
      alert('An error occurred. Please try again.');
    }
  };

  const handleResetPassword = () => {
    navigate('/reset-password');
  };

  return (
    <div className="youtube-container">
      <h2>Woah, you logged in! idk what should be here so heres a video</h2>
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
