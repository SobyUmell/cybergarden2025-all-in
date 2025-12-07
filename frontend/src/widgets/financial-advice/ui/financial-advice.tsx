"use client"

import { useState } from "react"
import { useMutation } from "@tanstack/react-query"
import { Lightbulb, Loader2, AlertCircle } from "lucide-react"
import { Button } from "@/shared/shadcn/ui/button"
import { transactionApi } from "@/features/transaction/api"

export const FinancialAdvice = () => {
  const [advice, setAdvice] = useState<string | null>(null)

  const mutation = useMutation({
    mutationFn: transactionApi.getAdvice,
    onSuccess: (data) => {
      setAdvice(data.advice)
    },
    retry: false,
  })

  const handleGetAdvice = () => {
    setAdvice(null)
    mutation.mutate()
  }

  const isTimeout = mutation.isError && mutation.error instanceof Error &&
    (mutation.error.message.includes("timeout") || mutation.error.message.includes("aborted"))

  return (
    <div className="mt-6 space-y-4">
      <Button
        onClick={handleGetAdvice}
        disabled={mutation.isPending}
        size="lg"
        className="w-full md:w-auto gap-1"
      >
        {mutation.isPending ? (
          <>
            <Loader2 className="w-5 h-5 mr-2 animate-spin" />
            Анализирую...
          </>
        ) : (
          <>
            <Lightbulb className="w-5 h-5 mr-2" />
            Получить рекомендации
          </>
        )}
      </Button>

      {advice && (
        <div className="bg-card rounded-lg border p-4 md:p-6 shadow-sm">
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-8 h-8 md:w-10 md:h-10 rounded-full bg-main-accent flex items-center justify-center">
              <Lightbulb className="w-4 h-4 md:w-5 md:h-5 text-primary-foreground" />
            </div>
            <div className="flex-1 space-y-2">
              <h3 className="font-semibold text-base md:text-lg">Финансовые рекомендации</h3>
              <div className="text-sm leading-relaxed whitespace-pre-wrap text-muted-foreground">
                {advice}
              </div>
            </div>
          </div>
        </div>
      )}

      {mutation.isError && (
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
          <div className="flex items-start gap-2">
            <AlertCircle className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-sm font-medium text-destructive">
                {isTimeout ? "Время ожидания истекло" : "Ошибка при получении рекомендаций"}
              </p>
              <p className="text-xs text-destructive/80 mt-1">
                {isTimeout
                  ? "Сервер не успел обработать запрос. Попробуйте позже."
                  : "Не удалось получить рекомендации. Попробуйте позже."}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
