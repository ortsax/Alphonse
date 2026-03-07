# ── Stage 1: download the pre-compiled release binary ─────────────────────────
FROM alpine:3 AS fetch

ARG VERSION=0.0.1
# TARGETARCH is set automatically by Docker BuildKit (amd64, arm64).
# Falls back to amd64 for plain `docker build` without --platform.
ARG TARGETARCH=amd64

RUN apk add --no-cache curl && \
    curl -fsSL \
      "https://github.com/ortsax/Alphonse/releases/download/v${VERSION}/alphonse_${VERSION}_linux_${TARGETARCH}.tar.gz" \
      | tar -xz -C /usr/local/bin/ && \
    chmod +x /usr/local/bin/alphonse

# ── Stage 2: lean runtime image ───────────────────────────────────────────────
FROM alpine:3

# ffmpeg  — media commands (mp3 extraction, video trim, etc.)
# ca-certificates — HTTPS for WhatsApp servers and AI integrations
# tzdata  — timezone-aware timestamps
RUN apk add --no-cache \
      ca-certificates \
      ffmpeg \
      tzdata

COPY --from=fetch /usr/local/bin/alphonse /usr/local/bin/alphonse

# /data holds database.db and .env — mount a volume here for persistence.
WORKDIR /data

ENTRYPOINT ["alphonse"]
