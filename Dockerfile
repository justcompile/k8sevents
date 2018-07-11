FROM scratch
ADD confs/cfg.json /
ADD dist/linux/k8sevents-linux-amd64 /
CMD ["/k8sevents-linux-amd64", "-config", "/cfg.json"]
