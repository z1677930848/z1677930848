#!/usr/bin/env bash

function build() {
	ROOT=$(dirname "$0")
	NAME="edge-api"
	DIST=$ROOT/"../dist/${NAME}"
	OS=${1}
	ARCH=${2}
	TAG=${3}
	NODE_ARCHITECTS=("amd64" "arm64")

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

	VERSION=$(lookup-version "$ROOT"/../internal/const/const.go)
	ZIP="${NAME}-${OS}-${ARCH}-${TAG}-v${VERSION}.zip"

	# build edge-node
	# 支持 EDGENODE_PATH 环境变量，默认使用 ../../EdgeNode
	if [ -z "${EDGENODE_PATH}" ]; then
		EDGENODE_PATH="$ROOT/../../EdgeNode"
	fi

	# 检查 EdgeNode 目录是否存在
	if [ ! -d "${EDGENODE_PATH}" ]; then
		echo "warning: EdgeNode directory not found at '${EDGENODE_PATH}', skipping edge-node build"
		NodeVersion="unknown"
	else
		NodeVersion=$(lookup-version "${EDGENODE_PATH}/internal/const/const.go")
		echo "building edge-node v${NodeVersion} ..."
		EDGE_NODE_BUILD_SCRIPT="${EDGENODE_PATH}/build/build.sh"
		if [ ! -f "$EDGE_NODE_BUILD_SCRIPT" ]; then
			echo "warning: edge-node build script not found at '${EDGE_NODE_BUILD_SCRIPT}', skipping edge-node build"
		else
			cd "${EDGENODE_PATH}/build" || exit
			echo "=============================="
			for arch in "${NODE_ARCHITECTS[@]}"; do
				if [ ! -f "${EDGENODE_PATH}/dist/edge-node-linux-${arch}-${TAG}-v${NodeVersion}.zip" ]; then
					./build.sh linux "$arch" $TAG
				else
					echo "use built node linux/$arch/v${NodeVersion}"
				fi
			done
			echo "=============================="
			cd - || exit
		fi
	fi

	rm -f "$ROOT"/deploy/*.zip
	# 复制 edge-node zip 文件（如果存在）
	for arch in "${NODE_ARCHITECTS[@]}"; do
		NODE_ZIP="${EDGENODE_PATH}/dist/edge-node-linux-${arch}-${TAG}-v${NodeVersion}.zip"
		if [ -f "$NODE_ZIP" ]; then
			cp "$NODE_ZIP" "$ROOT"/deploy/edge-node-linux-"${arch}"-v"${NodeVersion}".zip
		else
			echo "warning: edge-node zip not found at '$NODE_ZIP'"
		fi
	done

	# build edge-dns
	if [ "$TAG" = "plus" ]; then
		DNS_ROOT=$ROOT"/../../EdgeDNS"
		if [ -d "$DNS_ROOT"  ]; then
			DNSNodeVersion=$(lookup-version "$ROOT""/../../EdgeDNS/internal/const/const.go")
			echo "building edge-dns ${DNSNodeVersion} ..."
			EDGE_DNS_NODE_BUILD_SCRIPT=$ROOT"/../../EdgeDNS/build/build.sh"
			if [ ! -f "$EDGE_DNS_NODE_BUILD_SCRIPT" ]; then
				echo "unable to find edge-dns build script 'EdgeDNS/build/build.sh'"
				exit
			fi
			cd "$ROOT""/../../EdgeDNS/build" || exit
			echo "=============================="
			architects=("amd64" "arm64")
			for arch in "${architects[@]}"; do
				./build.sh linux "$arch" $TAG
			done
			echo "=============================="
			cd - || exit

			for arch in "${architects[@]}"; do
				cp "$ROOT""/../../EdgeDNS/dist/edge-dns-linux-${arch}-v${DNSNodeVersion}.zip" "$ROOT"/deploy/edge-dns-linux-"${arch}"-v"${DNSNodeVersion}".zip
			done
		fi
	fi

	# build sql
	# if [ $TAG = "plus" ]; then
	# 	echo "building sql ..."
	# 	"${ROOT}"/sql.sh
	# fi

	# copy files
	echo "copying ..."
	if [ ! -d "$DIST" ]; then
		mkdir "$DIST"
		mkdir "$DIST"/bin
		mkdir "$DIST"/configs
		mkdir "$DIST"/logs
		mkdir "$DIST"/data
	fi
	cp "$ROOT"/configs/api.template.yaml "$DIST"/configs/
	cp "$ROOT"/configs/db.template.yaml "$DIST"/configs/
	cp -R "$ROOT"/deploy "$DIST/"
	rm -f "$DIST"/deploy/.gitignore
	cp -R "$ROOT"/installers "$DIST"/

	# 确保 go.mod 是最新的
	cd "$ROOT"/.. || exit
	go mod tidy
	cd - || exit

	# building edge installer
	echo "building node installer ..."
	architects=("amd64" "arm64")
	for arch in "${architects[@]}"; do
		# TODO support arm, mips ...
		env GOOS=linux GOARCH="${arch}" CGO_ENABLED=0 go build -trimpath -tags $TAG --ldflags="-s -w" -o "$ROOT"/installers/edge-installer-helper-linux-"${arch}" "$ROOT"/../cmd/installer-helper/main.go
	done

	# building edge dns installer
	if [ $TAG = "plus" ]; then
		echo "building dns node installer ..."
		architects=("amd64" "arm64")
		for arch in "${architects[@]}"; do
			# TODO support arm, mips ...
			env GOOS=linux GOARCH="${arch}" CGO_ENABLED=0 go build -trimpath -tags $TAG --ldflags="-s -w" -o "$ROOT"/installers/edge-installer-dns-helper-linux-"${arch}" "$ROOT"/../cmd/installer-dns-helper/main.go
		done
	fi

	CC_PATH=""
	CXX_PATH=""
	if [[ "${OS}" == "linux" ]]; then
		if [ "${ARCH}" == "amd64" ]; then
			CC_PATH=$(command -v x86_64-linux-musl-gcc)
			CXX_PATH=$(command -v x86_64-linux-musl-g++)
		fi
		if [ "${ARCH}" == "arm64" ]; then
			CC_PATH=$(command -v aarch64-linux-musl-gcc)
			CXX_PATH=$(command -v aarch64-linux-musl-g++)
		fi
	fi
	# building api node
	# 如果找到了 musl 交叉编译工具链，使用 CGO 进行静态链接
	if [ -n "$CC_PATH" ] && [ -f "$CC_PATH" ]; then
		echo "using musl toolchain for static linking: $CC_PATH"
		env CC=$CC_PATH CXX=$CXX_PATH CGO_ENABLED=1 GOOS="$OS" GOARCH="$ARCH" go build -trimpath -tags $TAG --ldflags="-linkmode external -extldflags -static -s -w" -o "$DIST/bin/$NAME" "$ROOT"/../cmd/edge-api/main.go
	else
		# 如果没有找到 musl 工具链，使用 CGO_ENABLED=0 进行纯 Go 编译
		echo "musl toolchain not found, building without CGO"
		env GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=0 go build -trimpath -tags $TAG --ldflags="-s -w" -o "$DIST/bin/$NAME" "$ROOT"/../cmd/edge-api/main.go
	fi

	if [ ! -f "${DIST}/bin/${NAME}" ]; then
		echo "build failed!"
		exit
	fi

	# delete hidden files
	find "$DIST" -name ".DS_Store" -delete
	find "$DIST" -name ".gitignore" -delete

	echo "zip files"
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
