/**
 * @file    syslib.c
 * @copyright defined in aergo/LICENSE.txt
 */

#include "common.h"

#include "binaryen-c.h"
#include "ast_id.h"
#include "ast_exp.h"
#include "parse.h"
#include "ir_abi.h"
#include "ir_md.h"
#include "gen_util.h"
#include "iobuf.h"

#include "syslib.h"

char *lib_src =
"library system {\n"
"    func abs32(int32 i) int32 : \"__abs32\";\n"
"    func abs64(int64 i) int64 : \"__abs64\";\n"
"    func abs128(int128 i) int128 : \"__mpz_abs\";\n"

"    func pow32(int32 x, int32 y) int32 : \"__pow32\";\n"
"    func pow64(int64 x, int32 y) int64 : \"__pow64\";\n"
"    func pow128(int128 x, int32 y) int128 : \"__mpz_pow_ui\";\n"

"    func sign32(int32 x) int8 : \"__sign32\";\n"
"    func sign64(int64 x) int8 : \"__sign64\";\n"
"    func sign128(int128 x) int8 : \"__mpz_sign\";\n"

"    func sqrt32(int32 x) int32 : \"__sqrt32\";\n"
"    func sqrt64(int64 x) int64 : \"__sqrt64\";\n"
"    func sqrt128(int128 x) int128 : \"__mpz_sqrt\";\n"
"}";

