FROM maven:3-openjdk-18-slim
ARG PLANTUML_VERSION

RUN apt-get update && apt-get install -y --no-install-recommends wget graphviz fonts-symbola fonts-wqy-zenhei && rm -rf /var/lib/apt/lists/*
RUN wget -O /plantuml.jar http://sourceforge.net/projects/plantuml/files/plantuml.${PLANTUML_VERSION}.jar/download

ENTRYPOINT [ "java", "-XX:-UsePerfData", "-jar", "/plantuml.jar" ]
