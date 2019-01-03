/**
 * @file    gen_fn.c
 * @copyright defined in aergo/LICENSE.txt
 */

#include "common.h"

#include "gen_bb.h"
#include "gen_util.h"

#include "gen_fn.h"

void
fn_gen(gen_t *gen, ir_fn_t *fn)
{
    int i;
    BinaryenFunctionTypeRef spec;
    BinaryenExpressionRef body;

    spec = BinaryenAddFunctionType(gen->module, fn->name, BinaryenTypeNone(), fn->params,
                                   fn->param_cnt);

    /* 1st local for base stack address */
    gen_add_local(gen, TYPE_INT32);

    /* 2nd local for relooper */
    gen_add_local(gen, TYPE_INT32);

    /* generate local variables */
    array_foreach(&fn->locals, i) {
        gen_add_local(gen, array_get_id(&fn->locals, i)->meta.type);
    }

    gen->relooper = RelooperCreate();

    /* generate basic blocks */
    array_foreach(&fn->bbs, i) {
        bb_gen(gen, array_get_bb(&fn->bbs, i));
    }

    /* generate branches */
    array_foreach(&fn->bbs, i) {
        br_gen(gen, array_get_bb(&fn->bbs, i));
    }

    body = RelooperRenderAndDispose(gen->relooper, fn->entry_bb->rb, fn->param_cnt + 1,
                                    gen->module);

    BinaryenAddFunction(gen->module, fn->name, spec, gen->locals, gen->local_cnt,
                        BinaryenBlock(gen->module, NULL, &body, 1, BinaryenTypeNone()));

    gen->local_cnt = 0;
    gen->locals = NULL;
}

/* end of gen_fn.c */
