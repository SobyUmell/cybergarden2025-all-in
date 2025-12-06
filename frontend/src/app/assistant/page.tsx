import type { Metadata } from "next";
import Link from "next/link";
import { Button } from "@/shared/shadcn/ui/button";

export const metadata: Metadata = {
  title: "AI Ассистент",
  description: "Умный помощник для управления финансами",
};

export default function AssistantPage() {
  return (
    <div>
      <h1 className="text-2xl md:text-3xl font-bold mb-6">Ассистент</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">
      </div>
    </div>
  );
}
