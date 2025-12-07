"use client"

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"
import { useGetTransactionHistory } from "@/features/transaction/hooks"
import { useMemo } from "react"

const MONTH_NAMES = ["Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"]

export const IncomeExpenseChart = () => {
  const { data: transactionData, isLoading } = useGetTransactionHistory()

  const monthlyData = useMemo(() => {
    if (!transactionData?.data) return []

    const monthlyTotals = new Map<string, { income: number; expense: number }>()

    transactionData.data.forEach((transaction) => {
      const date = new Date(transaction.date)
      const monthKey = `${date.getFullYear()}-${date.getMonth()}`
      const monthLabel = `${MONTH_NAMES[date.getMonth()]} ${date.getFullYear()}`

      const current = monthlyTotals.get(monthKey) || { income: 0, expense: 0 }
      
      if (transaction.type === "Пополнение") {
        current.income += transaction.amount
      } else if (transaction.type === "Списание/Покупка") {
        current.expense += transaction.amount
      }

      monthlyTotals.set(monthKey, current)
    })

    return Array.from(monthlyTotals.entries())
      .map(([key, totals]) => {
        const [year, month] = key.split("-").map(Number)
        return {
          month: `${MONTH_NAMES[month]} ${year}`,
          sortKey: new Date(year, month).getTime(),
          income: totals.income,
          expense: totals.expense,
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

  if (monthlyData.length === 0) {
    return (
      <div className="w-full h-[300px] md:h-[350px]">
        <h3 className="text-lg font-semibold mb-4">Доходы и расходы</h3>
        <div className="flex items-center justify-center h-full">
          <p className="text-muted-foreground">Нет данных о транзакциях</p>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Доходы и расходы</h3>
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={monthlyData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey="income" fill="#82ca9d" name="Доходы" />
          <Bar dataKey="expense" fill="#ff8042" name="Расходы" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
