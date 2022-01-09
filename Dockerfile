FROM alpine

COPY valkontroller /valkontroller

ENTRYPOINT [ "./valkontroller" ]