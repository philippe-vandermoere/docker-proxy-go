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

    include /etc/nginx/ssl.conf;
    ssl_certificate {{ .Certificate.GetFileName }};
    ssl_certificate_key {{ .Certificate.GetPrivateKeyFileName }};
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
