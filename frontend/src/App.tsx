import { Routes, Route } from 'react-router-dom';
import { SeasonPage } from './pages/SeasonPage';
import { WeekPage } from './pages/WeekPage';
import { GamePage } from './pages/GamePage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<SeasonPage />} />
      <Route path="/season/:season" element={<SeasonPage />} />
      <Route path="/season/:season/week/:week" element={<WeekPage />} />
      <Route path="/game/:gameId" element={<GamePage />} />
    </Routes>
  );
}

export default App;