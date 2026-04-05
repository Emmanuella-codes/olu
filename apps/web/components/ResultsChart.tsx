import { TallyRow } from "@/types/types";

interface Props {
    tally: TallyRow[];
    totalVotes: number;
}

const COLORS = [
    "bg-brand-500",
    "bg-blue-500",
    "bg-purple-500",
    "bg-amber-500",
    "bg-pink-500",
    "bg-teal-500",
];

export default function ResultsChart({ tally, totalVotes }: Props) {
    if (!tally || tally.length === 0) {
        return (
          <div className="card text-center text-gray-400 py-10">
            No votes recorded yet.
          </div>
        );
    }
     
    const max = Math.max(...tally.map((r) => r.vote_count), 1);
      
    return (
        <div className="card space-y-4">
            <h2 className="text-sm font-semibold uppercase tracking-wider text-gray-400">
                Vote distribution
            </h2>
            {tally.map((row, i) => {
                const pct = totalVotes > 0 ? (row.vote_count / totalVotes) * 100 : 0;
                const barWidth = (row.vote_count / max) * 100;
                const color = COLORS[i % COLORS.length];
    
                return (
                    <div key={row.candidate_id}>
                        <div className="flex items-center justify-between mb-1">
                            <div className="flex items-center gap-2 min-w-0">
                                <span
                                    className={`w-2.5 h-2.5 rounded-full shrink-0 ${color}`}
                                />
                                <span className="text-sm font-medium text-gray-800 truncate">
                                    {row.name}
                                </span>
                                <span className="text-xs text-gray-400 font-mono shrink-0">
                                {   row.code}
                                </span>
                            </div>
                            <div className="flex items-center gap-3 shrink-0 ml-2">
                                <span className="text-sm font-semibold text-gray-900">
                                {row.vote_count.toLocaleString("en-NG")}
                                </span>
                                <span className="text-xs text-gray-400 w-10 text-right">
                                {pct.toFixed(1)}%
                                </span>
                            </div>
                        </div>
                        <div className="h-2.5 bg-gray-100 rounded-full overflow-hidden">
                            <div
                                className={`h-full rounded-full transition-all duration-700 ${color}`}
                                style={{ width: `${barWidth}%` }}
                            />
                        </div>
                    </div>
                );
             })}
        </div>
    );
}
