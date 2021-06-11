@REM Windows
SET GOOS=windows
SET GOARCH=amd64
SET CGO_ENABLED=1
go build -o rtmp2flv.exe main.go

rmdir /S /Q .\output\rtmp2flv

md .\output\rtmp2flv\output\live
md .\output\rtmp2flv\output\log
md .\output\rtmp2flv\conf

xcopy /S /Y /E .\static .\output\rtmp2flv\static\
xcopy /S /Y /E .\db .\output\rtmp2flv\db\
xcopy .\conf\conf.yml .\output\rtmp2flv\conf
xcopy .\rtmp2flv .\output\rtmp2flv
xcopy .\rtmp2flv.exe .\output\rtmp2flv
xcopy .\start.vbs .\output\rtmp2flv

pause