#!/bin/bash

install_sqlc() {
    echo "Installing sqlc: https://github.com/sqlc-dev/sqlc/cmd/sqlc"
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
}

install_templ() {
    echo -e "templ will be cloned into /tmp/templ, use `go install`, then /tmp/templ will be removed."
    mkdir -p /tmp/templ
    git clone https://github.com/a-h/templ.git /tmp/templ
    cd /tmp/templ
    go run ./get-version > .version
    cd cmd/templ && go install
    rm -rf /tmp/templ
}

install_air() {
    echo "Installing air: https://github.com/cosmtrek/air"
    go install github.com/cosmtrek/air@latest
}

install_tailwindcss() {
    echo -e "\nhttps://github.com/tailwindlabs/tailwindcss/releases \nReleases are in the link above. The files without extensions are binaries. Place whichever is relevant to your OS/platform in a PATH folder (e.g. /bin/)"
    echo -e "Alternatively, you can run `wget -O ~/.local/bin/tailwindcss <url of file> \& chmod +x ~/.local/bin/tailwindcss`"
}

#install_postgresql_ubuntu_debian
#install_sqlc
install_templ
#install_air
#install_tailwindcss

