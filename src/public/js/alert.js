window.onload = function () {
    let potentials = [
        { param: "unauthorized", component: "/c/alert/unauthorized" },
        { param: "signedout", component: "/c/alert/signout" },
        { param: "404", component: "/c/alert/404" },
        { param: "500", component: "/c/alert/500" },
    ];
    const params = new URLSearchParams(window.location.search);

    for (let key of potentials) {
        if (params.has(key.param)) {
            window.history.replaceState(null, "", params.get(key));
            htmx.ajax("GET", key.component, "#notification-target");
            break;
        }
    }
};
