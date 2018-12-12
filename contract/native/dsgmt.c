/**
 * @file    dsgmt.c
 * @copyright defined in aergo/LICENSE.txt
 */

#include "common.h"

#include "dsgmt.h"

static void
dsgmt_extend(dsgmt_t *dsgmt)
{
    dsgmt->cap += DSGMT_INIT_CAPACITY;

    dsgmt->lens = xrealloc(dsgmt->lens, sizeof(BinaryenIndex) * dsgmt->cap);
    dsgmt->addrs = xrealloc(dsgmt->addrs, sizeof(BinaryenExpressionRef) * dsgmt->cap);
    dsgmt->datas = xrealloc(dsgmt->datas, sizeof(char *) * dsgmt->cap);
}

int
dsgmt_add(dsgmt_t *dsgmt, BinaryenModuleRef module, void *ptr, uint32_t len)
{
    uint32_t offset = dsgmt->offset;

    if (ptr != NULL) {
        if (dsgmt->size >= dsgmt->cap)
            dsgmt_extend(dsgmt);

        dsgmt->lens[dsgmt->size] = (BinaryenIndex)len;
        dsgmt->addrs[dsgmt->size] = BinaryenConst(module, BinaryenLiteralInt32(offset));
        dsgmt->datas[dsgmt->size] = ptr;

        dsgmt->size++;
    }

    dsgmt->offset += ALIGN64(len);

    return offset;
}

/* end of dsgmt.c */
