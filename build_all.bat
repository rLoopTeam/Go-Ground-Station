SET OSLIST=darwin linux windows
SET ARCHLIST=386 amd64
FOR %%a in (%OSLIST%) do (
    SET GOOS=%%a
    FOR %%b in (%ARCHLIST% %%a) do (
        SET GOARCH=%%b
        ECHO go-groundstation-%GOOS%-%GOARCH%
    )
)
