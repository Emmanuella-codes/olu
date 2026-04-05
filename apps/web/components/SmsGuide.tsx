import { SMS_NUMBER } from "@/lib/api/constants";
import Link from "next/link";


export default function SmsGuide() {
    return (
        <div className="rounded-xl border border-gray-200 bg-gray-50 p-6 flex flex-col sm:flex-row items-start sm:items-center gap-4">
            <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0 text-xl">
                📱
            </div>
            <div className="flex-1 min-w-0">
                <h3 className="font-semibold text-gray-900">Using a feature phone?</h3>
                <p className="text-sm text-gray-500 mt-0.5">
                    You can still vote without internet. Send an SMS to{" "}
                    <strong className="font-mono text-gray-700">{SMS_NUMBER}</strong>.
                </p>
                <p className="text-sm font-mono bg-white border border-gray-200 rounded-lg px-3 py-2 mt-2 text-gray-700 inline-block">
                    VOTE A1 → {SMS_NUMBER}
                </p>
            </div>
            <Link
                href="/how-to-vote#sms"
                className="btn-secondary text-sm px-4 py-2 min-h-0 h-auto flex-shrink-0 whitespace-nowrap"
            >
                Learn more
            </Link>
        </div>
    )
}
