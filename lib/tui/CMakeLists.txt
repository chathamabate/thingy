
set(tui_Srcs 
    "src/private/Util.cpp"
)


add_library(
    tui
    STATIC ${tui_Srcs}
)

target_link_libraries(
    tui
    PUBLIC  ReqCompilerFlags
)

target_include_directories(
    tui
    PUBLIC include
)
