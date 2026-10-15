import { useState } from 'react';
import { ChevronIcon } from '../ui/ChevronIcon';
import { DriveItem } from './DriveItem';
import type { QuarterDriveSummary, DriveSummaryResponse } from '../../api/types';

interface DriveListProps {
  data: DriveSummaryResponse | undefined;
  homeTeam: string;
  awayTeam: string;
}

export function DriveList({ data, homeTeam, awayTeam }: DriveListProps) {
  if (!data) {
    return <div className="card p-8 text-center text-[var(--text-secondary)]">Loading drives...</div>;
  }

  const quarters = Object.entries(data.quarters)
    .sort(([a], [b]) => Number(a) - Number(b))
    .map(([quarter]) => quarter); // Keep as string key

  const hasOvertime = !!data.overtime;

  return (
    <div className="space-y-4">
      {quarters.map((quarter) => (
        <QuarterDriveAccordion
          key={quarter}
          quarter={quarter}
          summary={data.quarters[quarter as unknown as number]}
          homeTeam={homeTeam}
          awayTeam={awayTeam}
        />
      ))}
      {hasOvertime && data.overtime && (
        <QuarterDriveAccordion
          quarter={5}
          summary={data.overtime}
          homeTeam={homeTeam}
          awayTeam={awayTeam}
          isOvertime
        />
      )}
    </div>
  );
}

interface QuarterDriveAccordionProps {
  quarter: string | number;
  summary: QuarterDriveSummary;
  homeTeam: string;
  awayTeam: string;
  isOvertime?: boolean;
}

function QuarterDriveAccordion({ quarter, summary, homeTeam, awayTeam, isOvertime }: QuarterDriveAccordionProps) {
  const [open, setOpen] = useState(Number(quarter) <= 2);
  const label = isOvertime ? 'OT' : `Q${quarter}`;

  const homeTotal = summary.total?.home || { drives: 0, points: 0, time_of_possession: '0:00', yards: 0 };
  const awayTotal = summary.total?.away || { drives: 0, points: 0, time_of_possession: '0:00', yards: 0 };

  return (
    <details className="group" open={open}>
      <summary
        onClick={(e) => {
          e.preventDefault();
          setOpen(!open);
        }}
        className="flex items-center justify-between p-3 bg-gray-50 border-b cursor-pointer list-none dark:bg-gray-800 dark:border-gray-700"
      >
        <h3 className="font-semibold text-nfl-dark dark:text-nfl-light">
          {label} <span className="text-sm font-normal text-gray-500 dark:text-gray-400">({summary.drives.length} drives)</span>
        </h3>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 text-sm">
            <span className="font-medium text-nfl-dark dark:text-nfl-light">{homeTeam}</span>
            <span className="text-gray-500 dark:text-gray-400">{homeTotal.points} pts</span>
            <span className="text-gray-500 dark:text-gray-400">{homeTotal.time_of_possession}</span>
            <span className="text-gray-500 dark:text-gray-400">{homeTotal.yards} yds</span>
          </div>
          <span className="text-gray-400">-</span>
          <div className="flex items-center gap-2 text-sm">
            <span className="font-medium text-nfl-light dark:text-nfl-light">{awayTeam}</span>
            <span className="text-gray-500 dark:text-gray-400">{awayTotal.points} pts</span>
            <span className="text-gray-500 dark:text-gray-400">{awayTotal.time_of_possession}</span>
            <span className="text-gray-500 dark:text-gray-400">{awayTotal.yards} yds</span>
          </div>
          <ChevronIcon className={`transition-transform ${open ? 'rotate-180' : ''} text-gray-500 dark:text-gray-400`} />
        </div>
      </summary>
      <div className="p-4 space-y-2 bg-white dark:bg-gray-900">
        {summary.drives.map((drive, index) => (
          <DriveItem key={drive.drive_id} drive={drive} index={index} />
        ))}
      </div>
    </details>
  );
}