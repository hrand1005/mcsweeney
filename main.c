#include <stdio.h>
#include <stdlib.h>
#include <string.h>

const size_t MAX_LINES = 10;
const size_t LINE_BUFF = 256;

int process_each(char** lines, size_t n);

int main() {
    FILE* fp;
    fp = fopen("test.txt", "r");
    if (fp == NULL) {
        perror("Error opening file");
        return 1;
    }
    
    char buf[LINE_BUFF];
    char** lines = calloc(MAX_LINES, sizeof(char*));

    size_t i = 0;
    while (fgets(buf, LINE_BUFF, fp) && i < MAX_LINES) {
        lines[i] = strdup(buf);
        memset(buf, '\0', LINE_BUFF);
        i++;
    }

    if (ferror(fp) || i == 0) {
        perror("Error reading file");
        fclose(fp);
        return 1;
    }
    fclose(fp);

    process_each(lines, i);

    for (size_t j = 0; j < i; j++) {
        free(lines[j]);
    }
    free(lines);

    return 0;
}

int process_each(char** lines, size_t n) {
    for (int i = 0; i < n; i++) {
        printf("Processing file: %s\n", lines[i]);
    }
    return 0;
}
