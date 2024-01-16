This is the `no-magic-stack`. Build web applications with as little magic as possible. 

This is a contribution to the philosophy that abstractions should not hand-wave away extremely complicated software techniques. 
This is a stack where instead of learning about how someone else configured a library, you learn fundamental technologies.
This is a kind of template - not a proper one but maybe we'll get there.
This is not a library, and most certainly not another CLI tool. 
This is my way of learning some really cool pieces of software.

This is the stack:
- htmx        - `https://htmx.org/docs/`
- sqlc        - `https://github.com/sqlc-dev/sqlc/cmd/sqlc`
- tmpl        - `https://github.com/a-h/templ`
- Tailwindcss - `https://tailwindcss.com/blog/standalone-cli`
- air         - `https://github.com/cosmtrek/air`
- Postgresql  - `https://www.postgresql.org/download/`


Goals: 
1. Minimal abstractions: Server performance is important. Abstractions should compliment (or develop!) someone's skills in standard technologies (sql, css, http, etc), not get in the way of them.
2. Make the DX akin to developing on NextJS, Astro, or any of the other great JS frameworks. 
    - In other words: spend as little time as possible learning about how a library functions and instead focus on the fundamentals.
    - sqlc is a great choice for this. It makes you write actual, real-deal SQL but gives a great DX when you're inside your .go files.
4. Script the things that are annoying. This should include installing all of the relevant packages, and setting up the stack.
5. Examples: So far this includes a todo app and live chat app using websockets.
    - Including auth! Auth is very important, and I've seen too many tutorials and walkthroughs skip over it. Likely because it's actually kind of difficult to implement with confidence.
    - Testing: A cool thing about HTMX and sqlc is that they are highly testable. Testing in web applications might be overkill, but it should be included in any examples. 
    - Cool things that are a discussion for later: SEO, SSG & SSR, Serverless, Docker containers for development and production, an optional build step for optimizing for production (re: file sizes).
        - templ has stuff in their docs regarding this. Making this work and deployable (in self-hosted, monolithic cloud, and serverless) should be a priority.

Environment Variables:
`$DB_URI`: This is a combination of the user, password, host, port and database name. You can find it with the Supabase web interface at Settings > Datbase > Connection String > URI.
- We need this to get the schema of the remote database with `pg_dump $DB_URI -s -f src/sqlc/input/remote-schema.sqlc`
`JWT_SECRET`: The secret key to signing the JWTs.

Misc Notes:
- Turn off browser caching on your browser for your `localhost:3000` or whatever port you use. I've had issues with the browser caching tailwind's css file, not updating when tailwind updates said file.
- To run tests, use `cd tests/ && go test -v ./...`. This is because the tests are in the tests/ folder.

Problems to be ironed out:
- templ [claims to support hot reloading in the browser](https://templ.guide/commands-and-tools/hot-reload). Whenever I try to run `templ generate --proxy="localhost:3000"`, the program crashes. Maybe there's something else I should do. Even so, it's not entirely clear if the hot reloading extends to .go files. For now, we'll use air and refresh like cavemen.
  - Even using air, there is a compile time with templ or air (probably templ), where when you press refresh, the server is down until templ can compile, and air can compile that. Is there a browser setting to adjust the timeout period? Is there a misconfig somewhere? Maybe this is just the cost of using a compiled langauge, because air does compile to `tmp/` before throwing back up the server. Even so, why does it terminate? Does it know that it's not caught up yet?

VSCode
- For VSCode users, install the relevant extensions. Making a `.vscode/` folder isn't necessary because the setup is pretty easy. You can make them recommended extensions by putting the below json into `.vscode/extensions.json`, then search `@recommended` in the Extensions tab. Or just search for them. They're not hard to find. 

extensions.json
```json
{ "recommendations": [
    "golang.go",
    "a-h.templ",
    "tamasfe.even-better-toml",
    "bradlc.vscode-tailwindcss",
    "mads-hartmann.bash-ide-vscode",
]}
```

settings.json
```json
{ "gopls": { "ui.semanticTokens": true } }
```
Major TODOs:
1. Provide proper IDE integration for working with Postgresql. Sqlc provides excellent typesafety and static analysis between your schema and your queries, but the LSP integration for Postgresql is virtually non-existent from what I understand. I think there are solutions that work well if you literally connect to your live database - but in my opinion that's not the right answer. Surely this is a relatively solved problem right? At least for neovim, it is not. Supabase is working on [a proper Postgres LSP](https://github.com/supabase/postgres_lsp/), but as of January 2024, it's in the early development stages (although definitely being worked on, as evidenced by [this blog post from just a month ago](https://supabase.com/blog/postgres-language-server-implementing-parser))
