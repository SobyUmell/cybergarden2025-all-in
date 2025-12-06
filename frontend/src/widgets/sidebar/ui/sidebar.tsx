import Link from 'next/link'
import { Home, ArrowLeftRight, Settings, ChartPie, BotMessageSquare } from "lucide-react"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/shared/shadcn/ui/sidebar"

const items = [
  {
    title: "Главная",
    url: "/",
    icon: Home,
  },
  {
    title: "Аналитика",
    url: "/analytics",
    icon: ChartPie,
  },
  {
    title: "Транзакции",
    url: "/transactions",
    icon: ArrowLeftRight,
  },
  {
    title: "Ассистент",
    url: "/assistant",
    icon: BotMessageSquare,
  },
]

export function AppSidebar() {
  return (
    <Sidebar>
      <SidebarContent>
        <SidebarGroup className='space-y-8 p-4'>
          <SidebarGroupLabel className='text-3xl font-bold text-foreground'><span className='text-main-accent'>Cent</span>Keeper</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu className='gap-2'>
              {items.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton variant={"outline"} asChild>
                    <Link href={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  )
}
