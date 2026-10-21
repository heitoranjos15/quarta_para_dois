import { useState } from 'react';
import { ChevronIcon } from '../ui/ChevronIcon';
import type { DriveDetail, PlayDetail } from '../../api/types';
import { epaColor, formatEPA } from '../../utils/epa';

interface DriveItemProps {
  drive: DriveDetail;
  index: number;
}

export function DriveItem({ drive, index }: DriveItemProps) {
  const [expanded, setExpanded] = useState(false);
  const resultColor = getResultColor(drive.result);

  return (
    <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full p-4 flex items-center justify-between gap-4 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 text-left"
      >
        <div className="flex items-center gap-3 min-w-0 flex-1">
          <span className="text-sm font-mono text-gray-500 dark:text-gray-400 w-8">#{index + 1}</span>
          <span className={`px-2 py-1 rounded text-xs font-medium ${resultColor}`}>
            {drive.result.replace(/_/g, ' ')}
          </span>
          {drive.points > 0 && (
            <span className="text-xs font-mono font-medium text-green-700 dark:text-green-300">+{drive.points}</span>
          )}
          <span className="text-sm font-medium text-nfl-dark dark:text-nfl-light">
            {drive.team}
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400">
            {drive.start_yardline} → {drive.end_yardline}
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400 font-mono">
            {drive.yards} yds
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400 font-mono">
            {drive.plays_count} plays
          </span>
          <span className="text-xs text-gray-500 dark:text-gray-400 font-mono">
            {drive.time_of_possession}
          </span>
        </div>
        <ChevronIcon className={`transition-transform ${expanded ? 'rotate-180' : ''} text-gray-500 dark:text-gray-400`} />
      </button>
      {expanded && (
        <div className="p-4 bg-gray-50 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700">
          <div className="space-y-2 max-h-60 overflow-y-auto">
            {drive.plays.map((play) => (
              <DrivePlayRow key={play.play_id} play={play} />
            ))}
            <div className="pt-2 border-t border-gray-200 dark:border-gray-700">
              <p className="text-xs text-gray-600 dark:text-gray-400 font-medium">
                END PLAY: {drive.end_play_description}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function DrivePlayRow({ play }: { play: PlayDetail }) {
  const playTypeIcon = getPlayTypeIcon(play.play_type);

  return (
    <div className="flex items-start gap-2 text-xs">
      <span className="w-10 font-mono text-gray-500 dark:text-gray-400">{play.down}&{play.distance}</span>
      <span className="w-6">{playTypeIcon}</span>
      <span className="flex-1 text-gray-700 dark:text-gray-300">{play.description}</span>
      <span className={`${epaColor(play.epa)} font-mono whitespace-nowrap`}>EPA: {formatEPA(play.epa)}</span>
    </div>
  );
}

function getResultColor(result: string) {
  switch (result) {
    case 'TOUCHDOWN':
      return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
    case 'FIELD_GOAL':
      return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300';
    case 'PUNT':
      return 'bg-gray-100 text-gray-700 dark:bg-gray-900 dark:text-gray-300';
    case 'MISSED_FIELD_GOAL':
    case 'INTERCEPTION':
    case 'FUMBLE_LOST':
    case 'DOWNS':
      return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
    case 'SAFETY':
      return 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300';
    case 'END_OF_HALF':
    case 'END_OF_GAME':
      return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300';
    default:
      return 'bg-gray-100 text-gray-700 dark:bg-gray-900 dark:text-gray-300';
  }
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