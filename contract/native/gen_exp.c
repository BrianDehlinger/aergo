/**
 * @file    gen_exp.c
 * @copyright defined in aergo/LICENSE.txt
 */

#include "common.h"

#include "ast_id.h"
#include "ir_abi.h"
#include "ir_md.h"
#include "gen_util.h"
#include "syslib.h"

#include "gen_exp.h"

static BinaryenExpressionRef
exp_gen_lit(gen_t *gen, ast_exp_t *exp)
{
    value_t *val = &exp->u_lit.val;
    meta_t *meta = &exp->meta;
    ir_md_t *md = gen->md;

    switch (val->type) {
    case TYPE_BOOL:
        return i32_gen(gen, val_bool(val) ? 1 : 0);

    case TYPE_BYTE:
        return i32_gen(gen, val_byte(val));

    case TYPE_INT256:
        if (is_int256_meta(meta)) {
            char *str;

            if (value_fits_i32(val))
                return syslib_call(gen, FN_MPZ_SET_I32, 1, i32_gen(gen, val_i64(val)));

            if (value_fits_i64(val))
                return syslib_call(gen, FN_MPZ_SET_I64, 1, i64_gen(gen, val_i64(val)));

            str = mpz_get_str(NULL, 10, val_mpz(val));
            ASSERT(str != NULL && str[0] != '\0');

            return syslib_call(gen, FN_MPZ_SET_STR, 1, i32_gen(gen, sgmt_add_str(&md->sgmt, str)));
        }

        if (is_int64_meta(meta))
            return i64_gen(gen, val_i64(val));

        return i32_gen(gen, val_i64(val));

    case TYPE_OBJECT:
        return i32_gen(gen, sgmt_add_raw(&md->sgmt, val_ptr(val), val_size(val)));

    default:
        ASSERT2(!"invalid value", val->type, meta->type);
    }

    return NULL;
}

static BinaryenExpressionRef
exp_gen_array(gen_t *gen, ast_exp_t *exp, BinaryenExpressionRef value)
{
    uint32_t offset = 0;
    ast_exp_t *id_exp = exp->u_arr.id_exp;
    ast_exp_t *idx_exp = exp->u_arr.idx_exp;
    meta_t *meta = &exp->meta;
    BinaryenExpressionRef address;

    ASSERT1(is_array_meta(&id_exp->meta), id_exp->meta.type);

    address = exp_gen(gen, id_exp, NULL);

    if (is_fixed_meta(&id_exp->meta) && is_lit_exp(idx_exp)) {
        /* The total size of the subdimensions is required. */
        offset = meta_align(meta) + val_i64(&idx_exp->u_lit.val) * meta_memsz(meta);
    }
    else {
        fn_kind_t kind = FN_ARR_GET_I32;
        meta_t *elem_meta = id_exp->meta.elems[0];

        if (is_int64_meta(elem_meta))
            kind = FN_ARR_GET_I64;

        address = syslib_call(gen, kind, 4, address, i32_gen(gen, meta->arr_dim),
                             exp_gen(gen, idx_exp, NULL), i32_gen(gen, meta_typsz(elem_meta)));
    }

    ASSERT2(offset % meta_align(meta) == 0, offset, meta_align(meta));

    if (value != NULL) {
        ASSERT1(!is_array_meta(meta), meta->type);
        return BinaryenStore(gen->module, meta_iosz(meta), offset, 0, address, value,
                             meta_gen(meta));
    }

    if (is_array_meta(meta) || is_struct_meta(meta))
        /* When returning a middle element of a multidimensional array or returning a struct,
         * we must return an address value. */
        return BinaryenBinary(gen->module, BinaryenAddInt32(), address, i32_gen(gen, offset));

    return BinaryenLoad(gen->module, meta_iosz(meta), is_signed_meta(meta), offset, 0,
                        meta_gen(meta), address);
}

