export function getCurrentNFLWeek(season: number): number {
  const now = new Date();
  
  if (now.getFullYear() < season) {
    return 1;
  }
  
  if (now.getFullYear() > season) {
    return 18;
  }
  
  const seasonStart = new Date(season, 8, 1);
  const firstThursday = new Date(seasonStart);
  firstThursday.setDate(firstThursday.getDate() + ((4 - firstThursday.getDay() + 7) % 7));
  
  const week1Start = new Date(firstThursday);
  week1Start.setDate(week1Start.getDate() - 3);
  
  const diffTime = now.getTime() - week1Start.getTime();
  const diffWeeks = Math.floor(diffTime / (7 * 24 * 60 * 60 * 1000)) + 1;
  
  return Math.max(1, Math.min(18, diffWeeks));
}

export function getCurrentNFLSeason(): number {
  const now = new Date();
  const month = now.getMonth();
  const year = now.getFullYear();
  
  if (month >= 8) {
    return year;
  }
  
  return year - 1;
}