sys_fn_t sys_fntab_[FN_MAX] = {
    { "__udf", NULL, 0, { TYPE_NONE }, TYPE_NONE },
    { "__ctor", NULL, 0, { TYPE_NONE }, TYPE_NONE },
    { "__malloc", SYSLIB_MODULE".__malloc", 1, { TYPE_INT32 }, TYPE_INT32 },
    { "__memcpy", SYSLIB_MODULE".__memcpy", 3, { TYPE_INT32, TYPE_INT32, TYPE_INT32 }, TYPE_VOID },
    { "__assert", SYSLIB_MODULE".__assert", 6,
        { TYPE_BOOL, TYPE_STRING, TYPE_STRING, TYPE_INT32, TYPE_INT32, TYPE_INT32 }, TYPE_VOID },
    { "__strcat", SYSLIB_MODULE".__strcat", 2, { TYPE_STRING, TYPE_STRING }, TYPE_STRING },
    { "__strcmp", SYSLIB_MODULE".__strcmp", 2, { TYPE_STRING, TYPE_STRING }, TYPE_INT32 },
    { "__atoi32", SYSLIB_MODULE".__atoi32", 1, { TYPE_STRING }, TYPE_INT32 },
    { "__atoi64", SYSLIB_MODULE".__atoi64", 1, { TYPE_STRING }, TYPE_INT64 },
    { "__itoa32", SYSLIB_MODULE".__itoa32", 1, { TYPE_INT32 }, TYPE_STRING },
    { "__itoa64", SYSLIB_MODULE".__itoa64", 1, { TYPE_INT64 }, TYPE_STRING },
    { "__mpz_get_i32", SYSLIB_MODULE".__mpz_get_i32", 1, { TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_get_i64", SYSLIB_MODULE".__mpz_get_i64", 1, { TYPE_INT32 }, TYPE_INT64 },
    { "__mpz_get_str", SYSLIB_MODULE".__mpz_get_str", 1, { TYPE_INT32 }, TYPE_STRING },
    { "__mpz_set_i32", SYSLIB_MODULE".__mpz_set_i32", 1, { TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_set_i64", SYSLIB_MODULE".__mpz_set_i64", 1, { TYPE_INT64 }, TYPE_INT32 },
    { "__mpz_set_str", SYSLIB_MODULE".__mpz_set_str", 1, { TYPE_STRING }, TYPE_INT32 },
    { "__mpz_add", SYSLIB_MODULE".__mpz_add", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_sub", SYSLIB_MODULE".__mpz_sub", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_mul", SYSLIB_MODULE".__mpz_mul", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_div", SYSLIB_MODULE".__mpz_div", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_mod", SYSLIB_MODULE".__mpz_mod", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_and", SYSLIB_MODULE".__mpz_and", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_or", SYSLIB_MODULE".__mpz_or", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_xor", SYSLIB_MODULE".__mpz_xor", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_rshift", SYSLIB_MODULE".__mpz_rshift", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_lshift", SYSLIB_MODULE".__mpz_lshift", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_cmp", SYSLIB_MODULE".__mpz_cmp", 2, { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_neg", SYSLIB_MODULE".__mpz_neg", 1, { TYPE_INT32 }, TYPE_INT32 },
    { "__mpz_sign", SYSLIB_MODULE".__mpz_sign", 1, { TYPE_INT32 }, TYPE_INT8 },
    { "__array_get_i32", SYSLIB_MODULE".__array_get_i32", 2,
        { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__array_get_i64", SYSLIB_MODULE".__array_get_i64", 2,
        { TYPE_INT32, TYPE_INT32 }, TYPE_INT32 },
    { "__array_set_i32", SYSLIB_MODULE".__array_set_i32", 3,
        { TYPE_INT32, TYPE_INT32, TYPE_INT32 }, TYPE_VOID },
    { "__array_set_i64", SYSLIB_MODULE".__array_set_i64", 3,
        { TYPE_INT32, TYPE_INT32, TYPE_INT64 }, TYPE_VOID },
};

void
syslib_load(ast_t *ast)
{
    flag_t flag = { FLAG_NONE, 0, 0 };
    iobuf_t src;

    iobuf_init(&src, "system_library");
    iobuf_set(&src, strlen(lib_src), lib_src);

    parse(&src, flag, ast);
}

ir_abi_t *
syslib_abi(sys_fn_t *sys_fn)
{
    int i;
    ir_abi_t *abi = xcalloc(sizeof(ir_abi_t));

    ASSERT(sys_fn != NULL);

    abi->module = SYSLIB_MODULE;
    abi->name = sys_fn->name;

    abi->param_cnt = sys_fn->param_cnt;
    abi->params = xmalloc(sizeof(BinaryenType) * abi->param_cnt);

    for (i = 0; i < abi->param_cnt; i++) {
        abi->params[i] = type_gen(sys_fn->params[i]);
    }

    abi->result = type_gen(sys_fn->result);

    return abi;
}

ast_exp_t *
syslib_new_memcpy(trans_t *trans, ast_exp_t *dest_exp, ast_exp_t *src_exp,
                   uint32_t size, src_pos_t *pos)
{
    ast_exp_t *res_exp;
    ast_exp_t *param_exp;
    vector_t *param_exps = vector_new();
    sys_fn_t *sys_fn = SYS_FN(FN_MEMCPY);

    exp_add(param_exps, dest_exp);
    exp_add(param_exps, src_exp);

    param_exp = exp_new_lit_int(size, pos);
    meta_set_int32(&param_exp->meta);

    exp_add(param_exps, param_exp);

    res_exp = exp_new_call(FN_MEMCPY, NULL, param_exps, pos);

    res_exp->u_call.qname = sys_fn->qname;
    meta_set_void(&res_exp->meta);

    md_add_imp(trans->md, syslib_abi(sys_fn));

    return res_exp;
}

BinaryenExpressionRef
syslib_call_1(gen_t *gen, fn_kind_t kind, BinaryenExpressionRef argument)
{
    sys_fn_t *sys_fn = SYS_FN(kind);

    ASSERT2(sys_fn->param_cnt == 1, kind, sys_fn->param_cnt);

    md_add_imp(gen->md, syslib_abi(sys_fn));

    return BinaryenCall(gen->module, sys_fn->qname, &argument, 1, type_gen(sys_fn->result));
}

BinaryenExpressionRef
syslib_call_2(gen_t *gen, fn_kind_t kind, BinaryenExpressionRef argument1,
              BinaryenExpressionRef argument2)
{
    sys_fn_t *sys_fn = SYS_FN(kind);
    BinaryenExpressionRef arguments[2] = { argument1, argument2 };

    ASSERT2(sys_fn->param_cnt == 2, kind, sys_fn->param_cnt);

    md_add_imp(gen->md, syslib_abi(sys_fn));

    return BinaryenCall(gen->module, sys_fn->qname, arguments, 2, type_gen(sys_fn->result));
}

BinaryenExpressionRef
syslib_gen(gen_t *gen, fn_kind_t kind, int argc, ...)
{
    int i;
    va_list vargs;
    sys_fn_t *sys_fn = SYS_FN(kind);
    BinaryenExpressionRef arguments[SYSLIB_MAX_ARGS];

    ASSERT1(argc <= SYSLIB_MAX_ARGS, argc);
    ASSERT3(sys_fn->param_cnt == argc, kind, sys_fn->param_cnt, argc);

    va_start(vargs, argc);

    for (i = 0; i < argc; i++) {
        arguments[i] = va_arg(vargs, BinaryenExpressionRef);
    }

    va_end(vargs);

    md_add_imp(gen->md, syslib_abi(sys_fn));

    return BinaryenCall(gen->module, sys_fn->qname, arguments, argc, type_gen(sys_fn->result));
}

/* end of syslib.c */
