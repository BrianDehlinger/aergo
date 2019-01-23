/**
 * @file    ir.h
 * @copyright defined in aergo/LICENSE.txt
 */

#ifndef _IR_H
#define _IR_H

#include "common.h"

#include "array.h"
#include "ir_sgmt.h"

#ifndef _META_T
#define _META_T
typedef struct meta_s meta_t;
#endif /* ! _META_T */

#ifndef _IR_FN_T
#define _IR_FN_T
typedef struct ir_fn_s ir_fn_t;
#endif /* ! _IR_FN_T */

typedef struct ir_s {
    array_t abis;
    array_t fns;

    ir_sgmt_t sgmt;

    uint32_t offset;
} ir_t;

ir_t *ir_new(void);

void ir_add_heap(ir_t *ir, meta_t *meta, int idx);
void ir_add_fn(ir_t *ir, ir_fn_t *fn);

#endif /* no _IR_H */