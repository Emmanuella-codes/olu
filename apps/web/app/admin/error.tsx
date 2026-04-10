"use client";

export default function AdminError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4 text-center">
      <h2 className="text-xl font-bold text-gray-950">Something went wrong</h2>
      <p className="max-w-sm text-sm text-gray-500">{error.message || "An unexpected error occurred."}</p>
      <button onClick={reset} className="btn-primary">
        Try again
      </button>
    </div>
  );
}
