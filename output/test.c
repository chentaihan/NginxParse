//typedef struct {
//    ngx_http_upstream_conf_t       upstream;
//    ngx_http_proxy_headers_t       headers;
//#if (NGX_HTTP_CACHE)
//    ngx_http_proxy_headers_t       headers_cache;
//#endif
//    ngx_array_t                   *headers_source;
//    ngx_str_t                      url;
//
//#if (NGX_HTTP_CACHE)
//    ngx_http_complex_value_t       cache_key;
//#endif
//
//    ngx_http_proxy_vars_t          vars;
//
//    ngx_uint_t                     headers_hash_bucket_size;
//
//#if (NGX_HTTP_SSL)
//    ngx_uint_t                     ssl;
//
//    ngx_array_t                   *ssl_passwords;
//#endif
//} ngx_http_proxy_loc_conf_t;

struct ngx_event_s {
    void            *data;
    unsigned         cancelable:1;


#if 1 && (1 || 0)
    ngx_event_ovlp_t ovlp;
#endif

    ngx_queue_t      queue;

#if 1 && (1&&0) || (2||0) && (1+2)
    void            *thr_ctx;
#if 0
    uint32_t         padding[NGX_EVENT_T_PADDING];
#endif
#endif

#if 0 && 1 || 1 && 2
    unsigned         kq_vnode:1;
    int              kq_errno;
#endif

#if 0
    unsigned         kq_vnode1:1;
    int              kq_errno1;
#else
    unsigned         kq_vnode2:1;
    int              kq_errno2;
#endif

#if 1
    int              available;
#if 0
    int              available1;
#else
    unsigned         available2:1;
#endif
#else
    unsigned         available:1;
#endif

    ngx_event_handler_pt  handler;
};


typedef struct {
    union {
        ngx_http_geo_trees_t         trees;
        ngx_http_geo_high_ranges_t   high;
    } u;

    ngx_array_t                     *proxies;
    unsigned                         proxy_recursive:1;

    ngx_int_t                        index;
} ngx_http_geo_ctx_t;