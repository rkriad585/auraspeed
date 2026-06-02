cmake_minimum_required(VERSION 3.20)
project(auraspeed
  VERSION 3.1.5
  DESCRIPTION "Cross-platform terminal tool for network diagnostics"
  LANGUAGES NONE
)

set(AURASPEED_VERSION "${PROJECT_VERSION}")

if(CMAKE_GENERATOR_PLATFORM MATCHES "Win32|x86|ARM64")
  set(GOARCH "amd64")
elseif(CMAKE_GENERATOR_PLATFORM MATCHES "ARM64")
  set(GOARCH "arm64")
else()
  set(GOARCH "amd64")
endif()

find_program(GO_EXECUTABLE go REQUIRED)

add_custom_target(build
  COMMAND ${GO_EXECUTABLE} build -ldflags "-s -w -X 'auraspeed/cmd/root.Version=${AURASPEED_VERSION}'" -o auraspeed${CMAKE_EXECUTABLE_SUFFIX} ./cmd/main.go
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Building AuraSpeed ${AURASPEED_VERSION}"
)

add_custom_target(build-release
  COMMAND ${GO_EXECUTABLE} build -ldflags "-s -w -X 'auraspeed/cmd/root.Version=${AURASPEED_VERSION}' -X 'auraspeed/cmd/root.Commit=${AURASPEED_VERSION}' -X 'auraspeed/cmd/root.BuildTime=${BUILD_TIME}'" -o auraspeed${CMAKE_EXECUTABLE_SUFFIX} ./cmd/main.go
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Building AuraSpeed release ${AURASPEED_VERSION}"
)

add_custom_target(test
  COMMAND ${GO_EXECUTABLE} test ./... -v -short
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Running tests"
)

add_custom_target(clean
  COMMAND ${CMAKE_COMMAND} -E remove -f auraspeed
  COMMAND ${CMAKE_COMMAND} -E remove -f auraspeed.exe
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Cleaning build artifacts"
)

add_custom_target(lint
  COMMAND ${GO_EXECUTABLE} vet ./...
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Running go vet"
)

add_custom_target(format
  COMMAND ${GO_EXECUTABLE} fmt ./...
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Formatting source code"
)

add_custom_target(deps
  COMMAND ${GO_EXECUTABLE} mod tidy
  COMMAND ${GO_EXECUTABLE} mod download
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Tidying and downloading dependencies"
)

add_custom_target(distclean
  COMMAND ${CMAKE_COMMAND} -E remove_directory bin
  COMMAND ${CMAKE_COMMAND} -E remove_directory dist
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "Removing build output directories"
)

if(NOT CMAKE_BUILD_TYPE)
  set(CMAKE_BUILD_TYPE "Release" CACHE STRING "" FORCE)
endif()
