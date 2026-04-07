"use client";

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="card border-red-200 bg-red-50 py-12 text-center text-red-700">
      <p className="font-medium">Something went wrong.</p>
      {error.message && (
        <p className="mt-1 text-sm text-red-500">{error.message}</p>
      )}
      <button onClick={reset} className="btn-secondary mt-6">
        Try again
      </button>
    </div>
  );
}
