import { APP_NAME } from "@/lib/api/constants";

export default function Footer() {
  return (
    <footer>
      <div className="card">
        <p className="font-semibold text-gray-800">{APP_NAME} &mdash; Secure digital voting for Nigeria</p>
        <p className="mt-1">
          Feature phone? SMS <strong>VOTE &lt;CODE&gt;</strong> to <strong>4040</strong>
        </p>
      </div>
    </footer>
  );
}
