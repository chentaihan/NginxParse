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
//
//struct ngx_event_s {
//    void            *data;
//    unsigned         cancelable:1;
//
//
//#if 1 && (1 || 0)
//    ngx_event_ovlp_t ovlp;
//#endif
//
//    ngx_queue_t      queue;
//
//#if 1 && (1&&0) || (2||0) && (1+2)
//    void            *thr_ctx;
//#if 0
//    uint32_t         padding[NGX_EVENT_T_PADDING];
//#endif
//#endif
//
//#if 0 && 1 || 1 && 2
//    unsigned         kq_vnode:1;
//    int              kq_errno;
//#endif
//
//#if 0
//    unsigned         kq_vnode1:1;
//    int              kq_errno1;
//#else
//    unsigned         kq_vnode2:1;
//    int              kq_errno2;
//#endif
//
//#if 1
//    int              available;
//#if 0
//    int              available1;
//#else
//    unsigned         available2:1;
//#endif
//#else
//    unsigned         available:1;
//#endif
//
//    ngx_event_handler_pt  handler;
//    int a,b,c;
//    int (*f)(int a, int b)
//};
//
//
//typedef struct {
//    union {
//        ngx_http_geo_trees_t         trees;
//        ngx_http_geo_high_ranges_t   high;
//    } u;
//
//    ngx_array_t                     *proxies;
//    unsigned                         proxy_recursive:1;
//
//    ngx_int_t                        index;
//    int a,b,c;
//} ngx_http_geo_ctx_t;
//
//struct ngx_module_s {
//    ngx_uint_t            ctx_index;
//    ngx_uint_t            index;
//
//    char                 *name;
//
//    ngx_uint_t            spare0;
//    ngx_uint_t            spare1;
//
//    ngx_uint_t            version;
//    const char           *signature;
//
//    void                 *ctx;
//    ngx_command_t        *commands;
//    ngx_uint_t            type;
//
//    ngx_int_t           (*init_master)(ngx_log_t *log);
//
//    ngx_int_t           (*init_module)(ngx_cycle_t *cycle);
//
//    ngx_int_t           (*init_process)(ngx_cycle_t *cycle);
//    ngx_int_t           (*init_thread)(ngx_cycle_t *cycle);
//    void                (*exit_thread)(ngx_cycle_t *cycle);
//    void                (*exit_process)(ngx_cycle_t *cycle);
//
//    void                (*exit_master)(ngx_cycle_t *cycle);
//
//    uintptr_t             spare_hook0;
//    uintptr_t             spare_hook1;
//    uintptr_t             spare_hook2;
//    uintptr_t             spare_hook3;
//    uintptr_t             spare_hook4;
//    uintptr_t             spare_hook5;
//    uintptr_t             spare_hook6;
//    uintptr_t   spare_hook7;
//};

//
//typedef struct {
//    ngx_int_t   (*preconfiguration)(ngx_conf_t *cf);
//    ngx_int_t   (*postconfiguration)(ngx_conf_t *cf);
//
//    void       *(*create_main_conf)(ngx_conf_t *cf);
//    char       *(*init_main_conf)(ngx_conf_t *cf, void *conf);
//
//    void       *(*create_srv_conf)(ngx_conf_t *cf);
//    char       *(*merge_srv_conf)(ngx_conf_t *cf, void *prev, void *conf);
//
//    void       *(*create_loc_conf)(ngx_conf_t *cf);
//    char       *(*merge_loc_conf)(ngx_conf_t *cf, void *prev, void *conf);
//} ngx_http_module_t;
//
//static ngx_http_module_t  ngx_http_v2_module_ctx = {
//    ngx_http_v2_add_variables,             /* preconfiguration */
//    NULL,                                  /* postconfiguration */
//
//    ngx_http_v2_create_main_conf,          /* create main configuration */
//    ngx_http_v2_init_main_conf,            /* init main configuration */
//
//    ngx_http_v2_create_srv_conf,           /* create server configuration */
//    ngx_http_v2_merge_srv_conf,            /* merge server configuration */
//
//    ngx_http_v2_create_loc_conf,           /* create location configuration */
//    ngx_http_v2_merge_loc_conf             /* merge location configuration */
//};

struct ngx_command_t {
    ngx_str_t             name;
    ngx_uint_t            type;
    char               *(*set)(ngx_conf_t *cf, ngx_command_t *cmd, void *conf);
    ngx_uint_t            conf;
    ngx_uint_t            offset;
    void                 *post;
};

static ngx_command_t  ngx_http_secure_link_commands[] = {

    { ngx_string("secure_link"),
      NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_HTTP_LOC_CONF|NGX_CONF_TAKE1,
      ngx_http_set_complex_value_slot,
      NGX_HTTP_LOC_CONF_OFFSET,
      offsetof(ngx_http_secure_link_conf_t, variable),
      NULL },

    {
      ngx_string("secure_link_md5"),
      NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_HTTP_LOC_CONF|NGX_CONF_TAKE1,
      ngx_http_set_complex_value_slot,
      NGX_HTTP_LOC_CONF_OFFSET,
      offsetof(ngx_http_secure_link_conf_t, md5),
      NULL },

    { ngx_string("secure_link_secret"),
      NGX_HTTP_MAIN_CONF|NGX_HTTP_SRV_CONF|NGX_HTTP_LOC_CONF|NGX_CONF_TAKE1,
      ngx_conf_set_str_slot,
      NGX_HTTP_LOC_CONF_OFFSET,
      offsetof(ngx_http_secure_link_conf_t, secret),
      NULL },

      ngx_null_command
};

//typedef struct {
//    ngx_conf_post_handler_pt post_handler;
//} ngx_conf_post_t;

//ngx_str_t  ngx_http_html_default_types[] = {
//    ngx_string("text/html"),
//    ngx_null_string
//};

//typedef struct ngx_stream_session_s  ngx_stream_session_t;