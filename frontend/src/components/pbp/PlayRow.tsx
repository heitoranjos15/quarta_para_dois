import type { PlayDetail } from '../../api/types';

export function PlayRow({ play }: { play: PlayDetail }) {
  const epaColor = play.epa > 0 ? 'text-green-600 dark:text-green-400' : play.epa < 0 ? 'text-red-600 dark:text-red-400' : 'text-gray-500 dark:text-gray-400';
  const playTypeIcon = getPlayTypeIcon(play.play_type);

  return (
    <div className="p-3 bg-white border rounded-lg hover:bg-gray-50 transition-colors dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700">
      <div className="flex items-start gap-3">
        <span className="w-16 text-xs text-gray-500 dark:text-gray-400 font-mono">{play.down}&{play.distance}</span>
        <span className="w-8">{playTypeIcon}</span>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-nfl-dark dark:text-nfl-light">{play.description}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400">
            {play.passer && `QB: ${play.passer}`}
            {play.receiver && ` \u2192 ${play.receiver}`}
            {play.rusher && ` RB: ${play.rusher}`}
            {play.tackler && ` \u2014 tackled by ${play.tackler}`}
          </p>
        </div>
        <span className={`${epaColor} font-mono font-medium whitespace-nowrap`}>EPA: {play.epa > 0 ? '+' : ''}{play.epa.toFixed(1)}</span>
      </div>
      {play.is_scoring_play && (
        <span className="text-xs bg-nfl-gold text-nfl-dark px-2 py-0.5 rounded">SCORING PLAY</span>
      )}
    </div>
  );
}

function getPlayTypeIcon(playType: string) {
  switch (playType.toLowerCase()) {
    case 'pass':
      return '\u2197';
    case 'run':
      return '\u2192';
    case 'punt':
      return '\u2191';
    case 'field_goal':
      return '\u2601';
    case 'kickoff':
      return '\u26A1';
    case 'extra_point':
      return '\u2605';
    case 'qb_kneel':
      return '\u25BC';
    case 'qb_spike':
      return '\u23F9';
    default:
      return '\u2022';
  }
}