import type { Metadata } from "next";
import Link from "next/link";
import { Button } from "@/shared/shadcn/ui/button";

export const metadata: Metadata = {
  title: "AI Ассистент",
  description: "Умный помощник для управления финансами",
};

export default function AssistantPage() {
  return (
    <div className="flex min-h-screen flex-col p-8">
      <div className="mb-8 flex items-center justify-between">
        <div className="flex flex-col gap-2">
          <h1 className="text-3xl font-bold">AI Ассистент</h1>
          <p className="text-muted-foreground">
            Умный помощник для управления финансами
          </p>
        </div>
        <Link href="/">
          <Button variant="outline">На главную</Button>
        </Link>
      </div>

      <div className="flex flex-col gap-4">
        <p className="text-muted-foreground">
          Здесь будет интерфейс AI ассистента
        </p>
      </div>
    </div>
  );
}
