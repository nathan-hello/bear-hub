This is a template for using Go (with go-chi and air), HTMX, and Tailwind with Bun being our package manager and builder.


Instructions:
1. `bun install` - https://bun.sh/docs/installation
2. `go install` - only dependency so far is air, and that's only for dev anyways
3. `air` - from the root dir
4. `run.sh` is written so that you can have a tailwindcss and air run in the same terminal


Supabase: 
1. Supabase is installed through bun and is already in package.json
2. `bunx supabase init`
3. `sudo $(which bunx) supabase start` - needs to run as sudo because it's making a docker container
