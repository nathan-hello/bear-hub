This is a template for using Go (with go-chi and air), HTMX, and Tailwind with Bun being our package manager and builder.


Instructions:
0. `curl -fsSL https://bun.sh/install | bash` - install bun
1. `bun install` - bun install
2. `go install github.com/cosmtrek/air@latest` - is this possible to be put into go.mod? it's just a dev dependency
3. `bun run dev` - this is defined in `package.json`. it starts tailwindcss and air

Supabase: 
1. Supabase is installed through bun and is already in package.json
2. `bunx supabase init`
3. `sudo $(which bunx) supabase start` - needs to run as sudo because it's making a docker container
