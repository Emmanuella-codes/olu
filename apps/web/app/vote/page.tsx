import Steps from "@/components/Steps";

export default function VotePage() {
    return (
        <div className="">
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Cast your vote</h1>
            <p className="text-gray-500 mb-8">
                Your identity is verified by OTP. Each phone number can only vote once.
            </p>
            <Steps />
        </div>
    );
}
