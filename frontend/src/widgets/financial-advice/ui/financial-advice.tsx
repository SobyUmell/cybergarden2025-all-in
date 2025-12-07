"use client"

import { useState } from "react"
import { useMutation } from "@tanstack/react-query"
import { Lightbulb, Loader2 } from "lucide-react"
import { Button } from "@/shared/shadcn/ui/button"
import { transactionApi } from "@/features/transaction/api"

export const FinancialAdvice = () => {
  const [advice, setAdvice] = useState<string | null>(null)

  const mutation = useMutation({
    mutationFn: transactionApi.getAdvice,
    onSuccess: (data) => {
      setAdvice(data.advice)
    },
  })

  const handleGetAdvice = () => {
    mutation.mutate()
  }

  return (
    <div className="mt-6 space-y-4">
      <Button
        onClick={handleGetAdvice}
        disabled={mutation.isPending}
        size="lg"
        className="w-full md:w-auto"
      >
        {mutation.isPending ? (
          <>
            <Loader2 className="w-5 h-5 mr-2 animate-spin" />
            Анализирую...
          </>
        ) : (
          <>
            <Lightbulb className="w-5 h-5 mr-2" />
            Получить финансовые рекомендации
          </>
        )}
      </Button>

      {advice && (
        <div className="bg-card rounded-lg border p-6 shadow-sm">
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-10 h-10 rounded-full bg-main-accent flex items-center justify-center">
              <Lightbulb className="w-5 h-5 text-primary-foreground" />
            </div>
            <div className="flex-1 space-y-2">
              <h3 className="font-semibold text-lg">Финансовые рекомендации</h3>
              <div className="text-sm leading-relaxed whitespace-pre-wrap text-muted-foreground">
                {advice}
              </div>
            </div>
          </div>
        </div>
      )}

      {mutation.isError && (
        <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
          <p className="text-sm text-destructive">
            Не удалось получить рекомендации. Попробуйте позже.
          </p>
        </div>
      )}
    </div>
  )
}
