from ubuntu
maintainer Rodrigo Chacon <rochacon@gmail.com>
add s3stub /usr/local/bin/s3stub
run mkdir -p /srv/s3stub
cmd /usr/local/bin/s3stub -r /srv/s3stub -b :80
expose 80
