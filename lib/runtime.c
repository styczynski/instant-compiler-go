#include <string.h>
#include <stdlib.h>

int32_t AddStrings(int32_t a, int32_t b) {
    const int al = strlen(a);
    const int bl = strlen(b);
    char* result = (char*) malloc(sizeof(char) * (al+bl));
    strcpy(result, (char*) a);
    strcat(result, (char*) b);
    return (int32_t) result;
}