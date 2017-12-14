package main

const (
	NGX_COMMAND_T       = "ngx_command_t"
	NGX_STRING          = "ngx_string"
	NGX_STRING_LEN      = len(NGX_STRING)
	NGX_MODULE_T        = "ngx_module_t"
	NGX_HTTP_VARIABLE_T = "ngx_http_variable_t"
)

const (
	STRUCT_TYPE_COMMAND  = 1
	STRUCT_TYPE_MODULE   = 2
	STRUCT_TYPE_VARIABLE = 3
)

const (
	PATH_CONFIG        = "conf/"
	FILE_CONFIG_FORMAT = "config_format.html"
)

const (
	PATH_OUTPUT = "output/"
)
