"use client"

import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from "recharts"
import { useGetTransactionHistory } from "@/features/transaction/hooks"
import { useMemo } from "react"

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#8884D8", "#FF6B6B", "#4ECDC4", "#95E1D3"]

export const ExpenseByCategoryChart = () => {
  const { data: transactionData, isLoading } = useGetTransactionHistory()

  const categoryData = useMemo(() => {
    if (!transactionData?.data) return []

    const categoryTotals = new Map<string, number>()

    transactionData.data
      .filter((t) => t.type === "Списание/Покупка")
      .forEach((transaction) => {
        const current = categoryTotals.get(transaction.kategoria) || 0
        categoryTotals.set(transaction.kategoria, current + transaction.amount)
      })

    return Array.from(categoryTotals.entries())
      .map(([name, value]) => ({ name, value }))
      .sort((a, b) => b.value - a.value)
  }, [transactionData])

  if (isLoading) {
    return (
      <div className="w-full h-[250px] md:h-[350px] flex items-center justify-center">
        <p className="text-sm">Загрузка...</p>
      </div>
    )
  }

  if (categoryData.length === 0) {
    return (
      <div className="w-full h-[250px] md:h-[350px]">
        <h3 className="text-base md:text-lg font-semibold mb-2 md:mb-4">Расходы по категориям</h3>
        <div className="flex items-center justify-center h-full">
          <p className="text-sm text-muted-foreground">Нет данных о расходах</p>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full h-[250px] md:h-[350px]">
      <h3 className="text-base md:text-lg font-semibold mb-2 md:mb-4">Расходы по категориям</h3>
      <ResponsiveContainer width="100%" height="100%">
        <PieChart>
          <Pie
            data={categoryData}
            cx="50%"
            cy="50%"
            labelLine={false}
            label={({ name, percent }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}
            outerRadius={60}
            fill="#8884d8"
            dataKey="value"
          >
            {categoryData.map((_, index) => (
              <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip />
          <Legend wrapperStyle={{ fontSize: '12px' }} />
        </PieChart>
      </ResponsiveContainer>
    </div>
  )
}