static BinaryenExpressionRef
exp_gen_cast(gen_t *gen, ast_exp_t *exp)
{
    ast_exp_t *val_exp = exp->u_cast.val_exp;
    meta_t *from_meta = &val_exp->meta;
    meta_t *to_meta = &exp->meta;
    ir_md_t *md = gen->md;
    BinaryenOp op = 0;
    BinaryenExpressionRef value;

    value = exp_gen(gen, val_exp, NULL);

    switch (from_meta->type) {
    case TYPE_BOOL:
        ASSERT1(is_string_meta(to_meta), to_meta->type);
        return BinaryenSelect(gen->module, value, i32_gen(gen, sgmt_add_str(&md->sgmt, "true")),
                              i32_gen(gen, sgmt_add_str(&md->sgmt, "false")));

    case TYPE_BYTE:
        if (is_string_meta(to_meta))
            return syslib_call(gen, FN_CTOA, 1, value);
        /* fall through */

    case TYPE_INT8:
    case TYPE_INT16:
    case TYPE_INT32:
        if (is_string_meta(to_meta))
            return syslib_call(gen, FN_ITOA32, 1, value);

        if (is_int256_meta(to_meta))
            return syslib_call(gen, FN_MPZ_SET_I32, 1, value);

        if (is_int64_meta(to_meta))
            op = BinaryenExtendSInt32();
        else
            return value;
        break;

    case TYPE_INT64:
        if (is_string_meta(to_meta))
            return syslib_call(gen, FN_ITOA64, 1, value);

        if (is_int256_meta(to_meta))
            return syslib_call(gen, FN_MPZ_SET_I64, 1, value);

        if (!is_int64_meta(to_meta))
            op = BinaryenWrapInt64();
        else
            return value;
        break;

    case TYPE_INT256:
        if (is_string_meta(to_meta))
            return syslib_call(gen, FN_MPZ_GET_STR, 1, value);

        if (is_int64_meta(to_meta))
            return syslib_call_1(gen, FN_MPZ_GET_I64, value);

        return syslib_call_1(gen, FN_MPZ_GET_I32, value);

    case TYPE_STRING:
        if (is_int64_meta(to_meta))
            return syslib_call(gen, FN_ATOI64, 1, value);

        if (is_int256_meta(to_meta))
            return syslib_call(gen, FN_MPZ_SET_STR, 1, value);

        return syslib_call(gen, FN_ATOI32, 1, value);

    default:
        ASSERT2(!"invalid conversion", from_meta->type, to_meta->type);
    }

    return BinaryenUnary(gen->module, op, value);
}

static BinaryenExpressionRef
exp_gen_unary(gen_t *gen, ast_exp_t *exp)
{
    meta_t *meta = &exp->meta;
    BinaryenExpressionRef value;

    value = exp_gen(gen, exp->u_un.val_exp, NULL);

    switch (exp->u_un.kind) {
    case OP_NEG:
        if (is_int256_meta(meta))
            return syslib_call_1(gen, FN_MPZ_NEG, value);

        if (is_int64_meta(meta))
            return BinaryenBinary(gen->module, BinaryenSubInt64(), i64_gen(gen, 0), value);

        return BinaryenBinary(gen->module, BinaryenSubInt32(), i32_gen(gen, 0), value);

    case OP_NOT:
        return BinaryenUnary(gen->module, BinaryenEqZInt32(), value);

    case OP_BIT_NOT:
        if (is_int256_meta(meta))
            return syslib_call(gen, FN_MPZ_COM, 1, value);

        if (is_int64_meta(meta))
            return BinaryenBinary(gen->module, BinaryenXorInt64(), value, i64_gen(gen, -1));

        return BinaryenBinary(gen->module, BinaryenXorInt32(), value, i32_gen(gen, -1));

    default:
        ASSERT1(!"invalid operator", exp->u_un.kind);
    }

    return NULL;
}

