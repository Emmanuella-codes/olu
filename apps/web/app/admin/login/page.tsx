"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { adminLogin, setAdminToken } from "@/lib/admin-api/api";

function getErrorMessage(error: unknown, fallback: string) {
  return error instanceof Error && error.message ? error.message : fallback;
}

export default function AdminLoginPage() {
  const router = useRouter();
  const [error, setError] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError("");

    if (!email.trim() || !password) {
      setError("Please enter your email and password.");
      return;
    }

    setLoading(true);
    try {
      const { token } = await adminLogin(email.trim(), password);
      setAdminToken(token);
      router.push("/admin");
    } catch (error) {
      setError(getErrorMessage(error, "Login failed."));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto flex min-h-[70vh] max-w-md items-center">
      <div className="card w-full space-y-6">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-brand-600">Admin</p>
          <h1 className="mt-2 text-2xl font-bold text-gray-950">Sign in to Olu</h1>
          <p className="mt-1 text-sm text-gray-500">Manage candidates, results, and election settings.</p>
        </div>

        {error && <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="email" className="mb-1 block text-sm font-medium text-gray-700">
              Email
            </label>
            <input
              type="email"
              id="email"
              placeholder="admin@olu.com"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              className="input-field"
              autoComplete="email"
            />
          </div>

          <div>
            <label htmlFor="password" className="mb-1 block text-sm font-medium text-gray-700">
              Password
            </label>
            <input
              type="password"
              id="password"
              placeholder="Enter your password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="input-field"
              autoComplete="current-password"
            />
          </div>

          <button className="btn-primary w-full justify-center" type="submit" disabled={loading}>
            {loading ? "Signing in..." : "Sign in"}
          </button>
        </form>
      </div>
    </div>
  );
}
