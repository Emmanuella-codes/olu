import { faqs, smsSteps, webSteps } from "@/lib/data/steps";
import Link from "next/link";

export default function HowToVotePage() {
    return (
        <div className="max-w-2xl mx-auto">
            <h1 className="text-2xl font-bold text-gray-900 mb-2">How to vote</h1>
            <p className="text-gray-500 mb-10">
                Choose the method that works best for you — smartphone or feature phone.
            </p>
        
            {/* Web voting */}
            <section className="mb-12">
                <div className="flex items-center gap-3 mb-6">
                <div className="w-10 h-10 rounded-full bg-brand-500 text-white flex items-center justify-center font-bold shrink-0">
                    W
                </div>
                <h2 className="text-xl font-semibold text-gray-900">
                    Voting online (smartphone)
                </h2>
                </div>
        
                <ol className="space-y-4">
                {webSteps.map(({ step, title, body }) => (
                    <li key={step} className="flex gap-4">
                    <span className="w-8 h-8 rounded-full border-2 border-brand-300 text-brand-600 font-semibold text-sm flex items-center justify-center shrink-0 mt-0.5">
                        {step}
                    </span>
                    <div>
                        <h3 className="font-semibold text-gray-900">{title}</h3>
                        <p className="text-gray-500 text-sm mt-0.5">{body}</p>
                    </div>
                    </li>
                ))}
                </ol>
        
                <Link href="/vote" className="btn-primary mt-6 inline-flex">
                Vote now
                </Link>
            </section>
        
            {/* SMS voting */}
            <section className="mb-12">
                <div className="flex items-center gap-3 mb-6">
                <div className="w-10 h-10 rounded-full bg-gray-700 text-white flex items-center justify-center font-bold shrink-0">
                    S
                </div>
                <h2 className="text-xl font-semibold text-gray-900">
                    Voting by SMS (any phone)
                </h2>
                </div>
        
                <ol className="space-y-4">
                {smsSteps.map(({ step, title, body }) => (
                    <li key={step} className="flex gap-4">
                    <span className="w-8 h-8 rounded-full border-2 border-gray-300 text-gray-600 font-semibold text-sm flex items-center justify-center shrink-0 mt-0.5">
                        {step}
                    </span>
                    <div>
                        <h3 className="font-semibold text-gray-900">{title}</h3>
                        <p className="text-gray-500 text-sm mt-0.5">{body}</p>
                    </div>
                    </li>
                ))}
                </ol>
        
                <div className="mt-6 card bg-gray-50 border-gray-200">
                <p className="text-sm font-mono text-gray-700 text-center">
                    VOTE A1 → send to <strong>****</strong>
                </p>
                </div>
            </section>
        
            {/* FAQ */}
            <section>
                <h2 className="text-xl font-semibold text-gray-900 mb-4">
                Frequently asked questions
                </h2>
                <div className="space-y-4">
                {faqs.map(({ q, a }) => (
                    <details
                    key={q}
                    className="card cursor-pointer group"
                    >
                    <summary className="font-medium text-gray-900 list-none flex justify-between items-center">
                        {q}
                        <span className="text-gray-400 group-open:rotate-180 transition-transform ml-4 shrink-0">
                        ▾
                        </span>
                    </summary>
                    <p className="mt-3 text-gray-500 text-sm">{a}</p>
                    </details>
                ))}
                </div>
            </section>
        </div>
    );
}
