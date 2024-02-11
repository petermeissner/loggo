Remove-Item  .\dist -Recurse -Force -ErrorAction SilentlyContinue
New-Item -Path .\dist -ItemType Directory


go build .\src\loggo && Move-Item .\loggo.exe .\dist\loggo
go build .\src\loggo-serve && Move-Item .\loggo-serve.exe .\dist\loggo
go build .\src\loggo-connect && Move-Item .\loggo-connect.exe .\dist\loggo
