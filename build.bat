del /q parser\*

antlr4.exe -Dlanguage=Go Go.g4 -package parser -o parser
   
go build -o solution.exe