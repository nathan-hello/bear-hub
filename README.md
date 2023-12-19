This is a template for using Go (with go-chi and air), HTMX, and Tailwind with Bun being our package manager and builder.


Install:
bun
1. `curl -fsSL https://bun.sh/install | bash`
air
2. `go install github.com/cosmtrek/air@latest`
sqlc
3. `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
Tailwindcss and Supabase
3. `bun install` - bun install

Bun Scripts
- `bun run dev` - this is defined in `package.json`. it starts tailwindcss and air

Supabase: 

- First, you want to create a Supabase database: https://supabase.com/dashboard/projects
- You are going to need the following variables:
    Tip: After you have written the .env file, you can `source .env` to put those variables in your shell!
    - `$SUPA_DB_PASSWORD`: This is the password you used when you created the project.
    - `$SUPA_ACCESS_TOKEN`: Go to https://supabase.com/dashboard/account/tokens and make a token. This will be used alongside your database password for authenticating. 
    - `$SUPA_PROJECT_REF`: This is the Reference ID for your project. You can find it in the dashboard Settings > General > Reference ID
    - `$SUPA_DB_HOST`: You can find this in the dashboard Settings > Database > Host
 
    The following it put into .env.example because they are defaulted values.
    - `$SUPA_DB_PORT`: This is 5432 by default. Find it at Settings > Database > Port > Port
    - `$SUPA_DB_NAME`: This is `postgres` by default. Find it at Settings > Database > Database name
    - `$SUPA_DB_USER`: This is `postgres` by default. Find it at Settings > Database > User
    - `$SUPA_DB_URI`: This is a combination of the user, password, host, port and database name. You can find it at Settings > Datbase > Connection String > URI. Again, .env.example has this already filled out.     

1. Supabase is installed through bun and is already in package.json
2. `bunx supabase init`
3. `bunx supabase login --token $SUPA_ACCESS_TOKEN` 
    - The --token flag is optional. If not passed in, it will open the browser for you to log in.
    - If you can't open a browser from the terminal, you can pass `--no-browser` so it will give you a link that you can paste into a browser and log in elsewhere. See https://supabase.com/docs/reference/cli/supabase-login
4. `bunx supabase link --project-ref $SUPA_PROJECT_REF -p $SUPA_DB_PASSWORD`

