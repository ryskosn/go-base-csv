BINARY_NAME=basecsv
BINARY_WINDOWS=${BINARY_NAME}.exe
DEPLOY_PATH=/Users/USERNAME/Dropbox/foo/bar/base_csv_download

build:
	go build -o ${BINARY_NAME} -v

buildw:
	env GOOS=windows GOARCH=amd64 go build -o ${BINARY_WINDOWS} -v

rename:
	mv ${DEPLOY_PATH}/${BINARY_WINDOWS} ${DEPLOY_PATH}/_${BINARY_WINDOWS}

deploy: rename
	cp ./${BINARY_WINDOWS} ${DEPLOY_PATH}
