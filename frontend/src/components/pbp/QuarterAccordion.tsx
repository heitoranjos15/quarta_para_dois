import { useState } from 'react';
import { ChevronIcon } from '../ui/ChevronIcon';
import { PlayRow } from './PlayRow';
import type { PlayDetail } from '../../api/types';

interface QuarterAccordionProps {
  quarter: number;
  plays: PlayDetail[];
}

export function QuarterAccordion({ quarter, plays }: QuarterAccordionProps) {
  const [open, setOpen] = useState(quarter <= 2);

  return (
    <details className="group" open={open} onToggle={() => setOpen(!open)}>
      <summary className="flex items-center justify-between p-3 bg-gray-50 border-b cursor-pointer list-none dark:bg-gray-800 dark:border-gray-700">
        <h3 className="font-semibold text-nfl-dark dark:text-nfl-light">
          Q{quarter} <span className="text-sm font-normal text-gray-500 dark:text-gray-400">({plays.length} plays)</span>
        </h3>
        <ChevronIcon className="transition-transform group-open:rotate-180 text-gray-500 dark:text-gray-400" />
      </summary>
      <div className="p-2 space-y-1">
        {plays.map((play) => (
          <PlayRow key={play.play_id} play={play} />
        ))}
      </div>
    </details>
  );
}