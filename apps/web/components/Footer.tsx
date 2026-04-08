import { APP_NAME } from "@/lib/api/constants";

export default function Footer() {
  return (
    <footer className="border-t border-gray-200 bg-white">
      <div className="mx-auto max-w-4xl px-4 py-6 sm:px-6">
        <p className="font-semibold text-gray-800">{APP_NAME} &mdash; Secure digital voting for Nigeria</p>
        <p className="mt-1 text-sm text-gray-500">
          Feature phone? SMS <strong className="text-gray-700">VOTE &lt;CODE&gt;</strong> to{" "}
          <strong className="text-gray-700">****</strong>
        </p>
      </div>
    </footer>
  );
}
