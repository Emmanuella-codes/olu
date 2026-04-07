import { Metadata } from "next";
import { notFound } from "next/navigation";

import CandidateProfileView from "@/components/CandidateProfileView";
import { getCandidate, getCandidates } from "@/lib/api/api";
import { APP_NAME } from "@/lib/api/constants";

export const revalidate = 300;

interface Props {
  params: Promise<{ candidate_id: string }>;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { candidate_id } = await params;

  try {
    const candidate = await getCandidate(candidate_id);
    return {
      title: `${candidate.name} - ${APP_NAME}`,
      description: candidate.bio.slice(0, 155),
    };
  } catch {
    return { title: `Candidate - ${APP_NAME}` };
  }
}

export async function generateStaticParams() {
  try {
    const candidates = await getCandidates();
    return candidates.map((candidate) => ({ candidate_id: candidate.id }));
  } catch {
    return [];
  }
}

export default async function CandidatePage({ params }: Props) {
  const { candidate_id } = await params;
  let candidate;

  try {
    candidate = await getCandidate(candidate_id);
  } catch {
    notFound();
  }

  return <CandidateProfileView candidate={candidate} />;
}