static BinaryenExpressionRef
exp_gen_op_arith(gen_t *gen, ast_exp_t *exp, meta_t *meta)
{
    BinaryenOp op;
    BinaryenExpressionRef left, right;

    left = exp_gen(gen, exp->u_bin.l_exp, NULL);
    right = exp_gen(gen, exp->u_bin.r_exp, NULL);

    switch (exp->u_bin.kind) {
    case OP_ADD:
        if (is_string_meta(meta))
            return syslib_call_2(gen, FN_STRCAT, left, right);

        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_ADD, left, right);

        if (is_int64_meta(meta))
            op = BinaryenAddInt64();
        else
            op = BinaryenAddInt32();
        break;

    case OP_SUB:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_SUB, left, right);

        if (is_int64_meta(meta))
            op = BinaryenSubInt64();
        else
            op = BinaryenSubInt32();
        break;

    case OP_MUL:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_MUL, left, right);

        if (is_int64_meta(meta))
            op = BinaryenMulInt64();
        else
            op = BinaryenMulInt32();
        break;

    case OP_DIV:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_DIV, left, right);

        if (is_int64_meta(meta))
            op = BinaryenDivSInt64();
        else
            op = BinaryenDivSInt32();
        break;

    case OP_MOD:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_MOD, left, right);

        if (is_int64_meta(meta))
            op = BinaryenRemSInt64();
        else
            op = BinaryenRemSInt32();
        break;

    case OP_BIT_AND:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_AND, left, right);

        if (is_int64_meta(meta))
            op = BinaryenAndInt64();
        else
            op = BinaryenAndInt32();
        break;

    case OP_BIT_OR:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_OR, left, right);

        if (is_int64_meta(meta))
            op = BinaryenOrInt64();
        else
            op = BinaryenOrInt32();
        break;

    case OP_BIT_XOR:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_XOR, left, right);

        if (is_int64_meta(meta))
            op = BinaryenXorInt64();
        else
            op = BinaryenXorInt32();
        break;

    case OP_BIT_SHR:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_SHR, left, right);

        if (is_int64_meta(meta))
            op = BinaryenShrSInt64();
        else
            op = BinaryenShrSInt32();
        break;

    case OP_BIT_SHL:
        if (is_int256_meta(meta))
            return syslib_call_2(gen, FN_MPZ_SHL, left, right);

        if (is_int64_meta(meta))
            op = BinaryenShlInt64();
        else
            op = BinaryenShlInt32();
        break;

    default:
        ASSERT1(!"invalid operator", exp->u_bin.kind);
    }

    return BinaryenBinary(gen->module, op, left, right);
}

static BinaryenExpressionRef
exp_gen_op_cmp(gen_t *gen, ast_exp_t *exp, meta_t *meta)
{
    BinaryenOp op;
    BinaryenExpressionRef left, right;

    left = exp_gen(gen, exp->u_bin.l_exp, NULL);
    right = exp_gen(gen, exp->u_bin.r_exp, NULL);

    switch (exp->u_bin.kind) {
    case OP_AND:
        ASSERT(!is_int64_meta(meta));
        op = BinaryenAndInt32();
        break;

    case OP_OR:
        ASSERT(!is_int64_meta(meta));
        op = BinaryenOrInt32();
        break;

    case OP_EQ:
        if (is_int256_meta(meta)) {
            left = syslib_call_2(gen, FN_MPZ_CMP, left, right);
            right = i32_gen(gen, 0);
        }

        if (is_int64_meta(meta))
            op = BinaryenEqInt64();
        else
            op = BinaryenEqInt32();
        break;

    case OP_NE:
        if (is_int256_meta(meta)) {
            left = syslib_call_2(gen, FN_MPZ_CMP, left, right);
            right = i32_gen(gen, 0);
        }

        if (is_int64_meta(meta))
            op = BinaryenNeInt64();
        else
            op = BinaryenNeInt32();
        break;

    case OP_LT:
        if (is_int256_meta(meta)) {
            left = syslib_call_2(gen, FN_MPZ_CMP, left, right);
            right = i32_gen(gen, 0);
        }

        if (is_int64_meta(meta))
            op = BinaryenLtSInt64();
        else
            op = BinaryenLtSInt32();
        break;

    case OP_GT:
        if (is_int256_meta(meta)) {
            left = syslib_call_2(gen, FN_MPZ_CMP, left, right);
            right = i32_gen(gen, 0);
        }

        if (is_int64_meta(meta))
            op = BinaryenGtSInt64();
        else
            op = BinaryenGtSInt32();
        break;

    case OP_LE:
        if (is_int256_meta(meta)) {
            left = syslib_call_2(gen, FN_MPZ_CMP, left, right);
            right = i32_gen(gen, 0);
        }

        if (is_int64_meta(meta))
            op = BinaryenLeSInt64();
        else
            op = BinaryenLeSInt32();
        break;

    case OP_GE:
        if (is_int256_meta(meta)) {
            left = syslib_call_2(gen, FN_MPZ_CMP, left, right);
            right = i32_gen(gen, 0);
        }

        if (is_int64_meta(meta))
            op = BinaryenGeSInt64();
        else
            op = BinaryenGeSInt32();
        break;

    default:
        ASSERT1(!"invalid operator", exp->u_bin.kind);
    }

    if (is_string_meta(meta)) {
        left = syslib_call_2(gen, FN_STRCMP, left, right);
        right = i32_gen(gen, 0);
    }

    return BinaryenBinary(gen->module, op, left, right);
}

