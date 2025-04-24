FROM ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -y vim wget net-tools iputils-ping jq python3 tree python3-pip curl tzdata && \ 
ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN pip3 install boto3

COPY ./script/start.sh /feature-flag/
COPY ./build/feature-flag* /feature-flag/
