Dim WinScriptHost
Set WinScriptHost = CreateObject("WScript.Shell")
WinScriptHost.Run Chr(34) & "rtmp2flv.exe" & Chr(34), 0
Set WinScriptHost = Nothing