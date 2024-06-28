docker build -t nginx-file-server .

docker run -d -p 6080:80 -v C:\hiar_face\registe_path:/usr/share/nginx/files:ro --name nginx-file-server nginx-file-server
docker run -d -p 6080:80 -v ../libs/data:/usr/share/nginx/files --name nginx-file-server nginx-file-server

docker run -d –-restart=always -p 192.168.10.51:6081:6080 -p 192.168.10.51:80:80 -v C:\hiar_face\registe_path:/usr/share/nginx/files --name file-server nginx-file-server
