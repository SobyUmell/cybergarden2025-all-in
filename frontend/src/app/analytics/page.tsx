import type { Metadata } from "next";
import {
  ExpenseByCategoryChart,
  IncomeExpenseChart,
  BalanceTrendChart,
  MonthlySpendingChart,
  SavingsRateChart,
  SpendingComparisonChart,
} from "@/widgets/charts";

export const metadata: Metadata = {
  title: "Аналитика",
  description: "Финансовая аналитика и отчеты",
};

export default function AnalyticsPage() {
  return (
    <div className="container mx-auto p-4 md:p-6 lg:p-8">
      <h1 className="text-2xl md:text-3xl font-bold mb-6">Аналитика</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">
        <div className="bg-card rounded-lg border p-4 md:p-6">
          <ExpenseByCategoryChart />
        </div>
        
        <div className="bg-card rounded-lg border p-4 md:p-6">
          <IncomeExpenseChart />
        </div>
        
        <div className="bg-card rounded-lg border p-4 md:p-6">
          <BalanceTrendChart />
        </div>
        
        <div className="bg-card rounded-lg border p-4 md:p-6">
          <MonthlySpendingChart />
        </div>
        
        <div className="bg-card rounded-lg border p-4 md:p-6">
          <SavingsRateChart />
        </div>
        
        <div className="bg-card rounded-lg border p-4 md:p-6">
          <SpendingComparisonChart />
        </div>
      </div>
    </div>
  );
}
