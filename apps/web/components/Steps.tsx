"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { castVote, getCandidates, sendOTP, verifyOTP } from "@/lib/api/api";
import { VOTE_PHONE_KEY, VOTE_TOKEN_KEY } from "@/lib/api/constants";
import { Candidate, Step } from "@/types/types";

import ConfirmationModal from "./ConfirmationModal";
import OtpInput from "./OtpInput";

const STEP_ORDER: Step[] = ["phone", "otp", "confirm", "done"];

function getErrorMessage(error: unknown, fallback: string) {
  return error instanceof Error && error.message ? error.message : fallback;
}

export default function Steps() {
  const searchParams = useSearchParams();
  const router = useRouter();

  const [step, setStep] = useState<Step>("phone");
  const [phone, setPhone] = useState("");
  const [otpCode, setOtpCode] = useState("");
  const [token, setToken] = useState("");
  const [selectedCode, setSelectedCode] = useState(searchParams.get("code") ?? "");
  const [candidates, setCandidates] = useState<Candidate[]>([]);
  const [confirmationId, setConfirmationId] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [countdown, setCountdown] = useState(0);

  useEffect(() => {
    getCandidates()
      .then(setCandidates)
      .catch(() => setCandidates([]));
  }, []);

  useEffect(() => {
    const savedPhone = window.sessionStorage.getItem(VOTE_PHONE_KEY);
    const savedToken = window.sessionStorage.getItem(VOTE_TOKEN_KEY);

    if (savedPhone) {
      setPhone(savedPhone);
    }

    if (savedToken) {
      setToken(savedToken);
    }
  }, []);

  useEffect(() => {
    if (countdown <= 0) {
      return;
    }

    const id = window.setTimeout(() => setCountdown((current) => current - 1), 1000);
    return () => window.clearTimeout(id);
  }, [countdown]);

  const selectedCandidate = useMemo(
    () => candidates.find((candidate) => candidate.code === selectedCode.toUpperCase()) ?? null,
    [candidates, selectedCode]
  );

  const currentStepIndex = STEP_ORDER.indexOf(step);

  const handleSendOTP = useCallback(async () => {
    setError("");

    const trimmedPhone = phone.trim();
    if (!trimmedPhone) {
      setError("Please enter your phone number.");
      return;
    }

    if (!selectedCode) {
      setError("Please select a candidate.");
      return;
    }

    setLoading(true);
    try {
      await sendOTP(trimmedPhone);
      window.sessionStorage.setItem(VOTE_PHONE_KEY, trimmedPhone);
      setPhone(trimmedPhone);
      setStep("otp");
      setCountdown(60);
    } catch (error) {
      setError(getErrorMessage(error, "Failed to send OTP. Please try again."));
    } finally {
      setLoading(false);
    }
  }, [phone, selectedCode]);

  const handleVerifyOTP = useCallback(async () => {
    setError("");

    if (otpCode.length !== 6) {
      setError("Please enter the 6-digit code.");
      return;
    }

    setLoading(true);
    try {
      const response = await verifyOTP(phone.trim(), otpCode);
      setToken(response.token);
      window.sessionStorage.setItem(VOTE_TOKEN_KEY, response.token);
      setStep("confirm");
    } catch (error) {
      setError(getErrorMessage(error, "Invalid code. Please try again."));
    } finally {
      setLoading(false);
    }
  }, [otpCode, phone]);

  const handleCastVote = useCallback(async () => {
    if (!selectedCode || !token) {
      return;
    }

    setError("");
    setLoading(true);
    try {
      const response = await castVote(selectedCode, token);
      setConfirmationId(response.confirmation_id);
      setStep("done");
      window.sessionStorage.removeItem(VOTE_TOKEN_KEY);
      window.sessionStorage.removeItem(VOTE_PHONE_KEY);
    } catch (error) {
      const message = getErrorMessage(error, "Failed to cast vote. Please try again.");
      setError(message.includes("already") ? "You have already cast your vote." : message);
    } finally {
      setLoading(false);
    }
  }, [selectedCode, token]);

  return (
    <div>
      <div className="mb-8 flex items-center gap-2">
        {STEP_ORDER.map((item, index) => {
          const isCurrent = step === item;
          const isComplete = STEP_ORDER.indexOf(item) < currentStepIndex;

          return (
            <div key={item} className="flex items-center gap-2">
              <div
                className={`flex h-8 w-8 items-center justify-center rounded-full text-sm font-semibold transition-colors ${
                  isCurrent
                    ? "bg-brand-500 text-white"
                    : isComplete
                      ? "bg-brand-100 text-brand-700"
                      : "bg-gray-100 text-gray-400"
                }`}
              >
                {index + 1}
              </div>
              {index < STEP_ORDER.length - 1 && <div className="h-px w-8 flex-1 bg-gray-200" />}
            </div>
          );
        })}
      </div>

      {error && <div className="card mb-6 border-red-200 bg-red-50 text-red-700">{error}</div>}

      {step === "phone" && (
        <div className="card space-y-5">
          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700">
              Your Nigerian phone number
            </label>
            <input
              className="input-field"
              type="tel"
              placeholder="e.g. 08012345678"
              value={phone}
              onChange={(event) => setPhone(event.target.value)}
              onKeyDown={(event) => event.key === "Enter" && handleSendOTP()}
              autoComplete="tel"
            />
            <p className="mt-1 text-xs text-gray-400">
              We&apos;ll send a one-time code to verify your identity.
            </p>
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-gray-700">Candidate code</label>
            <select
              className="input-field"
              value={selectedCode}
              onChange={(event) => setSelectedCode(event.target.value)}
            >
              <option value="">Select a candidate...</option>
              {candidates.map((candidate) => (
                <option key={candidate.id} value={candidate.code}>
                  {candidate.code} — {candidate.name} ({candidate.party})
                </option>
              ))}
            </select>
          </div>

          <button
            className="btn-primary w-full justify-center"
            onClick={handleSendOTP}
            disabled={loading || !phone.trim() || !selectedCode}
          >
            {loading ? "Sending..." : "Send verification code"}
          </button>
        </div>
      )}

      {step === "otp" && (
        <div className="card space-y-5">
          <p className="text-sm text-gray-600">
            Enter the 6-digit code sent to <strong>{phone}</strong>.
          </p>

          <OtpInput value={otpCode} onChange={setOtpCode} />

          <button
            className="btn-primary w-full justify-center"
            onClick={handleVerifyOTP}
            disabled={loading || otpCode.length !== 6}
          >
            {loading ? "Verifying..." : "Verify code"}
          </button>

          <div className="text-center text-sm text-gray-400">
            {countdown > 0 ? (
              <span>Resend in {countdown}s</span>
            ) : (
              <button
                className="min-h-0 text-brand-600 hover:underline"
                onClick={() => {
                  setOtpCode("");
                  void handleSendOTP();
                }}
              >
                Resend code
              </button>
            )}
          </div>
        </div>
      )}

      {step === "confirm" && selectedCandidate && (
        <ConfirmationModal
          candidate={selectedCandidate}
          onConfirm={handleCastVote}
          onCancel={() => router.push("/")}
          loading={loading}
        />
      )}

      {step === "done" && (
        <div className="card space-y-4 text-center">
          <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
            <svg className="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h2 className="text-xl font-bold text-gray-900">Vote confirmed!</h2>
          <p className="text-gray-500">Your vote has been recorded securely.</p>
          <p className="rounded-lg bg-gray-50 px-4 py-2 font-mono text-xs text-gray-400">
            Confirmation ID: {confirmationId.slice(0, 8).toUpperCase()}
          </p>
          <button className="btn-secondary" onClick={() => router.push("/result")}>
            See live results
          </button>
        </div>
      )}
    </div>
  );
}
