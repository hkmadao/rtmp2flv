@REM Windows
SET GOOS=windows
SET GOARCH=amd64
SET CGO_ENABLED=1
go build -o rtmp2flv.exe main.go

rmdir /S /Q .\output\releases

md .\output\releases\output\live
md .\output\releases\output\log
md .\output\releases\conf

xcopy /S /Y /E .\static .\output\releases\static\
xcopy /S /Y /E .\db .\output\releases\db\
xcopy .\conf\conf.yml .\output\releases\conf
xcopy .\rtmp2flv .\output\releases
xcopy .\rtmp2flv.exe .\output\releases
xcopy .\start.vbs .\output\releases

pause