static BinaryenExpressionRef
exp_gen_binary(gen_t *gen, ast_exp_t *exp)
{
    switch (exp->u_bin.kind) {
    case OP_ADD:
    case OP_SUB:
    case OP_MUL:
    case OP_DIV:
    case OP_MOD:
    case OP_BIT_AND:
    case OP_BIT_OR:
    case OP_BIT_XOR:
    case OP_BIT_SHR:
    case OP_BIT_SHL:
        return exp_gen_op_arith(gen, exp, &exp->meta);

    case OP_AND:
    case OP_OR:
    case OP_EQ:
    case OP_NE:
    case OP_LT:
    case OP_GT:
    case OP_LE:
    case OP_GE:
        return exp_gen_op_cmp(gen, exp, &exp->u_bin.l_exp->meta);

    default:
        ASSERT1(!"invalid operator", exp->u_bin.kind);
    }

    return NULL;
}

static BinaryenExpressionRef
exp_gen_ternary(gen_t *gen, ast_exp_t *exp)
{
    return BinaryenSelect(gen->module, exp_gen(gen, exp->u_tern.pre_exp, NULL),
                          exp_gen(gen, exp->u_tern.in_exp, NULL),
                          exp_gen(gen, exp->u_tern.post_exp, NULL));
}

static BinaryenExpressionRef
exp_gen_access(gen_t *gen, ast_exp_t *exp, BinaryenExpressionRef value)
{
    uint32_t offset = exp->meta.rel_offset;
    meta_t *meta = &exp->meta;
    BinaryenExpressionRef address;

    address = exp_gen(gen, exp->u_acc.qual_exp, NULL);

    if (is_fn_id(exp->id))
        return address;

    if (value != NULL) {
        ASSERT1(!is_array_meta(meta) && !is_struct_meta(meta), meta->type);

        return BinaryenStore(gen->module, meta_iosz(meta), offset, 0, address, value,
                             meta_gen(meta));
    }

    if (is_struct_meta(meta) || is_array_meta(meta))
        return BinaryenBinary(gen->module, BinaryenAddInt32(), address, i32_gen(gen, offset));

    return BinaryenLoad(gen->module, meta_iosz(meta), is_signed_meta(meta), offset, 0,
                        meta_gen(meta), address);
}

