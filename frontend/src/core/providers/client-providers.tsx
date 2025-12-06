"use client";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "@/shared/lib";
import { SidebarProvider, SidebarTrigger } from "@/shared/shadcn/ui/sidebar";
import { AppSidebar } from "@/widgets/sidebar/ui/sidebar";

export const ClientProviders = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <SidebarProvider>
          <AppSidebar />
          <SidebarTrigger className="hidden sm:flex" />
          {children}
        </SidebarProvider>
      </QueryClientProvider>
    </>
  );
};
