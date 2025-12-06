"use client"

import { RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar, Legend, ResponsiveContainer } from "recharts"

const data = [
  { category: "Питание", thisMonth: 45000, lastMonth: 38000 },
  { category: "Транспорт", thisMonth: 12000, lastMonth: 15000 },
  { category: "Развлечения", thisMonth: 18000, lastMonth: 20000 },
  { category: "Здоровье", thisMonth: 8000, lastMonth: 5000 },
  { category: "Покупки", thisMonth: 25000, lastMonth: 30000 },
]

export const SpendingComparisonChart = () => {
  return (
    <div className="w-full h-[300px] md:h-[350px]">
      <h3 className="text-lg font-semibold mb-4">Сравнение расходов</h3>
      <ResponsiveContainer width="100%" height="100%">
        <RadarChart cx="50%" cy="50%" outerRadius="80%" data={data}>
          <PolarGrid />
          <PolarAngleAxis dataKey="category" />
          <PolarRadiusAxis />
          <Radar name="Этот месяц" dataKey="thisMonth" stroke="#8884d8" fill="#8884d8" fillOpacity={0.6} />
          <Radar name="Прошлый месяц" dataKey="lastMonth" stroke="#82ca9d" fill="#82ca9d" fillOpacity={0.6} />
          <Legend />
        </RadarChart>
      </ResponsiveContainer>
    </div>
  )
}
