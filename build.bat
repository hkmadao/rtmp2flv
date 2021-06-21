@echo off
chcp 65001
set /p ver=请输入版本：  
echo 版本：%ver% 打包开始
@REM window_amd64
@REM SET GOOS=windows
@REM SET GOARCH=amd64
SET CGO_ENABLED=1
go build -o rtmp2flv_%ver%_window_amd64.exe main.go

@REM window_amd64
@REM rmdir /S /Q .\output\rtmp2flv_%ver%_window_amd64

@REM md .\output\rtmp2flv_%ver%_window_amd64\output\live
@REM md .\output\rtmp2flv_%ver%_window_amd64\output\log
@REM md .\output\rtmp2flv_%ver%_window_amd64\conf

@REM xcopy /S /Y /E .\static .\output\rtmp2flv_%ver%_window_amd64\static\
@REM xcopy /S /Y /E .\db .\output\rtmp2flv_%ver%_window_amd64\db\
@REM xcopy .\conf\conf.yml .\output\rtmp2flv_%ver%_window_amd64\conf
@REM xcopy .\rtmp2flv_%ver%_window_amd64.exe .\output\rtmp2flv_%ver%_window_amd64\rtmp2flv.exe
@REM xcopy .\start.vbs .\output\rtmp2flv_%ver%_window_amd64\start.vbs

@REM linux_amd64
@REM rmdir /S /Q .\output\rtmp2flv_%ver%_linux_amd64

@REM md .\output\rtmp2flv_%ver%_linux_amd64\output\live
@REM md .\output\rtmp2flv_%ver%_linux_amd64\output\log
@REM md .\output\rtmp2flv_%ver%_linux_amd64\conf

@REM xcopy /S /Y /E .\static .\output\rtmp2flv_%ver%_linux_amd64\static\
@REM xcopy /S /Y /E .\db .\output\rtmp2flv_%ver%_linux_amd64\db\
@REM xcopy .\conf\conf.yml .\output\rtmp2flv_%ver%_linux_amd64\conf
@REM xcopy .\rtmp2flv_%ver%_linux_amd64 .\output\rtmp2flv_%ver%_linux_amd64\rtmp2flv

@REM linux_armv6
@REM rmdir /S /Q .\output\rtmp2flv_%ver%_linux_armv6

@REM md .\output\rtmp2flv_%ver%_linux_armv6\output\live
@REM md .\output\rtmp2flv_%ver%_linux_armv6\output\log
@REM md .\output\rtmp2flv_%ver%_linux_armv6\conf

@REM xcopy /S /Y /E .\static .\output\rtmp2flv_%ver%_linux_armv6\static\
@REM xcopy /S /Y /E .\db .\output\rtmp2flv_%ver%_linux_armv6\db\
@REM xcopy .\conf\conf.yml .\output\rtmp2flv_%ver%_linux_armv6\conf
@REM xcopy .\rtmp2flv_%ver%_linux_armv6 .\output\rtmp2flv_%ver%_linux_armv6\rtmp2flv

pause