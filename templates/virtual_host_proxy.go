package templates

const virtualHostProxy = `
{{- range $path := .GetPaths -}}
upstream {{ $.GetUpstream $path }} {
    {{- range $server := $.GetServers $path }}
    server {{ $server.Ip }}:{{ $server.Port }};
    {{- end }}
}

{{ end -}}

server {
    listen 80;
    listen [::]:80;

    server_name {{ .Domain }};
{{ if .IsHttps }}
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;

    server_name {{ .Domain }};

    ssl_certificate {{ .Certificate.GetFileName }};
    ssl_certificate_key {{ .Certificate.GetPrivateKeyFileName }};

    # https://ssl-config.mozilla.org/#server=nginx&server-version=1.17.0&config=modern
    ssl_certificate {{ .Certificate.GetFileName }};
    ssl_certificate_key {{ .Certificate.GetPrivateKeyFileName }};
{{- if .Certificate.HasChain }}
    ssl_trusted_certificate {{ .Certificate.GetChainFileName }};

    # OCSP stapling
    ssl_stapling on;
    ssl_stapling_verify on;
{{- end }}

    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:10m;
    ssl_session_tickets off;

    # modern configuration
    ssl_protocols TLSv1.3;
    ssl_prefer_server_ciphers off;

    # HSTS (ngx_http_headers_module is required) (63072000 seconds)
    add_header Strict-Transport-Security "max-age=63072000" always;

    resolver 1.1.1.1 1.0.0.1 valid=300s;
    resolver_timeout 5s;


{{- end }}
{{- range $path := .GetPaths }}

    location {{ $path }} {
        proxy_pass http://{{ $.GetUpstream $path }};
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_uid;
        proxy_read_timeout 900;
    }
{{- end }}
}
`
