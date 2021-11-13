FROM centos
WORKDIR /server
COPY . .
EXPOSE 8080
CMD ./envelope_rain_group10