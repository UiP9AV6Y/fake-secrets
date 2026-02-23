CHANGE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
CHANGE_ID ?= $(shell date +%s)
CHANGE_MODULE ?= $(shell $(GO) list -m)
CHANGE_AUTHOR_BASE_URL ?= $(CHANGE_MODULE)/
CHANGE_ISSUE_BASE_URL ?= $(CHANGE_MODULE)/issues/

$(CHANGE_DIR)/%.properties:
	@echo 'CHANGE="Describe your change here"' > $@
	@echo 'ISSUE=""' >> $@
	@echo 'AUTHOR=""' >> $@
	@echo 'BREAKING="false"' >> $@
	@echo "Edit $@ based on the Common Changelog conventions (https://common-changelog.org/#24-change-group)"

.PHONY: change-changed
change-changed: $(CHANGE_DIR)/changed-$(CHANGE_ID).properties

.PHONY: change-added
change-added: $(CHANGE_DIR)/added-$(CHANGE_ID).properties

.PHONY: change-removed
change-removed: $(CHANGE_DIR)/removed-$(CHANGE_ID).properties

.PHONY: change-fixed
change-fixed: $(CHANGE_DIR)/fixed-$(CHANGE_ID).properties

$(CHANGE_DIR)/CHANGELOG-v.md.tmp: CHANGE_VERSION = 0.0.0
$(CHANGE_DIR)/CHANGELOG-v.md.tmp: CHANGE_DATE = $(shell date +%F)
$(CHANGE_DIR)/CHANGELOG-v.md.tmp: $(CHANGE_DIR)/CHANGELOG-changed.md.tmp $(CHANGE_DIR)/CHANGELOG-added.md.tmp $(CHANGE_DIR)/CHANGELOG-removed.md.tmp $(CHANGE_DIR)/CHANGELOG-fixed.md.tmp
	printf "## $(CHANGE_VERSION) - $(CHANGE_DATE)\n" > $@
	cat $^ >> $@
	printf "\n" >> $@

$(CHANGE_DIR)/CHANGELOG-changed.md.tmp: $(wildcard $(CHANGE_DIR)/changed-*.properties)
	CHANGE_AUTHOR_BASE_URL=$(CHANGE_AUTHOR_BASE_URL) \
	  CHANGE_ISSUE_BASE_URL=$(CHANGE_ISSUE_BASE_URL) \
	  CHANGE_DIR=$(CHANGE_DIR) \
	  $(CHANGE_DIR)/common-change.sh Changed $^ \
	> $@

$(CHANGE_DIR)/CHANGELOG-added.md.tmp: $(wildcard $(CHANGE_DIR)/added-*.properties)
	CHANGE_AUTHOR_BASE_URL=$(CHANGE_AUTHOR_BASE_URL) \
	  CHANGE_ISSUE_BASE_URL=$(CHANGE_ISSUE_BASE_URL) \
	  CHANGE_DIR=$(CHANGE_DIR) \
	  $(CHANGE_DIR)/common-change.sh Added $^ \
	> $@

$(CHANGE_DIR)/CHANGELOG-removed.md.tmp: $(wildcard $(CHANGE_DIR)/removed-*.properties)
	CHANGE_AUTHOR_BASE_URL=$(CHANGE_AUTHOR_BASE_URL) \
	  CHANGE_ISSUE_BASE_URL=$(CHANGE_ISSUE_BASE_URL) \
	  CHANGE_DIR=$(CHANGE_DIR) \
	  $(CHANGE_DIR)/common-change.sh Removed $^ \
	> $@

$(CHANGE_DIR)/CHANGELOG-fixed.md.tmp: $(wildcard $(CHANGE_DIR)/fixed-*.properties)
	CHANGE_AUTHOR_BASE_URL=$(CHANGE_AUTHOR_BASE_URL) \
	  CHANGE_ISSUE_BASE_URL=$(CHANGE_ISSUE_BASE_URL) \
	  CHANGE_DIR=$(CHANGE_DIR) \
	  $(CHANGE_DIR)/common-change.sh Fixed $^ \
	> $@

.PHONY: clean-changelog
clean-changelog:
	$(RM) $(CHANGE_DIR)/CHANGELOG-*

.PHONY: clean-changes
clean-changes:
	$(RM) $(CHANGE_DIR)/changed-*.properties
	$(RM) $(CHANGE_DIR)/added-*.properties
	$(RM) $(CHANGE_DIR)/removed-*.properties
	$(RM) $(CHANGE_DIR)/fixed-*.properties
