import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  RadarChart,
  Radar,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Legend,
} from 'recharts';
import { TeamLogo } from '../ui';
import type { TeamComparisonStats, TeamUnitStats } from '../../api/types';

const COLORS = {
  home: '#002244',
  away: '#FB4F14',
  homeLight: '#1e3a5f',
  awayLight: '#ff7043',
};

interface TeamComparisonChartsProps {
  stats: TeamComparisonStats;
  homeTeam: string;
  awayTeam: string;
}

const radarData = [
  { metric: 'Red Zone', key: 'red_zone_pct' },
  { metric: '3rd Down', key: 'third_down_pct' },
  { metric: '4th Down', key: 'fourth_down_pct' },
  { metric: 'Explosive', key: 'explosive_plays' },
];

function getRadarValue(stats: TeamUnitStats, key: string): number {
  const value = stats[key as keyof TeamUnitStats];
  if (typeof value === 'number') {
    if (key.includes('_pct')) return value * 100;
    if (key === 'explosive_plays') return Math.min(value * 10, 100);
    return value;
  }
  return 0;
}

function OffenseRadar({ stats, team, color }: { stats: TeamUnitStats; team: string; color: string }) {
  const data = radarData.map((d) => ({
    metric: d.metric,
    value: getRadarValue(stats, d.key),
  }));

  return (
    <div className="w-64 h-64 mx-auto">
      <RadarChart cx="50%" cy="50%" outerRadius="70%" data={data}>
        <PolarGrid gridType="polygon" radialLines={false} stroke="rgba(0,0,0,0.1)" />
        <PolarAngleAxis dataKey="metric" tick={{ fontSize: 11, fill: '#666' }} />
        <PolarRadiusAxis angle={90} domain={[0, 100]} tick={false} />
        <Radar
          name={team}
          dataKey="value"
          stroke={color}
          fill={color}
          fillOpacity={0.15}
          strokeWidth={2}
          dot={false}
        />
        <Legend
          wrapperStyle={{ paddingTop: '10px' }}
          formatter={(): React.ReactNode => <span className="text-xs font-medium">{team}</span>}
        />
      </RadarChart>
    </div>
  );
}

function BarChartComparison({
  data,
  label,
}: {
  data: { name: string; home: number; away: number }[];
  label: string;
}) {
  return (
    <div className="w-full">
      <h4 className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">{label}</h4>
      <ResponsiveContainer width="100%" height={200}>
        <BarChart data={data} layout="vertical">
          <XAxis
            type="number"
            tick={{ fontSize: 11 }}
            axisLine={false}
            tickLine={false}
            domain={['auto', 'auto']}
          />
          <YAxis
            type="category"
            dataKey="name"
            width={60}
            tick={{ fontSize: 11, fontWeight: 500 }}
            axisLine={false}
            tickLine={false}
          />
          <Tooltip
            formatter={(value: number) => [value.toFixed(2), '']}
            labelFormatter={(name: string) => name}
            contentStyle={{ backgroundColor: '#fff', border: '1px solid #e5e5e5', borderRadius: '8px' }}
          />
          <Legend
            wrapperStyle={{ paddingTop: '10px' }}
            formatter={(value: string): React.ReactNode => (
              <span className="text-xs font-medium" style={{ color: value === 'Home' ? COLORS.home : COLORS.away }}>
                {value}
              </span>
            )}
          />
          <Bar dataKey="home" name="Home" fill={COLORS.home} radius={[0, 4, 4, 0]} maxBarSize={30} />
          <Bar dataKey="away" name="Away" fill={COLORS.away} radius={[4, 0, 0, 4]} maxBarSize={30} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}

export function TeamComparisonCharts({ stats, homeTeam, awayTeam }: TeamComparisonChartsProps) {
  const offenseData = [
    { name: 'EPA/Play', home: stats.home_offense.epa_per_play, away: stats.away_offense.epa_per_play },
    { name: 'Success Rate %', home: stats.home_offense.success_rate * 100, away: stats.away_offense.success_rate * 100 },
    { name: 'Pass EPA', home: stats.home_offense.pass_epa, away: stats.away_offense.pass_epa },
    { name: 'Rush EPA', home: stats.home_offense.rush_epa, away: stats.away_offense.rush_epa },
  ];

  const defenseData = [
    { name: 'EPA/Play Allowed', home: stats.home_defense.epa_per_play, away: stats.away_defense.epa_per_play },
    { name: 'Success Rate % Allowed', home: stats.home_defense.success_rate * 100, away: stats.away_defense.success_rate * 100 },
    { name: 'Pass EPA Allowed', home: stats.home_defense.pass_epa, away: stats.away_defense.pass_epa },
    { name: 'Rush EPA Allowed', home: stats.home_defense.rush_epa, away: stats.away_defense.rush_epa },
  ];

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-center gap-4 mb-4">
        <TeamLogo team={homeTeam} size="lg" />
        <span className="text-lg font-bold text-nfl-dark dark:text-nfl-light">{homeTeam}</span>
        <span className="text-gray-400">vs</span>
        <span className="text-lg font-bold text-nfl-light dark:text-nfl-light">{awayTeam}</span>
        <TeamLogo team={awayTeam} size="lg" />
      </div>

      <section className="space-y-6">
        <h3 className="text-lg font-semibold text-nfl-dark dark:text-nfl-light border-b border-gray-200 dark:border-gray-700 pb-2">
          Offense Comparison
        </h3>

        <div className="grid gap-6 md:grid-cols-2">
          <BarChartComparison data={offenseData} label="Key Offensive Metrics" />
          <OffenseRadar stats={stats.home_offense} team={homeTeam} color={COLORS.home} />
        </div>
      </section>

      <section className="space-y-6">
        <h3 className="text-lg font-semibold text-nfl-dark dark:text-nfl-light border-b border-gray-200 dark:border-gray-700 pb-2">
          Defense Comparison
        </h3>

        <div className="grid gap-6 md:grid-cols-2">
          <BarChartComparison data={defenseData} label="Key Defensive Metrics" />
          <OffenseRadar stats={stats.away_defense} team={awayTeam} color={COLORS.away} />
        </div>
      </section>

      <section className="space-y-6">
        <h3 className="text-lg font-semibold text-nfl-dark dark:text-nfl-light border-b border-gray-200 dark:border-gray-700 pb-2">
          Situational Performance
        </h3>

        <div className="grid gap-6 md:grid-cols-2">
          <OffenseRadar stats={stats.home_offense} team={`${homeTeam} Offense`} color={COLORS.home} />
          <OffenseRadar stats={stats.away_offense} team={`${awayTeam} Offense`} color={COLORS.away} />
        </div>
      </section>
    </div>
  );
}