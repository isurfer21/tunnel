@echo OFF
echo Cross build for MacOSX, Windows and Linux as 32bit and 64bit
rmdir /s /q bin
mkdir bin
call :xc darwin 386 tunnel
call :xc darwin amd64 tunnel
call :xc windows 386 tunnel .exe
call :xc windows amd64 tunnel .exe
call :xc linux 386 tunnel
call :xc linux amd64 tunnel
echo Done!
goto end

:xc
	echo %3_%1_%2
	set GOOS=%1
	set GOARCH=%2
	go build -o bin\%3_%1_%2%4 %3.go
	goto end

:end
	exit /b
