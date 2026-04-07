import { APP_NAME } from "@/lib/api/constants";
import Link from "next/link";

export default function Header() {
  return (
    <header>
      <div className="mx-auto flex max-w-4xl items-center justify-between px-4 py-4">
        <Link href="/" className="text-xl font-bold tracking-tight text-white">
          {APP_NAME}
        </Link>
        <nav className="flex gap-6 text-sm font-medium">
          <Link href="/" className="hover:text-brand-100">
            Candidates
          </Link>
          <Link href="/result" className="hover:text-brand-100">
            Results
          </Link>
          <Link href="/how-to-vote" className="hover:text-brand-100">
            How to vote
          </Link>
        </nav>
      </div>
    </header>
  );
}
