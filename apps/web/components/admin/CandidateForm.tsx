"use client";

import { FormEvent, useState } from "react";

import { CreateCandidatePayload } from "@/types/types";

interface Props {
  initial?: Partial<CreateCandidatePayload>;
  onSubmit: (candidate: CreateCandidatePayload) => void | Promise<void>;
  submitLabel?: string;
  loading?: boolean;
  error?: string;
}

const PARTY_OPTIONS = [
  { value: "apc",  label: "APC — All Progressives Congress" },
  { value: "pdp",  label: "PDP — Peoples Democratic Party" },
  { value: "lp",   label: "LP — Labour Party" },
  { value: "nnpp", label: "NNPP — New Nigeria Peoples Party" },
  { value: "apga", label: "APGA — All Progressives Grand Alliance" },
  { value: "adc",  label: "ADC — African Democratic Congress" },
  { value: "adp",  label: "ADP — Action Democratic Party" },
  { value: "apm",  label: "APM — Action Peoples Movement" },
  { value: "app",  label: "APP — Action Peoples Party" },
  { value: "a",    label: "A — Accord" },
  { value: "aac",  label: "AAC — African Action Congress" },
  { value: "aa",   label: "AA — Action Alliance" },
  { value: "bp",   label: "BP — Boot Party" },
  { value: "dla",  label: "DLA — Democratic Liberal Alliance" },
  { value: "nrm",  label: "NRM — National Rescue Movement" },
  { value: "prp",  label: "PRP — Peoples Redemption Party" },
  { value: "sdp",  label: "SDP — Social Democratic Party" },
  { value: "ypp",  label: "YPP — Young Progressives Party" },
  { value: "yp",   label: "YP — Youths Party" },
  { value: "zlp",  label: "ZLP — Zenith Labour Party" },
];

const emptyCandidate: CreateCandidatePayload = {
  code: "",
  name: "",
  party: "",
  bio: "",
  achievements: "",
  photo_url: "",
};

export default function CandidateForm({
  initial,
  onSubmit,
  submitLabel = "Save candidate",
  loading = false,
  error,
}: Props) {
  const [form, setForm] = useState<CreateCandidatePayload>({ ...emptyCandidate, ...initial });

  const handleChange = (field: keyof CreateCandidatePayload, value: string) => {
    setForm((current) => ({ ...current, [field]: value }));
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await onSubmit({
      ...form,
      code: form.code.trim().toUpperCase(),
      party: form.party.trim().toLowerCase(),
      name: form.name.trim(),
      bio: form.bio.trim(),
      achievements: form.achievements.trim(),
      photo_url: form.photo_url?.trim() || undefined,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      {error && <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">{error}</div>}

      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="code" className="mb-1 block text-sm font-medium text-gray-700">
            Candidate code
          </label>
          <input
            type="text"
            id="code"
            className="input-field"
            placeholder="A1"
            value={form.code}
            onChange={(event) => handleChange("code", event.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="party" className="mb-1 block text-sm font-medium text-gray-700">
            Political party
          </label>
          <select
            id="party"
            className="input-field"
            value={form.party}
            onChange={(event) => handleChange("party", event.target.value)}
            required
          >
            <option value="">Select a party...</option>
            {PARTY_OPTIONS.map(({ value, label }) => (
              <option key={value} value={value}>{label}</option>
            ))}
          </select>
        </div>
      </div>

      <div>
        <label htmlFor="name" className="mb-1 block text-sm font-medium text-gray-700">
          Full name
        </label>
        <input
          type="text"
          id="name"
          className="input-field"
          value={form.name}
          onChange={(event) => handleChange("name", event.target.value)}
          required
        />
      </div>

      <div>
        <label htmlFor="bio" className="mb-1 block text-sm font-medium text-gray-700">
          Bio
        </label>
        <textarea
          id="bio"
          className="input-field min-h-28"
          value={form.bio}
          onChange={(event) => handleChange("bio", event.target.value)}
          required
        />
      </div>

      <div>
        <label htmlFor="achievements" className="mb-1 block text-sm font-medium text-gray-700">
          Key achievements
        </label>
        <textarea
          id="achievements"
          className="input-field min-h-28"
          value={form.achievements}
          onChange={(event) => handleChange("achievements", event.target.value)}
          required
        />
      </div>

      <div>
        <label htmlFor="photo_url" className="mb-1 block text-sm font-medium text-gray-700">
          Photo URL
        </label>
        <input
          type="url"
          id="photo_url"
          className="input-field"
          value={form.photo_url ?? ""}
          onChange={(event) => handleChange("photo_url", event.target.value)}
        />
      </div>

      <button type="submit" className="btn-primary w-full justify-center" disabled={loading}>
        {loading ? "Saving..." : submitLabel}
      </button>
    </form>
  );
}
