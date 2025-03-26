ROOT_DIR_RELATIVE := ..
include $(ROOT_DIR_RELATIVE)/common.mk

DIAGRAM_SRCS := $(call rwildcard,.,*.md,*.plantuml)

# 生成 PlantUML 图表
diagrams: $(TOOLS_BIN_DIR)/plantuml-sentinal
    -docker run -u $(UID):$(GID) -v $(abspath .):/docs$(DOCKER_VOL_OPTS) plantuml-builder -tsvg /docs/**/*.md
    -docker run -u $(UID):$(GID) -v $(abspath .):/docs$(DOCKER_VOL_OPTS) plantuml-builder -tsvg /docs/**/*.plantuml
