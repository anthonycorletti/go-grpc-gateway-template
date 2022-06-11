#!/bin/sh -ex

go generate
statik -src=./docs -include='*.jpg,*.txt,*.html,*.css,*.js,*.json'
