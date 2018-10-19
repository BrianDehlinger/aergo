/**
 * @file    error.h
 * @copyright defined in aergo/LICENSE.txt
 */

#ifndef _ERROR_H
#define _ERROR_H

#include "common.h"

#include "enum.h"

#define DESC_MAX_LEN            512

#define FATAL(ec, ...)          error_exit((ec), LVL_FATAL, ## __VA_ARGS__)

#define ERROR(ec, pos, ...)     error_push((ec), LVL_ERROR, (pos), ## __VA_ARGS__)
#define INFO(ec, pos, ...)      error_push((ec), LVL_INFO, (pos), ## __VA_ARGS__)
#define WARN(ec, pos, ...)      error_push((ec), LVL_WARN, (pos), ## __VA_ARGS__)
#define DEBUG(ec, pos, ...)     error_push((ec), LVL_DEBUG, (pos), ## __VA_ARGS__)

#define error_first()           error_item(0)
#define error_last()            error_item(error_size())

#define is_no_error()           (error_size() == 0)

typedef struct error_s {
    ec_t code;
    errlvl_t level;
    char *path;
    int line;
    int col;
    char desc[DESC_MAX_LEN];
} error_t;

char *error_to_str(ec_t ec);
ec_t error_to_code(char *str);

int error_size(void);
ec_t error_item(int idx);

void error_push(ec_t ec, errlvl_t lvl, src_pos_t *pos, ...);
error_t *error_pop(void);

void error_clear(void);
void error_dump(void);

void error_exit(ec_t ec, errlvl_t lvl, ...);

static inline int
error_cmp(const void *x, const void *y)
{
    int res;
    error_t *e1 = (error_t *)x;
    error_t *e2 = (error_t *)y;

    ASSERT(e1->path != NULL);
    ASSERT(e2->path != NULL);

    res = strcmp(e1->path, e2->path);
    if (res != 0)
        return res;

    if (e1->line < e2->line)
        return -1;
    else if (e1->line > e2->line)
        return 1;

    if (e1->col < e2->col)
        return -1;
    else if (e1->col == e2->col)
        return 0;
    else
        return 1;
}

#endif /* ! _ERROR_H */