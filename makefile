.PHONY: modify-import clean compile

MAIN_BEFORE := main-to-modify
MAIN_AFTER := main.go

# creates a main with the desired active response
build:
	@if [ -z "$(IMPORT)" || -z $(OUTPUT)]; then \
		echo "Usage: make modify-import IMPORT=path/to/new-import OUTPUT=binary_name"; \
		exit 1; \
	fi
	# Copier MAIN_BEFORE dans MAIN_AFTER
	@cp $(MAIN_BEFORE) $(MAIN_AFTER)
	# Remplacer l'import par la valeur de IMPORT
	@sed -i "s|AR \"active-response/ACTIVE_RESPONSE_IMPORT\"|AR \"active-response/$(IMPORT)\"|g" $(MAIN_AFTER)
	@echo "Import modifié : AR \"active-response/$(IMPORT)\""

	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT).exe
	@echo "Compilation terminée pour Windows (amd64)."

clean:
	@if [ -f $(MAIN_AFTER) ]; then \
		rm $(MAIN_AFTER); \
		echo "$(MAIN_AFTER) supprimé."; \
	else \
		echo "$(MAIN_AFTER) n'existe pas."; \
	fi
	@rm -f myprogram.exe
	@echo "Fichiers temporaires supprimés."
	
