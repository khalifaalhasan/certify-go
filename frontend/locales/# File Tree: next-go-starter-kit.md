# File Tree: next-go-starter-kit

**Generated:** 3/23/2026, 12:38:50 AM
**Root Path:** `/home/khalifaalhasan/Projects/next-go-starter-kit`

```
├── backend
│   ├── cmd
│   │   └── api
│   │       └── main.go
│   ├── config
│   │   ├── config.development.yaml
│   │   ├── config.template.yaml
│   │   └── reflex.conf
│   ├── container
│   │   └── container.go
│   ├── internal
│   │   ├── auth
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── rbac
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── shared
│   │       └── models
│   │           ├── helpers.go
│   │           ├── rbac.go
│   │           └── user.go
│   ├── migrations
│   │   ├── seeders
│   │   │   └── rbac_seeder.sql
│   │   ├── 001_create_users_table.down.sql
│   │   ├── 001_create_users_table.up.sql
│   │   ├── 002_create_rbac_tables.down.sql
│   │   ├── 002_create_rbac_tables.up.sql
│   │   ├── 003_add_deleted_at_to_users.down.sql
│   │   └── 003_add_deleted_at_to_users.up.sql
│   ├── pkg
│   │   ├── config
│   │   │   ├── app.go
│   │   │   ├── auth.go
│   │   │   ├── cache.go
│   │   │   ├── config.go
│   │   │   └── database.go
│   │   ├── constants
│   │   │   └── common.go
│   │   ├── database
│   │   │   └── database.go
│   │   ├── middleware
│   │   │   ├── cors.go
│   │   │   ├── jwt.go
│   │   │   ├── limiter.go
│   │   │   └── rbac.go
│   │   ├── response
│   │   │   └── response.go
│   │   ├── router
│   │   │   └── router.go
│   │   ├── server
│   │   │   └── server.go
│   │   ├── token
│   │   │   └── jwt.go
│   │   ├── utils
│   │   │   ├── oauth.go
│   │   │   └── password.go
│   │   └── validator
│   │       └── validator.go
│   ├── tests
│   │   └── .gitkeep
│   ├── .env.example
│   ├── .gitignore
│   ├── Dockerfile
│   ├── Dockerfile.dev
│   ├── Makefile
│   ├── README.md
│   ├── docker-compose.dev.yaml
│   ├── docker-compose.yaml
│   ├── go.mod
│   └── go.sum
├── frontend
│   ├── app
│   │   ├── admin
│   │   │   ├── permissions
│   │   │   │   └── page.tsx
│   │   │   ├── roles
│   │   │   │   └── page.tsx
│   │   │   ├── users
│   │   │   │   └── page.tsx
│   │   │   └── layout.tsx
│   │   ├── auth
│   │   │   ├── login
│   │   │   │   └── page.tsx
│   │   │   ├── register
│   │   │   │   └── page.tsx
│   │   │   └── layout.tsx
│   │   ├── docs
│   │   │   └── page.tsx
│   │   ├── home
│   │   │   └── page.tsx
│   │   ├── settings
│   │   │   ├── preference
│   │   │   │   └── page.tsx
│   │   │   └── profile
│   │   │       └── page.tsx
│   │   ├── favicon.ico
│   │   ├── globals.css
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   └── providers.tsx
│   ├── locales
│   │   ├── en.json
│   │   └── id.json
│   ├── public
│   │   ├── file.svg
│   │   ├── globe.svg
│   │   ├── gns.png
│   │   ├── logo_icon.png
│   │   ├── next.svg
│   │   ├── vercel.svg
│   │   └── window.svg
│   ├── src
│   │   ├── application
│   │   │   └── hooks
│   │   │       ├── use-auth.ts
│   │   │       ├── use-mobile.tsx
│   │   │       ├── use-permission.ts
│   │   │       └── use-users-query.ts
│   │   ├── domain
│   │   │   └── services
│   │   │       ├── api-client.ts
│   │   │       ├── auth.service.ts
│   │   │       ├── rbac.service.ts
│   │   │       └── user.service.ts
│   │   ├── infrastructure
│   │   │   └── stores
│   │   │       ├── auth-store.ts
│   │   │       └── locale-store.ts
│   │   ├── lib
│   │   │   ├── query-client.ts
│   │   │   └── utils.ts
│   │   └── presentation
│   │       ├── components
│   │       │   ├── admin
│   │       │   │   └── permissions-container.tsx
│   │       │   ├── forms
│   │       │   │   ├── change-password-form.tsx
│   │       │   │   └── profile-info-form.tsx
│   │       │   ├── layout
│   │       │   │   ├── main-layout.tsx
│   │       │   │   ├── navbar.tsx
│   │       │   │   ├── protected-feature.tsx
│   │       │   │   ├── protected-route.tsx
│   │       │   │   ├── sidebar-item.tsx
│   │       │   │   └── sidebar.tsx
│   │       │   ├── pages
│   │       │   │   ├── docs-page.tsx
│   │       │   │   ├── home-page.tsx
│   │       │   │   ├── landing-page.tsx
│   │       │   │   ├── preference-settings-page.tsx
│   │       │   │   └── profile-settings-page.tsx
│   │       │   ├── ui
│   │       │   │   ├── accordion.tsx
│   │       │   │   ├── alert-dialog.tsx
│   │       │   │   ├── alert.tsx
│   │       │   │   ├── aspect-ratio.tsx
│   │       │   │   ├── avatar.tsx
│   │       │   │   ├── badge.tsx
│   │       │   │   ├── breadcrumb.tsx
│   │       │   │   ├── button-group.tsx
│   │       │   │   ├── button.tsx
│   │       │   │   ├── calendar.tsx
│   │       │   │   ├── card.tsx
│   │       │   │   ├── carousel.tsx
│   │       │   │   ├── chart.tsx
│   │       │   │   ├── checkbox.tsx
│   │       │   │   ├── code.tsx
│   │       │   │   ├── collapsible.tsx
│   │       │   │   ├── command.tsx
│   │       │   │   ├── context-menu.tsx
│   │       │   │   ├── dialog.tsx
│   │       │   │   ├── drawer.tsx
│   │       │   │   ├── dropdown-menu.tsx
│   │       │   │   ├── empty.tsx
│   │       │   │   ├── field.tsx
│   │       │   │   ├── form.tsx
│   │       │   │   ├── hover-card.tsx
│   │       │   │   ├── input-group.tsx
│   │       │   │   ├── input-otp.tsx
│   │       │   │   ├── input.tsx
│   │       │   │   ├── item.tsx
│   │       │   │   ├── kbd.tsx
│   │       │   │   ├── label.tsx
│   │       │   │   ├── menubar.tsx
│   │       │   │   ├── navigation-menu.tsx
│   │       │   │   ├── pagination.tsx
│   │       │   │   ├── popover.tsx
│   │       │   │   ├── progress.tsx
│   │       │   │   ├── radio-group.tsx
│   │       │   │   ├── resizable.tsx
│   │       │   │   ├── scroll-area.tsx
│   │       │   │   ├── select.tsx
│   │       │   │   ├── separator.tsx
│   │       │   │   ├── sheet.tsx
│   │       │   │   ├── sidebar.tsx
│   │       │   │   ├── skeleton.tsx
│   │       │   │   ├── slider.tsx
│   │       │   │   ├── sonner.tsx
│   │       │   │   ├── spinner.tsx
│   │       │   │   ├── switch.tsx
│   │       │   │   ├── table.tsx
│   │       │   │   ├── tabs.tsx
│   │       │   │   ├── textarea.tsx
│   │       │   │   ├── toggle-group.tsx
│   │       │   │   ├── toggle.tsx
│   │       │   │   └── tooltip.tsx
│   │       │   ├── theme-provider.tsx
│   │       │   └── theme-toggle.tsx
│   │       └── hooks
│   │           └── use-locale.ts
│   ├── .env.example
│   ├── .gitignore
│   ├── README.md
│   ├── bun.lock
│   ├── components.json
│   ├── eslint.config.mjs
│   ├── i18n.ts
│   ├── middleware.ts
│   ├── next-env.d.ts
│   ├── next.config.ts
│   ├── package-lock.json
│   ├── package.json
│   ├── postcss.config.mjs
│   └── tsconfig.json
└── README.md
```

---
*Generated by FileTree Pro Extension*