ifeq ($(DETECTED_OS),Windows)
	BIN := nty.exe
else
	BIN := nty
endif

install:
	@echo Installing...
ifeq ($(OS),Windows_NT)
	@go build -o nty.exe ./main.go
	@if not exist "%USERPROFILE%\bin\.nty" mkdir "%USERPROFILE%\bin\.nty"
	@copy nty.exe "%USERPROFILE%\bin\.nty\nty.exe"
	@setx PATH "%PATH%;%USERPROFILE%\bin\.nty"
	@echo Add complete (restart terminal)
else
	@go build -o nty ./main.go
	@cp nty ~/.local/bin/nty
endif
	@echo Done!