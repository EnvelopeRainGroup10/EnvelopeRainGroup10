FROM centos@sha256:dead07b4d8ed7e29e98de0f4504d87e8880d4347859d839686a31da35a3b532f
COPY main /root/server
EXPOSE 8080
CMD /root/server