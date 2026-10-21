import { TeamLogo } from '../ui';

export interface StatRowProps {
  label: string;
  home: number;
  away: number;
  unit?: string;
  format?: (n: number) => string;
  homeNote?: string;
  awayNote?: string;
  homeTeam?: string;
  awayTeam?: string;
  higherIsBetter?: boolean;
}

export function StatRow({
  label,
  home,
  away,
  unit = '',
  format,
  homeNote,
  awayNote,
  homeTeam,
  awayTeam,
  higherIsBetter = true,
}: StatRowProps) {
  const homeIsBetter = higherIsBetter ? home > away : home < away;
  const awayIsBetter = higherIsBetter ? away > home : away < home;
  const show = (n: number) => (format ? format(n) : `${n}${unit}`);

  return (
    <div className="flex items-center justify-between py-2 border-b border-gray-100 dark:border-gray-800">
      <span className="text-gray-600 dark:text-gray-400 text-sm">{label}</span>
      <div className="flex items-center gap-4">
        {homeTeam && <TeamLogo team={homeTeam} size="sm" />}
        <span
          className={`font-medium text-sm ${
            homeIsBetter ? 'text-green-600 dark:text-green-400' : 'text-nfl-dark dark:text-nfl-light'
          }`}
        >
          {show(home)}
        </span>
        {homeNote && <span className="text-xs text-gray-500 dark:text-gray-400">{homeNote}</span>}
        <span className="text-gray-400 text-sm">vs</span>
        <span
          className={`font-medium text-sm ${
            awayIsBetter ? 'text-green-600 dark:text-green-400' : 'text-nfl-light dark:text-nfl-light'
          }`}
        >
          {show(away)}
        </span>
        {awayNote && <span className="text-xs text-gray-500 dark:text-gray-400">{awayNote}</span>}
        {awayTeam && <TeamLogo team={awayTeam} size="sm" />}
      </div>
    </div>
  );
}
