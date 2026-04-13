export const webSteps = [
    {
        step: "1",
        title: "Browse candidates",
        body: "Read the bio and achievements of each candidate on the home page. Click a candidate card to see the full profile.",
    },
    {
        step: "2",
        title: "Enter your phone number",
        body: "Go to the Vote page and enter your Nigerian mobile number. We send a one-time verification code — your number is never stored in plain text.",
    },
    {
        step: "3",
        title: "Verify with OTP",
        body: "Enter the 6-digit code sent to your phone. The code expires in 10 minutes. You can request a resend after 60 seconds.",
    },
    {
        step: "4",
        title: "Select and confirm",
        body: "Choose your candidate and confirm your selection. Once submitted, your vote is final. Each phone number can only vote once.",
    },
    {
        step: "5",
        title: "Save your confirmation ID",
        body: "After voting you'll receive a unique confirmation ID. Keep it safe — it proves your vote was counted without revealing who you voted for.",
    },
];

export const smsSteps = [
    {
      step: "1",
      title: "Find your candidate's code",
      body: 'Each candidate has a short code (e.g. A1, B2). You can find the codes at idibo.ng or ask a registration officer.',
    },
    {
      step: "2",
      title: "Compose your SMS",
      body: 'Type: VOTE followed by a space and the candidate code. Example: VOTE A1',
    },
    {
      step: "3",
      title: "Send to ****",
      body: "Send your message to ****. No rates Applied!. You will receive a confirmation message once your vote is recorded.",
    },
];

export const faqs = [
    {
        q: "Can I change my vote after submitting?",
        a: "No. Once submitted, votes are final and cannot be changed or deleted. Please review your selection carefully before confirming.",
    },
    {
        q: "How do I know my vote was counted?",
        a: "After voting you receive a unique confirmation ID. You can cross-reference this with the public audit log after the election closes.",
    },
    {
        q: "Is my vote anonymous?",
        a: "Yes. We store a cryptographic hash of your phone number, never the number itself. Your vote is linked to a hash, not to your identity.",
    },
    {
        q: "What if I don't have a smartphone?",
        a: "Send VOTE <CODE> to **** from any mobile phone, no rates applied. The system works on all Nigerian networks.",
    },
    {
        q: "What if my SMS is not confirmed?",
        a: "If you don't receive a confirmation SMS within 10 minutes, your message may have been incorrectly formatted. Check the format (VOTE A1) and try again. If you have already voted, you'll be notified.",
    },
];
