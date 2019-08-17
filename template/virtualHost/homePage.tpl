server {
    listen 80 default_server;
    listen [::]:80 default_server;

    location / {
        root {{ .documentRoot }};
        index {{ .index }};
    }
}
