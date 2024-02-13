3. auth!
    - potentially two versions: one using supabase primatives and one using as much non-db specific apis (sql)
    - use email, github oauth, implement JWTs? again this might be something supabase can handle with their db api
    - email magic links, phone number support, oauth, and most importantly, the ability to mix/match these login methods to one account
        - a limitation could be that the alternate login method can't already have an account on it, OR we say that we can merge the two db entries
    - multiple layers of perms: `admin`, `user`

4. testing?
    - it would be really cool to have a test where you give it data, and test it against what it's supposed to be
    - kind of ui testing, but because htmx uses http to grab data and render it, it makes sense
    - this will only really make sense for something more complicated than todo, and testing html output sounds not fun so this is just a game theory
