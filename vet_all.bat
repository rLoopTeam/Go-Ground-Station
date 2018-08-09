@ECHO OFF
ECHO This script will vet all packages in this project
ECHO Vetting main package
go vet rloop/Go-Ground-Station-1/ 2> vet_results.txt
CALL :CheckError
ECHO Vetting constants package
go vet rloop/Go-Ground-Station-1/constants 2>> vet_results.txt
CALL :CheckError
ECHO Vetting datastore package
go vet rloop/Go-Ground-Station-1/datastore 2>> vet_results.txt
CALL :CheckError
ECHO Vetting gsgrpc package
go vet rloop/Go-Ground-Station-1/gsgrpc 2>> vet_results.txt
CALL :CheckError
ECHO Vetting gstypes package
go vet rloop/Go-Ground-Station-1/gstypes 2>> vet_results.txt
CALL :CheckError
ECHO Vetting helpers package
go vet rloop/Go-Ground-Station-1/helpers 2>> vet_results.txt
CALL :CheckError
ECHO Vetting logging package
go vet rloop/Go-Ground-Station-1/logging 2>> vet_results.txt
CALL :CheckError
ECHO Vetting parsing package
go vet rloop/Go-Ground-Station-1/parsing 2>> vet_results.txt
CALL :CheckError
ECHO Vetting server package
go vet rloop/Go-Ground-Station-1/server 2>> vet_results.txt
CALL :CheckError

EXIT /B %ERRORLEVEL%

:CheckError
if %errorlevel% equ 0 ( echo No warnings found ) else ( echo warnings found and written to file )