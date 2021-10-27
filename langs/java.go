package langs

func GetPOM() string {
	return pomFile;
}

func GetJava() string {
	return helloJavaSrcBoilerplate;
}

func GetDocker() string {
	return dockerFile;
}

const (
	dockerFile = `FROM openjdk:8-jdk-alpine
ARG JAR_FILE=target/%s-0.0.1-SNAPSHOT.jar
COPY ${JAR_FILE} app.jar
ENTRYPOINT ["java","-cp","/app.jar","com.naga.%s.Main"]
	`
	
	pomFile = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion> 
  <groupId>com.naga</groupId>
  <artifactId>%s</artifactId>
  <version>0.0.1-SNAPSHOT</version>
  <properties>
    <java.version>1.8</java.version>
  </properties> 
  <dependencies>
		 <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>4.11</version>
            <scope>test</scope>
        </dependency>	
  </dependencies>
  <build>
    <plugins>
      <plugin>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>
`

	helloJavaSrcBoilerplate = `package com.naga.%s;

import java.util.Scanner;

public class Main {

	
	private static String callFunction(String payload) {
		// Replace it with your own implementation
		return payload;
	}
	
	public static void main(String args[]) throws Exception {
		while(true) {
		 	try (Scanner scanner = new Scanner(System.in)) {
				while (scanner.hasNext()) {
				    String payload = scanner.nextLine();
				    System.out.println(callFunction(payload));
				}
				scanner.close();				
			}		
		}
	}
}
`
)