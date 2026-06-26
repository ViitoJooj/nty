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
	@powershell -NoProfile -Command "$$d=\"$$env:USERPROFILE\bin\.nty\"; $$p=[Environment]::GetEnvironmentVariable('PATH','User'); if([string]::IsNullOrEmpty($$p)){[Environment]::SetEnvironmentVariable('PATH',$$d,'User')}elseif($$p -notlike \"*$$d*\"){[Environment]::SetEnvironmentVariable('PATH',($$p.TrimEnd(';')+';'+$$d),'User')}"
	@echo Add complete (restart terminal)
else
	@go build -o nty ./main.go
	@cp nty ~/.local/bin/nty
endif
	@echo Done!