@echo off
chcp 65001
set /p ver=请输入版本：  
echo 版本：%ver% 打包开始
platforms="windows_amd64 linux_amd64 linux_arm"
for %%platform in (%platforms%) do (
    rmdir /S /Q .\resources\output\rtmp2flv_%ver%_window_amd64
    SET CGO_ENABLED=0
    if %platform%=="windows_amd64" (
        SET GOOS=windows
        SET GOARCH=amd64
        go build -o .\resources\output\rtmp2flv_%ver%_window_amd64\rtmp2flv.exe main.go
    )else if %platform%=="linux_amd64" (
        SET GOOS=linux
        SET GOARCH=amd64
        go build -o .\resources\output\rtmp2flv_%ver%_window_amd64\rtmp2flv main.go
    )else if %platform%=="linux_arm" (
        SET GOOS=linux
        SET GOARCH=arm
        go build -o .\resources\output\rtmp2flv_%ver%_window_amd64\rtmp2flv main.go
    )

    md .\resources\output\rtmp2flv_%ver%_window_amd64\output\live
    md .\resources\output\rtmp2flv_%ver%_window_amd64\output\log
    md .\resources\output\rtmp2flv_%ver%_window_amd64\conf

    xcopy /S /Y /E .\static .\resources\output\rtmp2flv_%ver%_window_amd64\static\
    xcopy /S /Y /E .\resources\db .\resources\output\rtmp2flv_%ver%_window_amd64\db\
    xcopy .\resources\conf .\resources\output\rtmp2flv_%ver%_window_amd64\conf

    cd .\resources\output\

    7z a -ttar -so rtmp2flv_%ver%_window_amd64.tar rtmp2flv_%ver%_window_amd64/ | 7z a -si rtmp2flv_%ver%_window_amd64.tar.gz
)

pause