FROM scratch

ADD api /

EXPOSE 3000
CMD ["/api"]