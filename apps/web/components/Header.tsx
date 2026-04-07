import { APP_NAME } from "@/lib/api/constants";
import Link from "next/link";

export default function Header() {
  return (
    <header className="border-b border-gray-200 bg-white shadow-sm">
      <div className="mx-auto flex max-w-4xl items-center justify-between px-4 py-4 sm:px-6">
        <Link href="/" className="text-xl font-bold tracking-tight text-brand-700">
          {APP_NAME}
        </Link>
        <nav className="flex gap-6 text-sm font-medium text-gray-600">
          <Link href="/" className="hover:text-brand-600 transition-colors">
            Candidates
          </Link>
          <Link href="/results" className="hover:text-brand-600 transition-colors">
            Results
          </Link>
          <Link href="/how-to-vote" className="hover:text-brand-600 transition-colors">
            How to vote
          </Link>
        </nav>
      </div>
    </header>
  );
}
