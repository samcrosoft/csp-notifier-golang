appPort := "3000"
appHost := "localhost"
appEndpoint := "http://" + appHost + ":" + appPort
dockerRepo := "samcrosoft/csp-notifier-golang"
dockerRepoTag := "v1"
dockerRepoWithTag := dockerRepo + ":" + dockerRepoTag
appBuildName := "collector"

default:
    @just --list --unsorted

# This will build the app as a linux binary
build_app build_os="linux" build_arch="amd64":
    env GOOS={{build_os}} \
        GOARCH={{build_arch}} \
        go build -o build/collector main.go    

# This will build the application into docker image
build_docker:
    #!/usr/bin/env bash
    docker build \
        -t {{dockerRepoWithTag}} \
        -f Dockerfile \
        .

# This will push the docker image to the docker hub
push_docker_image:
    #!/usr/bin/env bash
    docker image push {{dockerRepoWithTag}}

#  This will be used to test the endpoint
test_violation:
    #!/usr/bin/env bash
    TEST_DATE_TIME=$(date  +"%Y-%m-%d-%H-%M-%S");
    TEST_URL="https://example.com/at/${TEST_DATE_TIME}.html"
    TEST_REFERRER="https://sample-referrer.com"
    TEST_STATUS_CODE=$(( ( $RANDOM % 10 + 40)  + 1 ))
    TEST_DATA='{
          "csp-report": {
              "document-uri": $td_url,
              "referrer": $td_ref,
              "violated-directive": "script-src-elem",
              "effective-directive": "script-src-elem",
              "original-policy": "default-src '\''self'\''; img-src '\''self'\'' https://*.ytimg.com; script-src-elem '\''self'\'' https://storage.googleapis.com https://www.youtube.com; connect-src '\''self'\'' https://www.googleapis.com; frame-src '\''self'\'' https://www.youtube.com; base-uri '\''self'\''; frame-ancestors '\''none'\''; form-action '\''self'\''; block-all-mixed-content; report-uri https://reporting.example.com/;",
              "disposition": "report",
              "blocked-uri": "https://www.youtube.com/iframe_api",
              "line-number": 1,
              "column-number": 7982,
              "source-file": "://example.com/static/js/7.74a7cce6.chunk.js",
              "status-code": 0,
              "script-sample": ""
          }
    }';
    # generate sample payload
    JSON_TEST_DATA=$( jq -n \
                      --arg td_url "${TEST_URL}" \
                      --arg td_ref "${TEST_REFERRER}" \
                      "${TEST_DATA}")
    curl \
      --request POST \
      -H 'Content-Type: application/csp-report' \
      -d "${JSON_TEST_DATA}" \
      {{appEndpoint}}/report