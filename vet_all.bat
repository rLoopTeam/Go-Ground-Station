@ECHO OFF
ECHO This script will vet all packages in this project
ECHO Vetting main package
go vet rloop/Go-Ground-Station/
IF ERRORLEVEL 0 (
   ECHO No warnings found
)
ECHO Vetting constants package
go vet rloop/Go-Ground-Station/constants
ECHO Vetting datastore package
go vet rloop/Go-Ground-Station/datastore
ECHO Vetting gsgrpc package
go vet rloop/Go-Ground-Station/gsgrpc
ECHO Vetting gstypes package
go vet rloop/Go-Ground-Station/gstypes

