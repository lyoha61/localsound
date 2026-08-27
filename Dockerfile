FROM eclipse-temurin:25-jdk AS build

WORKDIR /app

COPY gradlew .
COPY gradle gradle
COPY build.gradle settings.gradle ./

RUN chmod +x gradlew
RUN ls -la /app
RUN head -n 1 gradlew
RUN ./gradlew dependencies --no-daemon

COPY src src

RUN ./gradlew bootJar --no-daemon


FROM eclipse-temurin:25-jre

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       curl \
       ffmpeg \
			 nodejs \
    && rm -rf /var/lib/apt/lists/*

ARG TARGETARCH

RUN if [ "$TARGETARCH" = "arm64" ]; then \
        curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux_aarch64 \
        -o /usr/local/bin/yt-dlp; \
    elif [ "$TARGETARCH" = "amd64" ]; then \
        curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux \
        -o /usr/local/bin/yt-dlp; \
    else \
        echo "Unsupported architecture: $TARGETARCH" && exit 1; \
    fi \
    && chmod +x /usr/local/bin/yt-dlp

COPY --from=build /app/build/libs/*.jar app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]