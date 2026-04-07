"use client";

import { useCallback, useEffect, useRef, useState } from "react";

import { getResults } from "@/lib/api/api";
import { Results, TallyRow } from "@/types/types";

import ResultsChart from "./ResultsChart";

const REFRESH_INTERVAL = 30_000;

export default function ResultsView() {
  const [results, setResults] = useState<Results | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const [pollingPaused, setPollingPaused] = useState(false);
  const failureCount = useRef(0);

  const fetchResults = useCallback(async () => {
    try {
      const data = await getResults();
      setResults(data);
      setLastUpdated(new Date());
      setError(false);
      failureCount.current = 0;
      setPollingPaused(false);
    } catch {
      failureCount.current += 1;
      setError(true);
      if (failureCount.current >= 3) {
        setPollingPaused(true);
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchResults();
    const id = window.setInterval(() => {
      if (failureCount.current < 3) {
        fetchResults();
      }
    }, REFRESH_INTERVAL);

    return () => window.clearInterval(id);
  }, [fetchResults]);

  const leader: TallyRow | null = results?.tally?.[0] ?? null;

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Live results</h1>
          {lastUpdated && (
            <p className="mt-0.5 text-xs text-gray-400">
              Updated {lastUpdated.toLocaleTimeString("en-NG")} &mdash; refreshes every 30s
            </p>
          )}
        </div>
        <button onClick={fetchResults} className="btn-secondary h-auto min-h-0 px-4 py-2 text-sm">
          Refresh
        </button>
      </div>

      {loading && <div className="card py-16 text-center text-gray-400">Loading results...</div>}

      {error && !loading && (
        <div className="card border-red-200 bg-red-50 text-red-700">
          Could not load results.{" "}
          {pollingPaused
            ? "Auto-refresh paused after repeated failures. "
            : "Please try again. "}
          <button onClick={fetchResults} className="underline font-medium">
            Retry now
          </button>
        </div>
      )}

      {results && !loading && (
        <>
          <div className="mb-6 grid grid-cols-2 gap-4 sm:grid-cols-3">
            <div className="card text-center">
              <p className="text-3xl font-bold text-brand-600">
                {results.total_votes.toLocaleString("en-NG")}
              </p>
              <p className="mt-1 text-sm text-gray-400">Total votes</p>
            </div>
            <div className="card text-center">
              <p className="text-3xl font-bold text-brand-600">{results.tally.length}</p>
              <p className="mt-1 text-sm text-gray-400">Candidates</p>
            </div>
            {leader && (
              <div className="card col-span-2 text-center sm:col-span-1">
                <p className="truncate text-lg font-bold text-gray-900">{leader.name}</p>
                <p className="text-sm font-medium text-brand-500">{leader.party.toLocaleUpperCase()}</p>
                <p className="mt-1 text-xs text-gray-400">Currently leading</p>
              </div>
            )}
          </div>

          <ResultsChart tally={results.tally} totalVotes={results.total_votes} />

          <div className="card mt-6 overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-100">
                  <th className="py-2 pr-4 text-left font-semibold text-gray-500">Rank</th>
                  <th className="py-2 pr-4 text-left font-semibold text-gray-500">Candidate</th>
                  <th className="py-2 pr-4 text-left font-semibold text-gray-500">Party</th>
                  <th className="py-2 text-right font-semibold text-gray-500">Votes</th>
                  <th className="py-2 pl-4 text-right font-semibold text-gray-500">Share</th>
                </tr>
              </thead>
              <tbody>
                {results.tally.map((row, index) => {
                  const percentage =
                    results.total_votes > 0
                      ? ((row.vote_count / results.total_votes) * 100).toFixed(1)
                      : "0.0";

                  return (
                    <tr
                      key={row.candidate_id}
                      className={`border-b border-gray-50 ${index === 0 ? "bg-brand-50" : ""}`}
                    >
                      <td className="py-3 pr-4 font-mono text-gray-400">{index + 1}</td>
                      <td className="py-3 pr-4 font-medium text-gray-900">
                        {row.name}
                        <span className="ml-2 rounded-full bg-gray-100 px-2 py-0.5 font-mono text-xs text-gray-500">
                          {row.code}
                        </span>
                      </td>
                      <td className="py-3 pr-4 text-gray-500">{row.party.toLocaleUpperCase()}</td>
                      <td className="py-3 text-right font-semibold text-gray-900">
                        {row.vote_count.toLocaleString("en-NG")}
                      </td>
                      <td className="py-3 pl-4 text-right font-medium text-brand-600">
                        {percentage}%
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}
