#include <pjnath.h>
#include <pjlib-util.h>
#include <pjlib.h>
#include "_cgo_export.h"

void  ice_cb(pj_ice_strans *ice_strans, pj_ice_strans_op op, pj_status_t status)
{
  go_ice_callback(ice_strans,op,status);
}

void  data_cb(pj_ice_strans *ice_st, unsigned comp_id, void *pkt, pj_size_t size, const pj_sockaddr_t *src_addr, unsigned src_addr_len)
{
  go_data_callback(ice_st, comp_id, pkt, size, (pj_sockaddr_t *)(src_addr), src_addr_len);
}

pj_ice_strans_cb *new_cb(void *ice, void *data)
{
  pj_ice_strans_cb *cb;
  cb = malloc(sizeof(pj_ice_strans_cb));
  cb->on_ice_complete = ice;
  cb->on_rx_data = data;
  return cb;
}
