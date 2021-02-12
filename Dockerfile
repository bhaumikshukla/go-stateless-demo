FROM golang:alpine AS builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM alpine AS final
WORKDIR /
COPY --from=builder /workspace/fibr .
CMD [ "./fibr" ]