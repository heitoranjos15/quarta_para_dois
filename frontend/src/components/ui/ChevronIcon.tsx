import { clsx } from 'clsx';

interface ChevronIconProps {
  className?: string;
}

export function ChevronIcon({ className }: ChevronIconProps) {
  return (
    <svg
      className={clsx('w-5 h-5 text-gray-400 transition-transform duration-200', className)}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
    </svg>
  );
}