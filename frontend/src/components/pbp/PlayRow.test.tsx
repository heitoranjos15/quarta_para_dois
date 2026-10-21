import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import type { PlayDetail } from '../../api/types';
import { PlayRow } from './PlayRow';

const play: PlayDetail = {
  play_id: 1,
  quarter: 1,
  down: 1,
  distance: 10,
  yardline: 75,
  play_type: 'pass',
  description: 'Love pass complete to Watson for 15 yards',
  epa: 0.42,
  passer: 'J.Love',
  receiver: 'C.Watson',
  rusher: '',
  tackler: '',
  is_scoring_play: false,
  drive_id: 1,
  drive_play_count: 4,
  drive_yards: 50,
  drive_result: 'TOUCHDOWN',
  drive_top: '2:35',
};

describe('PlayRow', () => {
  it('shows EPA with a sign and one decimal', () => {
    render(<PlayRow play={play} />);

    expect(screen.getByText('EPA: +0.4')).toBeInTheDocument();
  });

  it('shows a dash instead of EPA when nflverse has none', () => {
    render(<PlayRow play={{ ...play, epa: null }} />);

    expect(screen.getByText('EPA: —')).toBeInTheDocument();
  });
});
