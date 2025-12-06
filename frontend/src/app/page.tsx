import type { Metadata } from "next";
import {
  ExpenseByCategoryChart,
  IncomeExpenseChart,
  MonthlySpendingChart,
} from "@/widgets/charts";

export const metadata: Metadata = {
  title: "Аналитика",
  description: "Финансовая аналитика и отчеты",
};

export default function AnalyticsPage() {
  return (
    <div>
      <h1 className="text-2xl md:text-3xl font-bold mb-6">Аналитика</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-card rounded-lg border p-4 md:p-6">
          <ExpenseByCategoryChart />
        </div>

        <div className="bg-card rounded-lg border p-4 md:p-6">
          <IncomeExpenseChart />
        </div>

        <div className="md:col-span-2 bg-card rounded-lg border p-4 md:p-6">
          <MonthlySpendingChart />
        </div>
      </div>
    </div>
  );
}
