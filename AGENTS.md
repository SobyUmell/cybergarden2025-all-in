# Agent Guidelines

## Commands
- **Frontend (Next.js)**: Located in `frontend/`
  - Dev: `cd frontend && pnpm dev`
  - Build: `cd frontend && pnpm build`
  - No lint/test commands configured
- **Backend (Go)**: Located in `backend/` - check for go.mod and standard Go commands
- **ML**: Located in `ml/` - check for requirements.txt or other Python configs

## Code Style
- **TypeScript**: Strict mode enabled, no explicit any types
- **Imports**: Use `@/` alias for imports (maps to `src/`), e.g., `import { cn } from "@/lib/utils"`
- **Components**: Use shadcn/ui (New York style) stored in `@/shared/shadcn/ui/`
- **Styling**: Tailwind CSS v4 with `cn()` utility for class merging
- **Types**: Use React.ComponentProps<"element"> pattern, VariantProps for variants
- **Naming**: camelCase for functions/variables, PascalCase for components
- **State**: React Query (@tanstack/react-query) for server state, React Hook Form + Zod for forms
- **Function components**: Use standard function declarations, not arrow functions for exports
- **Error handling**: No specific patterns observed - implement appropriate try/catch
