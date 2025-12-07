"use client"

import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Legend, Line } from "recharts"
import { useGetTransactionHistory } from "@/features/transaction/hooks"
import { useMemo } from "react"

const MONTH_NAMES = ["Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"]

export const MonthlySpendingChart = () => {
  const { data: transactionData, isLoading } = useGetTransactionHistory()

  const monthlyExpenses = useMemo(() => {
    if (!transactionData?.data) return []

    const monthlyTotals = new Map<string, number>()

    transactionData.data
      .filter((t) => t.type === "Списание/Покупка")
      .forEach((transaction) => {
        const date = new Date(transaction.date)
        const monthKey = `${date.getFullYear()}-${date.getMonth()}`
        const current = monthlyTotals.get(monthKey) || 0
        monthlyTotals.set(monthKey, current + transaction.amount)
      })

    return Array.from(monthlyTotals.entries())
      .map(([key, amount]) => {
        const [year, month] = key.split("-").map(Number)
        return {
          month: `${MONTH_NAMES[month]} ${year}`,
          sortKey: new Date(year, month).getTime(),
          amount,
        }
      })
      .sort((a, b) => a.sortKey - b.sortKey)
      .slice(-6)
  }, [transactionData])

  if (isLoading) {
    return (
      <div className="w-full h-[300px] md:h-[350px] flex items-center justify-center">
        <p>Загрузка...</p>
      </div>
    )
  }

  if (monthlyExpenses.length === 0) {
    return (
      <div className="w-full h-[300px] md:h-[350px]">
        <h3 className="text-lg font-semibold mb-4">Ежемесячные расходы</h3>
        <div className="flex items-center justify-center h-full">
          <p className="text-muted-foreground">Нет данных о расходах</p>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Ежемесячные расходы</h3>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={monthlyExpenses}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="amount" stroke="#8884d8" strokeWidth={2} name="Расходы" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
