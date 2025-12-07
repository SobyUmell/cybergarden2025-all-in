"use client"

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"
import { useGetTransactionHistory } from "@/features/transaction/hooks"
import { useMemo } from "react"

const MONTH_NAMES = ["Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"]

export const BalanceTrendChart = () => {
  const { data: transactionData, isLoading } = useGetTransactionHistory()

  const balanceData = useMemo(() => {
    if (!transactionData?.data) return []

    const sortedTransactions = [...transactionData.data].sort((a, b) => a.date - b.date)
    
    const monthlyBalances = new Map<string, { balance: number; sortKey: number }>()
    let runningBalance = 0

    sortedTransactions.forEach((transaction) => {
      if (transaction.type === "Пополнение") {
        runningBalance += transaction.amount
      } else if (transaction.type === "Списание/Покупка") {
        runningBalance -= transaction.amount
      }

      const date = new Date(transaction.date)
      const monthKey = `${date.getFullYear()}-${date.getMonth()}`
      monthlyBalances.set(monthKey, {
        balance: runningBalance,
        sortKey: new Date(date.getFullYear(), date.getMonth()).getTime(),
      })
    })

    return Array.from(monthlyBalances.entries())
      .map(([key, data]) => {
        const [year, month] = key.split("-").map(Number)
        return {
          month: `${MONTH_NAMES[month]} ${year}`,
          sortKey: data.sortKey,
          balance: data.balance,
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

  if (balanceData.length === 0) {
    return (
      <div className="w-full h-[300px] md:h-[350px]">
        <h3 className="text-lg font-semibold mb-4">Динамика баланса</h3>
        <div className="flex items-center justify-center h-full">
          <p className="text-muted-foreground">Нет данных о балансе</p>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Динамика баланса</h3>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={balanceData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="balance" stroke="#8884d8" strokeWidth={2} name="Баланс" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
