import Link from "next/link";

export default function NotFound() {
  return (
    <div className="card py-16 text-center">
      <p className="text-5xl font-bold text-brand-500">404</p>
      <h1 className="mt-4 text-xl font-semibold text-gray-900">Page not found</h1>
      <p className="mt-2 text-gray-400">
        The page you&apos;re looking for doesn&apos;t exist.
      </p>
      <Link href="/" className="btn-primary mt-8 inline-flex">
        Back to candidates
      </Link>
    </div>
  );
}
