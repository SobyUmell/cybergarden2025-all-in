import { cn } from "@/shared/shadcn/lib/utils";
import { SidebarTrigger } from "@/shared/shadcn/ui/sidebar";

interface MobileMenuProps {
  className?: string;
}

export const MobileMenu = ({ className }: MobileMenuProps) => {
  return (
    <div className={cn("sticky top-0 z-20 flex bg-main-accent md:hidden justify-between gap-4 items-center w-full p-4", className)}>
      <h1 className="text-3xl font-bold text-foreground"><span className='text-background'>Cent</span>Keeper</h1>
      <SidebarTrigger variant={"secondary"} className="flex md:hidden text-main-accent text-2xl size-9" />
    </div>
  )
}

