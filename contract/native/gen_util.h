/**
 * @file    gen_util.h
 * @copyright defined in aergo/LICENSE.txt
 */

#ifndef _GEN_UTIL_H
#define _GEN_UTIL_H

#include "common.h"

#include "gen.h"
#include "meta.h"

uint32_t gen_add_local(gen_t *gen, meta_t *meta);

void gen_add_instr(gen_t *gen, BinaryenExpressionRef instr);

static inline BinaryenExpressionRef
gen_i32(gen_t *gen, int32_t v)
{
    return BinaryenConst(gen->module, BinaryenLiteralInt32(v));
}

static inline BinaryenExpressionRef
gen_i64(gen_t *gen, int64_t v)
{
    return BinaryenConst(gen->module, BinaryenLiteralInt64(v));
}

static inline BinaryenExpressionRef
gen_f32(gen_t *gen, float v)
{
    return BinaryenConst(gen->module, BinaryenLiteralFloat32(v));
}

static inline BinaryenExpressionRef
gen_f64(gen_t *gen, double v)
{
    return BinaryenConst(gen->module, BinaryenLiteralFloat64(v));
}

#endif /* no _GEN_UTIL_H */
