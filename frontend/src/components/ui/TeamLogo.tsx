import { clsx } from 'clsx';

const teamColors: Record<string, { primary: string; secondary: string; primaryDark: string; secondaryDark: string }> = {
  ARI: { primary: '#97233F', secondary: '#000000', primaryDark: '#C73659', secondaryDark: '#FFFFFF' },
  ATL: { primary: '#A71930', secondary: '#000000', primaryDark: '#D62D3F', secondaryDark: '#FFFFFF' },
  BAL: { primary: '#241773', secondary: '#9E7C0C', primaryDark: '#4A3FC7', secondaryDark: '#C8A000' },
  BUF: { primary: '#00338D', secondary: '#C60C30', primaryDark: '#3D6FE8', secondaryDark: '#FF3D3D' },
  CAR: { primary: '#0085CA', secondary: '#101820', primaryDark: '#4DD1FF', secondaryDark: '#FFFFFF' },
  CHI: { primary: '#0B162A', secondary: '#C83803', primaryDark: '#4A6FA5', secondaryDark: '#FF6B00' },
  CIN: { primary: '#FB4F14', secondary: '#000000', primaryDark: '#FF7840', secondaryDark: '#FFFFFF' },
  CLE: { primary: '#311D00', secondary: '#FF3C00', primaryDark: '#8B5A00', secondaryDark: '#FF6B3D' },
  DAL: { primary: '#003594', secondary: '#869397', primaryDark: '#3D6FE8', secondaryDark: '#B0BEC5' },
  DEN: { primary: '#FB4F14', secondary: '#002244', primaryDark: '#FF7840', secondaryDark: '#4A90D9' },
  DET: { primary: '#0076B6', secondary: '#B0B7BC', primaryDark: '#4DD1FF', secondaryDark: '#E0E0E0' },
  GB: { primary: '#203731', secondary: '#FFB612', primaryDark: '#4A7C6E', secondaryDark: '#FFD600' },
  HOU: { primary: '#03202F', secondary: '#A71930', primaryDark: '#00D4B8', secondaryDark: '#FF6B3D' },
  IND: { primary: '#002C5F', secondary: '#A2AAAD', primaryDark: '#3D6FE8', secondaryDark: '#D0D0D0' },
  JAX: { primary: '#006778', secondary: '#9F792C', primaryDark: '#00E5FF', secondaryDark: '#FFD600' },
  KC: { primary: '#E31837', secondary: '#FFB81C', primaryDark: '#FF4D4D', secondaryDark: '#FFD600' },
  LV: { primary: '#000000', secondary: '#A5ACAF', primaryDark: '#4A4A4A', secondaryDark: '#E0E0E0' },
  LAC: { primary: '#0080C6', secondary: '#FFC20E', primaryDark: '#4DD1FF', secondaryDark: '#FFD600' },
  LA: { primary: '#003594', secondary: '#FFA300', primaryDark: '#3D6FE8', secondaryDark: '#FFD600' },
  LAR: { primary: '#003594', secondary: '#FFA300', primaryDark: '#3D6FE8', secondaryDark: '#FFD600' },
  MIA: { primary: '#008E97', secondary: '#FC4C02', primaryDark: '#00E5FF', secondaryDark: '#FF6B3D' },
  MIN: { primary: '#4F2683', secondary: '#FFC62F', primaryDark: '#9B59D1', secondaryDark: '#FFD600' },
  NE: { primary: '#002244', secondary: '#C60C30', primaryDark: '#3D6FE8', secondaryDark: '#FF6B3D' },
  NO: { primary: '#D3BC8D', secondary: '#101820', primaryDark: '#FFF3C4', secondaryDark: '#000000' },
  NYG: { primary: '#0B2265', secondary: '#A71930', primaryDark: '#4A6FA5', secondaryDark: '#FF6B3D' },
  NYJ: { primary: '#125740', secondary: '#000000', primaryDark: '#4A9C7D', secondaryDark: '#FFFFFF' },
  PHI: { primary: '#004C54', secondary: '#A5ACAF', primaryDark: '#00D4B8', secondaryDark: '#E0E0E0' },
  PIT: { primary: '#FFB612', secondary: '#101820', primaryDark: '#FFD600', secondaryDark: '#000000' },
  SF: { primary: '#AA0000', secondary: '#B3995D', primaryDark: '#FF3D3D', secondaryDark: '#E0C870' },
  SEA: { primary: '#002244', secondary: '#69BE28', primaryDark: '#3D6FE8', secondaryDark: '#A8E670' },
  TB: { primary: '#D50A0A', secondary: '#FF7900', primaryDark: '#FF3D3D', secondaryDark: '#FFB84D' },
  TEN: { primary: '#0C2340', secondary: '#4B92DB', primaryDark: '#3D6FE8', secondaryDark: '#7DC7FF' },
  WAS: { primary: '#5A1414', secondary: '#FFB612', primaryDark: '#D64040', secondaryDark: '#FFD600' },
};

interface TeamLogoProps {
  team: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

const sizeClasses = {
  sm: 'w-8 h-8 text-[10px]',
  md: 'w-10 h-10 text-xs',
  lg: 'w-12 h-12 text-sm',
};

export function TeamLogo({ team, size = 'md', className }: TeamLogoProps) {
  const colors = teamColors[team] || { primary: '#6B7280', secondary: '#FFFFFF', primaryDark: '#9CA3AF', secondaryDark: '#FFFFFF' };

  return (
    <div
      className={clsx(
        'flex items-center justify-center rounded-lg font-bold text-white select-none',
        'bg-[var(--team-primary)] dark:bg-[var(--team-primary-dark)]',
        sizeClasses[size],
        className
      )}
      style={{
        '--team-primary': colors.primary,
        '--team-primary-dark': colors.primaryDark,
      } as React.CSSProperties}
      aria-label={team}
    >
      {team}
    </div>
  );
}