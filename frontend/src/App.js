import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './components/Login';
import Signup from './components/Signup';
import YoutubePage from './components/YoutubePage';
import ResetPasswordPage from './components/ResetPasswordPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/video" element={<YoutubePage />} />
        <Route path="/reset-password" element={<ResetPasswordPage />} />

      </Routes>
    </Router>
  );
}

export default App;
