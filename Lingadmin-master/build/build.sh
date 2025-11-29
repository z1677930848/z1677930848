#!/usr/bin/env bash

function build() {
	ROOT=$(dirname "$0")
	JS_ROOT=$ROOT/../web/public/js
	NAME="lingcdnadmin"
	DIST=$ROOT/"../dist/${NAME}"
	OS=${1}
	ARCH=${2}
	TAG=${3}

	if [ -z "$OS" ]; then
		echo "usage: build.sh OS ARCH"
		exit
	fi
	if [ -z "$ARCH" ]; then
		echo "usage: build.sh OS ARCH"
		exit
	fi
	if [ -z "$TAG" ]; then
		TAG="community"
	fi

	# checking environment
	echo "checking required commands ..."
	commands=("zip" "unzip" "go" "find" "sed")
	for cmd in "${commands[@]}"; do
		if [ "$(which "${cmd}")" ]; then
			echo "checking ${cmd}: ok"
		else
			echo "checking ${cmd}: not found"
			return
		fi
	done

	VERSION=$(lookup-version "$ROOT"/../internal/const/const.go)
	ZIP="${NAME}-${OS}-${ARCH}-${TAG}-v${VERSION}.zip"

	# build edge-api (optional)
	# If EDGEAPI_PATH is set, use it; otherwise default to ../../EdgeAPI
	if [ -z "${EDGEAPI_PATH}" ]; then
		EDGEAPI_PATH="$ROOT"/../../EdgeAPI
	fi

	if [ -d "${EDGEAPI_PATH}" ]; then
		APINodeVersion=$(lookup-version "${EDGEAPI_PATH}"/internal/const/const.go)
		echo "building edge-api v${APINodeVersion} ..."
		EDGE_API_BUILD_SCRIPT="${EDGEAPI_PATH}/build/build.sh"
		if [ ! -f "${EDGE_API_BUILD_SCRIPT}" ]; then
			echo "warning: edge-api build script not found at '${EDGE_API_BUILD_SCRIPT}', skipping edge-api build"
		else
			cd "${EDGEAPI_PATH}/build" || exit
			echo "=============================="
			./build.sh "$OS" "$ARCH" $TAG
			echo "=============================="
			cd - || exit
		fi
	else
		echo "warning: EdgeAPI directory not found at '${EDGEAPI_PATH}', skipping edge-api build"
		APINodeVersion="unknown"
	fi

    # generate files (optional, may fail in CI environment)
	echo "generating files ..."
	go run -tags $TAG "$ROOT"/../cmd/lingcdnadmin/main.go generate || echo "warning: generate failed, using existing components.src.js"

	# prefer npm-based build if package.json exists in web/
	if [ -f "$ROOT"/../web/package.json ] && [ "$(which npm)" ]; then
		echo "building web assets with npm (terser)..."
		npm --prefix "$ROOT"/../web run build
	elif [ "$(which uglifyjs)" ]; then
		echo "compress to component.js ..."
		uglifyjs --compress --mangle -- "${JS_ROOT}"/components.src.js > "${JS_ROOT}"/components.js
	else
		echo "copy to component.js ..."
		cp "${JS_ROOT}"/components.src.js "${JS_ROOT}"/components.js
	fi

	# create dir & copy files
	echo "copying ..."
	if [ ! -d "$DIST" ]; then
		mkdir "$DIST"
		mkdir "$DIST"/bin
		mkdir "$DIST"/configs
		mkdir "$DIST"/logs
	fi

	cp -R "$ROOT"/../web "$DIST"/
	rm -f "$DIST"/web/tmp/*
	rm -rf "$DIST"/web/public/js/components
	rm -f "$DIST"/web/public/js/components.src.js
	cp "$ROOT"/configs/server.template.yaml "$DIST"/configs/

	EDGE_API_ZIP_FILE="${EDGEAPI_PATH}/dist/edge-api-${OS}-${ARCH}-${TAG}-v${APINodeVersion}.zip"
	if [ -f "$EDGE_API_ZIP_FILE" ]; then
		cp "$EDGE_API_ZIP_FILE" "$DIST"/
		cd "$DIST"/ || exit
		unzip -q "$(basename "$EDGE_API_ZIP_FILE")"
		rm -f "$(basename "$EDGE_API_ZIP_FILE")"
		cd - || exit
	else
		echo "warning: EdgeAPI zip not found at '$EDGE_API_ZIP_FILE', skipping bundling"
	fi

	# build
	echo "building ${NAME} ..."
	env GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=1 go build -trimpath -tags $TAG -ldflags="-s -w" -o "$DIST"/bin/${NAME} "$ROOT"/../cmd/lingcdnadmin/main.go

	# delete hidden files
	find "$DIST" -name ".DS_Store" -delete
	find "$DIST" -name ".gitignore" -delete
	find "$DIST" -name "*.less" -delete
	find "$DIST" -name "*.css.map" -delete
	find "$DIST" -name "*.js.map" -delete

	# zip
	echo "zip files ..."
	cd "${DIST}/../" || exit
	if [ -f "${ZIP}" ]; then
		rm -f "${ZIP}"
	fi
	zip -r -X -q "${ZIP}" ${NAME}/
	rm -rf ${NAME}
	cd - || exit

	echo "[done]"
}

function lookup-version() {
	FILE=$1
	VERSION_DATA=$(cat "$FILE")
	re="Version[ ]+=[ ]+\"([0-9.]+)\""
	if [[ $VERSION_DATA =~ $re ]]; then
		VERSION=${BASH_REMATCH[1]}
		echo "$VERSION"
	else
		echo "could not match version"
		exit
	fi
}

build "$1" "$2" "$3"
