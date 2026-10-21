/** EPA is null when nflverse has no value for a play. */
export function formatEPA(epa: number | null): string {
  if (epa === null) return '—';
  return `${epa > 0 ? '+' : ''}${epa.toFixed(1)}`;
}

export function epaColor(epa: number | null): string {
  if (epa !== null && epa > 0) return 'text-green-600 dark:text-green-400';
  if (epa !== null && epa < 0) return 'text-red-600 dark:text-red-400';
  return 'text-gray-500 dark:text-gray-400';
}
