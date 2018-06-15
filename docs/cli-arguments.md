# Ingress Controller Command-line Arguments

```
Usage of ./nginx-ingress:
  -alsologtostderr
      log to standard error as well as files
  -default-server-tls-secret string
      Specifies a secret with a TLS certificate and key for SSL termination of the default server.
      The value must follow the following format: <namespace>/<name>.
      If not specified, the key and the cert from /etc/nginx/secrets/default is used.
  -health-status
      If present, the default server listening on port 80 with the health check location "/nginx-health"
      gets added to the main nginx configuration.
  -ingress-class string
      Specifies a class of ingress. Only processes Ingresses with this value in annotations.
      Can be used with --use-ingress-class-only. Default 'nginx' (default "nginx")
  -log_backtrace_at value
      when logging hits line file:N, emit a stack trace
  -log_dir string
      If non-empty, write log files in this directory
  -logtostderr
      log to standard error instead of files
  -nginx-configmaps string
      Specifies a configmaps resource that can be used to customize NGINX configuration.
      The value must follow the following format: <namespace>/<name>
  -nginx-plus
      Enables support for NGINX Plus.
  -proxy string
      If specified, the controller assumes a kubctl proxy server is running on the given url and creates a proxy client.
      Regenerated NGINX configuration files are not written to the disk, instead they are printed to stdout.
      Also NGINX is not getting invoked. This flag is for testing.
  -stderrthreshold value
      logs at or above this threshold go to stderr
  -use-ingress-class-only
      If true, ingress resource will handled by ingress controller with class which specifed
      by value of ingress-class. Default false
  -v value
      log level for V logs
  -version
      Print the version and git-commit hash and exit.
  -vmodule value
      comma-separated list of pattern=N settings for file-filtered logging
  -watch-namespace string
      Namespace to watch for Ingress/Services/Endpoints. By default the controller watches acrosss all namespaces
```