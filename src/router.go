package src

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/nathan-hello/htmx-template/src/routes"
)

func HandleSites() {

	type Site struct {
		route       string
		hfunc       http.HandlerFunc
		middlewares alice.Chain
	}

	sites := []Site{
		{route: "/",
			hfunc: routes.Root,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET"),
				RejectSubroute("/"),
			)},
		{route: "/todo",
			hfunc: routes.Todo,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET", "DELETE", "POST"),
			)},
		{route: "/signup",
			hfunc: routes.SignUp,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET", "POST"),
			)},
		{route: "/signin",
			hfunc: routes.SignIn,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET", "POST"),
			)},
		{route: "/profile",
			hfunc: routes.UserProfile,
			middlewares: alice.New(
				Logging,
				InjectClaimsOnValidToken,
				AllowMethods("GET"),
			)},
		{route: "/c",
			hfunc: routes.MicroComponents,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET"),
			)},
	}

	for _, v := range sites {
		http.Handle(v.route, v.middlewares.ThenFunc(v.hfunc))
	}

	http.ListenAndServe(":3000", nil)
}
