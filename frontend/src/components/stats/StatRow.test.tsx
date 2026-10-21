import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { StatRow } from './StatRow';

const mmss = (seconds: number) =>
  `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, '0')}`;

describe('StatRow', () => {
  it('formats values with the format prop and highlights the better one', () => {
    render(<StatRow label="Time of Possession" home={517} away={383} format={mmss} />);

    expect(screen.getByText('8:37').className).toContain('text-green-600');
    expect(screen.getByText('6:23').className).not.toContain('text-green-600');
  });

  it('shows a note next to a value when given', () => {
    render(<StatRow label="Red Zone TD %" home={50} away={25} unit="%" homeNote="2/4" awayNote="1/4" />);

    expect(screen.getByText('2/4')).toBeInTheDocument();
    expect(screen.getByText('1/4')).toBeInTheDocument();
  });
});
