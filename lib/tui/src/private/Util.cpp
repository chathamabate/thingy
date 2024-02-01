
#include <tui/Util.h>

using namespace tui;

void tui::printAt(uint8_t row, uint8_t col, const char *text) {
    char buf[20];
    sprintf(buf, CSI "%u;%uH%s", row, col, text);
    printf(SAVE_CURSOR "%s" RESTORE_CURSOR, buf);
}
