FROM maven:3.6.3-openjdk-17

WORKDIR /app

COPY pom.xml .

COPY src ./src

RUN mvn clean install

EXPOSE 8082

ENTRYPOINT ["java", "-jar", "target/cinemaabyss-events-service-1.0-SNAPSHOT.jar"]