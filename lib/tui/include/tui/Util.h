#pragma once

#include <stdint.h>
#include <stdio.h>


#define ESC "\033"
#define CSI ESC "["

#define RESET CSI "m"

#define BOLD CSI "1" "m"
#define ITALIC CSI "3" "m"
#define UNDERLINE CSI "4" "m"

#define FG_BLACK CSI "30" "m"
#define FG_RED  CSI "31" "m"
#define FG_GREEN CSI "32" "m"
#define FG_YELLOW CSI "33" "m"
#define FG_BLUE CSI "34" "m"
#define FG_MAGENTA CSI "35" "m"
#define FG_CYAN CSI "36" "m"
#define FG_WHITE CSI "37" "m"

#define FG_BRIGHT_BLACK CSI "90" "m"
#define FG_BRIGHT_RED CSI "91" "m"
#define FG_BRIGHT_GREEN CSI "92" "m"
#define FG_BRIGHT_YELLOW CSI "93" "m"
#define FG_BRIGHT_BLUE CSI "94" "m"
#define FG_BRIGHT_MAGENTA CSI "95" "m"
#define FG_BRIGHT_CYAN CSI "96" "m"
#define FG_BRIGHT_WHITE CSI "97" "m"

#define BG_BLACK CSI "40" "m"
#define BG_RED  CSI "41" "m"
#define BG_GREEN CSI "42" "m"
#define BG_YELLOW CSI "43" "m"
#define BG_BLUE CSI "44" "m"
#define BG_MAGENTA CSI "45" "m"
#define BG_CYAN CSI "46" "m"
#define BG_WHITE CSI "47" "m"

#define BG_BRIGHT_BLACK CSI "100" "m"
#define BG_BRIGHT_RED CSI "101" "m"
#define BG_BRIGHT_GREEN CSI "102" "m"
#define BG_BRIGHT_YELLOW CSI "103" "m"
#define BG_BRIGHT_BLUE CSI "104" "m"
#define BG_BRIGHT_MAGENTA CSI "105" "m"
#define BG_BRIGHT_CYAN CSI "106" "m"
#define BG_BRIGHT_WHITE CSI "107" "m"

#define SAVE_CURSOR CSI "s"
#define RESTORE_CURSOR CSI "u"

namespace tui {
    void printAt(uint8_t row, uint8_t col, const char *text);
}