static BinaryenExpressionRef
exp_gen_call(gen_t *gen, ast_exp_t *exp)
{
    int i;
    char *name;
    BinaryenIndex arg_cnt;
    BinaryenExpressionRef *arguments;

    if (exp->u_call.kind == FN_UDF || exp->u_call.kind == FN_NEW) {
        name = exp->u_call.qname;
    }
    else {
        sys_fn_t *sys_fn = SYS_FN(exp->u_call.kind);

        if (exp->u_call.kind != FN_ALLOCA)
            md_add_abi(gen->md, syslib_abi(sys_fn));

        name = sys_fn->qname;
    }

    ASSERT1(name != NULL, exp->u_call.kind);

    arg_cnt = vector_size(exp->u_call.arg_exps);
    arguments = xmalloc(sizeof(BinaryenExpressionRef) * arg_cnt);

    vector_foreach(exp->u_call.arg_exps, i) {
        arguments[i] = exp_gen(gen, vector_get_exp(exp->u_call.arg_exps, i), NULL);
    }

    return BinaryenCall(gen->module, name, arguments, arg_cnt, meta_gen(&exp->meta));
}

static BinaryenExpressionRef
exp_gen_sql(gen_t *gen, ast_exp_t *exp)
{
    /* TODO */
    return i32_gen(gen, 0);
}

static BinaryenExpressionRef
exp_gen_global(gen_t *gen, ast_exp_t *exp, BinaryenExpressionRef value)
{
    ASSERT(exp->u_glob.name != NULL);

    if (value != NULL)
        return BinaryenSetGlobal(gen->module, exp->u_glob.name, value);

    return BinaryenGetGlobal(gen->module, exp->u_glob.name, BinaryenTypeInt32());
}

static BinaryenExpressionRef
exp_gen_reg(gen_t *gen, ast_exp_t *exp, BinaryenExpressionRef value)
{
    BinaryenExpressionRef address;

    address = BinaryenGetLocal(gen->module, exp->meta.base_idx, meta_gen(&exp->meta));

    if (value != NULL)
        return BinaryenSetLocal(gen->module, exp->meta.base_idx, value);

    return address;
}

static BinaryenExpressionRef
exp_gen_mem(gen_t *gen, ast_exp_t *exp, BinaryenExpressionRef value)
{
    uint32_t offset;
    meta_t *meta = &exp->meta;
    BinaryenExpressionRef address;

    address = BinaryenGetLocal(gen->module, meta->base_idx, BinaryenTypeInt32());

    offset = meta->rel_addr + meta->rel_offset;

    if (value != NULL)
        return BinaryenStore(gen->module, meta_iosz(meta), offset, 0, address, value,
                             meta_gen(meta));

    if (is_raw_meta(meta))
        return BinaryenBinary(gen->module, BinaryenAddInt32(), address, i32_gen(gen, offset));

    return BinaryenLoad(gen->module, meta_iosz(meta), is_signed_meta(meta), offset, 0,
                        meta_gen(meta), address);
}

BinaryenExpressionRef
exp_gen(gen_t *gen, ast_exp_t *exp, BinaryenExpressionRef value)
{
    ASSERT(exp != NULL);

    switch (exp->kind) {
    case EXP_LIT:
        return exp_gen_lit(gen, exp);

    case EXP_ARRAY:
        return exp_gen_array(gen, exp, value);

    case EXP_CAST:
        return exp_gen_cast(gen, exp);

    case EXP_UNARY:
        return exp_gen_unary(gen, exp);

    case EXP_BINARY:
        return exp_gen_binary(gen, exp);

    case EXP_TERNARY:
        return exp_gen_ternary(gen, exp);

    case EXP_ACCESS:
        return exp_gen_access(gen, exp, value);

    case EXP_CALL:
        return exp_gen_call(gen, exp);

    case EXP_SQL:
        return exp_gen_sql(gen, exp);

    case EXP_GLOBAL:
        return exp_gen_global(gen, exp, value);

    case EXP_REG:
        return exp_gen_reg(gen, exp, value);

    case EXP_MEM:
        return exp_gen_mem(gen, exp, value);

    default:
        ASSERT1(!"invalid expression", exp->kind);
    }

    return NULL;
}

/* end of gen_exp.c */
