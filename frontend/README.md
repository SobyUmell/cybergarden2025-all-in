# Frontend

Personal finance management web application built as a Telegram Mini App with Next.js 16 and React 19.

## Overview

A responsive financial tracking application that integrates with Telegram's ecosystem, providing users with transaction management, analytics, and AI-powered financial assistance.

## Features

### 📊 Analytics Dashboard
Real-time financial analytics with interactive charts:
- Expense breakdown by category (pie chart)
- Income vs. Expense comparison
- Monthly spending trends
- Balance trend tracking
- Savings rate visualization
- Spending comparison analysis

### 💰 Transaction Management
Complete transaction lifecycle management:
- View transaction history with sortable/filterable data table
- Add new transactions with form validation
- Edit existing transactions
- Delete transactions
- Category-based organization
- Transaction type classification (income/expense)

### 🤖 AI Assistant
Intelligent financial assistant for:
- Budget advice and insights
- Spending pattern analysis
- Financial planning guidance
- Natural language interaction

### 🔐 Authentication
- Telegram WebApp authentication
- Secure init data validation
- User session management

## Tech Stack

### Core Framework
- **Next.js 16** - React framework with App Router
- **React 19** - Latest React with improved rendering
- **TypeScript 5** - Type-safe development
- **Tailwind CSS 4** - Utility-first styling

### UI Components
- **Radix UI** - Accessible component primitives
  - Dialog, Dropdown Menu, Select, Tooltip, Label, etc.
- **shadcn/ui** - Pre-styled components built on Radix
- **Lucide React** - Icon system
- **Recharts 3** - Chart library for analytics

### State Management & Data Fetching
- **TanStack Query 5** - Server state management
- **React Hook Form 7** - Form handling
- **Zod 4** - Schema validation

### Telegram Integration
- **@telegram-apps/sdk-react** - Telegram Mini Apps SDK
- Custom Telegram authentication hook
- Init data validation

### Utilities
- **date-fns** - Date manipulation
- **clsx** / **tailwind-merge** - Conditional styling
- **class-variance-authority** - Component variants

## Project Structure

```
src/
├── app/                    # Next.js App Router pages
│   ├── page.tsx           # Analytics dashboard (home)
│   ├── transactions/      # Transaction management
│   │   ├── page.tsx       # Transaction list
│   │   ├── create/[id]/   # Create transaction
│   │   ├── edit/[id]/     # Edit transaction
│   │   └── delete/[id]/   # Delete transaction
│   └── assistant/         # AI assistant chat
├── core/
│   └── providers/         # App-wide providers (Query, Telegram)
├── entities/              # Domain models
│   ├── category/          # Category schema
│   ├── transaction/       # Transaction schema
│   └── user/              # User schema
├── features/              # Feature-specific logic
│   └── transaction/
│       └── create-transaction/
├── shared/                # Shared utilities
│   ├── api/               # API client with TMA auth
│   ├── lib/               # Utility functions
│   ├── shadcn/            # shadcn/ui components
│   └── ui/                # Custom UI components
│       └── data-table/    # Reusable data table
└── widgets/               # Complex UI widgets
    ├── assistant-chat/    # AI chat interface
    ├── charts/            # All chart components
    ├── sidebar/           # Navigation sidebar + mobile menu
    └── transaction-table/ # Transaction data table
```

### Architecture Pattern

The project follows **Feature-Sliced Design (FSD)** principles:
- **app/** - Application layer (pages, layouts)
- **widgets/** - Complex UI compositions
- **features/** - User interactions and business logic
- **entities/** - Business entities and models
- **shared/** - Reusable code without business logic

## API Integration

The frontend communicates with the backend via REST API:

**Base Endpoint:** `/webapp`

**Authentication:** Telegram init data passed via `Authorization: tma <data>` header

**Endpoints:**
- `GET /webapp/datahistory` - Fetch transaction history
- `GET /webapp/datametrics` - Fetch analytics metrics
- `POST /webapp/addt` - Create transaction
- `POST /webapp/updatet` - Update transaction
- `POST /webapp/deletet` - Delete transaction

## Development

Install dependencies:
```bash
pnpm install
```

Run development server:
```bash
pnpm dev
```

Build for production:
```bash
pnpm build
```

Start production server:
```bash
pnpm start
```

Open [http://localhost:3000](http://localhost:3000) to view the app.

## Responsive Design

The application is fully responsive with:
- Mobile-first design approach
- Container queries for component-level responsiveness
- Mobile drawer menu for navigation
- Adaptive chart sizing and layout
- Touch-optimized interactions

## Environment

The app is designed to run within Telegram as a Mini App, leveraging Telegram's native authentication and user context.
