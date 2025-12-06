import type { Metadata } from "next";
import { AssistantChat } from "@/widgets/assistant-chat/ui/assistant-chat";

export const metadata: Metadata = {
  title: "AI Ассистент",
  description: "Умный помощник для управления финансами",
};

export default function AssistantPage() {
  return (
    <div className="flex flex-col h-full">
      <AssistantChat />
    </div>
  );
}
