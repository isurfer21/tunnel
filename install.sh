#!/usr/bin/bash env
# Copyright (c) 2019 Abhishek Kumar. All rights reserved. MIT license.

INSTALLATION_DIR="$HOME/.program/tunnel"
CURRENT_DIR="$PWD"

case $(uname -s) in
	Darwin) OS="darwin" ;;
	*) OS="linux" ;;
esac

case $(getconf LONG_BIT) in
	64) ARCH="amd64" ;;
	32) ARCH="386" ;;
esac

printf "Installing the Tunnel CLI at \n  $INSTALLATION_DIR \n"

if [ -d "$INSTALLATION_DIR" ]; then
	printf "Removing previous installation \n"
	rm -rf "$INSTALLATION_DIR"
fi
printf "Creating installation directory \n"
mkdir -p "$INSTALLATION_DIR"

if [ $# -eq 0 ]; then
	DOWNLOAD_URI_PATH=$(
		command curl -sSf https://github.com/isurfer21/tunnel/releases |
			command grep -o "/isurfer21/tunnel/releases/download/.*/tunnel_${OS}_${ARCH}" |
			command head -n 1
	)
	if [ ! "$DOWNLOAD_URI_PATH" ]; then exit 1; fi
	DOWNLOAD_LINK="https://github.com${DOWNLOAD_URI_PATH}"
else
	DOWNLOAD_LINK="https://github.com/isurfer21/tunnel/releases/download/${1}/tunnel_${OS}_${ARCH}"
fi
printf "Suitable build for ${OS} ${ARCH} \n  $DOWNLOAD_LINK \n"

cd "$INSTALLATION_DIR"
printf "Downloading the installation package \n"
curl -fL# $DOWNLOAD_LINK -o tunnel

printf "Allowing it to run as executable \n"
chmod +x tunnel

printf "Making it globally accessible \n"
printf "\nexport PATH=\"\$PATH:$INSTALLATION_DIR\"" >> "$HOME/.bash_profile"

printf "Tunnel was installed successfully \n"
if command -v tunnel >/dev/null; then
	printf "Run 'tunnel --help' to get started \n"
else
	printf "Manually add the directory to your \$HOME/.bash_profile (or similar) \n"
	printf "  export PATH=\"\$PATH:$INSTALLATION_DIR\" \n"
	printf "Run 'tunnel --help' to get started \n"
fi

printf "Done!\n"
exit 0