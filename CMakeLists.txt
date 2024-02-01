
cmake_minimum_required(VERSION 3.21)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Tine for project specs.
project(
    TiDoc
    VERSION 1.0
    DESCRIPTION "Document Creator for TI-84 PCE."
    LANGUAGES CXX 
)

# specify the C++ standard
add_library(ReqCompilerFlags INTERFACE)
target_compile_features(ReqCompilerFlags INTERFACE cxx_std_20)

# add compiler warning flags just when building this project via
# the BUILD_INTERFACE genex
set(gcc_like_cxx "$<COMPILE_LANG_AND_ID:CXX,ARMClang,AppleClang,Clang,GNU,LCC>")
set(msvc_cxx "$<COMPILE_LANG_AND_ID:CXX,MSVC>")
target_compile_options(ReqCompilerFlags INTERFACE
  "$<${gcc_like_cxx}:$<BUILD_INTERFACE:-Wall;-Wextra;-Wshadow;-Wformat=2;-Wunused>>"
  "$<${msvc_cxx}:$<BUILD_INTERFACE:-W3>>"
)

# We are not going to use vcpkg so that we can
# use github actions easily.

enable_testing()

include(FetchContent)
FetchContent_Declare(
  googletest
  URL https://github.com/google/googletest/archive/03597a01ee50ed33e9dfd640b249b4be3799d395.zip
)

include(GoogleTest)

add_executable(${PROJECT_NAME} src/main.cpp)
target_link_libraries(${PROJECT_NAME} 
    PUBLIC ReqCompilerFlags
)

