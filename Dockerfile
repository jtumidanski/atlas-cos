FROM maven:3.6.3-openjdk-14-slim AS build

COPY settings.xml /usr/share/maven/conf/

COPY pom.xml pom.xml
COPY cos-api/pom.xml cos-api/pom.xml
COPY cos-model/pom.xml cos-model/pom.xml
COPY cos-base/pom.xml cos-base/pom.xml
COPY cos-database/pom.xml cos-database/pom.xml

RUN mvn dependency:go-offline package -B

COPY cos-api/src cos-api/src
COPY cos-model/src cos-model/src
COPY cos-base/src cos-base/src
COPY cos-database/src cos-database/src

RUN mvn install

FROM openjdk:14-ea-jdk-alpine
USER root

RUN mkdir service

COPY --from=build /cos-base/target/ /service/

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait

RUN chmod +x /wait

ENV JAVA_TOOL_OPTIONS -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005

EXPOSE 5005

CMD /wait && java --enable-preview -jar /service/cos-base-1.0-SNAPSHOT.jar -Xdebug